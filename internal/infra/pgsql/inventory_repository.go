package pgsql

import (
	"context"
	"fmt"
	"irfanard27/incore-api/internal/domain/entity"
	"irfanard27/incore-api/internal/domain/repository"

	"github.com/jmoiron/sqlx"
)

type inventoryRepository struct {
	db *sqlx.DB
}

// NewInventoryRepository creates a new inventory repository with database connection.
func NewInventoryRepository(db *sqlx.DB) repository.InventoryRepository {
	return &inventoryRepository{db: db}
}

// All retrieves all inventory records ordered by creation date (newest first).
func (r *inventoryRepository) All(ctx context.Context) ([]entity.Inventory, error) {
	query := `
		SELECT *
		FROM inventories
		ORDER BY created_at DESC
	`
	var inventories []entity.Inventory
	err := r.db.SelectContext(ctx, &inventories, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all inventories: %w", err)
	}
	return inventories, nil
}

// Create inserts a new inventory record into the database.
func (r *inventoryRepository) Create(ctx context.Context, inventory *entity.Inventory) error {
	query := `
		INSERT INTO inventories (sku, name, price, quantity, created_at, updated_at)
		VALUES (:sku, :name, :price, :quantity, NOW(), NOW())
	`
	_, err := r.db.NamedExecContext(ctx, query, inventory)
	if err != nil {
		return fmt.Errorf("failed to create inventory: %w", err)
	}
	return nil
}

// Update modifies an existing inventory record in the database.
func (r *inventoryRepository) Update(ctx context.Context, inventory *entity.Inventory) error {
	query := `
		UPDATE inventories
		SET sku = :sku, name = :name, price = :price, quantity = :quantity, updated_at = NOW()
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, inventory)
	if err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}
	return nil
}

// Delete removes an inventory record from the database by its ID.
func (r *inventoryRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM inventories WHERE id = :id`
	params := map[string]any{
		"id": id,
	}
	_, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to delete inventory: %w", err)
	}
	return nil
}

// GetById retrieves a single inventory record from the database by its ID.
func (r *inventoryRepository) GetById(ctx context.Context, id string) (*entity.Inventory, error) {
	query := `
		SELECT *
		FROM inventories
		WHERE id = :id
	`
	params := map[string]any{
		"id": id,
	}
	rows, err := r.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory by id: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("inventory not found")
	}

	var inventory entity.Inventory
	err = rows.StructScan(&inventory)
	if err != nil {
		return nil, fmt.Errorf("failed to scan inventory: %w", err)
	}

	return &inventory, nil
}

func (r *inventoryRepository) Search(ctx context.Context, keyword string) ([]entity.Inventory, error) {
	query := `
		SELECT *
		FROM inventories
		WHERE sku ILIKE :keyword OR name ILIKE :keyword OR customer ILIKE :keyword
		ORDER BY created_at DESC
	`
	params := map[string]any{
		"keyword": "%" + keyword + "%",
	}

	rows, err := r.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to search inventories: %w", err)
	}
	defer rows.Close()

	var inventories []entity.Inventory
	for rows.Next() {
		var inventory entity.Inventory
		if err := rows.StructScan(&inventory); err != nil {
			return nil, fmt.Errorf("failed to scan inventory: %w", err)
		}
		inventories = append(inventories, inventory)
	}

	return inventories, nil
}
