package model

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}

type PublicUser struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

func (u *User) ToPublic() PublicUser {
	return PublicUser{
		ID:       u.ID,
		Name:     u.Name,
		Username: u.Username,
	}
}
