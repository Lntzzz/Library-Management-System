package db

import (
	"Library-Management-System/api/model"
	"Library-Management-System/api/util/db"
)

var User *UserDao

func InitUserDao(dbClient db.Client) {
	User = &UserDao{BaseDao: db.NewBaseDao[model.User](dbClient)}
}

type UserDao struct {
	*db.BaseDao[model.User]
}
