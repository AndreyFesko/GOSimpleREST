package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rest/models"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Get User function started")

	result := models.ListUsers()
	marshaled, _ := json.MarshalIndent(result, "", " ")
	w.Write(marshaled)
	log.Println("Users printed")
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Create User function started")
	var e models.Employee

	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal Server Error111")
		return
	}
	models.CUser(e)
	http.Redirect(w, r, "http://localhost:9090/users", 301)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete User function started")
	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID mustn't be empty")
		return
	}
	models.DUser(id)
	http.Redirect(w, r, "http://localhost:9090/users", 301)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Update User function started")
	var e models.Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal Server Error")
		return
	}
	id := mux.Vars(r)["id"]
	e.ID, _ = strconv.Atoi(id)
	models.UUser(e, id)
	http.Redirect(w, r, "http://localhost:9090/users", 301)
}

func ReadUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Read User function started")
	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID mustn't be empty")
		return
	}
	e := models.RUser(id)
	marshaled, _ := json.MarshalIndent(e, "", " ")
	w.Write(marshaled)
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	log.Println("Create Account function started")
	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID mustn't be empty")
		return
	}
	var a models.Account
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal Server Error111")
		return
		log.Println("Users printed")
	}
	models.CAccount(a, id)
	s := fmt.Sprintf("http://localhost:9090/users/%s", id)
	http.Redirect(w, r, s, 301)
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
	result := models.RAccount(id, accID)
	marshaled, _ := json.MarshalIndent(result, "", " ")
	w.Write(marshaled)
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete Account function started")
	id := mux.Vars(r)["id"]
	accID := mux.Vars(r)["acc_id"]
	models.DAccount(id, accID)
	s := fmt.Sprintf("http://localhost:9090/users/%s", id)
	http.Redirect(w, r, s, 301)
}

func ListTransactions(w http.ResponseWriter, r *http.Request) {
	log.Println("List Transactions function started")
	result := models.LTransactions()
	marshaled, _ := json.MarshalIndent(result, "", " ")
	w.Write(marshaled)
}

func CreateTransactions(w http.ResponseWriter, r *http.Request) {
	log.Println("Transactions function started")
	var t models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal Server Error111")
		return
	}
	models.Transactions(t)
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
	models.CTransaction(id)
	http.Redirect(w, r, "http://localhost:9090/transactions", 301)
}
