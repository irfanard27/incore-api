package repository

import (
	"context"
	"irfanard27/incore-api/internal/domain/entity"
)

type InventoryRepository interface {
	All(ctx context.Context) ([]entity.Inventory, error)
	Create(ctx context.Context, inventory *entity.Inventory) error
	Update(ctx context.Context, inventory *entity.Inventory) error
	Delete(ctx context.Context, id string) error
	GetById(ctx context.Context, id string) (*entity.Inventory, error)
	Search(ctx context.Context, kw string) ([]entity.Inventory, error)
}
