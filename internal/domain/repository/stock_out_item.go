package repository

import (
	"context"
	"irfanard27/incore-api/internal/domain/entity"
)

type StockOutItemRepository interface {
	Create(ctx context.Context, stockOutItem *entity.StockOutItem) error
	GetById(ctx context.Context, id string) (*entity.StockOutItem, error)
	GetByStockOutID(ctx context.Context, stockOutID string) ([]entity.StockOutItem, error)
	Update(ctx context.Context, stockOutItem *entity.StockOutItem) error
	Delete(ctx context.Context, id string) error
	All(ctx context.Context) ([]entity.StockOutItem, error)
}
