package api

import (
	"GoNews/pkg/storage"
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
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
	// api.router.HandleFunc("/add-post", api.addPostHandler).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/posts/{id}", api.updatePostHandler).Methods(http.MethodPut, http.MethodOptions)
	api.router.HandleFunc("/posts/{id}", api.deletePostHandler).Methods(http.MethodDelete, http.MethodOptions)

	api.router.HandleFunc("/", api.homeHandler).Methods(http.MethodGet)

	// Обработка статических файлов
	api.router.PathPrefix("/static/").HandlerFunc(api.staticFileHandler())

	//добавление пользователя
	api.router.HandleFunc("/add-user", api.addUserPageHandler).Methods("GET") // Для отображения формы
	api.router.HandleFunc("/add-user", api.addUserHandler).Methods("POST")    // Для обработки формы

	//добавление поста
	api.router.HandleFunc("/add-post", api.addPostPageHandler).Methods("GET") // Для отображения формы
	api.router.HandleFunc("/add-post", api.addPostHandler).Methods("POST")    // Для обработки формы
}

// Обработчик статических файлов
func (api *API) staticFileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Определяем путь к файлу
		filePath := filepath.Join("static", r.URL.Path[len("/static/"):])

		// Проверка существования файла
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.Error(w, "404 not found", http.StatusNotFound)
			return
		}

		// Логирование запроса
		log.Printf("Serving static file: %s", filePath)

		http.ServeFile(w, r, filePath)
	}
}

// обработчик ГЕТ юзер
func (api *API) addUserPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/add_user.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func (api *API) homeHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := api.db.Posts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, posts)
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
//VER1
// func (api *API) addPostHandler(w http.ResponseWriter, r *http.Request) {
// 	var p storage.Post

// 	err := json.NewDecoder(r.Body).Decode(&p)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

//		err = api.db.AddPost(p)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		w.WriteHeader(http.StatusOK)
//	}
//
// Добавление публикации.
func (api *API) addPostHandler(w http.ResponseWriter, r *http.Request) {
	// Парсим данные формы
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Ошибка парсинга формы", http.StatusInternalServerError)
		return
	}

	// Извлекаем значения из формы
	title := r.FormValue("title")
	content := r.FormValue("content")
	authorID := r.FormValue("author_id")

	// Проверяем, что authorID не пустой
	if authorID == "" {
		http.Error(w, "Автор не выбран", http.StatusBadRequest)
		return
	}

	fmt.Println("Заголовок:", title)
	fmt.Println("Контент:", content)
	fmt.Println("Автор:", authorID)

	// Преобразуем authorID в int
	authorIDInt, err := strconv.Atoi(authorID)
	if err != nil {
		http.Error(w, "Неверный формат author_id", http.StatusBadRequest)
		return
	}

	// Создаем новый объект Post
	p := storage.Post{
		Title:     title,
		Content:   content,
		AuthorID:  authorIDInt,       // Используем преобразованный int
		CreatedAt: time.Now().Unix(), // Устанавливаем текущую метку времени
	}

	// Добавляем публикацию в базу данных
	err = api.db.AddPost(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// для html
func (api *API) addPostPageHandler(w http.ResponseWriter, r *http.Request) {
	// Загружаем список авторов
	authors, err := api.db.GetAuthors()
	if err != nil {
		http.Error(w, "Ошибка загрузки авторов", http.StatusInternalServerError)
		return
	}

	// Создаём объект данных для шаблона
	data := storage.PageData{Authors: authors}

	// Парсим и рендерим шаблон
	tmpl, err := template.ParseFiles("templates/add_post.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
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

// обработка фоток
// Добавление пользователя с аватаркой
func (api *API) addUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем имя пользователя
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Имя не может быть пустым", http.StatusBadRequest)
		return
	}

	// Получаем файл аватарки
	file, header, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, "Ошибка загрузки файла", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Определяем расширение файла
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		http.Error(w, "Неподдерживаемый формат (только JPG, PNG)", http.StatusBadRequest)
		return
	}

	// Формируем путь к файлу: av_имя.jpg
	avatarPath := "static/avatars/av_" + name + ext

	// Декодируем изображение
	var img image.Image
	if ext == ".png" {
		img, err = png.Decode(file)
	} else {
		img, err = jpeg.Decode(file)
	}
	if err != nil {
		http.Error(w, "Ошибка декодирования изображения", http.StatusInternalServerError)
		return
	}

	// Сжимаем до 32x32
	img = resize.Resize(32, 32, img, resize.Lanczos3)

	// Создаём файл для сохранения
	outFile, err := os.Create(avatarPath)
	if err != nil {
		http.Error(w, "Ошибка сохранения файла", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	// Кодируем и сохраняем
	if ext == ".png" {
		err = png.Encode(outFile, img)
	} else {
		err = jpeg.Encode(outFile, img, nil)
	}
	if err != nil {
		http.Error(w, "Ошибка кодирования изображения", http.StatusInternalServerError)
		return
	}

	// Добавляем пользователя в БД
	err = api.db.AddAuthor(storage.Author{Name: name, AvatarURL: "/" + avatarPath})
	if err != nil {
		http.Error(w, "Ошибка сохранения пользователя", http.StatusInternalServerError)
		return
	}

	// Успех
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
