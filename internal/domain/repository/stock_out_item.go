package repository

import (
	"context"
	"irfanard27/incore-api/internal/domain/entity"
)

type StockOutItemRepository interface {
	Create(ctx context.Context, stockOutItem *entity.StockOutItem) error
	GetById(ctx context.Context, id string) (*entity.StockOutItem, error)
	GetByStockOutID(ctx context.Context, stockOutID string) ([]entity.StockOutItem, error)
	GetByStockOutIDWithInventory(ctx context.Context, stockOutID string) ([]entity.StockOutItem, []entity.Inventory, error)
	Update(ctx context.Context, stockOutItem *entity.StockOutItem) error
	Delete(ctx context.Context, id string) error
	All(ctx context.Context) ([]entity.StockOutItem, error)

	BatchCreate(ctx context.Context, stockOutItems []entity.StockOutItem) error
}
