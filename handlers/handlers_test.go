package handlers

import (
	"../models"
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
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

func TestAddOkK(t *testing.T) {
	body := []byte(`{"id": "1", "title": "test", "author": {"first_name": "first", "last_name": "last"}}`)
	request, _ := http.NewRequest("POST", "", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()
	handler.AddBook(recorder, request)

	result := recorder.Result()
	if result.StatusCode != http.StatusCreated {
		t.Error("Error while add model")
	}
}

func TestAddErrorOnUnmurshal(t *testing.T) {
	body := []byte(`{"id": "2", "title": "test", "author": "wrong_type"}`)
	request, _ := http.NewRequest("POST", "", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()
	handler.AddBook(recorder, request)

	result := recorder.Result()
	if result.StatusCode != http.StatusBadRequest {
		t.Error("No error while bad JSON")
	}
}

func TestAddWithStorageError(t *testing.T) {
	body := []byte(`{"id": "3", "title": "title_for_error", "author": {"first_name": "first", "last_name": "last"}}`)
	request, _ := http.NewRequest("POST", "", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()
	handler.AddBook(recorder, request)

	result := recorder.Result()
	if result.StatusCode != http.StatusInternalServerError {
		t.Error("No error while storage err")
	}
}

func TestGettAllForEmpty(t *testing.T) {
	for key := range storageMock.data {
		delete(storageMock.data, key)
	}

	request, _ := http.NewRequest("GET", "", nil)
	recorder := httptest.NewRecorder()
	handler.GetAllBooks(recorder, request)

	result := recorder.Result()
	if result.StatusCode != http.StatusNotFound {
		t.Error("No error while empty data")
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

func TestDeleteBadRequest(t *testing.T) {
	request, _ := http.NewRequest("DELETE", "", nil)
	recorder := httptest.NewRecorder()
	handler.DeleteBook(recorder, request)

	result := recorder.Result()
	if result.StatusCode != http.StatusBadRequest {
		t.Error("No error on bad request")
	}
}
