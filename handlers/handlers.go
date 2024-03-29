package handlers

import (
	"../models"
	"../services"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CrudHandler interface {
	GetAllBooks(w http.ResponseWriter, r *http.Request)
	GetBook(w http.ResponseWriter, r *http.Request)
	AddBook(w http.ResponseWriter, r *http.Request)
	EditBook(w http.ResponseWriter, r *http.Request)
	DeleteBook(w http.ResponseWriter, r *http.Request)
}

type CrudHandlerImpl struct {
	storageService services.Storage
}

func InitHandlers(storage services.Storage) CrudHandler {
	return &CrudHandlerImpl{storageService: storage}
}

// GET /books
func (c *CrudHandlerImpl) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := c.storageService.GetAllBooks()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// GET /books/{id}
func (c *CrudHandlerImpl) GetBook(w http.ResponseWriter, r *http.Request) {
	id, isExist := mux.Vars(r)["id"]
	if !isExist {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	castedID, _ := strconv.Atoi(id)
	bookFromStorage, err := c.storageService.GetBookById(castedID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookFromStorage)
	return
}

// POST /books
func (c *CrudHandlerImpl) AddBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := strconv.Itoa(rand.Intn(1000000))
	book.ID = id
	err = c.storageService.AddBook(&book)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// PUT /book/{id}
func (c *CrudHandlerImpl) EditBook(w http.ResponseWriter, r *http.Request) {
	id, isExist := mux.Vars(r)["id"]
	if !isExist {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	book.ID = id
	err = c.storageService.UpdateBook(&book)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DELETE /book/{id}
func (c *CrudHandlerImpl) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, isExist := mux.Vars(r)["id"]
	if !isExist {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	castedID, _ := strconv.Atoi(id)
	err := c.storageService.RemoveBook(castedID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
