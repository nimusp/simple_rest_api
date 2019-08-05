package main

import (
	"./models"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var books []models.Book

func main() {
	books = append(books, models.Book{
		ID:    "7",
		Title: "book title",
		Author: &models.Author{
			Firstname: "first",
			Lastname:  "second",
		},
	})

	books = append(books, models.Book{
		ID:    "9",
		Title: "another book title",
		Author: &models.Author{
			Firstname: "f",
			Lastname:  "s",
		},
	})

	router := mux.NewRouter()
	router.HandleFunc("/book/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", getAllBooks).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/book/{id}", editBook).Methods("PUT")
	router.HandleFunc("/book/{id}", deleteBook).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// GET /books
func getAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// GET /books/{id}
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, ok := mux.Vars(r)["id"]
	if !ok {
		json.NewEncoder(w).Encode(books)
		return
	}

	isFounded := false
	for _, item := range books {
		if item.ID == id {
			err := json.NewEncoder(w).Encode(item)
			isFounded = true
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
	}

	if !isFounded {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(&models.Book{})
}

// POST /books
func addBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book.ID = strconv.Itoa(rand.Intn(1000000))
	books = append(books, book)
	w.WriteHeader(http.StatusCreated)
}

// PUT /book/{id}
func editBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id := mux.Vars(r)["id"]
	isUpdated := false
	for index, item := range books {
		if item.ID == id {
			book.ID = id
			books[index] = book
			isUpdated = true
			break
		}
	}

	if !isUpdated {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DELETE /book/{id}
func deleteBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id := mux.Vars(r)["id"]
	isFinded := false
	for index, item := range books {
		if item.ID == id {
			books = append(books[:index], books[index+1:]...)
			isFinded = true
			break
		}
	}

	if !isFinded {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
