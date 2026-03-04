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
	// Not used currently but required by interface
	return fmt.Errorf("Create method not implemented")
}

func (r *stockOutItemRepository) BatchCreate(ctx context.Context, stockOutItems []entity.StockOutItem) error {
	if len(stockOutItems) == 0 {
		return fmt.Errorf("no stock out items provided")
	}

	query := `
		INSERT INTO stock_out_items (stock_out_id, inventory_id, quantity, created_at, updated_at)
		VALUES (:stock_out_id, :inventory_id, :quantity, NOW(), NOW())
	`
	_, err := r.db.NamedExecContext(ctx, query, stockOutItems)
	if err != nil {
		return fmt.Errorf("failed to batch create stock out items: %w", err)
	}
	return nil
}

func (r *stockOutItemRepository) GetById(ctx context.Context, id string) (*entity.StockOutItem, error) {
	// Not used currently but required by interface
	return nil, fmt.Errorf("GetById method not implemented")
}

func (r *stockOutItemRepository) GetByStockOutID(ctx context.Context, stockOutID string) ([]entity.StockOutItem, error) {
	query := `
		SELECT *
		FROM stock_out_items
		WHERE stock_out_id = :stock_out_id
		ORDER BY created_at DESC
	`
	params := map[string]any{
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

func (r *stockOutItemRepository) GetByStockOutIDWithInventory(ctx context.Context, stockOutID string) ([]entity.StockOutItem, []entity.Inventory, error) {
	query := `
		SELECT soi.*, i.id, i.sku, i.name
		FROM stock_out_items soi
		LEFT JOIN inventories i ON soi.inventory_id = i.id
		WHERE soi.stock_out_id = :stock_out_id
		ORDER BY soi.created_at DESC
	`
	params := map[string]any{
		"stock_out_id": stockOutID,
	}
	rows, err := r.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get stock out items by stock out id: %w", err)
	}
	defer rows.Close()

	var stockOutItems []entity.StockOutItem
	var inventories []entity.Inventory

	for rows.Next() {
		var stockOutItem entity.StockOutItem
		var inventoryID, inventorySku, inventoryName string

		err = rows.Scan(
			&stockOutItem.ID, &stockOutItem.StockOutID, &stockOutItem.InventoryID, &stockOutItem.Quantity, &stockOutItem.CreatedAt, &stockOutItem.UpdatedAt,
			&inventoryID, &inventorySku, &inventoryName,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan stock out item: %w", err)
		}

		stockOutItems = append(stockOutItems, stockOutItem)

		// Create inventory with only the fields we need
		inventory := entity.Inventory{
			ID:   inventoryID,
			Sku:  inventorySku,
			Name: inventoryName,
		}
		inventories = append(inventories, inventory)
	}

	return stockOutItems, inventories, nil
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
	// Not used currently but required by interface
	return nil, fmt.Errorf("All method not implemented")
}
