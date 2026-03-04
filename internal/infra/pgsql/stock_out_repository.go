package pgsql

import (
	"context"
	"fmt"
	"irfanard27/incore-api/internal/domain/entity"
	"irfanard27/incore-api/internal/domain/repository"

	"github.com/jmoiron/sqlx"
)

type stockOutRepository struct {
	db *sqlx.DB
}

type stockOutWithTotalItem struct {
	entity.StockOut
	TotalItem int `db:"total_item"`
}

func NewStockOutRepository(db *sqlx.DB) repository.StockOutRepository {
	return &stockOutRepository{db: db}
}

func (r *stockOutRepository) All(ctx context.Context) ([]entity.StockOut, []int, error) {
	query := `
		SELECT s.*, COALESCE(COUNT(soi.id), 0) as total_item
		FROM stock_out s 
		LEFT JOIN stock_out_items soi ON s.id = soi.stock_out_id
		GROUP BY s.id
		ORDER BY s.created_at DESC
	`
	var stockOutsWithTotal []stockOutWithTotalItem
	err := r.db.SelectContext(ctx, &stockOutsWithTotal, query)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get all stock outs: %w", err)
	}

	// Convert to entity.StockOut (without total_item) and collect totals
	stockOuts := make([]entity.StockOut, len(stockOutsWithTotal))
	totalItems := make([]int, len(stockOutsWithTotal))
	for i, s := range stockOutsWithTotal {
		stockOuts[i] = s.StockOut
		totalItems[i] = s.TotalItem
	}
	return stockOuts, totalItems, nil
}

func (r *stockOutRepository) GetById(ctx context.Context, id string) (*entity.StockOut, error) {
	query := `
		SELECT s.*
		FROM stock_out s 
		WHERE s.id = $1
	`
	var stockOut entity.StockOut
	err := r.db.GetContext(ctx, &stockOut, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock out by id: %w", err)
	}

	return &stockOut, nil
}

func (r *stockOutRepository) GetByTransactionID(ctx context.Context, transactionID string) (*entity.StockOut, error) {
	query := `
		SELECT *
		FROM stock_out
		WHERE transaction_id = :transaction_id
	`
	params := map[string]any{
		"transaction_id": transactionID,
	}
	rows, err := r.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock out by transaction id: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("stock out not found")
	}

	var stockOut entity.StockOut
	err = rows.StructScan(&stockOut)
	if err != nil {
		return nil, fmt.Errorf("failed to scan stock out: %w", err)
	}

	return &stockOut, nil
}

func (r *stockOutRepository) Create(ctx context.Context, stockOut *entity.StockOut) (string, error) {
	query := `
		INSERT INTO stock_out (transaction_id, status, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id
	`
	var id string
	err := r.db.QueryRowxContext(ctx, query, stockOut.TransactionID, stockOut.Status).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to create stock out: %w", err)
	}
	return id, nil
}

func (r *stockOutRepository) Update(ctx context.Context, stockOut *entity.StockOut) error {
	query := `
		UPDATE stock_out
		SET transaction_id = :transaction_id, status = :status, updated_at = NOW()
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, stockOut)
	if err != nil {
		return fmt.Errorf("failed to update stock out: %w", err)
	}
	return nil
}

func (r *stockOutRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM stock_out WHERE id = :id`
	params := map[string]any{
		"id": id,
	}
	_, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to delete stock out: %w", err)
	}
	return nil
}

func (r *stockOutRepository) AllocateStock(ctx context.Context, stockOutID string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get stock out items with inventory details
	query := `
		SELECT soi.id, soi.stock_out_id, soi.inventory_id, soi.quantity, i.quantity as available_quantity, i.reserved_quantity
		FROM stock_out_items soi
		JOIN inventories i ON soi.inventory_id = i.id
		WHERE soi.stock_out_id = $1
		FOR UPDATE
	`
	var items []entity.StockOutItemWithAvailability
	err = tx.SelectContext(ctx, &items, query, stockOutID)
	if err != nil {
		return fmt.Errorf("failed to get stock out items: %w", err)
	}

	for _, item := range items {
		// Debug: Log current state before allocation
		fmt.Printf("DEBUG: Before allocation - InventoryID: %s, Available: %d, Reserved: %d, ToReserve: %d\n",
			item.InventoryID, item.AvailableQuantity, item.ReservedQuantity, item.Quantity)

		// Check if stock is sufficient
		if item.AvailableQuantity-item.ReservedQuantity < item.Quantity {
			return fmt.Errorf("insufficient stock for inventory %s. Available: %d, Reserved: %d, Requested: %d",
				item.InventoryID, item.AvailableQuantity, item.ReservedQuantity, item.Quantity)
		}

		// Reserve the stock
		updateQuery := `
			UPDATE inventories 
			SET reserved_quantity = reserved_quantity + $1, updated_at = NOW()
			WHERE id = $2
		`
		_, err = tx.ExecContext(ctx, updateQuery, item.Quantity, item.InventoryID)
		if err != nil {
			return fmt.Errorf("failed to reserve stock: %w", err)
		}
	}

	// Update stock out status to ALLOCATED
	updateStatusQuery := `
		UPDATE stock_out 
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err = tx.ExecContext(ctx, updateStatusQuery, entity.StockOutStatusAllocated, stockOutID)
	if err != nil {
		return fmt.Errorf("failed to update stock out status: %w", err)
	}

	return tx.Commit()
}

func (r *stockOutRepository) ExecuteStockOut(ctx context.Context, stockOutID string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get stock out items with inventory details
	query := `
		SELECT soi.id, soi.stock_out_id, soi.inventory_id, soi.quantity, i.quantity as available_quantity, i.reserved_quantity
		FROM stock_out_items soi
		JOIN inventories i ON soi.inventory_id = i.id
		WHERE soi.stock_out_id = $1
		FOR UPDATE
	`
	var items []entity.StockOutItemWithAvailability
	err = tx.SelectContext(ctx, &items, query, stockOutID)
	if err != nil {
		return fmt.Errorf("failed to get stock out items: %w", err)
	}

	for _, item := range items {

		actualDeduct := item.Quantity
		fmt.Println("actualDeduct", actualDeduct)
		updateQuery := `
			UPDATE inventories 
			SET quantity = quantity - $1, 
				reserved_quantity = reserved_quantity - $1,
				updated_at = NOW()
			WHERE id = $2
		`
		_, err = tx.ExecContext(ctx, updateQuery, actualDeduct, item.InventoryID)
		if err != nil {
			return fmt.Errorf("failed to execute stock out: %w", err)
		}
	}

	// Update stock out status to IN_PROGRESS
	updateStatusQuery := `
		UPDATE stock_out 
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err = tx.ExecContext(ctx, updateStatusQuery, entity.StockOutStatusInProgress, stockOutID)
	if err != nil {
		return fmt.Errorf("failed to update stock out status: %w", err)
	}

	return tx.Commit()
}

func (r *stockOutRepository) RollbackStockOut(ctx context.Context, stockOutID string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get stock out items with inventory details
	query := `
		SELECT soi.id, soi.stock_out_id, soi.inventory_id, soi.quantity, i.reserved_quantity
		FROM stock_out_items soi
		JOIN inventories i ON soi.inventory_id = i.id
		WHERE soi.stock_out_id = $1
		FOR UPDATE
	`
	var items []entity.StockOutItemWithReserved
	err = tx.SelectContext(ctx, &items, query, stockOutID)
	if err != nil {
		return fmt.Errorf("failed to get stock out items: %w", err)
	}

	for _, item := range items {
		// Debug: Log current state before rollback
		fmt.Printf("DEBUG: Before rollback - InventoryID: %s, Quantity: %d, Reserved: %d, ToReturn: %d\n",
			item.InventoryID, item.Quantity, item.ReservedQuantity, item.Quantity)

		// Return reserved stock back to available (increase reserved_quantity)
		updateQuery := `
			UPDATE inventories 
			SET reserved_quantity = reserved_quantity - $1, updated_at = NOW()
			WHERE id = $2
		`
		_, err = tx.ExecContext(ctx, updateQuery, item.Quantity, item.InventoryID)
		if err != nil {
			return fmt.Errorf("failed to rollback stock: %w", err)
		}
	}

	// Update stock out status to CANCELLED
	updateStatusQuery := `
		UPDATE stock_out 
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err = tx.ExecContext(ctx, updateStatusQuery, entity.StockOutStatusCancelled, stockOutID)
	if err != nil {
		return fmt.Errorf("failed to update stock out status: %w", err)
	}

	return tx.Commit()
}

func (r *stockOutRepository) CompleteStockOut(ctx context.Context, stockOutID string) error {
	query := `
		UPDATE stock_out 
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, entity.StockOutStatusDone, stockOutID)
	if err != nil {
		return fmt.Errorf("failed to complete stock out: %w", err)
	}
	return nil
}
