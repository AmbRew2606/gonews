package postgres

import (
	"database/sql"
	"fmt"

	"GoNews/pkg/storage"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Store представляет собой хранилище PostgreSQL.
type Store struct {
	db *sql.DB
}

// New создаёт новое подключение к PostgreSQL.
func New(dsn string) (*Store, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к PostgreSQL: %w", err)
	}

	// Проверим подключение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка проверки соединения: %w", err)
	}

	return &Store{db: db}, nil
}

// Получение всех публикаций
func (s *Store) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query("SELECT id, title, content, author_id, author_name, created_at, published_at FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []storage.Post
	for rows.Next() {
		var p storage.Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.AuthorID, &p.AuthorName, &p.CreatedAt, &p.PublishedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

// Добавление публикации
func (s *Store) AddPost(p storage.Post) error {
	_, err := s.db.Exec("INSERT INTO posts (title, content, author_id, author_name, created_at, published_at) VALUES ($1, $2, $3, $4, $5, $6)",
		p.Title, p.Content, p.AuthorID, p.AuthorName, p.CreatedAt, p.PublishedAt)
	return err
}

// Обновление публикации
func (s *Store) UpdatePost(p storage.Post) error {
	_, err := s.db.Exec("UPDATE posts SET title=$1, content=$2, author_id=$3, author_name=$4, created_at=$5, published_at=$6 WHERE id=$7",
		p.Title, p.Content, p.AuthorID, p.AuthorName, p.CreatedAt, p.PublishedAt, p.ID)
	return err
}

// Удаление публикации
func (s *Store) DeletePost(p storage.Post) error {
	_, err := s.db.Exec("DELETE FROM posts WHERE id=$1", p.ID)
	return err
}
