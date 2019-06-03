package models

import (
	"github.com/golang/glog"
	_ "github.com/lib/pq"
	"github.com/rest/config"
	"gopkg.in/guregu/null.v3"
)

type Transaction struct {
	ID                    int      `json:"id"`
	Type                  string   `json:"type"`
	AccountForOperationID int      `json:"acc_for_operation"`
	RecievedIDAccount     null.Int `json:"recieved_id"`
	Value                 int      `json:"value"`
	Canceled              bool     `json:"canceled"`
	CanceledID            null.Int `json:"canceled_id"`
}

var commands = map[string]string{
	"Deposit":  "Deposit",
	"Withdraw": "Withdraw",
	"Transfer": "Transfer",
	"Cancel":   "Cancel",
}

func LTransactions() map[int]Transaction {
	glog.V(3).Info("LTransactions function started")
	rows, err := config.DB.Query("SELECT * FROM transactions ORDER BY id DESC")
	if err != nil {
		glog.Fatal(err)
	}
	defer rows.Close()
	result := make(map[int]Transaction)
	for rows.Next() {
		t := new(Transaction)
		err := rows.Scan(&t.ID, &t.Type, &t.AccountForOperationID, &t.RecievedIDAccount, &t.Value, &t.Canceled, &t.CanceledID)
		if err != nil {
			glog.Fatal(err)
		}
		result[t.ID] = *t
	}
	return result
}

func Transactions(t Transaction) {
	glog.V(3).Info("Transactions function started")
	tx, err := config.DB.Begin()
	if err != nil {
		glog.Fatal(err)
	}
	switch command := t.Type; command {
	case commands["Withdraw"]:
		sqlStatement := "INSERT INTO transactions (type, account_for_operation_id," +
			"value)	VALUES ($1, $2, $3)"
		_, err := config.DB.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.Value)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		sqlStatement = "UPDATE accounts SET value = value - $1	WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, t.Value, t.AccountForOperationID)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
	case commands["Deposit"]:
		sqlStatement := "INSERT INTO transactions (type, account_for_operation_id," +
			"value)	VALUES ($1, $2, $3)"
		_, err := config.DB.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.Value)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		sqlStatement = "UPDATE accounts SET value = value + $1	WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, t.Value, t.AccountForOperationID)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
	case commands["Transfer"]:
		sqlStatement := "INSERT INTO transactions (type, account_for_operation_id," +
			"recieve_id_account, value)	VALUES ($1, $2, $3, $4)"
		_, err = config.DB.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.RecievedIDAccount, t.Value)
		if err != nil {
			glog.Fatal(err)
		}
		stmt, err := tx.Prepare("UPDATE accounts SET value = value + $1	WHERE id = $2")
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		defer stmt.Close()
		if _, err := stmt.Exec(t.Value, t.RecievedIDAccount); err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		stmt, err = tx.Prepare("UPDATE accounts SET value = value - $1	WHERE id = $2")
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		defer stmt.Close()
		if _, err := stmt.Exec(t.Value, t.AccountForOperationID); err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
	}
	tx.Commit()
}

func CTransaction(id string) {
	glog.V(3).Info("CTransaction function started")
	tx, err := config.DB.Begin()
	if err != nil {
		glog.Fatal(err)
	}
	sqlStatement := "SELECT * FROM transactions WHERE id = $1"
	rows, err := config.DB.Query(sqlStatement, id)
	if err != nil {
		tx.Rollback()
		glog.Fatal(err)
	}
	defer rows.Close()
	t := new(Transaction)
	for rows.Next() {
		err = rows.Scan(&t.ID, &t.Type, &t.AccountForOperationID, &t.RecievedIDAccount, &t.Value, &t.Canceled, &t.CanceledID)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
	}
	if t.Canceled == true {
		tx.Rollback()
		glog.Fatal("Transaction was canceled before")
	}
	switch command := t.Type; command {
	case commands["Withdraw"]:
		sqlStatement = "UPDATE accounts SET value = value + $1	WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, t.Value, t.AccountForOperationID)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		sqlStatement = "UPDATE transactions SET canceled = $1 WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, true, t.ID)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		t.Type = commands["Cancel"]
		sqlStatement = "INSERT INTO transactions (type, account_for_operation_id," +
			"value, canceled_id)	VALUES ($1, $2, $3, $4)"
		_, err = config.DB.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.Value, t.ID)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		tl := GetLastTransaction()
		SendCancelMessage(tl)
	case commands["Deposit"]:
		sqlStatement = "UPDATE accounts SET value = value - $1	WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, t.Value, t.AccountForOperationID)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		sqlStatement = "UPDATE transactions SET canceled = $1 WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, true, t.ID)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		t.Type = commands["Cancel"]
		sqlStatement = "INSERT INTO transactions (type, account_for_operation_id," +
			"value, canceled_id)	VALUES ($1, $2, $3, $4)"
		_, err = config.DB.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.Value, t.ID)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		tl := GetLastTransaction()
		SendCancelMessage(tl)
	case commands["Transfer"]:
		sqlStatement = "UPDATE transactions SET canceled = $1 WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, true, t.ID)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		t.Type = commands["Cancel"]
		sqlStatement = "INSERT INTO transactions (type, account_for_operation_id," +
			"recieve_id_account, value, canceled_id)	VALUES ($1, $2, $3, $4, $5)"
		_, err = config.DB.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.RecievedIDAccount, t.Value, t.ID)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		tl := GetLastTransaction()
		SendCancelMessage(tl)
		sqlStatement = "UPDATE accounts SET value = value - $1	WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, t.Value, t.RecievedIDAccount)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		sqlStatement = "UPDATE accounts SET value = value + $1	WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, t.Value, t.AccountForOperationID)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
	}
	tx.Commit()
}

func GetLastTransaction() (t Transaction) {
	glog.V(3).Info("GetLastTransaction function started")
	rows, err := config.DB.Query("SELECT * FROM transactions WHERE id = (SELECT MAX(id) FROM transactions)")
	if err != nil {
		glog.Fatal(err)
	}
	defer rows.Close()
	t = Transaction{}
	for rows.Next() {
		err = rows.Scan(&t.ID, &t.Type, &t.AccountForOperationID, &t.RecievedIDAccount, &t.Value, &t.Canceled, &t.CanceledID)
		if err != nil {
			glog.Fatal(err)
		}
	}
	return t
}

func GetTransaction(id int64) (t Transaction) {
	glog.V(3).Info("GetTransactions function started")
	sqlStatement := "SELECT * FROM transactions WHERE id = $1"
	rows, err := config.DB.Query(sqlStatement, id)
	if err != nil {
		glog.Fatal(err)
	}
	defer rows.Close()
	t = Transaction{}
	for rows.Next() {
		err = rows.Scan(&t.ID, &t.Type, &t.AccountForOperationID, &t.RecievedIDAccount, &t.Value, &t.Canceled, &t.CanceledID)
		if err != nil {
			glog.Fatal(err)
		}
	}
	return t
}
