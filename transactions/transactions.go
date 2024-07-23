package transactions

import (
	"strconv"
	"time"
)

var TransactionList []Transaction

type Transaction struct {
	Transaction_id    string    `json:"transaction_id"`
	Transaction_type  string    `json:"transaction_type"`
	Amount            int       `json:"amount"`
	Timestamp         time.Time `json:"timestamp"`
	Account_id        string    `json:"account_id"`
	SenderAccountId   string    `json:"sender_account_id"`
	ReceiverAccountId string    `json:"receiver_account_id"`
}

func GetNextId(trans []Transaction) string {
	i := len(trans)
	return strconv.Itoa(i + 1)
}
