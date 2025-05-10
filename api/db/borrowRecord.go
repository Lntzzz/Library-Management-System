package db

import (
	"Library-Management-System/api/model"
	"database/sql"
)

var BorrowRecord *BorrowRecordDao

type BorrowRecordDao struct{}

func (b *BorrowRecordDao) Add(record model.BorrowRecord) error {
	query := "INSERT INTO borrow_record (id, user_id, book_id, borrowed_at, status) VALUES (?, ?, ?, now(), ?)"
	_, err := Db.Exec(query, record.Id, record.UserId, record.BookId, record.Status)
	if err != nil {
		return err
	}
	return nil
}

func (b *BorrowRecordDao) Get(recordId string) (*model.BorrowRecord, error) {
	query := "SELECT id, user_id, book_id, borrowed_at, status FROM borrow_record WHERE id = ?"
	row := Db.QueryRow(query, recordId)

	var record model.BorrowRecord
	err := row.Scan(&record.Id, &record.UserId, &record.BookId, &record.BorrowedAt, &record.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &record, nil
}

func (b *BorrowRecordDao) Update(record model.BorrowRecord) error {
	query := "UPDATE borrow_record SET user_id = ?, book_id = ?, borrowed_at = ?, status = ? WHERE id = ?"
	_, err := Db.Exec(query, record.UserId, record.BookId, record.BorrowedAt, record.Status, record.Id)
	if err != nil {
		return err
	}
	return nil
}

func (b *BorrowRecordDao) Delete(recordId string) error {
	query := "DELETE FROM borrow_record WHERE id = ?"
	_, err := Db.Exec(query, recordId)
	if err != nil {
		return err
	}
	return nil
}

func (b *BorrowRecordDao) GetByBookIdAndUserId(bookId, userId string) (*model.BorrowRecord, error) {
	query := "SELECT id, user_id, book_id, borrowed_at, status FROM borrow_record WHERE book_id = ? AND user_id = ?"
	row := Db.QueryRow(query, bookId, userId)

	var record model.BorrowRecord
	err := row.Scan(&record.Id, &record.UserId, &record.BookId, &record.BorrowedAt, &record.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &record, nil
}
