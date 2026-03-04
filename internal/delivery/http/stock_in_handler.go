package http

import (
	"fmt"
	"irfanard27/incore-api/internal/domain/dto"
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
