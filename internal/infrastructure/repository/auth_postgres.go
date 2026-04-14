package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/TTekmii/todo-list-app/internal/domain/model"
	"github.com/TTekmii/todo-list-app/internal/domain/repo"
	"github.com/jmoiron/sqlx"
)

var _ repo.Authorization = (*AuthPostgres)(nil)

type dbUser struct {
	ID           int    `db:"id"`
	Name         string `db:"name"`
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
}

func toDomainUser(du dbUser) model.User {
	return model.User{
		ID:           du.ID,
		Name:         du.Name,
		Username:     du.Username,
		PasswordHash: du.PasswordHash,
	}
}

func fromDomainUser(u model.User) dbUser {
	return dbUser{
		Name:         u.Name,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
	}
}

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(ctx context.Context, user model.User) (int, error) {
	dbUser := fromDomainUser(user)

	var id int

	query := fmt.Sprintf(
		"INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id",
		usersTable,
	)

	row := r.db.QueryRowContext(ctx, query, dbUser.Name, dbUser.Username, dbUser.PasswordHash)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to scan user id: %w", err)
	}

	return id, nil
}

func (r *AuthPostgres) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	var du dbUser

	query := fmt.Sprintf(
		"SELECT id, name, username, password_hash FROM %s WHERE username=$1",
		usersTable,
	)

	if err := r.db.GetContext(ctx, &du, query, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, fmt.Errorf("user not found: %w", err)
		}
		return model.User{}, fmt.Errorf("failed to fetch user: %w", err)
	}

	return toDomainUser(du), nil
}
