package models

import (
	"fmt"
	"strconv"

	"github.com/golang/glog"
)

func SendMessage(t Transaction) {
	switch command := t.Type; command {
	case commands["Deposit"], commands["Withdraw"]:
		userID := GetUserIDFromAccount(t.AccountForOperationID)
		user := RUser(strconv.Itoa(userID))
		str := fmt.Sprintf("EMAIL: send to %s. \nYour account is %d %s by the amount %d", user.Email, t.AccountForOperationID, t.Type, t.Value)
		glog.Info(str)
		if user.Phone.Valid != false {
			str := fmt.Sprintf("SMS: send to %s.\nYour account is %d %s by the amount %d", user.Phone.String, t.AccountForOperationID, t.Type, t.Value)
			glog.Info(str)
		}
	case commands["Transfer"]:
		userOperationID := GetUserIDFromAccount(t.AccountForOperationID)
		userRecievedID := GetUserIDFromAccount(int(t.RecievedIDAccount.Int64))
		userOperation := RUser(strconv.Itoa(userOperationID))
		userRecieved := RUser(strconv.Itoa(userRecievedID))
		str := fmt.Sprintf("EMAIL: send to %s.\nYour account has been transferred to account %d for the amount of %d", userOperation.Email, t.RecievedIDAccount.Int64, t.Value)
		glog.Info(str)
		if userOperation.Phone.Valid != false {
			str := fmt.Sprintf("SMS: send to %s.\nYour account has been transferred to account %d for the amount of %d", userOperation.Phone.String, t.RecievedIDAccount.Int64, t.Value)
			glog.Info(str)
		}
		str = fmt.Sprintf("EMAIL: send to %s.\nYour account has been transferred from account %d for the amount of %d", userRecieved.Email, t.AccountForOperationID, t.Value)
		glog.Info(str)
		if userRecieved.Phone.Valid != false {
			str := fmt.Sprintf("SMS: send to %s.\nYour account has been transferred from account %d for the amount of %d", userRecieved.Phone.String, t.AccountForOperationID, t.Value)
			glog.Info(str)
		}
	}
}

func SendCancelMessage(t Transaction) {
	canceledID := t.CanceledID.Int64
	tl := GetTransaction(canceledID)
	switch command := tl.Type; command {
	case commands["Deposit"]:
		userID := GetUserIDFromAccount(tl.AccountForOperationID)
		user := RUser(strconv.Itoa(userID))
		str := fmt.Sprintf("EMAIL: send to %s.\nDeposit operation canceled. Your account is %d withdraw by the amount %d", user.Email, tl.AccountForOperationID, tl.Value)
		glog.Info(str)
		if user.Phone.Valid != false {
			str := fmt.Sprintf("SMS: send to %s.\nDeposit operation canceled. Your account is %d withdraw by the amount %d", user.Phone.String, tl.AccountForOperationID, tl.Value)
			glog.Info(str)
		}
	case commands["Withdraw"]:
		userID := GetUserIDFromAccount(tl.AccountForOperationID)
		user := RUser(strconv.Itoa(userID))
		str := fmt.Sprintf("EMAIL: send to %s.\nWithdraw operation canceled. Your account is %d deposit by the amount %d", user.Email, tl.AccountForOperationID, tl.Value)
		glog.Info(str)
		if user.Phone.Valid != false {
			str := fmt.Sprintf("SMS: send to %s.\nWithdraw operation canceled. Your account is %d deposit by the amount %d", user.Phone.String, tl.AccountForOperationID, tl.Value)
			glog.Info(str)
		}
	case commands["Transfer"]:
		userOperationID := GetUserIDFromAccount(t.AccountForOperationID)
		userRecievedID := GetUserIDFromAccount(int(t.RecievedIDAccount.Int64))
		userOperation := RUser(strconv.Itoa(userOperationID))
		userRecieved := RUser(strconv.Itoa(userRecievedID))
		str := fmt.Sprintf("EMAIL: send to %s.\nTransfer operation canceled. Your account has been transferred from account %d for the amount of %d", userOperation.Email, tl.RecievedIDAccount.Int64, tl.Value)
		glog.Info(str)
		if userOperation.Phone.Valid != false {
			str := fmt.Sprintf("SMS: send to %s.\nTransfer operation canceled. Your account has been transferred from account %d for the amount of %d", userOperation.Phone.String, tl.RecievedIDAccount.Int64, tl.Value)
			glog.Info(str)
		}
		str = fmt.Sprintf("EMAIL: send to %s.\nTransfer operation canceled. Your account has been transferred to account %d for the amount of %d", userRecieved.Email, tl.AccountForOperationID, tl.Value)
		glog.Info(str)
		if userRecieved.Phone.Valid != false {
			str := fmt.Sprintf("SMS: send to %s.\nTransfer operation canceled. Your account has been transferred to account %d for the amount of %d", userRecieved.Phone.String, tl.AccountForOperationID, tl.Value)
			glog.Info(str)
		}
	}
}
