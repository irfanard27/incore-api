package dto

import "irfanard27/incore-api/internal/domain/entity"

type StockOutDTO struct {
	ID            string            `json:"id"`
	TransactionID string            `json:"transaction_id"`
	Status        string            `json:"status"`
	TotalItem     int               `json:"total_item"`
	Items         []StockOutItemDTO `json:"items"`
	CreatedAt     string            `json:"created_at"`
	UpdatedAt     string            `json:"updated_at"`
}

type StockOutItemDTO struct {
	InventoryID string               `json:"inventory_id"`
	Quantity    int                  `json:"quantity"`
	Inventory   *MinimalInventoryDTO `json:"inventory"`
}

type StocksOutDTO struct {
	ID            string `json:"id"`
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	TotalItems    int    `json:"total_items"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type CreateStockOutDTO struct {
	Items []StockOutItemDTO `json:"items"`
}

type UpdateStockOutStatusDTO struct {
	Status string `json:"status" binding:"required"`
}

func (c CreateStockOutDTO) ToEntity() []entity.StockOutItem {
	items := []entity.StockOutItem{}
	for _, item := range c.Items {
		items = append(items, entity.StockOutItem{
			InventoryID: item.InventoryID,
			Quantity:    item.Quantity,
		})
	}
	return items
}

func ToStocksOutDTO(stock *entity.StockOut, totalItem int) StocksOutDTO {
	return StocksOutDTO{
		ID:            stock.ID,
		TransactionID: stock.TransactionID,
		Status:        stock.Status,
		TotalItems:    totalItem,
		CreatedAt:     stock.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     stock.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToStockOutDTOWithItems(stock *entity.StockOut, items []entity.StockOutItem, inventories []entity.Inventory) StockOutDTO {
	stockOutItems := make([]StockOutItemDTO, len(items))
	for i, item := range items {
		var inventory *MinimalInventoryDTO
		if i < len(inventories) {
			inventory = &MinimalInventoryDTO{
				ID:   inventories[i].ID,
				Sku:  inventories[i].Sku,
				Name: inventories[i].Name,
			}
		}

		stockOutItems[i] = StockOutItemDTO{
			InventoryID: item.InventoryID,
			Quantity:    item.Quantity,
			Inventory:   inventory,
		}
	}

	return StockOutDTO{
		ID:            stock.ID,
		TransactionID: stock.TransactionID,
		Status:        stock.Status,
		TotalItem:     len(items),
		Items:         stockOutItems,
		CreatedAt:     stock.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     stock.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToStocksOutDTOList(stocks []entity.StockOut, totalItems []int) []StocksOutDTO {
	res := make([]StocksOutDTO, len(stocks))
	for i, stock := range stocks {
		totalItem := 0
		if i < len(totalItems) {
			totalItem = totalItems[i]
		}
		res[i] = ToStocksOutDTO(&stock, totalItem)
	}
	return res
}
