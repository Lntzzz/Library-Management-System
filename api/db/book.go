package db

import (
	"Library-Management-System/api/model"
	"database/sql"
	"time"
)

var Book *BookDao

type BookDao struct{}

func (b *BookDao) CountByName(bookName string) (total int64, err error) {
	query := "select count(id) as total from book where name like ?"
	row := Db.QueryRow(query, "%"+bookName+"%")
	err = row.Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (b *BookDao) ListByName(bookName string, pageNumber int64, pageSize int64) (books []model.Book, err error) {
	query := "select id, name, author, stock from book where name like ? limit ?, ?"
	offset := (pageNumber - 1) * pageSize
	rows, err := Db.Query(query, "%"+bookName+"%", offset, pageSize)
	if err != nil {
		return []model.Book{}, err
	}
	for rows.Next() {
		var book model.Book
		err = rows.Scan(&book.Id, &book.Name, &book.Author, &book.Stock)
		if err != nil {
			panic(err)
		}
		books = append(books, book)
	}
	return books, nil
}

func (b *BookDao) List(pageNumber int64, pageSize int64) (books []model.Book, err error) {
	query := "select id, name, author, stock from book limit ?, ?"
	offset := (pageNumber - 1) * pageSize
	rows, err := Db.Query(query, offset, pageSize)
	if err != nil {
		return []model.Book{}, err
	}
	for rows.Next() {
		var book model.Book
		err = rows.Scan(&book.Id, &book.Name, &book.Author, &book.Stock)
		if err != nil {
			panic(err)
		}
		books = append(books, book)
	}
	return books, nil
}

func (b *BookDao) Get(bookId string) (*model.Book, error) {
	query := "select id, name, author, stock from book where id = ?"
	rows, err := Db.Query(query, bookId)
	if err != nil {
		return &model.Book{}, err
	}
	defer rows.Close()

	var book model.Book
	if rows.Next() {
		err = rows.Scan(&book.Id, &book.Name, &book.Author, &book.Stock)
		if err != nil {
			return &model.Book{}, err
		}
	} else {
		return &model.Book{}, sql.ErrNoRows
	}
	time.Sleep(10 * time.Millisecond)
	return &book, nil
}

func (b *BookDao) Add(book model.Book) error {
	query := "insert into book (id, name, author, stock) values (?, ?, ?, ?)"
	_, err := Db.Exec(query, book.Id, book.Name, book.Author, book.Stock)
	if err != nil {
		return err
	}
	return nil
}

func (b *BookDao) Delete(bookId string) error {
	query := "delete from book where id = ?"
	_, err := Db.Exec(query, bookId)
	if err != nil {
		return err
	}
	return nil
}

func (b *BookDao) Update(bookId string, bookName string, bookAuthor string) error {
	query := "update book set name = ?, author = ? where id = ?"
	_, err := Db.Exec(query, bookName, bookAuthor, bookId)
	if err != nil {
		return err
	}
	return nil
}

func (b *BookDao) IncreaseBookCount(id string) error {
	query := "update book set stock = stock + 1 where id = ?"
	_, err := Db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
