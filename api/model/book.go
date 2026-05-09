package model

type Book struct {
	Id     string `db:"id"`
	Name   string `db:"name"`
	Author string `db:"author"`
	Stock  int    `db:"stock"`
}

func (book Book) GetId() string {
	return book.Id
}

func (book Book) GetName() string {
	return book.Name
}

func (book Book) SetName(name string) Book {
	book.Name = name
	return book
}

func (book Book) GetAuthor() string {
	return book.Author
}

func (book Book) SetAuthor(author string) Book {
	book.Author = author
	return book
}

func (book Book) GetStock() int {
	return book.Stock
}

func (book Book) SetStock(stock int) Book {
	book.Stock = stock
	return book
}
