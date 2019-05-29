package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gopkg.in/guregu/null.v3"
)

type Employee struct {
	User
	Accounts []Account `json:"accounts"`
}

type User struct {
	ID        int         `json:"id"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Email     string      `json:"email"`
	Phone     null.String `json:"phone"`
}

type Account struct {
	ID       int  `json:"id"`
	UserID   int  `json:"user_id"`
	Value    int  `json:"value"`
	Currency int  `json:"currency"`
	Active   bool `json:"active"`
}

type Transaction struct {
	ID                    int      `json:"id"`
	Type                  string   `json:"type"`
	AccountForOperationID int      `json:"acc_for_operation"`
	RecievedIDAccount     null.Int `json:"recieved_id"`
	Value                 int      `json:"value"`
	Canceled              bool     `json:"canceled"`
	CanceledID            null.Int `json:"canceled_id"`
}

const (
	host     = "localhost"
	port     = "5432"
	user     = "employee"
	password = "employee"
	dbname   = "dbank"
)

var commands = map[string]string{
	"Deposit":  "Deposit",
	"Withdraw": "Withdraw",
	"Transfer": "Transfer",
	"Cancel":   "Cancel",
}
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

	s.HandleFunc("/users/{id}/accounts", CreateAccount).Methods(http.MethodPost)
	s.HandleFunc("/users/{id}/accounts/{acc_id}/balance", ReadAccount).Methods(http.MethodGet)
	s.HandleFunc("/users/{id}/accounts/{acc_id}", DeleteAccount).Methods(http.MethodDelete)

	s.HandleFunc("/transactions", ListTransactions).Methods(http.MethodGet)
	s.HandleFunc("/transactions", Transactions).Methods(http.MethodPost)
	s.HandleFunc("/transactions/{id}", CancelTransaction).Methods(http.MethodDelete)

	http.Handle("/", r)
	if err := http.ListenAndServe(":9090", loggingMiddleware(r)); err != nil {
		log.Fatal(err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		log.Println(r.RequestURI, r.Method)
		next.ServeHTTP(w, r)
	})
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("List User function started")
	rows, err := db.Query("SELECT * FROM users ORDER BY id DESC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	result := make(map[int]User)
	for rows.Next() {
		e := new(User)
		err := rows.Scan(&e.ID, &e.FirstName, &e.LastName, &e.Email, &e.Phone)
		if err != nil {
			log.Fatal(err)
		}
		result[e.ID] = *e
	}
	marshaled, _ := json.MarshalIndent(result, "", " ")
	w.Write(marshaled)
	log.Println("Users printed")
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Create User function started")
	var e Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal Server Error111")
		return
	}

	sqlStatement := "INSERT INTO users (first_name, last_name, email, phone)	VALUES ($1, $2, $3, $4)"
	_, err := db.Exec(sqlStatement, e.FirstName, e.LastName, e.Email, e.Phone)
	if err != nil {
		log.Fatal(err)
	}
	// w.WriteHeader(statusCode, 201)
	http.Redirect(w, r, "http://localhost:9090/users", 301)
	log.Printf("User %s %s created", e.FirstName, e.LastName)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete User function started")
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
	http.Redirect(w, r, "http://localhost:9090/users", 301)
	log.Printf("User with id %s deleted", id)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Update User function started")
	var e Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal Server Error")
		return
	}
	id := mux.Vars(r)["id"]
	e.ID, _ = strconv.Atoi(id)

	sqlStatement := "UPDATE users SET first_name = $2, last_name = $3, email = $4, phone = $5 WHERE id = $1"
	_, err := db.Exec(sqlStatement, e.ID, e.FirstName, e.LastName, e.Email, e.Phone)
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "http://localhost:9090/users", 301)
	log.Printf("User with id %s updated", id)
}

func ReadUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Read User function started")
	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID mustn't be empty")
		return
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	sqlStatement := "SELECT * FROM users WHERE id = $1"
	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	defer rows.Close()
	e := new(Employee)
	for rows.Next() {
		err = rows.Scan(&e.ID, &e.FirstName, &e.LastName, &e.Email, &e.Phone)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}
	sqlStatement = "SELECT * FROM accounts WHERE id_user = $1 and active = $2"
	rows, err = db.Query(sqlStatement, id, false)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		a := new(Account)
		err = rows.Scan(&a.ID, &a.UserID, &a.Value, &a.Currency, &a.Active)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		e.Accounts = append(e.Accounts, *a)
	}
	tx.Commit()
	marshaled, _ := json.MarshalIndent(e, "", " ")
	w.Write(marshaled)
	log.Printf("Data user with id %s printed", id)
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	log.Println("Create Account function started")
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
		log.Println("Users printed")
	}

	sqlStatement := "INSERT INTO accounts (id_user, value, currency) VALUES ($1, $2, $3)"
	_, err := db.Exec(sqlStatement, id, a.Value, a.Currency)
	if err != nil {
		log.Fatal(err)
	}
	s := fmt.Sprintf("http://localhost:9090/users/%s", id)
	http.Redirect(w, r, s, 301)
	log.Printf("Account for user with id %s created", id)
}

func ReadAccount(w http.ResponseWriter, r *http.Request) {
	log.Println("Read Account function started")
	id := mux.Vars(r)["id"]
	accID := mux.Vars(r)["acc_id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID mustn't be empty")
		return
	}
	sqlStatement := "SELECT value, currency FROM accounts WHERE id = $1 and active = $2"
	rows, err := db.Query(sqlStatement, accID, false)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var value, currency int
	for rows.Next() {
		err = rows.Scan(&value, &currency)
		if err != nil {
			log.Fatal(err)
		}
	}
	result := make(map[string]int)
	result["id"], _ = strconv.Atoi(accID)
	result["value"] = value
	result["currency"] = currency
	marshaled, _ := json.MarshalIndent(result, "", " ")
	w.Write(marshaled)
	log.Printf("Account with id %s printed", accID)
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete Account function started")
	id := mux.Vars(r)["id"]
	accID := mux.Vars(r)["acc_id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID mustn't be empty")
		return
	}
	sqlStatement := "DELETE FROM accounts WHERE id = $1 and active = $2"
	_, err := db.Exec(sqlStatement, accID, false)
	if err != nil {
		log.Fatal(err)
	}

	s := fmt.Sprintf("http://localhost:9090/users/%s", id)
	http.Redirect(w, r, s, 301)
	log.Printf("Account %s users with id %s deleted", accID, id)
}

func ListTransactions(w http.ResponseWriter, r *http.Request) {
	log.Println("List Transactions function started")
	rows, err := db.Query("SELECT * FROM transactions ORDER BY id DESC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	result := make(map[int]Transaction)
	for rows.Next() {
		t := new(Transaction)
		err := rows.Scan(&t.ID, &t.Type, &t.AccountForOperationID, &t.RecievedIDAccount, &t.Value, &t.Canceled, &t.CanceledID)
		if err != nil {
			log.Fatal(err)
		}
		result[t.ID] = *t
	}
	marshaled, _ := json.MarshalIndent(result, "", " ")
	w.Write(marshaled)
	log.Println("List Transactions printed")
}

func Transactions(w http.ResponseWriter, r *http.Request) {
	log.Println("Transactions function started")
	var t Transaction
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal Server Error111")
		return
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	switch command := t.Type; command {
	case commands["Withdraw"]:
		sqlStatement := "INSERT INTO transactions (type, account_for_operation_id," +
			"value)	VALUES ($1, $2, $3)"
		_, err := db.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.Value)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}

		sqlStatement = "UPDATE accounts SET value = value - $1	WHERE id = $2"
		_, err = db.Exec(sqlStatement, t.Value, t.AccountForOperationID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		log.Println("Transactions Withdraw completed")
	case commands["Deposit"]:
		sqlStatement := "INSERT INTO transactions (type, account_for_operation_id," +
			"value)	VALUES ($1, $2, $3)"
		_, err := db.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.Value)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		sqlStatement = "UPDATE accounts SET value = value + $1	WHERE id = $2"
		_, err = db.Exec(sqlStatement, t.Value, t.AccountForOperationID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		log.Println("Transactions Deposit completed")
	case commands["Transfer"]:
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
		log.Println("Transactions Transfer completed")
	}
	tx.Commit()
	http.Redirect(w, r, r.URL.RequestURI(), 301)
}

func CancelTransaction(w http.ResponseWriter, r *http.Request) {
	log.Println("Cancel transaction function started")
	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID mustn't be empty")
		return
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	sqlStatement := "SELECT * FROM transactions WHERE id = $1"
	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	defer rows.Close()
	t := new(Transaction)
	for rows.Next() {
		err = rows.Scan(&t.ID, &t.Type, &t.AccountForOperationID, &t.RecievedIDAccount, &t.Value, &t.Canceled, &t.CanceledID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}
	if t.Canceled == true {
		tx.Rollback()
		log.Fatal("Transaction was canceled before")
	}
	switch command := t.Type; command {
	case commands["Withdraw"]:
		log.Println("Cancel withdraw")
		sqlStatement = "UPDATE accounts SET value = value + $1	WHERE id = $2"
		_, err = db.Exec(sqlStatement, t.Value, t.AccountForOperationID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		sqlStatement = "UPDATE transactions SET canceled = $1 WHERE id = $2"
		_, err = db.Exec(sqlStatement, true, t.ID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		t.Type = commands["Cancel"]
		sqlStatement = "INSERT INTO transactions (type, account_for_operation_id," +
			"value, canceled_id)	VALUES ($1, $2, $3, $4)"
		_, err = db.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.Value, t.ID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	case commands["Deposit"]:
		log.Println("Cancel deposit")
		sqlStatement = "UPDATE accounts SET value = value - $1	WHERE id = $2"
		_, err = db.Exec(sqlStatement, t.Value, t.AccountForOperationID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		sqlStatement = "UPDATE transactions SET canceled = $1 WHERE id = $2"
		_, err = db.Exec(sqlStatement, true, t.ID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		t.Type = commands["Cancel"]
		sqlStatement = "INSERT INTO transactions (type, account_for_operation_id," +
			"value, canceled_id)	VALUES ($1, $2, $3, $4)"
		_, err = db.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.Value, t.ID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	case commands["Transfer"]:
		log.Println("Cancel transfer")
		sqlStatement = "UPDATE transactions SET canceled = $1 WHERE id = $2"
		_, err = db.Exec(sqlStatement, true, t.ID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		t.Type = commands["Cancel"]
		sqlStatement = "INSERT INTO transactions (type, account_for_operation_id," +
			"recieve_id_account, value, canceled_id)	VALUES ($1, $2, $3, $4, $5)"
		_, err = db.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.RecievedIDAccount, t.Value, t.ID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		sqlStatement = "UPDATE accounts SET value = value - $1	WHERE id = $2"
		_, err = db.Exec(sqlStatement, t.Value, t.RecievedIDAccount)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		sqlStatement = "UPDATE accounts SET value = value + $1	WHERE id = $2"
		_, err = db.Exec(sqlStatement, t.Value, t.AccountForOperationID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}
	tx.Commit()
	log.Println("Transactions canceled")
	http.Redirect(w, r, "http://localhost:9090/transactions", 301)
}
