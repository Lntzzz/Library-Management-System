package db

import (
	"Library-Management-System/api/model"
	"testing"
)

func TestAddBook(t *testing.T) {
	book := model.Book{
		Id:     "book-01",
		Name:   "Test Book",
		Author: "Alice",
		Stock:  1,
	}
	err := Book.Add(book)
	if err != nil {
		t.Errorf("Failed to add book: %v", err)
	}
}

func TestAddBorrowRecord(t *testing.T) {
	record := model.BorrowRecord{
		Id:     "record-01",
		UserId: "user-01",
		BookId: "book-01",
		Status: 1,
	}
	err := BorrowRecord.Add(record)
	if err != nil {
		t.Errorf("Failed to add borrow record: %v", err)
	}
}
