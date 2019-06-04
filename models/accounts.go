package models

import (
	"strconv"

	"github.com/golang/glog"
	_ "github.com/lib/pq"
	"github.com/rest/config"
)

// Account struct
type Account struct {
	ID       int  `json:"id"`
	UserID   int  `json:"user_id"`
	Value    int  `json:"value"`
	Currency int  `json:"currency"`
	Active   bool `json:"active"`
}

// CAccount create account
func CAccount(a Account, id string) {
	glog.V(3).Info("CAccount function started")
	sqlStatement := "INSERT INTO accounts (id_user, value, currency) VALUES ($1, $2, $3)"
	_, err := config.DB.Exec(sqlStatement, id, a.Value, a.Currency)
	if err != nil {
		glog.Fatal(err)
	}
}

// RAccount takes balance account
func RAccount(id, accID string) (map[string]int, error) {
	glog.V(3).Info("RAccount function started")
	result := make(map[string]int)
	sqlStatement := "SELECT value, currency FROM accounts WHERE id = $1 and active = $2"
	rows, err := config.DB.Query(sqlStatement, accID, false)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	var value, currency int
	for rows.Next() {
		err = rows.Scan(&value, &currency)
		if err != nil {
			return result, err
		}
	}
	result["id"], _ = strconv.Atoi(accID)
	result["value"] = value
	result["currency"] = currency
	return result, nil
}

// DAccount delete account
func DAccount(id, accID string) {
	glog.V(3).Info("DAccount function started")
	sqlStatement := "DELETE FROM accounts WHERE id = $1 and active = $2"
	_, err := config.DB.Exec(sqlStatement, accID, false)
	if err != nil {
		glog.Fatal(err)
	}
}

// GetUserIDFromAccount takes user ID by account ID
func GetUserIDFromAccount(accID int) (id int, err error) {
	glog.V(3).Info("GetUserIDFromAccount function started")
	a := Account{}
	sqlStatement := "SELECT * FROM accounts WHERE id = $1"
	rows, err := config.DB.Query(sqlStatement, accID)
	if err != nil {
		return a.UserID, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&a.ID, &a.UserID, &a.Value, &a.Currency, &a.Active)
		if err != nil {
			return a.UserID, err
		}
	}
	return a.UserID, nil
}
