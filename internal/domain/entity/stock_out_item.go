package entity

import "time"

type StockOutItem struct {
	ID           string    `db:"id"`
	StockOutID   string    `db:"stock_out_id"`
	InventoryID  string    `db:"inventory_id"`
	Quantity     int       `db:"quantity"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
