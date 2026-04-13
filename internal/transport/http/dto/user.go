package dto

import "github.com/TTekmii/todo-list-app/internal/domain/model"

type RegisterInput struct {
	Name     string `json:"name" binding:"required,min=2"`
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

func (r *RegisterInput) ToDomain() model.User {
	return model.User{
		Name:     r.Name,
		Username: r.Username,
	}
}

type UserResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

func UserFromDomain(u model.User) UserResponse {
	pub := u.ToPublic()
	return UserResponse{
		ID:       pub.ID,
		Name:     pub.Name,
		Username: pub.Username,
	}
}
