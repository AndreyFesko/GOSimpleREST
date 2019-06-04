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

	result, err := models.ListUsers()
	if err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "StatusInternalServerError")
		return
	}
	marshaled, _ := json.MarshalIndent(result, "", " ")
	w.Write(marshaled)
	w.WriteHeader(http.StatusOK)
	glog.V(2).Info("Users printed")
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	glog.V(3).Info("Handler CreateUser")
	var e models.Employee

	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Request")
		return
	}
	models.CUser(e)
	w.WriteHeader(http.StatusCreated)
	glog.V(2).Info("User created")
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	glog.V(3).Info("Handler DeleteUser")
	id := mux.Vars(r)["id"]
	if id == "" {
		glog.Error("ID mustn't be empty")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Request")
		return
	}
	models.DUser(id)
	w.WriteHeader(http.StatusOK)
	glog.V(2).Info("User removed")
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	glog.V(3).Info("Handler UpdateUser")
	var e models.Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Request")
		return
	}
	id := mux.Vars(r)["id"]
	e.ID, _ = strconv.Atoi(id)
	models.UUser(e, id)
	w.WriteHeader(http.StatusOK)
	glog.V(2).Info("User updated")
}

func ReadUser(w http.ResponseWriter, r *http.Request) {
	glog.V(3).Info("Handler ReadUser")
	id := mux.Vars(r)["id"]
	if id == "" {
		glog.Error("ID mustn't be empty")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Request")
		return
	}
	e, err := models.RUser(id)
	if err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "StatusInternalServerError")
		return
	}
	marshaled, _ := json.MarshalIndent(e, "", " ")
	w.Write(marshaled)
	w.WriteHeader(http.StatusOK)
	glog.V(2).Info("Data user printed")
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	glog.V(3).Info("Handler CreateAccount")
	id := mux.Vars(r)["id"]
	if id == "" {
		glog.Error("ID mustn't be empty")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Request")
		return
	}
	var a models.Account
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Request")
		return
	}
	models.CAccount(a, id)
	w.WriteHeader(http.StatusCreated)
	glog.V(2).Info("Account created")
}

func ReadAccount(w http.ResponseWriter, r *http.Request) {
	glog.V(3).Info("Handler ReadAccount")
	id := mux.Vars(r)["id"]
	accID := mux.Vars(r)["acc_id"]
	if id == "" {
		glog.Error("ID mustn't be empty")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Request")
		return
	}
	result, err := models.RAccount(id, accID)
	if err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "StatusInternalServerError")
		return
	}
	marshaled, _ := json.MarshalIndent(result, "", " ")
	w.Write(marshaled)
	w.WriteHeader(http.StatusOK)
	glog.V(2).Info("Data account printed")
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	glog.V(3).Info("Handler DeleteAccount")
	id := mux.Vars(r)["id"]
	accID := mux.Vars(r)["acc_id"]
	models.DAccount(id, accID)
	w.WriteHeader(http.StatusOK)
	glog.V(2).Info("Account removed")
}

func ListTransactions(w http.ResponseWriter, r *http.Request) {
	glog.V(3).Info("Handler ListTransactions")
	result, err := models.LTransactions()
	if err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "StatusInternalServerError")
		return
	}
	marshaled, _ := json.MarshalIndent(result, "", " ")
	w.Write(marshaled)
	w.WriteHeader(http.StatusOK)
	glog.V(2).Info("List transactions printed")
}

func CreateTransactions(w http.ResponseWriter, r *http.Request) {
	glog.V(3).Info("Handler CreateTransactions")
	var t models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Request")
		return
	}
	models.Transactions(t)
	models.SendMessage(t)
	w.WriteHeader(http.StatusCreated)
	glog.V(2).Info("Transaction created")
}

func CancelTransaction(w http.ResponseWriter, r *http.Request) {
	glog.V(3).Info("Handler CancelTransaction")
	id := mux.Vars(r)["id"]
	if id == "" {
		glog.Error("ID mustn't be empty")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Request")
		return
	}
	models.CTransaction(id)
	w.WriteHeader(http.StatusOK)
	glog.V(2).Info("Transaction canceled")
}
