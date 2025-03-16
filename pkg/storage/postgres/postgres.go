package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"GoNews/pkg/storage"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Store представляет собой хранилище PostgreSQL.
type Store struct {
	db *sql.DB
}

// New - подключение к БД PostgreSQL
func New() (*Store, error) {
	// повторная проверка переменных
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Проверка переменных .env
	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		return nil, fmt.Errorf("не все переменные окружения загружены")
	}

	// Экранизация пароля (если требуется)
	escapedPassword := url.QueryEscape(password)

	// Строка подключения
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, escapedPassword, dbname)

	// Открытие соединения с БД
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("Ошибка подключения к PostgreSQL: %v", err)
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		log.Printf("Ошибка пинга PostgreSQL: %v", err)
		return nil, fmt.Errorf("не удалось подключиться к БД: %w", err)
	}

	log.Println("Подключение к PostgreSQL успешно")
	return &Store{db: db}, nil
}

// Close - закрытие соединения с БД
func (s *Store) Close() {
	if err := s.db.Close(); err != nil {
		log.Printf("Ошибка при закрытии соединения с PostgreSQL: %v", err)
	} else {
		log.Println("Соединение с PostgreSQL закрыто")
	}
}

func (s *Store) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query(`
        SELECT posts.id, posts.title, posts.content, posts.created_at, 
               authors.id, authors.name, authors.avatar_url 
        FROM posts 
        JOIN authors ON posts.author_id = authors.id
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []storage.Post
	for rows.Next() {
		var p storage.Post
		var a storage.Author //объект автора
		var createdAtUnix int64

		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &createdAtUnix, &a.ID, &a.Name, &a.AvatarURL); err != nil {
			return nil, err
		}
		// Конвертируем Unix timestamp в строку с форматом даты
		p.CreatedAt = createdAtUnix
		p.FormattedDate = time.Unix(createdAtUnix, 0).Format("02.01.2006 15:04")

		p.Author = a // Присваиваем автора в структуру поста
		posts = append(posts, p)
	}

	return posts, nil
}

// Добавление публикации
func (s *Store) AddPost(p storage.Post) error {
	// p.CreatedAt = time.Now().Unix()
	_, err := s.db.Exec("INSERT INTO posts (title, content, author_id) VALUES ($1, $2, $3)",
		p.Title, p.Content, p.AuthorID)
	return err
}

// Обновление публикации
func (s *Store) UpdatePost(p storage.Post) error {
	_, err := s.db.Exec("UPDATE posts SET title=$1, content=$2, author_id=$3, created_at=$4  WHERE id=$5",
		p.Title, p.Content, p.AuthorID, p.CreatedAt, p.ID)
	return err
}

// Удаление публикации
func (s *Store) DeletePost(p storage.Post) error {
	_, err := s.db.Exec("DELETE FROM posts WHERE id=$1", p.ID)
	return err
}

func (s *Store) AddAuthor(a storage.Author) error {
	_, err := s.db.Exec(`INSERT INTO authors (name, avatar_url) VALUES ($1, $2)`, a.Name, a.AvatarURL)
	return err
}

func (s *Store) GetAuthorByID(id int) (storage.Author, error) {
	var a storage.Author
	err := s.db.QueryRow(`SELECT id, name, avatar_url FROM authors WHERE id = $1`, id).
		Scan(&a.ID, &a.Name, &a.AvatarURL)
	if err != nil {
		return storage.Author{}, err
	}
	return a, nil
}

func (s *Store) GetAuthors() ([]storage.Author, error) {
	rows, err := s.db.Query("SELECT id, name, avatar_url FROM authors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []storage.Author
	for rows.Next() {
		var a storage.Author
		if err := rows.Scan(&a.ID, &a.Name, &a.AvatarURL); err != nil {
			return nil, err
		}
		authors = append(authors, a)
	}

	return authors, nil
}
