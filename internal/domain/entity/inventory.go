package entity

type Inventory struct {
	ID               string  `db:"id"`
	Sku              string  `db:"sku"`
	Name             string  `db:"name"`
	Price            int     `db:"price"`
	Quantity         int     `db:"quantity"`
	Customer         *string `db:"customer"`
	ReservedQuantity int     `db:"reserved_quantity"`
	CreatedAt        string  `db:"created_at"`
	UpdatedAt        string  `db:"updated_at"`
}
