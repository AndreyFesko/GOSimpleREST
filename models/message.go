package models

import (
	"fmt"
	"strconv"

	"github.com/golang/glog"
)

// SendMessage func
func SendMessage(t Transaction) {
	switch command := t.Type; command {
	case commands["Deposit"], commands["Withdraw"]:
		userID, _ := GetUserIDFromAccount(t.AccountForOperationID)
		user, _ := RUser(strconv.Itoa(userID))
		strGen := "%s: send to %s. \nYour account is %d %s by the amount %d"
		str := fmt.Sprintf(strGen, "EMAIL", user.Email, t.AccountForOperationID, t.Type, t.Value)
		glog.Info(str)
		if user.Phone.Valid != false {
			str := fmt.Sprintf(strGen, "SMS", user.Phone.String, t.AccountForOperationID, t.Type, t.Value)
			glog.Info(str)
		}
	case commands["Transfer"]:
		userOperationID, _ := GetUserIDFromAccount(t.AccountForOperationID)
		userRecievedID, _ := GetUserIDFromAccount(int(t.RecievedIDAccount.Int64))
		userOperation, _ := RUser(strconv.Itoa(userOperationID))
		userRecieved, _ := RUser(strconv.Itoa(userRecievedID))
		strGen := "%s: send to %s.\nYour account has been transferred to account %d for the amount of %d"
		str := fmt.Sprintf(strGen, "EMAIL", userOperation.Email, t.RecievedIDAccount.Int64, t.Value)
		glog.Info(str)
		if userOperation.Phone.Valid != false {
			str := fmt.Sprintf(strGen, "SMS", userOperation.Phone.String, t.RecievedIDAccount.Int64, t.Value)
			glog.Info(str)
		}
		str = fmt.Sprintf(strGen, "EMAIL", userRecieved.Email, t.AccountForOperationID, t.Value)
		glog.Info(str)
		if userRecieved.Phone.Valid != false {
			str := fmt.Sprintf(strGen, "SMS", userRecieved.Phone.String, t.AccountForOperationID, t.Value)
			glog.Info(str)
		}
	}
}

// SendCancelMessage func
func SendCancelMessage(t Transaction) {
	canceledID := t.CanceledID.Int64
	tl := GetTransaction(canceledID)
	strGen := "%s: send to %s.\n%s operation canceled. Your account is %d withdraw by the amount %d"
	switch command := tl.Type; command {
	case commands["Deposit"]:
		userID, _ := GetUserIDFromAccount(tl.AccountForOperationID)
		user, _ := RUser(strconv.Itoa(userID))
		str := fmt.Sprintf(strGen, "EMAIL", "Deposit", user.Email, tl.AccountForOperationID, tl.Value)
		glog.Info(str)
		if user.Phone.Valid != false {
			str := fmt.Sprintf(strGen, "SMS", "Deposit", user.Phone.String, tl.AccountForOperationID, tl.Value)
			glog.Info(str)
		}
	case commands["Withdraw"]:
		userID, _ := GetUserIDFromAccount(tl.AccountForOperationID)
		user, _ := RUser(strconv.Itoa(userID))
		str := fmt.Sprintf(strGen, "EMAIL", "Withdraw", user.Email, tl.AccountForOperationID, tl.Value)
		glog.Info(str)
		if user.Phone.Valid != false {
			str := fmt.Sprintf(strGen, "SMSL", "Withdraw", user.Phone.String, tl.AccountForOperationID, tl.Value)
			glog.Info(str)
		}
	case commands["Transfer"]:
		userOperationID, _ := GetUserIDFromAccount(t.AccountForOperationID)
		userRecievedID, _ := GetUserIDFromAccount(int(t.RecievedIDAccount.Int64))
		userOperation, _ := RUser(strconv.Itoa(userOperationID))
		userRecieved, _ := RUser(strconv.Itoa(userRecievedID))
		str := fmt.Sprintf(strGen, "EMAIL", "Transfer", userOperation.Email, tl.RecievedIDAccount.Int64, tl.Value)
		glog.Info(str)
		if userOperation.Phone.Valid != false {
			str := fmt.Sprintf(strGen, "SMS", "Transfer", userOperation.Phone.String, tl.RecievedIDAccount.Int64, tl.Value)
			glog.Info(str)
		}
		str = fmt.Sprintf(strGen, "EMAIL", "Transfer", userRecieved.Email, tl.AccountForOperationID, tl.Value)
		glog.Info(str)
		if userRecieved.Phone.Valid != false {
			str := fmt.Sprintf(strGen, "SMS", "Transfer", userRecieved.Phone.String, tl.AccountForOperationID, tl.Value)
			glog.Info(str)
		}
	}
}
