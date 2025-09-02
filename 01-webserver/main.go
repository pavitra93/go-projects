package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed.", http.StatusNotFound)
		return
	}

	_, err := fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])

	if err != nil {
		log.Println(err)
	}
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/form" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	fmt.Fprintln(w, "Form requesting handling initialized")
	if err := r.ParseForm(); err != nil {
		log.Println(err)
	}

	_, err := fmt.Fprintf(w, "Name: %s\n", r.FormValue("name"))
	if err != nil {
		log.Println(err)
	}

	_, err = fmt.Fprintf(w, "Email: %s", r.FormValue("email"))
	if err != nil {
		log.Println(err)
	}
}

func main() {
	fileServer := http.FileServer(http.Dir("./public"))

	http.Handle("/", fileServer)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/hello", helloHandler)

	fmt.Printf("Listening on port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
