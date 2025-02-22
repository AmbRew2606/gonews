package mongo

import (
	"context"
	"fmt"
	"os"
	"time"

	"GoNews/pkg/storage"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Store представляет собой хранилище MongoDB.
type Store struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// New создаёт новое подключение к MongoDB.
// func New(uri string) (*Store, error) {
// 	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
// 	if err != nil {
// 		return nil, fmt.Errorf("ошибка создания клиента MongoDB: %w", err)
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	if err := client.Connect(ctx); err != nil {
// 		return nil, fmt.Errorf("ошибка подключения к MongoDB: %w", err)
// 	}

//		collection := client.Database("gonews").Collection("posts")
//		return &Store{client: client, collection: collection}, nil
//	}
func New() (*Store, error) {
	// Загружаем переменные окружения
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить .env файл: %w", err)
	}

	uri := os.Getenv("MONGO_URI") // Получаем строку подключения из .env

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания клиента MongoDB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		return nil, fmt.Errorf("ошибка подключения к MongoDB: %w", err)
	}

	collection := client.Database("gonews").Collection("posts")
	return &Store{client: client, collection: collection}, nil
}

// Получение всех публикаций
func (s *Store) Posts() ([]storage.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []storage.Post
	for cursor.Next(ctx) {
		var p storage.Post
		if err := cursor.Decode(&p); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

// Добавление публикации
func (s *Store) AddPost(p storage.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.collection.InsertOne(ctx, p)
	return err
}

// Обновление публикации
func (s *Store) UpdatePost(p storage.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": p.ID}
	update := bson.M{"$set": p}

	_, err := s.collection.UpdateOne(ctx, filter, update)
	return err
}

// Удаление публикации
func (s *Store) DeletePost(p storage.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.collection.DeleteOne(ctx, bson.M{"id": p.ID})
	return err
}
