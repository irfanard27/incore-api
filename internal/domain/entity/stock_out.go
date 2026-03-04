package entity

import "time"

type StockOut struct {
	ID            string    `db:"id"`
	TransactionID string    `db:"transaction_id"`
	Status        string    `db:"status"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
