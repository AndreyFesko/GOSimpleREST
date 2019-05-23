package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gopkg.in/guregu/null.v3"
)

type Employee struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type Account struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`
	Value  int `json:"value"`
}

type Transaction struct {
	ID                    int      `json:"id"`
	Type                  string   `json:"type"`
	AccountForOperationID int      `json:"acc_for_operation"`
	RecievedIDAccount     null.Int `json:"recieved_id"`
	Value                 int      `json:"value"`
}

const (
	host     = "localhost"
	port     = "5432"
	user     = "employee"
	password = "employee"
	dbname   = "dbank"
)

var db *sql.DB

func init() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := mux.NewRouter()
	s := r.PathPrefix("/").Subrouter()

	s.HandleFunc("/users", ListUsers).Methods(http.MethodGet)
	s.HandleFunc("/users", CreateUser).Methods(http.MethodPost)

	s.HandleFunc("/users/{id}", ReadUser).Methods(http.MethodGet)
	s.HandleFunc("/users/{id}", UpdateUser).Methods(http.MethodPatch)
	s.HandleFunc("/users/{id}", DeleteUser).Methods(http.MethodDelete)

	s.HandleFunc("/users/{id}", CreateAccount).Methods(http.MethodPost)
	s.HandleFunc("/users/{id}/{AccountID}", DeleteAccount).Methods(http.MethodDelete)

	s.HandleFunc("/transactions", ListTransactions).Methods(http.MethodGet)
	s.HandleFunc("/transactions", Transactions).Methods(http.MethodPost)

	http.Handle("/", r)
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatal(err)
	}
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM users ORDER BY id DESC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	result := []Employee{}
	for rows.Next() {
		e := new(Employee)
		err := rows.Scan(&e.ID, &e.FirstName, &e.LastName, &e.Email)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, *e)
	}
	marshaled, _ := json.MarshalIndent(result, "", " ")
	w.Write(marshaled)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var e Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal Server Error111")
		return
	}

	sqlStatement := "INSERT INTO users (first_name, last_name, email)	VALUES ($1, $2, $3)"
	_, err := db.Exec(sqlStatement, e.FirstName, e.LastName, e.Email)
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "http://localhost:9090/users/", 301)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID mustn't be empty")
		return
	}
	sqlStatement := "DELETE FROM users WHERE id = $1"
	_, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "http://localhost:9090/users/", 301)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var e Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal Server Error")
		return
	}
	id := mux.Vars(r)["id"]
	e.ID, _ = strconv.Atoi(id)

	sqlStatement := "UPDATE users SET first_name = $2, last_name = $3, email = $4 WHERE id = $1"
	_, err := db.Exec(sqlStatement, e.ID, e.FirstName, e.LastName, e.Email)
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "http://localhost:9090/users/", 301)
}

func ReadUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID mustn't be empty")
		return
	}
	sqlStatement := "SELECT * FROM users WHERE id = $1"
	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	e := new(Employee)
	for rows.Next() {
		err = rows.Scan(&e.ID, &e.FirstName, &e.LastName, &e.Email)
		if err != nil {
			log.Fatal(err)
		}
	}
	marshaled, _ := json.MarshalIndent(e, "", " ")
	w.Write(marshaled)
	sqlStatement = "SELECT * FROM accounts WHERE id_user = $1"
	rows, err = db.Query(sqlStatement, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	result := []Account{}
	for rows.Next() {
		a := new(Account)
		err = rows.Scan(&a.ID, &a.UserID, &a.Value)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, *a)
	}
	marshaled, _ = json.MarshalIndent(result, "", " ")
	w.Write(marshaled)
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID mustn't be empty")
		return
	}
	var a Account
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal Server Error111")
		return
	}

	sqlStatement := "INSERT INTO accounts (id_user, value)	VALUES ($1, $2)"
	_, err := db.Exec(sqlStatement, id, a.Value)
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, r.URL.RequestURI(), 301)
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	account_id := mux.Vars(r)["AccountID"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID mustn't be empty")
		return
	}
	sqlStatement := "DELETE FROM accounts WHERE id = $1"
	_, err := db.Exec(sqlStatement, account_id)
	if err != nil {
		log.Fatal(err)
	}

	s := fmt.Sprintf("http://localhost:9090/users/%s", id)
	http.Redirect(w, r, s, 301)
}

func ListTransactions(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM transactions ORDER BY id DESC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	result := []Transaction{}
	for rows.Next() {
		t := new(Transaction)
		err := rows.Scan(&t.ID, &t.Type, &t.AccountForOperationID, &t.RecievedIDAccount, &t.Value)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, *t)
	}
	marshaled, _ := json.MarshalIndent(result, "", " ")
	w.Write(marshaled)
}

func Transactions(w http.ResponseWriter, r *http.Request) {
	var t Transaction
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal Server Error111")
		return
	}
	switch command := t.Type; command {
	case "WriteOff":
		sqlStatement := "INSERT INTO transactions (type, account_for_operation_id," +
			"value)	VALUES ($1, $2, $3)"
		_, err := db.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.Value)
		if err != nil {
			log.Fatal(err)
		}

		sqlStatement = "UPDATE accounts SET value = value - $1	WHERE id = $2"
		_, err = db.Exec(sqlStatement, t.Value, t.AccountForOperationID)
		if err != nil {
			log.Fatal(err)
		}
	case "Refill":
		sqlStatement := "INSERT INTO transactions (type, account_for_operation_id," +
			"value)	VALUES ($1, $2, $3)"
		_, err := db.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.Value)
		if err != nil {
			log.Fatal(err)
		}
		sqlStatement = "UPDATE accounts SET value = value + $1	WHERE id = $2"
		_, err = db.Exec(sqlStatement, t.Value, t.AccountForOperationID)
		if err != nil {
			log.Fatal(err)
		}
	case "Transfer":
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		sqlStatement := "INSERT INTO transactions (type, account_for_operation_id," +
			"recieve_id_account, value)	VALUES ($1, $2, $3, $4)"
		_, err = db.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.RecievedIDAccount, t.Value)
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare("UPDATE accounts SET value = value + $1	WHERE id = $2")
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		defer stmt.Close()
		if _, err := stmt.Exec(t.Value, t.RecievedIDAccount); err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		stmt, err = tx.Prepare("UPDATE accounts SET value = value - $1	WHERE id = $2")
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		defer stmt.Close()
		if _, err := stmt.Exec(t.Value, t.AccountForOperationID); err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		tx.Commit()
	}
	http.Redirect(w, r, r.URL.RequestURI(), 301)
}
