package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"

	"GoNews/pkg/storage"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Store представляет собой хранилище PostgreSQL.
type Store struct {
	db *sql.DB
}

// New - подключение к БД PostgreSQL
func New() (*Store, error) {
	// Читаем переменные окружения (они уже загружены в server.go)
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Проверяем, что все переменные загружены
	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		return nil, fmt.Errorf("не все переменные окружения загружены")
	}

	// Экранируем пароль (если есть спецсимволы)
	escapedPassword := url.QueryEscape(password)

	// Формируем строку подключения
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, escapedPassword, dbname)

	// Открываем соединение с БД
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("Ошибка подключения к PostgreSQL: %v", err)
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// Проверяем подключение
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

// // New создаёт новое подключение к PostgreSQL.
// func New(dsn string) (*Store, error) {
// 	db, err := sql.Open("postgres", dsn)
// 	if err != nil {
// 		return nil, fmt.Errorf("ошибка подключения к PostgreSQL: %w", err)
// 	}

// 	// Проверим подключение
// 	if err := db.Ping(); err != nil {
// 		return nil, fmt.Errorf("ошибка проверки соединения: %w", err)
// 	}

// 	return &Store{db: db}, nil
// }

// Получение всех публикаций
// VER 1
// func (s *Store) Posts() ([]storage.Post, error) {
// 	rows, err := s.db.Query("SELECT id, title, content, author_id, created_at FROM posts")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var posts []storage.Post
// 	for rows.Next() {
// 		var p storage.Post
// 		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.AuthorID, &p.CreatedAt); err != nil {
// 			return nil, err
// 		}
// 		posts = append(posts, p)
// 	}

// 	return posts, nil
// }

// VER 2
// func (db *Store) Posts() ([]storage.Post, error) {

// 	rows, err := s.db.Query("SELECT posts.id, posts.title, posts.content, posts.created_at, authors.name FROM posts JOIN authors ON posts.author_id = authors.id")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var posts []storage.Post
// 	for rows.Next() {
// 		var p storage.Post
// 		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.AuthorName)
// 		if err != nil {
// 			return nil, err
// 		}
// 		posts = append(posts, p)
// 	}

//		return posts, nil
//	}
//
// VER 3
func (s *Store) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query("SELECT posts.id, posts.title, posts.content, posts.created_at, authors.id, authors.name, authors.avatar_url FROM posts JOIN authors ON posts.author_id = authors.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []storage.Post
	for rows.Next() {
		var p storage.Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.AuthorID, &p.AuthorName, &p.AuthorAvatar); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

// Добавление публикации
func (s *Store) AddPost(p storage.Post) error {
	_, err := s.db.Exec("INSERT INTO posts (title, content, author_id, created_at) VALUES ($1, $2, $3, $4)",
		p.Title, p.Content, p.AuthorID, p.CreatedAt)
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

// Добавление публикации
// func (s *Store) AddPost(p storage.Post) error {
// 	_, err := s.db.Exec("INSERT INTO posts (title, content, author_id, author_name, created_at, published_at) VALUES ($1, $2, $3, $4, $5, $6)",
// 		p.Title, p.Content, p.AuthorID, p.AuthorName, p.CreatedAt, p.PublishedAt)
// 	return err
// }

// // Обновление публикации
// func (s *Store) UpdatePost(p storage.Post) error {
// 	_, err := s.db.Exec("UPDATE posts SET title=$1, content=$2, author_id=$3, author_name=$4, created_at=$5, published_at=$6 WHERE id=$7",
// 		p.Title, p.Content, p.AuthorID, p.AuthorName, p.CreatedAt, p.PublishedAt, p.ID)
// 	return err
// }

// // Удаление публикации
// func (s *Store) DeletePost(p storage.Post) error {
// 	_, err := s.db.Exec("DELETE FROM posts WHERE id=$1", p.ID)
// 	return err
// }
