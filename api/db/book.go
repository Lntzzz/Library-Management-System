package db

import (
	"Library-Management-System/api/model"
	"Library-Management-System/api/util/db"
)

var Book *BookDao

func InitBookDao(dbClient db.Client) {
	Book = &BookDao{BaseDao: db.NewBaseDao[model.Book](dbClient)}
}

type BookDao struct {
	*db.BaseDao[model.Book]
}
