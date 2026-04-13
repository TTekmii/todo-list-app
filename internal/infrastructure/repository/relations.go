package repository

type userListRow struct {
	ID     int `db:"id"`
	UserID int `db:"user_id"`
	ListID int `db:"list_id"`
}

type listItemRow struct {
	ID     int `db:"id"`
	ListID int `db:"list_id"`
	ItemID int `db:"item_id"`
}
