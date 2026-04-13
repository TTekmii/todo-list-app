package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/TTekmii/todo-list-app/internal/domain/model"
	"github.com/TTekmii/todo-list-app/internal/domain/repo"
	"github.com/jmoiron/sqlx"
)

var _ repo.TodoItem = (*TodoItemPostgres)(nil)

type dbItem struct {
	ID          int    `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Done        bool   `db:"done"`
}

func toDomainItem(di dbItem) model.TodoItem {
	return model.TodoItem{
		ID:          di.ID,
		Title:       di.Title,
		Description: di.Description,
		Done:        di.Done,
	}
}

func fromDomainItem(i model.TodoItem) dbItem {
	return dbItem{
		ID:          i.ID,
		Title:       i.Title,
		Description: i.Description,
		Done:        i.Done,
	}
}

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Create(ctx context.Context, listId int, item model.TodoItem) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	dbItem := fromDomainItem(item)
	var itemId int

	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoItemsTable)

	row := tx.QueryRowContext(ctx, createItemQuery, dbItem.Title, dbItem.Description)
	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to create item: %w", err)
	}

	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2)", listsItemsTable)
	_, err = tx.ExecContext(ctx, createListItemsQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to link item to list: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return itemId, nil
}

func (r *TodoItemPostgres) GetAll(ctx context.Context, userId, listId int) ([]model.TodoItem, error) {
	var dbItems []dbItem
	query := fmt.Sprintf(`
			SELECT ti.id, ti.title, ti.description, ti.done 
			FROM %s ti
			INNER JOIN %s li ON ti.id = li.item_id
			INNER JOIN %s ul ON li.list_id = ul.list_id
			WHERE li.list_id = $1
			AND ul.user_id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)
	if err := r.db.SelectContext(ctx, &dbItems, query, listId, userId); err != nil {
		return nil, err
	}

	items := make([]model.TodoItem, len(dbItems))
	for i, di := range dbItems {
		items[i] = toDomainItem(di)
	}
	return items, nil
}

func (r *TodoItemPostgres) GetById(ctx context.Context, userId, itemId int) (model.TodoItem, error) {
	var dbItem dbItem
	query := fmt.Sprintf(`
			SELECT ti.id, ti.title, ti.description, ti.done 
			FROM %s ti
			INNER JOIN %s li ON ti.id = li.item_id
			INNER JOIN %s ul ON li.list_id = ul.list_id
			WHERE ti.id = $1
			AND ul.user_id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)
	if err := r.db.GetContext(ctx, &dbItem, query, itemId, userId); err != nil {
		return model.TodoItem{}, err
	}

	return toDomainItem(dbItem), nil
}

func (r *TodoItemPostgres) Delete(ctx context.Context, userId, itemId int) error {
	query := fmt.Sprintf(`
			DELETE FROM %s ti USING %s li, %s ul
			WHERE ti.id = li.item_id
			AND li.list_id = ul.list_id
			AND ul.user_id = $1
			AND ti.id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)
	_, err := r.db.ExecContext(ctx, query, userId, itemId)
	return err
}

func (r *TodoItemPostgres) Update(ctx context.Context, userId, itemId int, input model.UpdateItemInput) error {
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

	if input.Done != nil {
		setValue = append(setValue, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	setQuery := strings.Join(setValue, ", ")

	query := fmt.Sprintf(`
			UPDATE %s ti SET %s
			FROM %s li, %s ul
			WHERE ti.id = li.item_id
			AND li.list_id = ul.list_id
			AND ul.user_id = $%d 
			AND ti.id = $%d`,
		todoItemsTable, setQuery, listsItemsTable, usersListsTable, argId, argId+1)
	args = append(args, userId, itemId)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}
