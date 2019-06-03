package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rest/models"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	glog.V(3).Info("Handler GetUsers")

	result := models.ListUsers()
	marshaled, _ := json.MarshalIndent(result, "", " ")
	w.Write(marshaled)
	glog.V(2).Info("Users printed")
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	glog.V(3).Info("Handler CreateUser")
	var e models.Employee

	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		glog.Fatal(err)
	}
	models.CUser(e)
	http.Redirect(w, r, "http://localhost:9090/users", 301)
	glog.V(2).Info("User created")
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	glog.V(3).Info("Handler DeleteUser")
	id := mux.Vars(r)["id"]
	if id == "" {
		glog.Fatal("ID mustn't be empty")
	}
	models.DUser(id)
	http.Redirect(w, r, "http://localhost:9090/users", 301)
	glog.V(2).Info("User removed")
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	glog.V(3).Info("Handler UpdateUser")
	var e models.Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		glog.Fatal(err)
	}
	id := mux.Vars(r)["id"]
	e.ID, _ = strconv.Atoi(id)
	models.UUser(e, id)
	http.Redirect(w, r, "http://localhost:9090/users", 301)
	glog.V(2).Info("User updated")
}

func ReadUser(w http.ResponseWriter, r *http.Request) {
	glog.V(3).Info("Handler ReadUser")
	id := mux.Vars(r)["id"]
	if id == "" {
		glog.Fatal("ID must'n be empty")
	}
	e := models.RUser(id)
	marshaled, _ := json.MarshalIndent(e, "", " ")
	w.Write(marshaled)
	glog.V(2).Info("Data user printed")
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	glog.V(3).Info("Handler CreateAccount")
	id := mux.Vars(r)["id"]
	if id == "" {
		glog.Fatal("ID mustn't be empty")
	}
	var a models.Account
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		glog.Fatal(err)
	}
	models.CAccount(a, id)
	s := fmt.Sprintf("http://localhost:9090/users/%s", id)
	http.Redirect(w, r, s, 301)
	glog.V(2).Info("Account created")
}

func ReadAccount(w http.ResponseWriter, r *http.Request) {
	glog.V(3).Info("Handler ReadAccount")
	id := mux.Vars(r)["id"]
	accID := mux.Vars(r)["acc_id"]
	if id == "" {
		glog.Fatal("ID mustn't be empty")
	}
	result := models.RAccount(id, accID)
	marshaled, _ := json.MarshalIndent(result, "", " ")
	w.Write(marshaled)
	glog.V(2).Info("Data account printed")
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	glog.V(3).Info("Handler DeleteAccount")
	id := mux.Vars(r)["id"]
	accID := mux.Vars(r)["acc_id"]
	models.DAccount(id, accID)
	s := fmt.Sprintf("http://localhost:9090/users/%s", id)
	http.Redirect(w, r, s, 301)
	glog.V(2).Info("Account removed")
}

func ListTransactions(w http.ResponseWriter, r *http.Request) {
	glog.V(3).Info("Handler ListTransactions")
	result := models.LTransactions()
	marshaled, _ := json.MarshalIndent(result, "", " ")
	w.Write(marshaled)
	glog.V(2).Info("List transactions printed")
}

func CreateTransactions(w http.ResponseWriter, r *http.Request) {
	glog.V(3).Info("Handler CreateTransactions")
	var t models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		glog.Fatal(err)
	}
	models.Transactions(t)
	models.SendMessage(t)
	http.Redirect(w, r, r.URL.RequestURI(), 301)
	glog.V(2).Info("Transaction created")
}

func CancelTransaction(w http.ResponseWriter, r *http.Request) {
	glog.V(3).Info("Handler CancelTransaction")
	id := mux.Vars(r)["id"]
	if id == "" {
		glog.Fatal("ID mustn't be empty")
	}
	models.CTransaction(id)
	http.Redirect(w, r, "http://localhost:9090/transactions", 301)
	glog.V(2).Info("Transaction canceled")
}
