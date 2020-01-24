package main

import (
	"./config"
	"./entities"
	"./models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Transaction struct {
	UserID            int64 `json:"userid"`
	TransactionAmount int64 `json:"transactionamount"`
}

func creditRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var t Transaction
	_ = json.NewDecoder(r.Body).Decode(&t)
	doCredit(t.UserID, t.TransactionAmount)
	json.NewEncoder(w).Encode(t)

	return
}

func debitRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var t Transaction
	_ = json.NewDecoder(r.Body).Decode(&t)
	doDebit(t.UserID, t.TransactionAmount)

	json.NewEncoder(w).Encode(t)
	return
}

func main() {
	// Init Router
	r := mux.NewRouter()

	// Route Handlers / Endpoints
	r.HandleFunc("/api/credit", creditRequest).Methods("PUT")
	r.HandleFunc("/api/debit", debitRequest).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8082", r))
}

func doCredit(uid int64, tranamount int64) {
	db, err := config.GetMySQLDB()
	if err != nil {
		fmt.Println(err)
	} else {
		tranModel := models.TransactionModel{Db: db}
		tran := entities.Transaction{
			UserId:           uid,
			TransactionAmout: tranamount,
		}
		rows, err := tranModel.CreditTransaction(tran)
		if err != nil {
			fmt.Println("my_error: ", err)
		} else {
			if rows > 0 {
				fmt.Println("Done")
			}
		}
	}

}

func doDebit(uid int64, tranamount int64) entities.Transaction {
	db, err := config.GetMySQLDB()
	if err != nil {
		fmt.Println(err)
	} else {
		tranModel := models.TransactionModel{Db: db}
		tran := entities.Transaction{
			UserId:           uid,
			TransactionAmout: tranamount,
		}
		result, err := tranModel.DebitTransaction(tran)
		if err != nil {
			fmt.Println("my_error: ", err)
		} else {
			fmt.Println(result)
			return result
		}
	}

	result := entities.Transaction{uid, tranamount}
	return result
}
