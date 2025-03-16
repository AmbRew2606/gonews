package storage

// Post - публикация.
type Post struct {
	ID            int
	Title         string
	Content       string
	AuthorID      int
	Author        Author
	CreatedAt     int64
	FormattedDate string // Дополнительное поле для вывода в шаблон
	// PublishedAt int64
}

// Author - автор публикаций.
type Author struct {
	ID        int
	Name      string
	AvatarURL string
}

// PageData - структура для передачи данных в шаблоны.
type PageData struct {
	Authors []Author // Список авторов
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	Posts() ([]Post, error) // получение всех публикаций
	AddPost(Post) error     // создание новой публикации
	UpdatePost(Post) error  // обновление публикации
	DeletePost(Post) error  // удаление публикации по ID

	// Новый метод для работы с авторами
	AddAuthor(Author) error            // создание нового автора
	GetAuthorByID(int) (Author, error) // получение автора по ID
	GetAuthors() ([]Author, error)     // получение всех авторов
}
