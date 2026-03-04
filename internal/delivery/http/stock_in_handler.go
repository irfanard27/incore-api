package http

import (
	"fmt"
	"irfanard27/incore-api/internal/domain/dto"
	"irfanard27/incore-api/internal/domain/entity"
	"irfanard27/incore-api/internal/usecase"
	httpresponse "irfanard27/incore-api/pkg/http_response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StockInHandler struct {
	stockInUsecase usecase.StockInUsecase
}

func NewStockInHandler(stockInUsecase usecase.StockInUsecase) *StockInHandler {
	return &StockInHandler{
		stockInUsecase: stockInUsecase,
	}
}

func (h *StockInHandler) GetAllStockIn(c *gin.Context) {
	stocks, err := h.stockInUsecase.GetAllStockIn(c.Request.Context())
	if err != nil {
		httpresponse.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Convert entities to DTOs
	res := dto.ToStocksInDTOList(stocks)

	httpresponse.SuccessResponse(c, "Stock ins retrieved successfully", res)
}

func (h *StockInHandler) GetStockInById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		httpresponse.ErrorResponse(c, fmt.Errorf("id is required"), http.StatusBadRequest)
		return
	}

	stock, items, inventories, err := h.stockInUsecase.GetStockInById(c.Request.Context(), id)
	if err != nil {
		httpresponse.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Convert entity to DTO with items
	res := dto.ToStockInDTOWithItems(stock, items, inventories)

	httpresponse.SuccessResponse(c, "Stock in retrieved successfully", res)
}

func (h *StockInHandler) CreateStockIn(c *gin.Context) {
	var req dto.CreateStockInDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	if len(req.Items) == 0 {
		httpresponse.ErrorResponse(c, fmt.Errorf("items must be provided"), http.StatusBadRequest)
		return
	}

	stockin, err := h.stockInUsecase.CreateStockIn(c.Request.Context(), req.ToEntity())

	if err != nil {
		httpresponse.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	httpresponse.SuccessResponse(c, "Stock in created successfully", stockin)
}

func (h *StockInHandler) UpdateStatus(c *gin.Context) {
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

	var stock *entity.StockIn
	var err error

	switch req.Status {
	case "in_progress":
		stock, err = h.stockInUsecase.UpdateStatusToInProgress(c.Request.Context(), id)
	case "done":
		stock, err = h.stockInUsecase.UpdateStatusToDone(c.Request.Context(), id)
	default:
		httpresponse.ErrorResponse(c, fmt.Errorf("invalid status '%s'. Valid statuses: 'in_progress', 'done'", req.Status), http.StatusBadRequest)
		return
	}

	if err != nil {
		httpresponse.ErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("Stock in status updated to '%s' successfully", req.Status)
	if req.Status == "done" {
		message = "Stock in status updated to 'done' successfully and inventory quantities updated"
	}

	httpresponse.SuccessResponse(c, message, stock)
}
