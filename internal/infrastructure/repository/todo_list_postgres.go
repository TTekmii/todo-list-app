package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/TTekmii/todo-list-app/internal/domain/model"
	"github.com/TTekmii/todo-list-app/internal/domain/repo"
	"github.com/jmoiron/sqlx"
)

var _ repo.TodoList = (*TodoListPostgres)(nil)

type dbList struct {
	ID          int    `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
}

func toDomainList(dl dbList) model.TodoList {
	return model.TodoList{
		ID:          dl.ID,
		Title:       dl.Title,
		Description: dl.Description,
	}
}

func fromDomainList(l model.TodoList) dbList {
	return dbList{
		ID:          l.ID,
		Title:       l.Title,
		Description: l.Description,
	}
}

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) Create(ctx context.Context, userId int, list model.TodoList) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	dbList := fromDomainList(list)
	var listID int

	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	row := tx.QueryRowContext(ctx, createListQuery, dbList.Title, dbList.Description)
	if err := row.Scan(&listID); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to create list: %w", err)
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	_, err = tx.ExecContext(ctx, createUsersListQuery, userId, listID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to link user to list: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return listID, nil
}

func (r *TodoListPostgres) GetAll(ctx context.Context, userId int) ([]model.TodoList, error) {
	var dbLists []dbList

	query := fmt.Sprintf(`
			SELECT tl.id, tl.title, tl.description 
			FROM %s tl 
			INNER JOIN %s ul ON tl.id = ul.list_id 
			WHERE ul.user_id = $1`,
		todoListsTable, usersListsTable,
	)

	if err := r.db.SelectContext(ctx, &dbLists, query, userId); err != nil {
		return nil, fmt.Errorf("failed to fetch lists: %w", err)
	}

	lists := make([]model.TodoList, len(dbLists))
	for i, dl := range dbLists {
		lists[i] = toDomainList(dl)
	}

	return lists, nil
}

func (r *TodoListPostgres) GetById(ctx context.Context, userId, listId int) (model.TodoList, error) {
	var dbList dbList

	query := fmt.Sprintf(`
			SELECT tl.id, tl.title, tl.description 
			FROM %s tl
			INNER JOIN %s ul ON tl.id = ul.list_id 
			WHERE ul.user_id = $1 AND ul.list_id = $2`,
		todoListsTable, usersListsTable,
	)

	if err := r.db.GetContext(ctx, &dbList, query, userId, listId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.TodoList{}, fmt.Errorf("list not found: %w", err)
		}
		return model.TodoList{}, fmt.Errorf("failed to fetch list: %w", err)
	}

	return toDomainList(dbList), nil
}

func (r *TodoListPostgres) Delete(ctx context.Context, userId, listId int) error {
	query := fmt.Sprintf(`
			DELETE FROM %s tl USING %s ul 
			WHERE tl.id = ul.list_id 
			AND ul.user_id = $1 AND ul.list_id = $2`,
		todoListsTable, usersListsTable,
	)

	result, err := r.db.ExecContext(ctx, query, userId, listId)

	if err != nil {
		return fmt.Errorf("failed to delete list: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("list not found or access denied")
	}

	return err
}

func (r *TodoListPostgres) Update(ctx context.Context, userId, listId int, input model.UpdateListInput) error {
	if !input.HasChanges() {
		return nil
	}

	setValue := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValue = append(setValue, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValue = append(setValue, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValue, ", ")

	query := fmt.Sprintf(`
			UPDATE %s tl SET %s
			FROM %s ul
			WHERE tl.id = ul.list_id 
			AND ul.list_id = $%d 
			AND ul.user_id = $%d`,
		todoListsTable, setQuery, usersListsTable, argId, argId+1,
	)
	args = append(args, listId, userId)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update list: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("list not found or access denied")
	}

	return err
}
