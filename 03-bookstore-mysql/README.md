# Go Bookstore API (Gorilla Mux + GORM + MySQL)

A simple RESTful API for managing books, built with **Go**, **Gorilla Mux** for routing, and **GORM** for ORM over **MySQL**.

---

## 📦 Tech Stack

* Go 1.20+
* Gorilla Mux (`github.com/gorilla/mux`)
* GORM (`gorm.io/gorm`) + MySQL driver (`gorm.io/driver/mysql`)
* net/http, encoding/json

---

## 🗂 Project Structure (suggested)

```
03-bookstore-mysql/
├─ main.go                        # app entrypoint
├─ pkg/
│  ├─ config/                     # DB connection, env
│  ├─ controllers/                # HTTP handlers
│  ├─ models/                     # GORM models & queries
│  ├─ routes/                     # mux routes
│  └─ utils/                      # helpers (e.g., ParseBody)
├─ go.mod / go.sum
└─ README.md
```
---
## 🔧 Setup

### 1) Initialize module & get deps

```bash
# inside project root
go mod init github.com/pavitra93/03-bookstore-mysql

# deps
go get github.com/gorilla/mux@latest
go get gorm.io/gorm@latest
go get gorm.io/driver/mysql@latest
```

### 2) Configure database

Create a **MySQL** database, e.g. `bookstore`.

Set a DSN in your config (example):

```
user:password@tcp(127.0.0.1:3306)/bookstore?charset=utf8mb4&parseTime=True&loc=Local
```

**.env example** (if your `config.Connect()` loads from env):

```env
DB_USER=root
DB_PASS=secret
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=bookstore
DB_PARAMS=charset=utf8mb4&parseTime=True&loc=Local
```

Then in `config.Connect()` compose:

```
<user>:<pass>@tcp(<host>:<port>)/<name>?<DB_PARAMS>
```

> The model’s `init()` runs `AutoMigrate(&Book{})` which will create the `books` table on startup if missing.

### 3) Run the server

```bash
go run ./cmd/server   # if you have that structure
# or, with your current layout
go run .
```

Server starts at **[http://localhost:8080](http://localhost:8080)**.

### 4) Build a binary (optional)

```bash
go build -o bookstore-api
./bookstore-api
```

---

## 📘 Data Model

```go
// pkg/models/book.go

type Book struct {
    gorm.Model         // ID, CreatedAt, UpdatedAt, DeletedAt
    Name        string `json:"name"`
    Author      string `json:"author"`
    Publication string `json:"publication"`
    Year        int    `json:"year"`
}
```

---

## 🔗 API Endpoints

Base URL: `http://localhost:8080`

### Get all books

**GET** `/books`

```bash
curl -s http://localhost:8080/books | jq
```

**200 OK** → `[]Book`

### Get book by ID

**GET** `/books/{id}`

```bash
curl -s http://localhost:8080/books/1 | jq
```

**200 OK** → `Book`

**404 Not Found** if missing (your controller returns 404 when `GetBookById` is nil).

### Create a book

**POST** `/books`

```bash
curl -s -X POST http://localhost:8080/books \
  -H 'Content-Type: application/json' \
  -d '{
        "name": "Clean Architecture",
        "author": "Robert C. Martin",
        "publication": "Pearson",
        "year": 2017
      }' | jq
```

**200 OK** → Created `Book` (consider `201 Created` as an improvement).

### Update a book

**PUT** `/books/{id}`

```bash
curl -s -X PUT http://localhost:8080/books/1 \
  -H 'Content-Type: application/json' \
  -d '{ "name": "Clean Architecture (2nd Ed)", "year": 2020 }' | jq
```

**200 OK** → Updated `Book`

**404 Not Found** if missing.

### Delete a book

**DELETE** `/books/{id}`

```bash
curl -s -X DELETE http://localhost:8080/books/1
```

**200 OK** → text: `Deleted book with ID 1.`
*(Consider `204 No Content` for production.)*

## 📜 License

MIT (or choose your own)

---

## 🙌 Credits

Built on top of Gorilla Mux & GORM. Nice work assembling your bookstore API!
