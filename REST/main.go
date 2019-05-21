package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	//"sync"

	//"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Employee struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type Account struct {
	ID      int `json:"id"`
	ID_User int `json:"id_user"`
	Value   int `json:"value"`
}

const (
	host     = "localhost"
	port     = "5432"
	user     = "employee"
	password = "employee"
	dbname   = "dbank"
)

func dbConnect() (db *sql.DB) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func main() {
	r := mux.NewRouter()
	s := r.PathPrefix("/users").Subrouter()

	s.HandleFunc("/", ListUsers).Methods(http.MethodGet)
	s.HandleFunc("/", CreateUser).Methods(http.MethodPost)

	s.HandleFunc("/{id}", ReadUser).Methods(http.MethodGet)
	s.HandleFunc("/{id}", UpdateUser).Methods(http.MethodPatch)
	s.HandleFunc("/{id}", DeleteUser).Methods(http.MethodDelete)
	s.HandleFunc("/{id}", CreateAccount).Methods(http.MethodPost)
	//s.HandleFunc("/{id}", DeleteAccount).Methods(http.MethodDelete)

	http.Handle("/", r)
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatal(err)
	}
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	db := dbConnect()
	rows, err := db.Query("SELECT * FROM users ORDER BY id DESC")
	if err != nil {
		log.Fatal(err)
	}

	employee := Employee{}
	result := []Employee{}

	for rows.Next() {
		var id int
		var first_name, last_name, email string
		err := rows.Scan(&id, &first_name, &last_name, &email)
		if err != nil {
			log.Fatal(err)
		}
		employee.ID = id
		employee.FirstName = first_name
		employee.LastName = last_name
		employee.Email = email
		result = append(result, employee)
	}
	marshaled, _ := json.MarshalIndent(result, "", " ")
	w.Write(marshaled)
	defer db.Close()
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var employee Employee
	db := dbConnect()
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal Server Error111")
		return
	}

	sqlStatement := "INSERT INTO users (first_name, last_name, email)	VALUES ($1, $2, $3)"
	_, err := db.Exec(sqlStatement, employee.FirstName, employee.LastName, employee.Email)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, "http://localhost:9090/users/", 301)
	defer db.Close()
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	db := dbConnect()
	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID mustn't be empty")
		return
	}
	sqlStatement := "DELETE FROM users WHERE id = $1"
	_, err := db.Exec(sqlStatement, id)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, "http://localhost:9090/users/", 301)
	defer db.Close()
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var employee Employee

	db := dbConnect()

	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal Server Error")
		return
	}
	id := mux.Vars(r)["id"]
	employee.ID, _ = strconv.Atoi(id)

	sqlStatement := "UPDATE users SET first_name = $2, last_name = $3, email = $4 WHERE id = $1"
	_, err := db.Exec(sqlStatement, employee.ID, employee.FirstName, employee.LastName, employee.Email)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "http://localhost:9090/users/", 301)
	defer db.Close()
}

func ReadUser(w http.ResponseWriter, r *http.Request) {
	db := dbConnect()
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

	employee := Employee{}

	for rows.Next() {
		var id int
		var first_name, last_name, email string
		err = rows.Scan(&id, &first_name, &last_name, &email)
		if err != nil {
			log.Fatal(err)
		}
		employee.ID = id
		employee.FirstName = first_name
		employee.LastName = last_name
		employee.Email = email
	}
	marshaled, _ := json.MarshalIndent(employee, "", " ")
	w.Write(marshaled)

	sqlStatement = "SELECT * FROM accounts WHERE id_user = $1"
	rows, err = db.Query(sqlStatement, id)
	if err != nil {
		log.Fatal(err)
	}

	account := Account{}
	result := []Account{}

	for rows.Next() {
		var id, id_user, value int
		err = rows.Scan(&id, &id_user, &value)
		if err != nil {
			log.Fatal(err)
		}
		account.ID = id
		account.ID_User = id_user
		account.Value = value
		result = append(result, account)
	}
	marshaled, _ = json.MarshalIndent(result, "", " ")
	w.Write(marshaled)
	defer db.Close()
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	db := dbConnect()
	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID mustn't be empty")
		return
	}
	var account Account

	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal Server Error111")
		return
	}

	sqlStatement := "INSERT INTO accounts (id_user, value)	VALUES ($1, $2)"
	_, err := db.Exec(sqlStatement, id, account.Value)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, r.URL.RequestURI(), 301)
	defer db.Close()
}
