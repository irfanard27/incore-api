package repository

import (
	"context"
	"irfanard27/incore-api/internal/domain/entity"
)

type StockInRepository interface {
	Create(ctx context.Context, stockIn *entity.StockIn) (string, error)
	GetById(ctx context.Context, id string) (*entity.StockIn, error)
	GetByTransactionID(ctx context.Context, transactionID string) (*entity.StockIn, error)
	Update(ctx context.Context, stockIn *entity.StockIn) error
	Delete(ctx context.Context, id string) error
	All(ctx context.Context) ([]entity.StockIn, error)
}
