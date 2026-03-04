package pgsql

import (
	"context"
	"fmt"
	"irfanard27/incore-api/internal/domain/entity"
	"irfanard27/incore-api/internal/domain/repository"

	"github.com/jmoiron/sqlx"
)

type stockInItemRepository struct {
	db *sqlx.DB
}

func NewStockInItemRepository(db *sqlx.DB) repository.StockInItemRepository {
	return &stockInItemRepository{db: db}
}

func (r *stockInItemRepository) Create(ctx context.Context, stockInItem *entity.StockInItem) error {
	query := `
		INSERT INTO stock_in_items (id, stock_in_id, inventory_id, quantity, created_at, updated_at)
		VALUES (:id, :stock_in_id, :inventory_id, :quantity, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, stockInItem)
	if err != nil {
		return fmt.Errorf("failed to create stock in item: %w", err)
	}
	return nil
}

func (r *stockInItemRepository) GetById(ctx context.Context, id string) (*entity.StockInItem, error) {
	query := `
		SELECT id, stock_in_id, inventory_id, quantity, created_at, updated_at
		FROM stock_in_items
		WHERE id = :id
	`
	params := map[string]interface{}{
		"id": id,
	}
	rows, err := r.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock in item by id: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("stock in item not found")
	}

	var stockInItem entity.StockInItem
	err = rows.StructScan(&stockInItem)
	if err != nil {
		return nil, fmt.Errorf("failed to scan stock in item: %w", err)
	}

	return &stockInItem, nil
}

func (r *stockInItemRepository) GetByStockInID(ctx context.Context, stockInID string) ([]entity.StockInItem, error) {
	query := `
		SELECT id, stock_in_id, inventory_id, quantity, created_at, updated_at
		FROM stock_in_items
		WHERE stock_in_id = :stock_in_id
		ORDER BY created_at DESC
	`
	params := map[string]interface{}{
		"stock_in_id": stockInID,
	}
	rows, err := r.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock in items by stock in id: %w", err)
	}
	defer rows.Close()

	var stockInItems []entity.StockInItem
	for rows.Next() {
		var stockInItem entity.StockInItem
		err = rows.StructScan(&stockInItem)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stock in item: %w", err)
		}
		stockInItems = append(stockInItems, stockInItem)
	}

	return stockInItems, nil
}

func (r *stockInItemRepository) Update(ctx context.Context, stockInItem *entity.StockInItem) error {
	query := `
		UPDATE stock_in_items
		SET stock_in_id = :stock_in_id, inventory_id = :inventory_id, quantity = :quantity, updated_at = :updated_at
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, stockInItem)
	if err != nil {
		return fmt.Errorf("failed to update stock in item: %w", err)
	}
	return nil
}

func (r *stockInItemRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM stock_in_items WHERE id = :id`
	params := map[string]interface{}{
		"id": id,
	}
	_, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to delete stock in item: %w", err)
	}
	return nil
}

func (r *stockInItemRepository) All(ctx context.Context) ([]entity.StockInItem, error) {
	query := `
		SELECT id, stock_in_id, inventory_id, quantity, created_at, updated_at
		FROM stock_in_items
		ORDER BY created_at DESC
	`
	var stockInItems []entity.StockInItem
	err := r.db.SelectContext(ctx, &stockInItems, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all stock in items: %w", err)
	}
	return stockInItems, nil
}

func (r *stockInItemRepository) BatchCreate(ctx context.Context, stockInItems []entity.StockInItem) error {
	if len(stockInItems) == 0 {
		return fmt.Errorf("no stock in items provided")
	}

	query := `
		INSERT INTO stock_in_items (stock_in_id, inventory_id, quantity, created_at, updated_at)
		VALUES (:stock_in_id, :inventory_id, :quantity, NOW(), NOW())
	`
	_, err := r.db.NamedExecContext(ctx, query, stockInItems)
	if err != nil {
		return fmt.Errorf("failed to batch create stock in items: %w", err)
	}
	return nil
}
