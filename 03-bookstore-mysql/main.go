package main

import (
	"github.com/gorilla/mux"
	"github.com/pavitra93/go-projects/03-bookstore-mysql/pkg/routes"
	_ "gorm.io/driver/mysql"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	routes.RegisterBookStoreRoutes(r)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
