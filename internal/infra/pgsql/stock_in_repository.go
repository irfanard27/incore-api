package pgsql

import (
	"context"
	"fmt"
	"irfanard27/incore-api/internal/domain/entity"
	"irfanard27/incore-api/internal/domain/repository"

	"github.com/jmoiron/sqlx"
)

type stockInRepository struct {
	db *sqlx.DB
}

func NewStockInRepository(db *sqlx.DB) repository.StockInRepository {
	return &stockInRepository{db: db}
}

func (r *stockInRepository) All(ctx context.Context) ([]entity.StockIn, error) {
	query := `
		SELECT s.*, COALESCE(COUNT(si.id), 0) as total_item
		FROM stock_in s 
		LEFT JOIN stock_in_items si ON s.id = si.stock_in_id
		GROUP BY s.id
		ORDER BY s.created_at DESC
	`
	var stockIns []entity.StockIn
	err := r.db.SelectContext(ctx, &stockIns, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all stock ins: %w", err)
	}
	return stockIns, nil
}

func (r *stockInRepository) GetById(ctx context.Context, id string) (*entity.StockIn, error) {
	query := `
		SELECT s.*, COALESCE(COUNT(si.id), 0) as total_item
		FROM stock_in s 
		LEFT JOIN stock_in_items si ON s.id = si.stock_in_id
		WHERE s.id = :id
		GROUP BY s.id
	`
	params := map[string]any{
		"id": id,
	}

	rows, err := r.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock in by id: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("stock in not found")
	}

	var stockIn entity.StockIn
	err = rows.StructScan(&stockIn)
	if err != nil {
		return nil, fmt.Errorf("failed to scan stock in: %w", err)
	}

	return &stockIn, nil
}

func (r *stockInRepository) GetByTransactionID(ctx context.Context, transactionID string) (*entity.StockIn, error) {
	query := `
		SELECT *
		FROM stock_in
		WHERE transaction_id = :transaction_id
	`
	params := map[string]any{
		"transaction_id": transactionID,
	}
	rows, err := r.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock in by transaction id: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("stock in not found")
	}

	var stockIn entity.StockIn
	err = rows.StructScan(&stockIn)
	if err != nil {
		return nil, fmt.Errorf("failed to scan stock in: %w", err)
	}

	return &stockIn, nil
}

func (r *stockInRepository) Create(ctx context.Context, stockIn *entity.StockIn) (string, error) {
	query := `
		INSERT INTO stock_in (transaction_id, created_at, updated_at)
		VALUES ($1, NOW(), NOW())
		RETURNING id
	`
	var id string
	err := r.db.QueryRowxContext(ctx, query, stockIn.TransactionID).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to create stock in: %w", err)
	}
	return id, nil
}

func (r *stockInRepository) Update(ctx context.Context, stockIn *entity.StockIn) error {
	query := `
		UPDATE stock_in
		SET transaction_id = :transaction_id, status = :status, updated_at = :updated_at
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, stockIn)
	if err != nil {
		return fmt.Errorf("failed to update stock in: %w", err)
	}
	return nil
}

func (r *stockInRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM stock_in WHERE id = :id`
	params := map[string]any{
		"id": id,
	}
	_, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to delete stock in: %w", err)
	}
	return nil
}
