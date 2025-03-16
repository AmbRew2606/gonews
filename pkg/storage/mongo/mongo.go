package mongo

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"GoNews/pkg/storage"
)

// Store представляет собой хранилище MongoDB
type Store struct {
	client *mongo.Client
	db     *mongo.Database
}

// New - подключение к MongoDB
func New() (*Store, error) {
	// Получение переменных окружения для подключения
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		return nil, fmt.Errorf("MONGO_URI не установлен")
	}

	// Подключаемся к MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Printf("Ошибка подключения к MongoDB: %v", err)
		return nil, fmt.Errorf("ошибка подключения к MongoDB: %w", err)
	}

	// Проверяем соединение
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Printf("Ошибка пинга MongoDB: %v", err)
		return nil, fmt.Errorf("не удалось подключиться к MongoDB: %w", err)
	}

	// Получаем доступ к базе данных
	db := client.Database("GoNews")

	log.Println("Подключение к MongoDB успешно")
	return &Store{client: client, db: db}, nil
}

// Close - закрытие соединения с MongoDB
func (s *Store) Close() {
	if err := s.client.Disconnect(context.Background()); err != nil {
		log.Printf("Ошибка при закрытии соединения с MongoDB: %v", err)
	} else {
		log.Println("Соединение с MongoDB закрыто")
	}
}

// get post
func (s *Store) Posts() ([]storage.Post, error) {
	collection := s.db.Collection("posts")
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var posts []storage.Post
	for cursor.Next(context.Background()) {
		var p storage.Post
		if err := cursor.Decode(&p); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// add post
func (s *Store) AddPost(p storage.Post) error {
	collection := s.db.Collection("posts")
	_, err := collection.InsertOne(context.Background(), p)
	return err
}

// get authors
func (s *Store) GetAuthors() ([]storage.Author, error) {
	collection := s.db.Collection("authors")
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var authors []storage.Author
	for cursor.Next(context.Background()) {
		var a storage.Author
		if err := cursor.Decode(&a); err != nil {
			return nil, err
		}
		authors = append(authors, a)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return authors, nil
}

// add authors
func (s *Store) AddAuthor(a storage.Author) error {
	collection := s.db.Collection("authors")
	_, err := collection.InsertOne(context.Background(), a)
	return err
}
