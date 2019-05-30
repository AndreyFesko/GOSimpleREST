package models

import (
	"log"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/rest/config"
)

type Account struct {
	ID       int  `json:"id"`
	UserID   int  `json:"user_id"`
	Value    int  `json:"value"`
	Currency int  `json:"currency"`
	Active   bool `json:"active"`
}

func CAccount(a Account, id string) {
	log.Println("CAccount function started")
	sqlStatement := "INSERT INTO accounts (id_user, value, currency) VALUES ($1, $2, $3)"
	_, err := config.DB.Exec(sqlStatement, id, a.Value, a.Currency)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Account for user with id %s created", id)
}

func RAccount(id, accID string) map[string]int {
	log.Println("RAccount function started")
	sqlStatement := "SELECT value, currency FROM accounts WHERE id = $1 and active = $2"
	rows, err := config.DB.Query(sqlStatement, accID, false)
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
	log.Printf("Account with id %s printed", accID)
	return result
}

func DAccount(id, accID string) {
	log.Println("DAccount function started")
	sqlStatement := "DELETE FROM accounts WHERE id = $1 and active = $2"
	_, err := config.DB.Exec(sqlStatement, accID, false)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Account %s users with id %s deleted", accID, id)
}
