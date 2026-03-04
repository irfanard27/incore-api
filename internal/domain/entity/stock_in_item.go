package entity

import "time"

type StockInItem struct {
	ID          string    `db:"id"`
	StockInID   string    `db:"stock_in_id"`
	InventoryID string    `db:"inventory_id"`
	Quantity    int       `db:"quantity"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
