package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

type Transaction struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Type      string    `json:"type"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

var db *sql.DB

func main() {
	initDB()
	r := mux.NewRouter()
	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/transactions", getTransactions).Methods("GET")
	r.HandleFunc("/transactions/{id}", getTransaction).Methods("GET")
	r.HandleFunc("/transactions", createTransaction).Methods("POST")

	fmt.Println("Server started at port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./bank.db")
	if err != nil {
		log.Fatal(err)
	}
	createTables()
}

func createTables() {
	createUsersTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"username" TEXT,
		"password" TEXT
	);`

	createTransactionsTableSQL := `
	CREATE TABLE IF NOT EXISTS transactions (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"user_id" INTEGER,
		"type" TEXT,
		"amount" REAL,
		"timestamp" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(createUsersTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createTransactionsTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Password)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	row := db.QueryRow("SELECT * FROM users WHERE id = ?", params["id"])

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, user.Password)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User created successfully")
}

func getTransactions(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM transactions")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var transaction Transaction
		err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.Type, &transaction.Amount, &transaction.Timestamp)
		if err != nil {
			log.Fatal(err)
		}
		transactions = append(transactions, transaction)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func getTransaction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	row := db.QueryRow("SELECT * FROM transactions WHERE id = ?", params["id"])

	var transaction Transaction
	err := row.Scan(&transaction.ID, &transaction.UserID, &transaction.Type, &transaction.Amount, &transaction.Timestamp)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

func createTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction Transaction
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO transactions (user_id, type, amount) VALUES (?, ?, ?)", transaction.UserID, transaction.Type, transaction.Amount)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Transaction created successfully")
}
