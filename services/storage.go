package services

import (
	"../models"
	"database/sql"
	_ "github.com/lib/pq"
)

const driverName = "postgres"

type Storage interface {
	GetAllBooks() ([]*models.Book, error)
	GetBookById(id int) (*models.Book, error)
	AddBook(book *models.Book) error
	UpdateBook(book *models.Book) error
	RemoveBook(id int) error
}

type StorageImpl struct {
	db *sql.DB
}

func InitDB(dbURL string) Storage {
	db, err := sql.Open(driverName, dbURL)
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	return &StorageImpl{db: db}
}

func (s *StorageImpl) GetAllBooks() ([]*models.Book, error) {
	rows, err := s.db.Query("SELECT * FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := make([]*models.Book, 0)
	for rows.Next() {
		var book models.Book
		var author models.Author
		rows.Scan(&book.ID, &book.Title, &author.Firstname, &author.Lastname)
		book.Author = &author
		books = append(books, &book)
	}
	return books, err
}

func (s *StorageImpl) GetBookById(id int) (*models.Book, error) {
	row := s.db.QueryRow(
		`SELECT id, title, author_name, author_last_name
		FROM books
		WHERE id = $1`, id)

	var bookID, title, firstName, lastName string
	err := row.Scan(&bookID, &title, &firstName, &lastName)
	if err != nil {
		return nil, err
	}
	author := &models.Author{Firstname: firstName, Lastname: lastName}
	println(title)
	return &models.Book{
		ID:     bookID,
		Title:  title,
		Author: author,
	}, nil
}

func (s *StorageImpl) AddBook(book *models.Book) error {
	statement, err := s.db.Prepare(
		`INSERT INTO books (title, author_name, author_last_name)
		VALUES ($1, $2, $3)`)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(book.Title, book.Author.Firstname, book.Author.Lastname)
	return err
}

func (s *StorageImpl) UpdateBook(book *models.Book) error {
	statement, err := s.db.Prepare(
		`UPDATE books
		SET title = $1, author_name = $2, author_last_name = $3
		WHERE id = $4`)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(book.Title, book.Author.Firstname, book.Author.Lastname, book.ID)
	return err
}

func (s *StorageImpl) RemoveBook(id int) error {
	statement, err := s.db.Prepare(
		`DELETE FROM books
		WHERE id = $1`)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(id)
	return err
}
