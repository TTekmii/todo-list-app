package repo

import "github.com/TTekmii/todo-list-app/internal/domain/model"

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUser(username, password string) (model.User, error)
}
