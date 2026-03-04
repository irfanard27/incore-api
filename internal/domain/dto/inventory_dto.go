package dto

import (
	"irfanard27/incore-api/internal/domain/entity"
)

// ToEntity converts CreateInventoryRequest to Inventory entity
func (req *CreateInventoryRequest) ToEntity() *entity.Inventory {
	return &entity.Inventory{
		Sku:      req.Sku,
		Name:     req.Name,
		Price:    req.Price,
		Quantity: req.Quantity,
		Customer: &req.Customer,
	}
}

// ToEntity converts UpdateInventoryRequest to Inventory entity
func (req *UpdateInventoryRequest) ToEntity() *entity.Inventory {
	return &entity.Inventory{
		ID:       req.ID,
		Sku:      req.Sku,
		Name:     req.Name,
		Price:    req.Price,
		Quantity: req.Quantity,
		Customer: &req.Customer,
	}
}

type InventoryDTO struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Sku              string `json:"sku"`
	Quantity         int    `json:"quantity"`
	ReservedQuantity int    `json:"reserved_quantity"`
	Customer         string `json:"customer"`
	Price            int    `json:"price"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

type CreateInventoryRequest struct {
	Name     string `json:"name" binding:"required"`
	Sku      string `json:"sku" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"`
	Customer string `json:"customer"`
	Price    int    `json:"price" binding:"required"`
}

type UpdateInventoryRequest struct {
	ID       string `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Sku      string `json:"sku" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"`
	Customer string `json:"customer"`
	Price    int    `json:"price" binding:"required"`
}

type DeleteInventoryRequest struct {
	ID string `json:"id" binding:"required"`
}

type InventoriesResponse struct {
	Data    []InventoryDTO `json:"data"`
	Message string         `json:"message"`
}

type InventoryResponse struct {
	Data    InventoryDTO `json:"data"`
	Message string       `json:"message"`
}

//
