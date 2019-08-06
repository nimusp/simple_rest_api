package main

import (
	"./handlers"
	"./services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var storage services.Storage

func main() {
	handler := handlers.InitHandlers()

	router := mux.NewRouter()
	router.HandleFunc("/book/{id}", handler.GetBook).Methods("GET")
	router.HandleFunc("/books", handler.GetAllBooks).Methods("GET")
	router.HandleFunc("/books", handler.AddBook).Methods("POST")
	router.HandleFunc("/book/{id}", handler.EditBook).Methods("PUT")
	router.HandleFunc("/book/{id}", handler.DeleteBook).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
