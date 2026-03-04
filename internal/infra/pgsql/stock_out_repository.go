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

func NewStockOutRepository(db *sqlx.DB) repository.StockOutRepository {
	return &stockOutRepository{db: db}
}

func (r *stockOutRepository) Create(ctx context.Context, stockOut *entity.StockOut) error {
	query := `
		INSERT INTO stock_out (id, transaction_id, status, created_at, updated_at)
		VALUES (:id, :transaction_id, :status, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, stockOut)
	if err != nil {
		return fmt.Errorf("failed to create stock out: %w", err)
	}
	return nil
}

func (r *stockOutRepository) GetById(ctx context.Context, id string) (*entity.StockOut, error) {
	query := `
		SELECT id, transaction_id, status, created_at, updated_at
		FROM stock_out
		WHERE id = :id
	`
	params := map[string]interface{}{
		"id": id,
	}
	rows, err := r.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock out by id: %w", err)
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

func (r *stockOutRepository) GetByTransactionID(ctx context.Context, transactionID string) (*entity.StockOut, error) {
	query := `
		SELECT id, transaction_id, status, created_at, updated_at
		FROM stock_out
		WHERE transaction_id = :transaction_id
	`
	params := map[string]interface{}{
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

func (r *stockOutRepository) Update(ctx context.Context, stockOut *entity.StockOut) error {
	query := `
		UPDATE stock_out
		SET transaction_id = :transaction_id, status = :status, updated_at = :updated_at
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
	params := map[string]interface{}{
		"id": id,
	}
	_, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to delete stock out: %w", err)
	}
	return nil
}

func (r *stockOutRepository) All(ctx context.Context) ([]entity.StockOut, error) {
	query := `
		SELECT id, transaction_id, status, created_at, updated_at
		FROM stock_out
		ORDER BY created_at DESC
	`
	var stockOuts []entity.StockOut
	err := r.db.SelectContext(ctx, &stockOuts, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all stock outs: %w", err)
	}
	return stockOuts, nil
}
