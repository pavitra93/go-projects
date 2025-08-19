# Go Movies REST API (Gorilla Mux)

A minimal RESTful API for managing a list of movies using **Go**, **net/http**, and **Gorilla Mux**.

---

## 🧭 Overview

This project exposes CRUD endpoints for a `Movie` resource stored in-memory. It’s great for learning Go’s HTTP server, routing with Gorilla Mux, and basic JSON handling.

**Tech stack:**

* Go 1.20+
* net/http
* github.com/gorilla/mux (router)

---

## 📦 Project Structure

```
.
├── main.go        # application entrypoint & handlers
└── README.md      # this file
```

---

## ✅ Features

* List all movies
* Get a single movie by ID
* Create a movie
* Update (replace) a movie
* Delete a movie

> ⚠️ Data is stored in a process-level slice (`[]Movie`). All changes are ephemeral and non-persistent.

---

## 🧪 Data Model

```json
Movie {
  id: string,
  isbn: string,
  title: string,
  year: number,
  director: {
    first_name: string,
    last_name: string
  }
}
```

Example object:

```json
{
  "id": "1",
  "isbn": "123",
  "title": "Movie",
  "year": 2000,
  "director": { "first_name": "James", "last_name": "Bond" }
}
```

---

## 🛠 Setup & Run

### 1) Clone and init module

```bash
# inside project folder
go mod init github.com/you/movies-api
```

### 2) Get dependencies

```bash
go get github.com/gorilla/mux@latest
```

### 3) Run the server

```bash
go run .
```

Server starts on **[http://localhost:8080](http://localhost:8080)**

---

## 🔗 API Reference

Base URL: `http://localhost:8080`

### 1) Get all movies

**GET** `/movies`

```bash
curl -s http://localhost:8080/movies | jq
```

Response: `200 OK`

```json
[
  {
    "id": "1",
    "isbn": "123",
    "title": "Movie",
    "year": 2000,
    "director": { "first_name": "James", "last_name": "Bond" }
  },
  {
    "id": "2",
    "isbn": "345",
    "title": "Movie 2",
    "year": 2002,
    "director": { "first_name": "John", "last_name": "Doe" }
  }
]
```

### 2) Get movie by ID

**GET** `/movies/{id}`

```bash
curl -s http://localhost:8080/movies/1 | jq
```

Response: `200 OK` with movie JSON, or `200 OK` with empty body if not found (current behavior).

### 3) Create a movie

**POST** `/movies`

```bash
curl -s -X POST http://localhost:8080/movies \
  -H 'Content-Type: application/json' \
  -d '{
        "isbn": "9999",
        "title": "Interstellar",
        "year": 2014,
        "director": {"first_name": "Christopher", "last_name": "Nolan"}
      }' | jq
```

Response: `200 OK` (currently returns the full list).

### 4) Update (replace) a movie

**PUT** `/movies/{id}`

```bash
curl -s -X PUT http://localhost:8080/movies/1 \
  -H 'Content-Type: application/json' \
  -d '{
        "isbn": "1111",
        "title": "The Batman",
        "year": 2022,
        "director": {"first_name": "Matt", "last_name": "Reeves"}
      }' | jq
```

Response: `200 OK` (currently returns the full list). See **Known quirks** about IDs changing.

### 5) Delete a movie

**DELETE** `/movies/{id}`

```bash
curl -s -X DELETE http://localhost:8080/movies/2 | jq
```

Response: `200 OK` (returns remaining list).

---

## 📜 License

MIT (or choose your own)

---

## 🙌 Credits

Starter inspired by standard Gorilla Mux tutoring examples. You rock for building and iterating on this! 🚀
