package db

import (
	"Library-Management-System/api/model"
	"Library-Management-System/api/util/db"
)

var LoanRecord *LoanRecordDao

func InitLoanRecordDao(dbClient db.Client) {
	LoanRecord = &LoanRecordDao{BaseDao: db.NewBaseDao[model.BorrowRecord](dbClient)}
}

type LoanRecordDao struct {
	*db.BaseDao[model.BorrowRecord]
}
