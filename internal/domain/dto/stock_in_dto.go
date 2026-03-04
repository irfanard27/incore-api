package dto

import "irfanard27/incore-api/internal/domain/entity"

type StockInDTO struct {
	TransactionID string           `json:"transaction_id"`
	Status        string           `json:"status"`
	Items         []StockInItemDTO `json:"items"`
}

type StockInItemDTO struct {
	InventoryID string        `json:"inventory_id"`
	Quantity    int           `json:"quantity"`
	Inventory   *InventoryDTO `json:"inventory"`
}

type CreateStockInDTO struct {
	Items []StockInItemDTO `json:"items"`
}

func (c CreateStockInDTO) ToEntity() []entity.StockInItem {

	items := []entity.StockInItem{}
	for _, item := range c.Items {
		items = append(items, entity.StockInItem{
			InventoryID: item.InventoryID,
			Quantity:    item.Quantity,
		})
	}
	return items
}
