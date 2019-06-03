package models

import (
	"strconv"

	"github.com/golang/glog"
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
	glog.V(3).Info("CAccount function started")
	sqlStatement := "INSERT INTO accounts (id_user, value, currency) VALUES ($1, $2, $3)"
	_, err := config.DB.Exec(sqlStatement, id, a.Value, a.Currency)
	if err != nil {
		glog.Fatal(err)
	}
}

func RAccount(id, accID string) map[string]int {
	glog.V(3).Info("RAccount function started")
	sqlStatement := "SELECT value, currency FROM accounts WHERE id = $1 and active = $2"
	rows, err := config.DB.Query(sqlStatement, accID, false)
	if err != nil {
		glog.Fatal(err)
	}
	defer rows.Close()
	var value, currency int
	for rows.Next() {
		err = rows.Scan(&value, &currency)
		if err != nil {
			glog.Fatal(err)
		}
	}
	result := make(map[string]int)
	result["id"], _ = strconv.Atoi(accID)
	result["value"] = value
	result["currency"] = currency
	return result
}

func DAccount(id, accID string) {
	glog.V(3).Info("DAccount function started")
	sqlStatement := "DELETE FROM accounts WHERE id = $1 and active = $2"
	_, err := config.DB.Exec(sqlStatement, accID, false)
	if err != nil {
		glog.Fatal(err)
	}
}

func GetUserIDFromAccount(acc_id int) (id int) {
	glog.V(3).Info("GetUserIDFromAccount function started")
	sqlStatement := "SELECT * FROM accounts WHERE id = $1"
	rows, err := config.DB.Query(sqlStatement, acc_id)
	if err != nil {
		glog.Fatal(err)
	}
	defer rows.Close()
	a := Account{}
	for rows.Next() {
		err = rows.Scan(&a.ID, &a.UserID, &a.Value, &a.Currency, &a.Active)
		if err != nil {
			glog.Fatal(err)
		}
	}
	return a.UserID
}
