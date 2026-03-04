package repository

import (
	"context"
	"irfanard27/incore-api/internal/domain/entity"
)

type StockInItemRepository interface {
	Create(ctx context.Context, stockInItem *entity.StockInItem) error
	GetById(ctx context.Context, id string) (*entity.StockInItem, error)
	GetByStockInID(ctx context.Context, stockInID string) ([]entity.StockInItem, error)
	GetByStockInIDWithInventory(ctx context.Context, stockInID string) ([]entity.StockInItem, []entity.Inventory, error)
	Update(ctx context.Context, stockInItem *entity.StockInItem) error
	Delete(ctx context.Context, id string) error
	All(ctx context.Context) ([]entity.StockInItem, error)

	BatchCreate(ctx context.Context, stockInItems []entity.StockInItem) error
}
