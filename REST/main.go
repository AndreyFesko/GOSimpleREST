package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	s.HandleFunc("/", CreateNewUser).Methods(http.MethodPost)

	http.Handle("/", r)
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatal(err)
	}
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	db := dbConnect()
	rows, err := db.Query("select * from users order by id desc")
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

func CreateNewUser(w http.ResponseWriter, r *http.Request) {
	var user Employee
	db := dbConnect()
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal Server Error111")
		return
	}

	sqlStatement := `INSERT INTO users (id, first_name, last_name, email)
	VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(sqlStatement, user.ID, user.FirstName, user.LastName, user.Email)
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("select * from users order by id desc")
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
