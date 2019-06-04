package models

import (
	"github.com/golang/glog"
	_ "github.com/lib/pq"
	"github.com/rest/config"
	"gopkg.in/guregu/null.v3"
)

// Transaction struct
type Transaction struct {
	ID                    int      `json:"id"`
	Type                  string   `json:"type"`
	AccountForOperationID int      `json:"acc_for_operation"`
	RecievedIDAccount     null.Int `json:"recieved_id"`
	Value                 int      `json:"value"`
	Canceled              bool     `json:"canceled"`
	CanceledID            null.Int `json:"canceled_id"`
}

// Commands bank operations
var commands = map[string]string{
	"Deposit":  "Deposit",
	"Withdraw": "Withdraw",
	"Transfer": "Transfer",
	"Cancel":   "Cancel",
}

// LTransactions takes all transactions
func LTransactions() (map[int]Transaction, error) {
	glog.V(3).Info("LTransactions function started")
	result := make(map[int]Transaction)
	rows, err := config.DB.Query("SELECT * FROM transactions ORDER BY id DESC")
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		t := new(Transaction)
		err := rows.Scan(&t.ID, &t.Type, &t.AccountForOperationID, &t.RecievedIDAccount, &t.Value, &t.Canceled, &t.CanceledID)
		if err != nil {
			return result, err
		}
		result[t.ID] = *t
	}
	return result, nil
}

// Transactions create
func Transactions(t Transaction) {
	glog.V(3).Info("Transactions function started")
	tx, err := config.DB.Begin()
	if err != nil {
		glog.Fatal(err)
	}
	switch command := t.Type; command {
	case commands["Withdraw"]:
		glog.V(4).Info("Withdraw transaction")
		// Saved transaction
		sqlStatement := "INSERT INTO transactions (type, account_for_operation_id," +
			"value)	VALUES ($1, $2, $3)"
		_, err := config.DB.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.Value)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		glog.V(4).Info("Transaction saved")
		// Withdraw account
		sqlStatement = "UPDATE accounts SET value = value - $1	WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, t.Value, t.AccountForOperationID)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		glog.V(4).Info("Account ", t.AccountForOperationID, " withdrawn ", t.Value)
	case commands["Deposit"]:
		glog.V(4).Info("Deposit transaction")
		// Saved transaction
		sqlStatement := "INSERT INTO transactions (type, account_for_operation_id," +
			"value)	VALUES ($1, $2, $3)"
		_, err := config.DB.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.Value)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		glog.V(4).Info("Transaction saved")
		// Deposit account
		sqlStatement = "UPDATE accounts SET value = value + $1	WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, t.Value, t.AccountForOperationID)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		glog.V(4).Info("Account ", t.AccountForOperationID, " deposit ", t.Value)
	case commands["Transfer"]:
		glog.V(4).Info("Transfer transaction")
		// Saved transaction
		sqlStatement := "INSERT INTO transactions (type, account_for_operation_id," +
			"recieve_id_account, value)	VALUES ($1, $2, $3, $4)"
		_, err = config.DB.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.RecievedIDAccount, t.Value)
		if err != nil {
			glog.Fatal(err)
		}
		glog.V(4).Info("Transaction saved")
		// Deposit account
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
		glog.V(4).Info("Account ", t.RecievedIDAccount, " deposit ", t.Value)
		// Withdraw account
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
		glog.V(4).Info("Account ", t.AccountForOperationID, " withdrawn ", t.Value)
	}
	tx.Commit()
}

// CTransaction cancel
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
	glog.V(4).Info("Get transaction ", id, " from database")
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
		// Deposit account
		sqlStatement = "UPDATE accounts SET value = value + $1	WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, t.Value, t.AccountForOperationID)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		glog.V(4).Info("Account ", t.AccountForOperationID, " deposit ", t.Value)
		// Marked cancel
		sqlStatement = "UPDATE transactions SET canceled = $1 WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, true, t.ID)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		glog.V(4).Info("Transaction ", t.ID, " marked as canceled")
		t.Type = commands["Cancel"]
		// Saved transaction
		sqlStatement = "INSERT INTO transactions (type, account_for_operation_id," +
			"value, canceled_id)	VALUES ($1, $2, $3, $4)"
		_, err = config.DB.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.Value, t.ID)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		glog.V(4).Info("Cancel transaction saved")
		tl := GetLastTransaction()
		SendCancelMessage(tl)
	case commands["Deposit"]:
		// Withdraw account
		sqlStatement = "UPDATE accounts SET value = value - $1	WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, t.Value, t.AccountForOperationID)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		glog.V(4).Info("Account ", t.AccountForOperationID, " withdrawn ", t.Value)
		// Marked cancel
		sqlStatement = "UPDATE transactions SET canceled = $1 WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, true, t.ID)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		glog.V(4).Info("Transaction ", t.ID, " marked as canceled")
		t.Type = commands["Cancel"]
		// Saved transaction
		sqlStatement = "INSERT INTO transactions (type, account_for_operation_id," +
			"value, canceled_id)	VALUES ($1, $2, $3, $4)"
		_, err = config.DB.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.Value, t.ID)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		glog.V(4).Info("Cancel transaction saved")
		tl := GetLastTransaction()
		SendCancelMessage(tl)
	case commands["Transfer"]:
		// Marked cancel
		sqlStatement = "UPDATE transactions SET canceled = $1 WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, true, t.ID)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		glog.V(4).Info("Transaction ", t.ID, " marked as canceled")
		t.Type = commands["Cancel"]
		// Saved transaction
		sqlStatement = "INSERT INTO transactions (type, account_for_operation_id," +
			"recieve_id_account, value, canceled_id)	VALUES ($1, $2, $3, $4, $5)"
		_, err = config.DB.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.RecievedIDAccount, t.Value, t.ID)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		glog.V(4).Info("Cancel transaction saved")
		tl := GetLastTransaction()
		SendCancelMessage(tl)
		// Withdraw account
		sqlStatement = "UPDATE accounts SET value = value - $1	WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, t.Value, t.RecievedIDAccount)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		glog.V(4).Info("Account ", t.RecievedIDAccount, " withdrawn ", t.Value)
		sqlStatement = "UPDATE accounts SET value = value + $1	WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, t.Value, t.AccountForOperationID)
		if err != nil {
			tx.Rollback()
			glog.Fatal(err)
		}
		glog.V(4).Info("Account ", t.AccountForOperationID, " deposit ", t.Value)
	}
	tx.Commit()
}

// GetLastTransaction takes last inserted transaction
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

// GetTransaction takes transaction
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
