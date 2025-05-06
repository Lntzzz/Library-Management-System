package service

import (
	"Library-Management-System/api/constant"
	"Library-Management-System/api/db"
	"Library-Management-System/api/dto/request"
	"Library-Management-System/api/dto/response"
	"Library-Management-System/api/idgen"
	"Library-Management-System/api/model"
	"Library-Management-System/api/util/xerror"
)

var Book *BookService

func init() {
	Book = NewBookService()
}

type BookService struct{}

func NewBookService() *BookService {
	return &BookService{}
}

func (b *BookService) Add(option *request.AddBookOption) (*response.AddBookResponse, xerror.OpenError) {
	bookByName, err := db.Book.ListByName(option.Name, 1, 1)
	if err != nil {
		return nil, xerror.ErrAddBookFailed
	}
	var bookId string
	if len(bookByName) != 0 {
		bookId = bookByName[0].Id
		stockErr := db.Book.IncreaseBookCount(bookByName[0].Id)
		if stockErr != nil {
			return nil, xerror.ErrAddBookFailed
		}
	} else {
		newBookId := idgen.NewBookId()
		bookId = newBookId
		book := model.Book{
			Id:     newBookId,
			Name:   option.Name,
			Author: option.Author,
			Stock:  constant.DefaultStock,
		}
		if insertErr := db.Book.Add(book); insertErr != nil {
			return nil, xerror.ErrAddBookFailed
		}
	}
	return &response.AddBookResponse{BookId: bookId}, nil
}

func (b *BookService) Delete(option *request.DeleteBookOption) (*response.DeleteBookResponse, xerror.OpenError) {
	book, err := db.Book.Get(option.Id)
	if book == nil {
		return nil, xerror.ErrBookNotFound
	}
	if err != nil {
		return nil, xerror.ErrQueryBookFailed
	}
	delErr := db.Book.Delete(option.Id)
	if delErr != nil {
		return nil, xerror.ErrDeleteBookFailed
	}
	return &response.DeleteBookResponse{BookId: book.Id}, nil
}

func (b *BookService) Update(option *request.UpdateBookOption) (*response.UpdateBookResponse, xerror.OpenError) {
	book, err := db.Book.Get(option.Id)
	if book == nil {
		return nil, xerror.ErrBookNotFound
	}
	if err != nil {
		return nil, xerror.ErrQueryBookFailed
	}
	err = db.Book.Update(option.Id, option.Name, option.Author)
	if err != nil {
		return nil, xerror.ErrUpdateBookFailed.SetMessage(err.Error())
	}
	return &response.UpdateBookResponse{BookId: book.Id}, nil
}

func (b *BookService) Describe(option *request.DescribeBookOption) (*response.DescribeBookResponse, xerror.OpenError) {
	book, err := db.Book.Get(option.Id)
	if book == nil {
		return nil, xerror.ErrBookNotFound
	}
	if err != nil {
		return nil, xerror.ErrQueryBookFailed
	}
	return &response.DescribeBookResponse{Book: *book}, nil
}

func (b *BookService) Describes(option *request.DescribeBooksOption) (*response.DescribesBooksResponse, xerror.OpenError) {
	books, err := db.Book.List(option.PageNumber, option.PageSize)
	if err != nil {
		return nil, xerror.ErrQueryBookFailed
	}
	return &response.DescribesBooksResponse{Books: books}, nil
}
