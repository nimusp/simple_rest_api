package handlers

import (
	"../models"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
)

type StorageMock struct {
	data map[int]*models.Book
}

func (s *StorageMock) GetAllBooks() ([]*models.Book, error) {
	if len(s.data) == 0 {
		return nil, errors.New("Empty data")
	}

	res := make([]*models.Book, 0, len(s.data))
	for _, v := range s.data {
		res = append(res, v)
	}
	return res, nil
}

func (s *StorageMock) GetBookById(id int) (*models.Book, error) {
	book, isExist := s.data[id]
	if !isExist {
		return nil, errors.New("Model by id not found")
	}
	return book, nil
}

func (s *StorageMock) AddBook(book *models.Book) error {
	castedID, _ := strconv.Atoi(book.ID)

	if book.Title == "title_for_error" {
		return errors.New("Storage error")
	}

	s.data[castedID] = book
	return nil
}

func (s *StorageMock) UpdateBook(book *models.Book) error {
	castedID, _ := strconv.Atoi(book.ID)
	_, isExist := s.data[castedID]
	if !isExist {
		return errors.New("Model does not exist")
	}
	s.data[castedID] = book
	return nil
}

func (s *StorageMock) RemoveBook(id int) error {
	_, isExist := s.data[id]
	if !isExist {
		return errors.New("Model does not exist")
	}
	delete(s.data, id)
	return nil
}

var storageMock *StorageMock
var handler CrudHandler

func init() {
	storageMock = &StorageMock{data: make(map[int]*models.Book)}
	handler = InitHandlers(storageMock)
}

// ------------------------------------- GET ALL -------------------------------------
func TestGettAllForEmpty(t *testing.T) {
	cleanData()

	request, _ := http.NewRequest("GET", "", nil)
	recorder := httptest.NewRecorder()
	handler.GetAllBooks(recorder, request)

	result := recorder.Result()
	if result.StatusCode != http.StatusNotFound {
		t.Error("No error with empty data")
	}
}

func TestGettAllOK(t *testing.T) {
	body := []byte(`{"id": "1", "title": "test", "author": {"first_name": "first", "last_name": "last"}}`)
	request, _ := http.NewRequest("POST", "", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()
	handler.AddBook(recorder, request)
	recorder.Flush()

	request, _ = http.NewRequest("GET", "", nil)
	handler.GetAllBooks(recorder, request)

	result := recorder.Result()
	data, err := ioutil.ReadAll(result.Body)
	if err != nil || len(data) == 0 {
		t.Error("Something going wrong")
	}
}

// ------------------------------------- GET -------------------------------------
func TestGetBadRequest(t *testing.T) {
	request := httptest.NewRequest("GET", "/book/", nil)
	recorder := httptest.NewRecorder()
	handler.GetBook(recorder, request)

	result := recorder.Result()

	if result.StatusCode != http.StatusBadRequest {
		t.Error("No 400 with request without id")
	}
}

func TestGetNotFound(t *testing.T) {
	cleanData()

	router := mux.NewRouter()
	router.HandleFunc("/book/{id}", handler.GetBook).Methods("GET")

	server := httptest.NewServer(router)
	defer server.Close()

	response, err := http.Get(server.URL + "/book/1")
	if err != nil {
		t.Error("Error while request to GET /book/1", err)
	}
	if response.StatusCode != http.StatusNotFound {
		t.Error("No 404 on not exist id")
	}
}

func TestGetOK(t *testing.T) {
	bookToSend, id := prepareMockData()

	router := mux.NewRouter()
	router.HandleFunc("/book/{id}", handler.GetBook).Methods("GET")

	server := httptest.NewServer(router)
	defer server.Close()

	response, err := http.Get(server.URL + "/book/" + id)
	if err != nil {
		t.Error("Error on request GET /book/"+id, err)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Error with decoding item from server")
	}
	defer response.Body.Close()

	var book models.Book
	err = json.Unmarshal(data, &book)
	if err != nil {
		t.Error("Error with parse JSON")
	}

	if !(book.Title == bookToSend.Title &&
		book.Author.Firstname == bookToSend.Author.Firstname &&
		book.Author.Lastname == bookToSend.Author.Lastname) {
		t.Error("Received wrong data")
	}
}

// ------------------------------------- ADD -------------------------------------
func TestAddBadRequest(t *testing.T) {
	body := []byte(`{"id": "2", "title": "test", "author": "wrong_type"}`)
	request, _ := http.NewRequest("POST", "", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()
	handler.AddBook(recorder, request)

	result := recorder.Result()
	if result.StatusCode != http.StatusBadRequest {
		t.Error("No error with bad JSON")
	}
}

func TestAddWithStorageError(t *testing.T) {
	body := []byte(`{"id": "3", "title": "title_for_error", "author": {"first_name": "first", "last_name": "last"}}`)
	request, _ := http.NewRequest("POST", "", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()
	handler.AddBook(recorder, request)

	result := recorder.Result()
	if result.StatusCode != http.StatusInternalServerError {
		t.Error("No error with storage err")
	}
}

func TestAddOkK(t *testing.T) {
	body := []byte(`{"id": "1", "title": "test", "author": {"first_name": "first", "last_name": "last"}}`)
	request, _ := http.NewRequest("POST", "", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()
	handler.AddBook(recorder, request)

	result := recorder.Result()
	if result.StatusCode != http.StatusCreated {
		t.Error("Error on add model")
	}
}

// ------------------------------------- EDIT -------------------------------------
func TestEditBadRequest(t *testing.T) {
	request, _ := http.NewRequest("PUT", "", nil)
	recorder := httptest.NewRecorder()
	handler.EditBook(recorder, request)

	result := recorder.Result()
	if result.StatusCode != http.StatusBadRequest {
		t.Error("No error with bad request")
	}
}

func TestEditNotFound(t *testing.T) {
	prepareMockData()

	router := mux.NewRouter()
	router.HandleFunc("/book/{id}", handler.EditBook).Methods("PUT")

	server := httptest.NewServer(router)
	defer server.Close()

	bookToSend := models.Book{
		ID:    "2",
		Title: "test 2",
		Author: &models.Author{
			Firstname: "first 2",
			Lastname:  "last 2",
		},
	}
	body, _ := json.Marshal(bookToSend)

	request, _ := http.NewRequest("PUT", server.URL+"/book/999", bytes.NewBuffer(body))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error("Error on PUT /book/999", err)
	}
	if response.StatusCode != http.StatusNotFound {
		t.Error("No 404 on request with not exist id")
	}
}

func TestEditWithWrongJSON(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/book/{id}", handler.EditBook).Methods("PUT")

	server := httptest.NewServer(router)
	defer server.Close()

	wrongJSON := []byte(`{"id": "1", "title": "wrong_json", "author": "error_here"}`)
	request, _ := http.NewRequest("PUT", server.URL+"/book/1", bytes.NewBuffer(wrongJSON))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error("Error on PUT /book/1", err)
	}
	if response.StatusCode != http.StatusInternalServerError {
		t.Error("No 500 with wrong JSON")
	}
}

func TestEditOK(t *testing.T) {
	initBook, id := prepareMockData()

	router := mux.NewRouter()
	router.HandleFunc("/book/{id}", handler.EditBook).Methods("PUT")
	router.HandleFunc("/book/{id}", handler.GetBook).Methods("GET")

	server := httptest.NewServer(router)
	defer server.Close()

	bookToSend := models.Book{
		ID:    "2",
		Title: "test 2",
		Author: &models.Author{
			Firstname: "first 2",
			Lastname:  "last 2",
		},
	}
	body, _ := json.Marshal(bookToSend)

	request, _ := http.NewRequest("PUT", server.URL+"/book/"+id, bytes.NewBuffer(body))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error("Error while PUT /book/"+id, err)
	}
	if response.StatusCode != http.StatusOK {
		t.Error("No 200 on PUT /book/" + id)
	}

	response, err = http.Get(server.URL + "/book/" + id)
	if err != nil {
		t.Error("Error on GET /book/"+id, err)
	}

	var bookFromServer *models.Book
	data, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(data, &bookFromServer)
	if err != nil {
		t.Error("Error with parse JSON from server")
	}
	defer response.Body.Close()

	if initBook.Title == bookToSend.Title ||
		initBook.Author.Firstname == bookToSend.Author.Firstname ||
		initBook.Author.Lastname == bookToSend.Author.Lastname {
		t.Error("Model not replaced")
	}

	if !(bookFromServer.Title == bookToSend.Title &&
		bookFromServer.Author.Firstname == bookToSend.Author.Firstname &&
		bookFromServer.Author.Lastname == bookToSend.Author.Lastname) {
		t.Error("Error on replace")
	}
}

// ------------------------------------- DELETE -------------------------------------
func TestDeleteBadRequest(t *testing.T) {
	request, _ := http.NewRequest("DELETE", "", nil)
	recorder := httptest.NewRecorder()
	handler.DeleteBook(recorder, request)

	result := recorder.Result()
	if result.StatusCode != http.StatusBadRequest {
		t.Error("No error on bad request")
	}
}

func TestDeleteNotFound(t *testing.T) {
	cleanData()

	router := mux.NewRouter()
	router.HandleFunc("/book/{id}", handler.DeleteBook).Methods("DELETE")

	server := httptest.NewServer(router)
	defer server.Close()

	request, _ := http.NewRequest("DELETE", server.URL+"/book/42", nil)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error("Error on request DELETE /book/42", err)
	}
	if response.StatusCode != http.StatusNotFound {
		t.Error("No 404 on DELETE with not exist id")
	}
}

func TestDeleteOK(t *testing.T) {
	_, id := prepareMockData()

	router := mux.NewRouter()
	router.HandleFunc("/book/{id}", handler.DeleteBook).Methods("DELETE")

	server := httptest.NewServer(router)
	defer server.Close()

	request, _ := http.NewRequest("DELETE", server.URL+"/book/"+id, nil)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error("Error on request to DELETE /book/"+id, err)
	}
	if response.StatusCode != http.StatusNoContent {
		t.Error("No 204 on DELETE with correct id")
	}
}

// ------------------------------------- helpers -------------------------------------

func cleanData() {
	for key := range storageMock.data {
		delete(storageMock.data, key)
	}
}

func prepareMockData() (*models.Book, string) {
	cleanData()

	bookToSend := models.Book{
		ID:    "1",
		Title: "test",
		Author: &models.Author{
			Firstname: "first",
			Lastname:  "last",
		},
	}
	body, _ := json.Marshal(bookToSend)

	request, _ := http.NewRequest("POST", "", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()
	handler.AddBook(recorder, request)
	recorder.Flush()

	var item *models.Book
	for k := range storageMock.data {
		item = storageMock.data[k]
		break
	}

	return &bookToSend, item.ID
}
