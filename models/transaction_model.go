package models

import (
	"../entities"
	"database/sql"
)

type TransactionModel struct {
	Db *sql.DB
}

func (tranModel TransactionModel) DebitTransaction(transaction entities.Transaction) (entities.Transaction, error) {
	rows1, err1 := tranModel.Db.Query("select * from wallet where userid = ?", transaction.UserId)
	if err1 != nil {
		return entities.Transaction{}, err1
	} else {
		tran := entities.Transaction{}
		for rows1.Next() {
			var id int64
			var userid int64
			var transactionamount int64
			err2 := rows1.Scan(&id, &userid, &transactionamount)
			if err2 != nil {
				return entities.Transaction{}, err2
			} else {
				tran = entities.Transaction{userid, transactionamount}
				if transactionamount-transaction.TransactionAmout < 0 {
					return tran, nil
				} else {
					_, err3 := tranModel.Db.Exec(
						"update wallet set balance = balance - ? where userid = ?",
						transaction.TransactionAmout, transaction.UserId)
					if err3 != nil {
						return tran, err3
					} else {
						tran.TransactionAmout = tran.TransactionAmout - transaction.TransactionAmout
						return tran, nil
					}
				}
			}
		}
		return tran, nil

	}

}

func (tranModel TransactionModel) CreditTransaction(tran entities.Transaction) (int64, error) {
	result, err := tranModel.Db.Exec(
		"update wallet set balance = balance + ? where userid = ?",
		tran.TransactionAmout, tran.UserId)
	if err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}
