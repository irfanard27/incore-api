package entity

import "time"

type StockOutItem struct {
	ID          string    `db:"id"`
	StockOutID  string    `db:"stock_out_id"`
	InventoryID string    `db:"inventory_id"`
	Quantity    int       `db:"quantity"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type StockOutItemWithAvailability struct {
	ID                string `db:"id"`
	StockOutID        string `db:"stock_out_id"`
	InventoryID       string `db:"inventory_id"`
	Quantity          int    `db:"quantity"`
	AvailableQuantity int    `db:"available_quantity"`
	ReservedQuantity  int    `db:"reserved_quantity"`
}

type StockOutItemWithReserved struct {
	ID               string `db:"id"`
	StockOutID       string `db:"stock_out_id"`
	InventoryID      string `db:"inventory_id"`
	Quantity         int    `db:"quantity"`
	ReservedQuantity int    `db:"reserved_quantity"`
}
