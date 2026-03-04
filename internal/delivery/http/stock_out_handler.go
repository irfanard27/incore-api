package http

import (
	"context"
	"fmt"
	"irfanard27/incore-api/internal/domain/dto"
	"irfanard27/incore-api/internal/domain/entity"
	"irfanard27/incore-api/internal/usecase"
	httpresponse "irfanard27/incore-api/pkg/http_response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StockOutHandler struct {
	stockOutUsecase usecase.StockOutUsecase
}

func NewStockOutHandler(stockOutUsecase usecase.StockOutUsecase) *StockOutHandler {
	return &StockOutHandler{
		stockOutUsecase: stockOutUsecase,
	}
}

func (h *StockOutHandler) GetAllStockOut(c *gin.Context) {
	stocks, totalItems, err := h.stockOutUsecase.GetAllStockOut(c.Request.Context())
	if err != nil {
		httpresponse.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Convert entities to DTOs with total items
	res := dto.ToStocksOutDTOList(stocks, totalItems)

	httpresponse.SuccessResponse(c, "Stock outs retrieved successfully", res)
}

func (h *StockOutHandler) GetStockOutById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		httpresponse.ErrorResponse(c, fmt.Errorf("id is required"), http.StatusBadRequest)
		return
	}

	stock, items, inventories, err := h.stockOutUsecase.GetStockOutById(c.Request.Context(), id)
	if err != nil {
		httpresponse.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Convert entity to DTO with items
	res := dto.ToStockOutDTOWithItems(stock, items, inventories)

	httpresponse.SuccessResponse(c, "Stock out retrieved successfully", res)
}

func (h *StockOutHandler) CreateStockOut(c *gin.Context) {
	var req dto.CreateStockOutDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	if len(req.Items) == 0 {
		httpresponse.ErrorResponse(c, fmt.Errorf("items must be provided"), http.StatusBadRequest)
		return
	}

	stockout, err := h.stockOutUsecase.CreateStockOut(c.Request.Context(), req.ToEntity())

	if err != nil {
		httpresponse.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	httpresponse.SuccessResponse(c, "Stock out created successfully", stockout)
}

func (h *StockOutHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		httpresponse.ErrorResponse(c, fmt.Errorf("id is required"), http.StatusBadRequest)
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	statusHandlers := map[string]func(context.Context, string) (*entity.StockOut, error){
		"allocated":   h.stockOutUsecase.AllocateStock,
		"in_progress": h.stockOutUsecase.ExecuteStockOut,
		"done":        h.stockOutUsecase.CompleteStockOut,
		"cancelled":   h.stockOutUsecase.CancelStockOut,
	}

	var stock *entity.StockOut
	var err error

	if handler, exists := statusHandlers[req.Status]; exists {
		stock, err = handler(c.Request.Context(), id)
	} else {
		httpresponse.ErrorResponse(c, fmt.Errorf("invalid status '%s'. Valid statuses: 'allocated', 'in_progress', 'done', 'cancelled'", req.Status), http.StatusBadRequest)
		return
	}

	if err != nil {
		httpresponse.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("Stock out status updated to '%s' successfully", req.Status)
	switch req.Status {
	case "cancelled":
		message = "Stock out cancelled and stock rolled back successfully"
	case "done":
		message = "Stock out completed successfully"
	}

	httpresponse.SuccessResponse(c, message, stock)
}

func (h *StockOutHandler) DeleteStockOut(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		httpresponse.ErrorResponse(c, fmt.Errorf("id is required"), http.StatusBadRequest)
		return
	}

	err := h.stockOutUsecase.DeleteStockOut(c.Request.Context(), id)
	if err != nil {
		httpresponse.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	httpresponse.SuccessResponse(c, "Stock out deleted successfully", nil)
}
