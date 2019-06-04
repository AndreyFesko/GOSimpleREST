package models

import (
	"github.com/golang/glog"
	_ "github.com/lib/pq"
	"github.com/rest/config"
	"gopkg.in/guregu/null.v3"
)

// Employee Struct has accounts information
type Employee struct {
	User
	Accounts []Account `json:"accounts"`
}

// User Struct
type User struct {
	ID        int         `json:"id"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Email     string      `json:"email"`
	Phone     null.String `json:"phone"`
}

// ListUsers takes data from db about all users
func ListUsers() (map[int]User, error) {
	glog.V(3).Info("Function list users started")
	result := make(map[int]User)
	rows, err := config.DB.Query("SELECT * FROM users ORDER BY id DESC")
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		e := new(User)
		err := rows.Scan(&e.ID, &e.FirstName, &e.LastName, &e.Email, &e.Phone)
		if err != nil {
			return result, err
		}
		result[e.ID] = *e
	}
	return result, nil
}

// CUser create user and save to db
func CUser(e Employee) {
	glog.V(3).Info("CUser function started")
	sqlStatement := "INSERT INTO users (first_name, last_name, email, phone)	VALUES ($1, $2, $3, $4)"
	_, err := config.DB.Exec(sqlStatement, e.FirstName, e.LastName, e.Email, e.Phone)
	if err != nil {
		glog.Fatal(err)
	}
}

// DUser delete user from db
func DUser(id string) {
	glog.V(3).Info("DUser function started")
	sqlStatement := "DELETE FROM users WHERE id = $1"
	_, err := config.DB.Exec(sqlStatement, id)
	if err != nil {
		glog.Fatal(err)
	}
}

// UUser update user
func UUser(e Employee, id string) {
	glog.V(3).Info("UUser function started")
	sqlStatement := "UPDATE users SET first_name = $2, last_name = $3, email = $4, phone = $5 WHERE id = $1"
	_, err := config.DB.Exec(sqlStatement, e.ID, e.FirstName, e.LastName, e.Email, e.Phone)
	if err != nil {
		glog.Fatal(err)
	}
}

// RUser takes data specified user
func RUser(id string) (e Employee, err error) {
	glog.V(3).Info("RUser function started")
	tx, err := config.DB.Begin()
	if err != nil {
		return e, err
	}
	// takes data about user
	sqlStatement := "SELECT * FROM users WHERE id = $1"
	rows, err := config.DB.Query(sqlStatement, id)
	if err != nil {
		tx.Rollback()
		return e, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&e.ID, &e.FirstName, &e.LastName, &e.Email, &e.Phone)
		if err != nil {
			tx.Rollback()
			return e, err
		}
	}
	// takes data about users accounts
	sqlStatement = "SELECT * FROM accounts WHERE id_user = $1 and active = $2"
	rows, err = config.DB.Query(sqlStatement, id, false)
	if err != nil {
		tx.Rollback()
		return e, err
	}
	defer rows.Close()
	for rows.Next() {
		a := new(Account)
		err = rows.Scan(&a.ID, &a.UserID, &a.Value, &a.Currency, &a.Active)
		if err != nil {
			tx.Rollback()
			return e, err
		}
		e.Accounts = append(e.Accounts, *a)
	}
	tx.Commit()
	return e, nil
}
