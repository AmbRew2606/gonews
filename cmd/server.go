package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/storage"
	"GoNews/pkg/storage/postgres"
	"log"
	"net/http"
)

// Сервер GoNews.
type server struct {
	db  storage.Interface
	api *api.API
}

func main() {
	var srv server

	// Используем БД в памяти (для тестов)
	// srv.db = memdb.New()

	// Подключение к PostgreSQL
	db2, err := postgres.New("postgres://postgres:password@localhost:5432/gonews?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	srv.db = db2

	// Подключение к MongoDB (альтернативный вариант)
	/*
		db3, err := mongo.New("mongodb://localhost:27017/")
		if err != nil {
			log.Fatal(err)
		}
		srv.db = db3
	*/

	// Запуск API
	srv.api = api.New(srv.db)
	log.Println("Сервер запущен на :8080")
	http.ListenAndServe(":8080", srv.api.Router())
}
