package repo

import "github.com/TTekmii/todo-list-app/internal/domain/model"

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUser(username, password string) (model.User, error)
}

// type Authorization interface {
//     CreateUser(ctx context.Context, user model.User) (int, error)
//     GetUserByUsername(ctx context.Context, username string) (model.User, error)
// }
