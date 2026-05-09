package model

import "time"

type User struct {
	Id        string    `db:"id"`
	Name      string    `db:"name"`
	Password  string    `db:"password"`
	Type      int       `db:"type"` //区分管理员和用户
	CreatedAt time.Time `db:"created_at"`
}

func (user User) GetId() string {
	return user.Id
}

func (user User) GetName() string {
	return user.Name
}

func (user User) SetName(name string) User {
	user.Name = name
	return user
}

func (user User) GetType() int {
	return user.Type
}

func (user User) SetType(t int) User {
	user.Type = t
	return user
}

func (user User) GetPassword() string {
	return user.Password
}

func (user User) SetPassword(password string) User {
	user.Password = password
	return user
}

func (user User) GetCreatedAt() time.Time {
	return user.CreatedAt
}

func (user User) SetCreatedAt(createdAt time.Time) User {
	user.CreatedAt = createdAt
	return user
}
