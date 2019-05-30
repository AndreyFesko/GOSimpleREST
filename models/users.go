package models

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/rest/config"
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

func ListUsers() map[int]User {
	log.Println("List User function started")
	rows, err := config.DB.Query("SELECT * FROM users ORDER BY id DESC")
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
	return result
}

func CUser(e Employee) {
	log.Println("CUser function started")
	sqlStatement := "INSERT INTO users (first_name, last_name, email, phone)	VALUES ($1, $2, $3, $4)"
	_, err := config.DB.Exec(sqlStatement, e.FirstName, e.LastName, e.Email, e.Phone)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("User %s %s created", e.FirstName, e.LastName)
}

func DUser(id string) {
	log.Println("DUser function started")
	sqlStatement := "DELETE FROM users WHERE id = $1"
	_, err := config.DB.Exec(sqlStatement, id)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("User with id %s deleted", id)
}

func UUser(e Employee, id string) {
	log.Println("UUser function started")
	sqlStatement := "UPDATE users SET first_name = $2, last_name = $3, email = $4, phone = $5 WHERE id = $1"
	_, err := config.DB.Exec(sqlStatement, e.ID, e.FirstName, e.LastName, e.Email, e.Phone)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("User with id %s updated", id)
}

func RUser(id string) (e Employee) {
	log.Println("RUser function started")
	tx, err := config.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	sqlStatement := "SELECT * FROM users WHERE id = $1"
	rows, err := config.DB.Query(sqlStatement, id)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&e.ID, &e.FirstName, &e.LastName, &e.Email, &e.Phone)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}
	sqlStatement = "SELECT * FROM accounts WHERE id_user = $1 and active = $2"
	rows, err = config.DB.Query(sqlStatement, id, false)
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
	log.Printf("Data user with id %s printed", id)
	return e
}
