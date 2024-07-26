package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"

	account "api/accounts"
	"api/transactions"
)

const (
	host     = "postgres"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "postgres"
)

func main() {
	//Setup Go-Chi router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	//Setup database connection string
	postgreSQL := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//Get connection pointer to database
	db, err_db := sql.Open("postgres", postgreSQL)
	if err_db != nil {
		panic(err_db)
	}
	defer db.Close() // Delay closing database connection

	//Open a connection to database
	err_open := db.Ping()
	if err_open != nil {
		panic(err_open)
	}
	fmt.Println("Open connection succes")

	//Mockup Data
	accList := account.AccountList
	acc1 := account.Account{}
	acc1.SetAccount("1234567890", "John Doe", "john@example.com", 500)
	acc2 := account.Account{}
	acc2.SetAccount("1234567891", "Jane Smith", "jane@example.com", 1000)
	accList = append(accList, acc1)
	accList = append(accList, acc2)
	transList := transactions.TransactionList
	/* transList = append(transList, ) */

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	//--------------------Account--------------------//
	r.Get("/accounts", func(w http.ResponseWriter, r *http.Request) {
		// db.Query()
		val, _ := json.Marshal(accList)
		w.Write([]byte(val))
	})

	r.Get("/accounts/{account_id}", func(w http.ResponseWriter, r *http.Request) {
		accountId := chi.URLParam(r, "account_id")
		var selectAcc account.Account

		for _, acc := range accList {
			if acc.Account_id == accountId {
				selectAcc = acc
			}
		}

		val, _ := json.Marshal(selectAcc)

		w.Write([]byte(val))
	})

	r.Post("/accounts", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err == nil {
			postAcc := account.Account{}
			postAcc.Account_id = account.GetNextId(accList)
			json.Unmarshal(body, &postAcc)

			accList = append(accList, postAcc)

			val, _ := json.Marshal(postAcc)

			w.Write([]byte(val))
		} else {
			w.Write([]byte("Failed"))
		}
	})

	r.Put("/accounts/{account_id}", func(w http.ResponseWriter, r *http.Request) {
		accountId := chi.URLParam(r, "account_id")
		body, err := io.ReadAll(r.Body)

		var acc1 account.Account
		var selectAcc account.Account
		var selectIndex int

		for i, acc := range accList {
			if acc.Account_id == accountId {
				selectAcc = acc
				selectIndex = i
			}
		}

		if err != nil {
			w.Write([]byte("Failed"))
		} else {
			json.Unmarshal(body, &acc1)
			selectAcc.Name = acc1.Name
			selectAcc.Email = acc1.Email
			selectAcc.Balance = acc1.Balance

			accList[selectIndex] = selectAcc
			w.Write([]byte("Status code 200"))
		}
	})

	r.Patch("/accounts/{account_id}", func(w http.ResponseWriter, r *http.Request) {
		accountId := chi.URLParam(r, "account_id")
		body, err := io.ReadAll(r.Body)

		var acc1 account.Account
		var selectAcc account.Account
		var selectIndex int

		for i, acc := range accList {
			if acc.Account_id == accountId {
				selectAcc = acc
				selectIndex = i
			}
		}

		if err != nil {
			w.Write([]byte("Failed"))
		} else {
			json.Unmarshal(body, &acc1)
			selectAcc.Name = acc1.Name

			accList[selectIndex] = selectAcc
			w.Write([]byte("Status code 200"))
		}
	})

	r.Delete("/accounts/{account_id}", func(w http.ResponseWriter, r *http.Request) {
		accountId := chi.URLParam(r, "account_id")
		var selectIndex int

		for i, acc := range accList {
			if acc.Account_id == accountId {
				selectIndex = i
			}
		}

		accListLength := len(accList)
		deleteList := account.AccountList
		slice1 := accList[:selectIndex]
		slice2 := accList[selectIndex+1 : accListLength]
		deleteList = append(deleteList, slice1...)
		deleteList = append(deleteList, slice2...)
		accList = deleteList
		w.Write([]byte("Status code 204"))
	})

	//--------------------Transaction--------------------//
	//Deposit
	r.Post("/transactions/deposit", func(w http.ResponseWriter, r *http.Request) {
		currentTime := time.Now()
		body, err := io.ReadAll(r.Body)

		if err != nil {
			w.Write([]byte("EOF Error"))
		}

		trans := transactions.Transaction{}
		body_err := json.Unmarshal(body, &trans)

		if body_err != nil {
			w.Write([]byte("Request Error : Invalid Body"))
		}

		hasId, accIndex := account.CheckId(accList, trans.Account_id)

		if !hasId {
			w.Write([]byte("Request Error : Account not found"))
		}

		trans.Transaction_id = transactions.GetNextId(transList)
		trans.Transaction_type = "deposit"
		trans.Timestamp = currentTime
		accList[accIndex].Balance = accList[accIndex].Balance + int64(trans.Amount)
		transList = append(transList, trans)
		res, _ := json.Marshal(trans)
		w.Write([]byte(res))

	})
	//Withdraw
	r.Post("/transactions/withdraw", func(w http.ResponseWriter, r *http.Request) {
		currentTime := time.Now()
		body, err := io.ReadAll(r.Body)

		if err != nil {
			w.Write([]byte("EOF Error"))
		}

		trans := transactions.Transaction{}
		body_err := json.Unmarshal(body, &trans)

		if body_err != nil {
			w.Write([]byte("Request Error : Invalid Body"))
		}

		hasId, accIndex := account.CheckId(accList, trans.Account_id)

		if !hasId {
			w.Write([]byte("Request Error : Account not found"))
		}

		if accList[accIndex].Balance >= int64(trans.Amount) {
			trans.Transaction_id = transactions.GetNextId(transList)
			trans.Transaction_type = "withdraw"
			trans.Timestamp = currentTime
			accList[accIndex].Balance = accList[accIndex].Balance - int64(trans.Amount)
			transList = append(transList, trans)
			res, _ := json.Marshal(trans)
			w.Write([]byte(res))
		}
	})
	//Transfer
	r.Post("/transactions/transfer", func(w http.ResponseWriter, r *http.Request) {
		currentTime := time.Now()
		body, err := io.ReadAll(r.Body)

		if err != nil {
			w.Write([]byte("EOF Error"))
		}

		trans := transactions.Transaction{}
		body_err := json.Unmarshal(body, &trans)

		if body_err != nil {
			w.Write([]byte("Request Error : Invalid Body"))
		}

		hasSenderId, accSenderIndex := account.CheckId(accList, trans.SenderAccountId)
		hasReceiverId, accReceiveIndex := account.CheckId(accList, trans.ReceiverAccountId)

		if !hasSenderId || !hasReceiverId {
			w.Write([]byte("Request Error : Account not found"))
		}

		if accList[accSenderIndex].Balance >= int64(trans.Amount) {
			trans.Transaction_id = transactions.GetNextId(transList)
			trans.Transaction_type = "transfer"
			trans.Timestamp = currentTime
			accList[accSenderIndex].Balance = accList[accSenderIndex].Balance - int64(trans.Amount)
			accList[accReceiveIndex].Balance = accList[accReceiveIndex].Balance + int64(trans.Amount)
			transList = append(transList, trans)
			res, _ := json.Marshal(trans)
			w.Write([]byte(res))
		}
	})

	fmt.Println("Server Started")
	http.ListenAndServe(":3000", r)
}
