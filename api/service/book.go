package service

import (
	"Library-Management-System/api/constant"
	"Library-Management-System/api/db"
	"Library-Management-System/api/dto/request"
	"Library-Management-System/api/dto/response"
	"Library-Management-System/api/idgen"
	"Library-Management-System/api/model"
	"Library-Management-System/api/util/xerror"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

var Book *BookService

func init() {
	Book = NewBookService()
}

type BookService struct {
	redisClient *redis.Client
}

func NewBookService() *BookService {
	return &BookService{
		redisClient: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379", // Redis 地址
			Password: "",               // Redis 密码
			DB:       0,                // Redis 数据库索引
		}),
	}
}

func (b *BookService) Add(option *request.AddBookOption) (*response.AddBookResponse, xerror.OpenError) {
	bookByName, err := db.Book.ListByName(option.Name, 1, 1)
	if err != nil {
		return nil, xerror.ErrAddBookFailed
	}
	var bookId string
	if len(bookByName) != 0 {
		if option.Author != bookByName[0].Author {
			return nil, xerror.ErrBookConflict.SetMessage("Book name already exists with a different author")
		}
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
	// 更新 Redis 缓存
	redisKey := "book:" + option.Id
	bookJSON, jsonErr := json.Marshal(book)
	if jsonErr == nil {
		b.redisClient.Set(context.Background(), redisKey, bookJSON, time.Hour*1) // 设置缓存过期时间为 1 小时
	} else {
		return nil, xerror.ErrUpdateBookFailed.SetMessage("Failed to update Redis cache")
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

func (b *BookService) DescribeWithCache(option *request.DescribeBookOption) (*response.DescribeBookResponse, xerror.OpenError) {
	// 定义 Redis 键
	redisKey := "book:" + option.Id

	// 从 Redis 获取数据
	cachedBook, err := b.redisClient.Get(context.Background(), redisKey).Result()
	if err == nil && cachedBook != "" {
		// 如果 Redis 中有数据，反序列化并返回
		var book model.Book
		if jsonErr := json.Unmarshal([]byte(cachedBook), &book); jsonErr == nil {
			return &response.DescribeBookResponse{Book: book}, nil
		}
	}

	// 如果 Redis 中没有数据，查询 MySQL
	book, err := db.Book.Get(option.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, xerror.ErrBookNotFound
		}
		return nil, xerror.ErrQueryBookFailed
	}

	// 将查询结果写入 Redis
	bookJSON, jsonErr := json.Marshal(book)
	if jsonErr == nil {
		b.redisClient.Set(context.Background(), redisKey, bookJSON, time.Hour*1) // 设置缓存过期时间为 1 小时
	}

	return &response.DescribeBookResponse{Book: *book}, nil
}
