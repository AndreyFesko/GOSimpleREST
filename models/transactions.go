package models

import (
	"log"

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
	log.Println("LTransactions function started")
	rows, err := config.DB.Query("SELECT * FROM transactions ORDER BY id DESC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	result := make(map[int]Transaction)
	for rows.Next() {
		t := new(Transaction)
		err := rows.Scan(&t.ID, &t.Type, &t.AccountForOperationID, &t.RecievedIDAccount, &t.Value, &t.Canceled, &t.CanceledID)
		if err != nil {
			log.Fatal(err)
		}
		result[t.ID] = *t
	}
	log.Println("List Transactions printed")
	return result
}

func Transactions(t Transaction) {
	log.Println("Transactions function started")
	tx, err := config.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	switch command := t.Type; command {
	case commands["Withdraw"]:
		sqlStatement := "INSERT INTO transactions (type, account_for_operation_id," +
			"value)	VALUES ($1, $2, $3)"
		_, err := config.DB.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.Value)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}

		sqlStatement = "UPDATE accounts SET value = value - $1	WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, t.Value, t.AccountForOperationID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		log.Println("Transactions Withdraw completed")
	case commands["Deposit"]:
		sqlStatement := "INSERT INTO transactions (type, account_for_operation_id," +
			"value)	VALUES ($1, $2, $3)"
		_, err := config.DB.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.Value)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		sqlStatement = "UPDATE accounts SET value = value + $1	WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, t.Value, t.AccountForOperationID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		log.Println("Transactions Deposit completed")
	case commands["Transfer"]:
		sqlStatement := "INSERT INTO transactions (type, account_for_operation_id," +
			"recieve_id_account, value)	VALUES ($1, $2, $3, $4)"
		_, err = config.DB.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.RecievedIDAccount, t.Value)
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare("UPDATE accounts SET value = value + $1	WHERE id = $2")
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		defer stmt.Close()
		if _, err := stmt.Exec(t.Value, t.RecievedIDAccount); err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		stmt, err = tx.Prepare("UPDATE accounts SET value = value - $1	WHERE id = $2")
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		defer stmt.Close()
		if _, err := stmt.Exec(t.Value, t.AccountForOperationID); err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		log.Println("Transactions Transfer completed")
	}
	tx.Commit()
}

func CTransaction(id string) {
	log.Println("CTransaction function started")
	tx, err := config.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	sqlStatement := "SELECT * FROM transactions WHERE id = $1"
	rows, err := config.DB.Query(sqlStatement, id)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	defer rows.Close()
	t := new(Transaction)
	for rows.Next() {
		err = rows.Scan(&t.ID, &t.Type, &t.AccountForOperationID, &t.RecievedIDAccount, &t.Value, &t.Canceled, &t.CanceledID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}
	if t.Canceled == true {
		tx.Rollback()
		log.Fatal("Transaction was canceled before")
	}
	switch command := t.Type; command {
	case commands["Withdraw"]:
		log.Println("Cancel withdraw")
		sqlStatement = "UPDATE accounts SET value = value + $1	WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, t.Value, t.AccountForOperationID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		sqlStatement = "UPDATE transactions SET canceled = $1 WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, true, t.ID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		t.Type = commands["Cancel"]
		sqlStatement = "INSERT INTO transactions (type, account_for_operation_id," +
			"value, canceled_id)	VALUES ($1, $2, $3, $4)"
		_, err = config.DB.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.Value, t.ID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	case commands["Deposit"]:
		log.Println("Cancel deposit")
		sqlStatement = "UPDATE accounts SET value = value - $1	WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, t.Value, t.AccountForOperationID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		sqlStatement = "UPDATE transactions SET canceled = $1 WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, true, t.ID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		t.Type = commands["Cancel"]
		sqlStatement = "INSERT INTO transactions (type, account_for_operation_id," +
			"value, canceled_id)	VALUES ($1, $2, $3, $4)"
		_, err = config.DB.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.Value, t.ID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	case commands["Transfer"]:
		log.Println("Cancel transfer")
		sqlStatement = "UPDATE transactions SET canceled = $1 WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, true, t.ID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		t.Type = commands["Cancel"]
		sqlStatement = "INSERT INTO transactions (type, account_for_operation_id," +
			"recieve_id_account, value, canceled_id)	VALUES ($1, $2, $3, $4, $5)"
		_, err = config.DB.Exec(sqlStatement, t.Type, t.AccountForOperationID, t.RecievedIDAccount, t.Value, t.ID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		sqlStatement = "UPDATE accounts SET value = value - $1	WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, t.Value, t.RecievedIDAccount)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		sqlStatement = "UPDATE accounts SET value = value + $1	WHERE id = $2"
		_, err = config.DB.Exec(sqlStatement, t.Value, t.AccountForOperationID)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}
	tx.Commit()
	log.Println("Transactions canceled")
}
