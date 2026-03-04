package usecase

import (
	"context"
	"fmt"
	"irfanard27/incore-api/internal/domain/entity"
	"irfanard27/incore-api/internal/domain/repository"
	"time"
)

type StockOutUsecase interface {
	// Stage 1: Allocation (Draft)
	CreateStockOut(ctx context.Context, items []entity.StockOutItem) (*entity.StockOut, error)
	AllocateStock(ctx context.Context, stockOutID string) (*entity.StockOut, error)

	// Stage 2: Execution (In Progress/Done)
	ExecuteStockOut(ctx context.Context, stockOutID string) (*entity.StockOut, error)
	CompleteStockOut(ctx context.Context, stockOutID string) (*entity.StockOut, error)

	// Rollback
	CancelStockOut(ctx context.Context, stockOutID string) (*entity.StockOut, error)

	// Queries
	GetStockOutById(ctx context.Context, id string) (*entity.StockOut, []entity.StockOutItem, []entity.Inventory, error)
	GetStockOutByTransactionID(ctx context.Context, transactionID string) (*entity.StockOut, error)
	GetAllStockOut(ctx context.Context) ([]entity.StockOut, []int, error)
	DeleteStockOut(ctx context.Context, id string) error
}

type stockOutUsecase struct {
	stockOutRepo     repository.StockOutRepository
	stockOutItemRepo repository.StockOutItemRepository
	inventoryRepo    repository.InventoryRepository
}

func NewStockOutUsecase(stockOutRepo repository.StockOutRepository, stockOutItemRepo repository.StockOutItemRepository, inventoryRepo repository.InventoryRepository) StockOutUsecase {
	return &stockOutUsecase{
		stockOutRepo:     stockOutRepo,
		stockOutItemRepo: stockOutItemRepo,
		inventoryRepo:    inventoryRepo,
	}
}

func (s *stockOutUsecase) GetAllStockOut(ctx context.Context) ([]entity.StockOut, []int, error) {
	stockOuts, totalItems, err := s.stockOutRepo.All(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get all stock outs: %w", err)
	}
	return stockOuts, totalItems, nil
}

// CreateStockOut creates a new stock out with DRAFT status (Stage 1)
func (s *stockOutUsecase) CreateStockOut(ctx context.Context, items []entity.StockOutItem) (*entity.StockOut, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("items must be provided")
	}

	// Validate that all inventories exist and have sufficient stock
	for i := range items {
		inventory, err := s.inventoryRepo.GetById(ctx, items[i].InventoryID)
		if err != nil {
			return nil, fmt.Errorf("inventory %s not found: %w", items[i].InventoryID, err)
		}

		availableStock := inventory.Quantity - inventory.ReservedQuantity
		if availableStock < items[i].Quantity {
			return nil, fmt.Errorf("insufficient stock for inventory %s. Available: %d, Requested: %d",
				items[i].InventoryID, availableStock, items[i].Quantity)
		}
	}

	// Create new stock out with DRAFT status
	stockOut := &entity.StockOut{
		TransactionID: generateStockOutTransactionID(),
		Status:        entity.StockOutStatusDraft,
	}

	id, err := s.stockOutRepo.Create(ctx, stockOut)
	if err != nil {
		return nil, fmt.Errorf("failed to create stock out: %w", err)
	}

	// Set the ID returned from repository
	stockOut.ID = id

	// Set stock_out_id for each item before batch create
	for i := range items {
		items[i].StockOutID = id
	}

	err = s.stockOutItemRepo.BatchCreate(ctx, items)
	if err != nil {
		return nil, fmt.Errorf("failed to batch create stock out items: %w", err)
	}

	return stockOut, nil
}

// AllocateStock performs stock allocation and changes status to ALLOCATED (Stage 1 completion)
func (s *stockOutUsecase) AllocateStock(ctx context.Context, stockOutID string) (*entity.StockOut, error) {
	// Check if stock out exists and is in DRAFT status
	stockOut, err := s.stockOutRepo.GetById(ctx, stockOutID)
	if err != nil {
		return nil, fmt.Errorf("stock out not found: %w", err)
	}

	if stockOut.Status != entity.StockOutStatusDraft {
		return nil, fmt.Errorf("cannot allocate stock: current status is '%s', expected '%s'",
			stockOut.Status, entity.StockOutStatusDraft)
	}

	// Perform stock allocation in repository (with transaction handling)
	err = s.stockOutRepo.AllocateStock(ctx, stockOutID)
	if err != nil {
		return nil, fmt.Errorf("failed to allocate stock: %w", err)
	}

	// Get updated stock out
	stockOut, err = s.stockOutRepo.GetById(ctx, stockOutID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated stock out: %w", err)
	}

	return stockOut, nil
}

// ExecuteStockOut changes status to IN_PROGRESS and decreases actual stock (Stage 2 start)
func (s *stockOutUsecase) ExecuteStockOut(ctx context.Context, stockOutID string) (*entity.StockOut, error) {
	// Check if stock out exists and is in ALLOCATED status
	stockOut, err := s.stockOutRepo.GetById(ctx, stockOutID)
	if err != nil {
		return nil, fmt.Errorf("stock out not found: %w", err)
	}

	if stockOut.Status != entity.StockOutStatusAllocated {
		return nil, fmt.Errorf("cannot execute stock out: current status is '%s', expected '%s'",
			stockOut.Status, entity.StockOutStatusAllocated)
	}

	// Execute stock out (decrease actual quantities)
	err = s.stockOutRepo.ExecuteStockOut(ctx, stockOutID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute stock out: %w", err)
	}

	// Get updated stock out
	stockOut, err = s.stockOutRepo.GetById(ctx, stockOutID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated stock out: %w", err)
	}

	return stockOut, nil
}

// CompleteStockOut changes status to DONE (Stage 2 completion)
func (s *stockOutUsecase) CompleteStockOut(ctx context.Context, stockOutID string) (*entity.StockOut, error) {
	// Check if stock out exists and is in IN_PROGRESS status
	stockOut, err := s.stockOutRepo.GetById(ctx, stockOutID)
	if err != nil {
		return nil, fmt.Errorf("stock out not found: %w", err)
	}

	if stockOut.Status != entity.StockOutStatusInProgress {
		return nil, fmt.Errorf("cannot complete stock out: current status is '%s', expected '%s'",
			stockOut.Status, entity.StockOutStatusInProgress)
	}

	// Complete stock out
	err = s.stockOutRepo.CompleteStockOut(ctx, stockOutID)
	if err != nil {
		return nil, fmt.Errorf("failed to complete stock out: %w", err)
	}

	// Get updated stock out
	stockOut, err = s.stockOutRepo.GetById(ctx, stockOutID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated stock out: %w", err)
	}

	return stockOut, nil
}

// CancelStockOut performs rollback and changes status to CANCELLED
func (s *stockOutUsecase) CancelStockOut(ctx context.Context, stockOutID string) (*entity.StockOut, error) {
	// Check if stock out exists
	stockOut, err := s.stockOutRepo.GetById(ctx, stockOutID)
	if err != nil {
		return nil, fmt.Errorf("stock out not found: %w", err)
	}

	// Can only cancel from ALLOCATED or IN_PROGRESS status
	if stockOut.Status != entity.StockOutStatusAllocated && stockOut.Status != entity.StockOutStatusInProgress {
		return nil, fmt.Errorf("cannot cancel stock out: current status is '%s', can only cancel from '%s' or '%s'",
			stockOut.Status, entity.StockOutStatusAllocated, entity.StockOutStatusInProgress)
	}

	// Perform rollback
	err = s.stockOutRepo.RollbackStockOut(ctx, stockOutID)
	if err != nil {
		return nil, fmt.Errorf("failed to rollback stock out: %w", err)
	}

	// Get updated stock out
	stockOut, err = s.stockOutRepo.GetById(ctx, stockOutID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated stock out: %w", err)
	}

	return stockOut, nil
}

func (s *stockOutUsecase) GetStockOutById(ctx context.Context, id string) (*entity.StockOut, []entity.StockOutItem, []entity.Inventory, error) {
	stockOut, err := s.stockOutRepo.GetById(ctx, id)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get stock out by id: %w", err)
	}

	items, inventories, err := s.stockOutItemRepo.GetByStockOutIDWithInventory(ctx, id)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get stock out items: %w", err)
	}

	return stockOut, items, inventories, nil
}

func (s *stockOutUsecase) GetStockOutByTransactionID(ctx context.Context, transactionID string) (*entity.StockOut, error) {
	stockOut, err := s.stockOutRepo.GetByTransactionID(ctx, transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock out by transaction id: %w", err)
	}
	return stockOut, nil
}

func (s *stockOutUsecase) DeleteStockOut(ctx context.Context, id string) error {
	// First check if stock out exists
	stockOut, err := s.stockOutRepo.GetById(ctx, id)
	if err != nil {
		return fmt.Errorf("stock out not found: %w", err)
	}

	// Can only delete from DRAFT status
	if stockOut.Status != entity.StockOutStatusDraft {
		return fmt.Errorf("cannot delete stock out: current status is '%s', can only delete from '%s'",
			stockOut.Status, entity.StockOutStatusDraft)
	}

	err = s.stockOutRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete stock out: %w", err)
	}

	return nil
}

func generateStockOutTransactionID() string {
	return "STO-" + time.Now().Format("20060102150405")
}
