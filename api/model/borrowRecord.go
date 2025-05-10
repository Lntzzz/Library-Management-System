package model

import "time"

type BorrowRecord struct {
	Id         string    `db:"id" json:"id"`
	UserId     string    `db:"user_id" json:"userId"`
	BookId     string    `db:"book_id" json:"bookId"`
	BorrowedAt time.Time `db:"borrowed_at" json:"borrowedAt"`
	ReturnedAt time.Time `db:"returned_at" json:"returnedAt"`
	Status     int       `db:"status" json:"status"`
	// 可选：续借次数
	// RenewCount int        `db:"renew_count" json:"renewCount"`
}

func (record BorrowRecord) GetId() string {
	return record.Id
}

func (record BorrowRecord) GetUserId() string {
	return record.UserId
}

func (record BorrowRecord) SetUserId(userId string) BorrowRecord {
	record.UserId = userId
	return record
}

func (record BorrowRecord) GetBookId() string {
	return record.BookId
}

func (record BorrowRecord) SetBookId(bookId string) BorrowRecord {
	record.BookId = bookId
	return record
}

func (record BorrowRecord) GetBorrowedAt() time.Time {
	return record.BorrowedAt
}

func (record BorrowRecord) SetBorrowedAt(borrowedAt time.Time) BorrowRecord {
	record.BorrowedAt = borrowedAt
	return record
}

func (record BorrowRecord) GetReturnedAt() time.Time {
	return record.ReturnedAt
}

func (record BorrowRecord) SetReturnedAt(returnedAt time.Time) BorrowRecord {
	record.ReturnedAt = returnedAt
	return record
}

func (record BorrowRecord) GetStatus() int {
	return record.Status
}

func (record BorrowRecord) SetStatus(status int) BorrowRecord {
	record.Status = status
	return record
}
