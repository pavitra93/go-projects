package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Year     int       `json:"year"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

var movies []Movie

// Get All Movies
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movies)
}

// delete Movie
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for idx, movie := range movies {
		if movie.ID == params["id"] {
			movies = append(movies[:idx], movies[idx+1:]...)
			break
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movies)
}

// get single movie
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, movie := range movies {
		if movie.ID == params["id"] {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	movie.ID = strconv.Itoa(rand.Int())
	movie.Isbn = movie.Isbn
	movie.Title = movie.Title
	movie.Year = movie.Year
	movie.Director = &Director{
		FirstName: movie.Director.FirstName,
		LastName:  movie.Director.LastName,
	}
	movies = append(movies, movie)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movies)
	return
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for idx, movie := range movies {
		if movie.ID == params["id"] {
			movies = append(movies[:idx], movies[idx+1:]...)
			var movie Movie
			if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			movie.ID = strconv.Itoa(rand.Int())
			movie.Isbn = movie.Isbn
			movie.Title = movie.Title
			movie.Year = movie.Year
			movie.Director = &Director{
				FirstName: movie.Director.FirstName,
				LastName:  movie.Director.LastName,
			}
			movies = append(movies, movie)
			break
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movies)
	return
}

func main() {
	// Sample initiale data
	movies = append(movies, Movie{
		ID:    "1",
		Isbn:  "123",
		Title: "Movie",
		Year:  2000,
		Director: &Director{
			FirstName: "James",
			LastName:  "Bond",
		},
	},
		Movie{
			ID:    "2",
			Isbn:  "345",
			Title: "Movie 2",
			Year:  2002,
			Director: &Director{
				FirstName: "John",
				LastName:  "Doe",
			},
		})

	// Initializing routes
	r := mux.NewRouter()
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	// Starting server
	fmt.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
