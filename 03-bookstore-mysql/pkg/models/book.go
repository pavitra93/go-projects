package models

import (
	"github.com/pavitra93/go-projects/03-bookstore-mysql/pkg/config"
	"gorm.io/gorm"
)

var db *gorm.DB

type Book struct {
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Publication string `json:"publication"`
	Year        int    `json:"year"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	err := db.AutoMigrate(&Book{})
	if err != nil {
		return
	}
}

func (b *Book) CreateBook() *Book {
	db.Create(&b)
	return b
}

func GetAllBook() []Book {
	var Books []Book
	db.Find(&Books)
	return Books
}

func GetBookById(id int) (*Book, *gorm.DB) {
	var Book Book
	res := db.Where("id = ?", id).First(&Book)
	if res.Error != nil {
		return nil, db
	}
	return &Book, db
}

func DeleteBookById(id int) {
	var Book Book
	db.Where("id = ?", id).Delete(&Book)
}
