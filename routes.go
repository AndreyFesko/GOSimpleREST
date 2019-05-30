package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{"GetUsers", "GET", "/users", GetUsers},
	Route{"CreateUser", "POST", "/users", CreateUser},
	Route{"UpdateUser", "PATCH", "/users/{id}", UpdateUser},
	Route{"DeleteUser", "DELETE", "/users/{id}", DeleteUser},
	Route{"ReadUsers", "GET", "/users/{id}", ReadUser},

	Route{"CreateAccount", "POST", "/users/{id}/accounts", CreateAccount},
	Route{"ReadAccount", "GET", "/users/{id}/accounts/{acc_id}/balance", ReadAccount},
	Route{"DeleteAccount", "DELETE", "/users/{id}/accounts/{acc_id}", DeleteAccount},

	Route{"ListTransactions", "GET", "/transactions", ListTransactions},
	Route{"CreateTransactions", "POST", "/transactions", CreateTransactions},
	Route{"CancelTransaction", "DELETE", "/transactions/{id}", CancelTransaction},
}
