package http

import (
	"fmt"
	"irfanard27/incore-api/internal/domain/dto"
	"irfanard27/incore-api/internal/usecase"
	httpresponse "irfanard27/incore-api/pkg/http_response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	inventoryUsecase usecase.InventoryUsecase
}

func NewInventoryHandler(inventoryUsecase usecase.InventoryUsecase) *InventoryHandler {
	return &InventoryHandler{
		inventoryUsecase: inventoryUsecase,
	}
}

// CreateInventory handles inventory creation
func (h *InventoryHandler) CreateInventory(c *gin.Context) {
	var req dto.CreateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	_, err := h.inventoryUsecase.CreateInventory(c.Request.Context(), req.Sku, req.Name, req.Price, req.Quantity)
	if err != nil {
		httpresponse.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	httpresponse.SuccessResponse(c, "Inventory created successfully", nil)
}

// GetInventoryById handles getting inventory by ID
func (h *InventoryHandler) GetInventoryById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		httpresponse.ErrorResponse(c, fmt.Errorf("id parameter is required"), http.StatusBadRequest)
		return
	}

	inventory, err := h.inventoryUsecase.GetInventoryById(c.Request.Context(), id)
	if err != nil {
		httpresponse.ErrorResponse(c, err, http.StatusNotFound)
		return
	}

	httpresponse.SuccessResponse(c, "Inventory retrieved successfully", inventory)
}

// GetAllInventories handles getting all inventories
func (h *InventoryHandler) GetAllInventories(c *gin.Context) {
	inventories, err := h.inventoryUsecase.GetAllInventories(c.Request.Context())
	if err != nil {
		httpresponse.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	var inventoryDTOs []dto.InventoryDTO
	for _, inventory := range inventories {
		inventoryDTOs = append(inventoryDTOs, dto.InventoryDTO{
			ID:               inventory.ID,
			Name:             inventory.Name,
			Sku:              inventory.Sku,
			Quantity:         inventory.Quantity,
			ReservedQuantity: inventory.ReservedQuantity,
			Price:            inventory.Price,
			CreatedAt:        inventory.CreatedAt,
			UpdatedAt:        inventory.UpdatedAt,
		})
	}

	httpresponse.SuccessResponse(c, "Inventories retrieved successfully", inventoryDTOs)
}

func (h *InventoryHandler) SearchInventory(c *gin.Context) {

	keyword := c.Query("keyword")

	inventories, err := h.inventoryUsecase.SearchInventory(c.Request.Context(), keyword)
	if err != nil {
		httpresponse.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	var inventoryDTOs []dto.InventoryDTO
	for _, inventory := range inventories {
		inventoryDTOs = append(inventoryDTOs, dto.InventoryDTO{
			ID:               inventory.ID,
			Name:             inventory.Name,
			Sku:              inventory.Sku,
			Quantity:         inventory.Quantity,
			ReservedQuantity: inventory.ReservedQuantity,
			Price:            inventory.Price,
			CreatedAt:        inventory.CreatedAt,
			UpdatedAt:        inventory.UpdatedAt,
		})
	}

	httpresponse.SuccessResponse(c, "Inventories retrieved successfully", inventoryDTOs)
}

// UpdateInventory handles inventory update
func (h *InventoryHandler) UpdateInventory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		httpresponse.ErrorResponse(c, fmt.Errorf("id parameter is required"), http.StatusBadRequest)
		return
	}

	var req dto.UpdateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	// Use ID from URL parameter, not from request body
	_, err := h.inventoryUsecase.UpdateInventory(c.Request.Context(), req.ToEntity())
	if err != nil {
		httpresponse.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	httpresponse.SuccessResponse(c, "Inventory updated successfully", nil)
}

// DeleteInventory handles inventory deletion
func (h *InventoryHandler) DeleteInventory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		httpresponse.ErrorResponse(c, fmt.Errorf("id parameter is required"), http.StatusBadRequest)
		return
	}

	err := h.inventoryUsecase.DeleteInventory(c.Request.Context(), id)
	if err != nil {
		httpresponse.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	httpresponse.SuccessResponse(c, "Inventory deleted successfully", nil)
}
