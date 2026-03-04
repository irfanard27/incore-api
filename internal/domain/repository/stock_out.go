package repository

import (
	"context"
	"irfanard27/incore-api/internal/domain/entity"
)

type StockOutRepository interface {
	Create(ctx context.Context, stockOut *entity.StockOut) error
	GetById(ctx context.Context, id string) (*entity.StockOut, error)
	GetByTransactionID(ctx context.Context, transactionID string) (*entity.StockOut, error)
	Update(ctx context.Context, stockOut *entity.StockOut) error
	Delete(ctx context.Context, id string) error
	All(ctx context.Context) ([]entity.StockOut, error)
}
