package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pavitra93/go-projects/03-bookstore-mysql/pkg/models"
	"github.com/pavitra93/go-projects/03-bookstore-mysql/pkg/utils"
	"net/http"
	"strconv"
)

var NewBook models.Book

func GetBook(w http.ResponseWriter, r *http.Request) {
	NewBooks := models.GetAllBook()
	res, _ := json.Marshal(NewBooks)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(res)
	if err != nil {
		return
	}
}

func GetBookByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["id"]
	ID, _ := strconv.Atoi(bookID)
	BookDetails, _ := models.GetBookById(ID)
	if BookDetails == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	res, _ := json.Marshal(BookDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
	return
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	NewBookDetails := &models.Book{}
	utils.ParseBody(r, NewBookDetails)
	b := NewBookDetails.CreateBook()
	res, _ := json.Marshal(b)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
	return
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["id"]
	ID, _ := strconv.Atoi(bookID)
	models.DeleteBookById(ID)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Deleted book with ID " + vars["id"] + "."))
	return
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["id"]
	ID, err := strconv.Atoi(bookID)
	if err != nil {
		fmt.Println(err)
	}

	// Get book from DB
	BookDetails, db := models.GetBookById(ID)
	if BookDetails == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// create new book from request body
	UpdateBookDetails := &models.Book{}
	utils.ParseBody(r, UpdateBookDetails)

	if UpdateBookDetails.Name != "" {
		BookDetails.Name = UpdateBookDetails.Name
	}

	if UpdateBookDetails.Author != "" {
		BookDetails.Author = UpdateBookDetails.Author
	}

	if UpdateBookDetails.Publication != "" {
		BookDetails.Publication = UpdateBookDetails.Publication
	}

	if UpdateBookDetails.Year > 0 {
		BookDetails.Year = UpdateBookDetails.Year
	}

	db.Save(&BookDetails)
	res, _ := json.Marshal(BookDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}
