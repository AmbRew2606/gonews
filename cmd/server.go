package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/storage"
	"GoNews/pkg/storage/postgres"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

// Сервер GoNews.
type server struct {
	db  storage.Interface
	api *api.API
}

func main() {

	// Используем БД в памяти (для тестов)
	// srv.db = memdb.New()

	// Загружаем переменные окружения из .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Не удалось загрузить .env файл, используем переменные окружения")
	}

	var srv server

	// Подключение к PostgreSQL
	db, err := postgres.New() // Теперь New() не принимает строку, а берёт данные из os.Getenv()
	if err != nil {
		log.Fatal("Ошибка подключения к PostgreSQL:", err)
	}
	srv.db = db

	// Подключение к MongoDB (альтернативный вариант)
	/*
		db3, err := mongo.New("mongodb://localhost:27017/")
		if err != nil {
			log.Fatal(err)
		}
		srv.db = db3
	*/
	////////////////////////
	// mongoURI := os.Getenv("MONGO_URI") // Получаем строку подключения из .env
	// db3, err := mongo.New(mongoURI) // Используем переменную окружения для подключения
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// srv.db = db3

	// Запуск API
	srv.api = api.New(srv.db)
	log.Println("Сервер запущен на :8080")
	http.ListenAndServe(":8080", srv.api.Router())
}
