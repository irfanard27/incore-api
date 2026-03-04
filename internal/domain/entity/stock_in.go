package entity

import "time"

type StockIn struct {
	ID            string    `db:"id"`
	TransactionID string    `db:"transaction_id"`
	Status        string    `db:"status"`
	TotalItem     int       `db:"total_item"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
