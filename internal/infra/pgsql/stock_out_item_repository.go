package pgsql

import (
	"context"
	"fmt"
	"irfanard27/incore-api/internal/domain/entity"
	"irfanard27/incore-api/internal/domain/repository"

	"github.com/jmoiron/sqlx"
)

type stockOutItemRepository struct {
	db *sqlx.DB
}

func NewStockOutItemRepository(db *sqlx.DB) repository.StockOutItemRepository {
	return &stockOutItemRepository{db: db}
}

func (r *stockOutItemRepository) Create(ctx context.Context, stockOutItem *entity.StockOutItem) error {
	query := `
		INSERT INTO stock_out_items (id, stock_out_id, inventory_id, quantity, created_at, updated_at)
		VALUES (:id, :stock_out_id, :inventory_id, :quantity, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, stockOutItem)
	if err != nil {
		return fmt.Errorf("failed to create stock out item: %w", err)
	}
	return nil
}

func (r *stockOutItemRepository) GetById(ctx context.Context, id string) (*entity.StockOutItem, error) {
	query := `
		SELECT id, stock_out_id, inventory_id, quantity, created_at, updated_at
		FROM stock_out_items
		WHERE id = :id
	`
	params := map[string]interface{}{
		"id": id,
	}
	rows, err := r.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock out item by id: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("stock out item not found")
	}

	var stockOutItem entity.StockOutItem
	err = rows.StructScan(&stockOutItem)
	if err != nil {
		return nil, fmt.Errorf("failed to scan stock out item: %w", err)
	}

	return &stockOutItem, nil
}

func (r *stockOutItemRepository) GetByStockOutID(ctx context.Context, stockOutID string) ([]entity.StockOutItem, error) {
	query := `
		SELECT id, stock_out_id, inventory_id, quantity, created_at, updated_at
		FROM stock_out_items
		WHERE stock_out_id = :stock_out_id
		ORDER BY created_at DESC
	`
	params := map[string]interface{}{
		"stock_out_id": stockOutID,
	}
	rows, err := r.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock out items by stock out id: %w", err)
	}
	defer rows.Close()

	var stockOutItems []entity.StockOutItem
	for rows.Next() {
		var stockOutItem entity.StockOutItem
		err = rows.StructScan(&stockOutItem)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stock out item: %w", err)
		}
		stockOutItems = append(stockOutItems, stockOutItem)
	}

	return stockOutItems, nil
}

func (r *stockOutItemRepository) Update(ctx context.Context, stockOutItem *entity.StockOutItem) error {
	query := `
		UPDATE stock_out_items
		SET stock_out_id = :stock_out_id, inventory_id = :inventory_id, quantity = :quantity, updated_at = :updated_at
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, stockOutItem)
	if err != nil {
		return fmt.Errorf("failed to update stock out item: %w", err)
	}
	return nil
}

func (r *stockOutItemRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM stock_out_items WHERE id = :id`
	params := map[string]interface{}{
		"id": id,
	}
	_, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to delete stock out item: %w", err)
	}
	return nil
}

func (r *stockOutItemRepository) All(ctx context.Context) ([]entity.StockOutItem, error) {
	query := `
		SELECT id, stock_out_id, inventory_id, quantity, created_at, updated_at
		FROM stock_out_items
		ORDER BY created_at DESC
	`
	var stockOutItems []entity.StockOutItem
	err := r.db.SelectContext(ctx, &stockOutItems, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all stock out items: %w", err)
	}
	return stockOutItems, nil
}
