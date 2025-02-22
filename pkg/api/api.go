package api

import (
	"GoNews/pkg/storage"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Программный интерфейс сервера GoNews
type API struct {
	db     storage.Interface
	router *mux.Router
}

// Конструктор объекта API
func New(db storage.Interface) *API {
	api := API{
		db: db,
	}
	api.router = mux.NewRouter()
	api.endpoints()
	return &api
}

// Регистрация обработчиков API.
func (api *API) endpoints() {
	api.router.HandleFunc("/posts/all", api.postsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/posts", api.addPostHandler).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/posts/{id}", api.updatePostHandler).Methods(http.MethodPut, http.MethodOptions)
	api.router.HandleFunc("/posts/{id}", api.deletePostHandler).Methods(http.MethodDelete, http.MethodOptions)
}

// Получение маршрутизатора запросов.
// Требуется для передачи маршрутизатора веб-серверу.
func (api *API) Router() *mux.Router {
	return api.router
}

// Получение всех публикаций.
func (api *API) postsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := api.db.Posts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

// Добавление публикации.
func (api *API) addPostHandler(w http.ResponseWriter, r *http.Request) {
	var p storage.Post
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.db.AddPost(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Обновление публикации.
// VER 1
//
//	func (api *API) updatePostHandler(w http.ResponseWriter, r *http.Request) {
//		var p storage.Post
//		err := json.NewDecoder(r.Body).Decode(&p)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		err = api.db.UpdatePost(p)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		w.WriteHeader(http.StatusOK)
//	}
//
// VER 2
//
//	func (api *API) updatePostHandler(w http.ResponseWriter, r *http.Request) {
//		vars := mux.Vars(r) // Получаем параметры из URL
//		id := vars["id"]    // Получаем ID из URL
//		var p storage.Post
//		err := json.NewDecoder(r.Body).Decode(&p)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		p.ID = id // Присваиваем ID из URL
//		err = api.db.UpdatePost(p)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		w.WriteHeader(http.StatusOK)
//	}
//
// VER 3
func (api *API) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Получаем параметры из URL
	idStr := vars["id"] // Получаем ID как строку

	// Преобразуем строку ID в int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var p storage.Post
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p.ID = id // Присваиваем ID из URL
	err = api.db.UpdatePost(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Удаление публикации.
// VER 1
//
//	func (api *API) deletePostHandler(w http.ResponseWriter, r *http.Request) {
//		var p storage.Post
//		err := json.NewDecoder(r.Body).Decode(&p)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		err = api.db.DeletePost(p)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		w.WriteHeader(http.StatusOK)
//	}
//
// VAR 2
func (api *API) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Получаем параметры из URL
	idStr := vars["id"] // Получаем ID из URL

	// Преобразуем строку в целое число
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	err = api.db.DeletePost(storage.Post{ID: id}) // Передаем ID для удаления
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
