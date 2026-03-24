package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"pharmacy-storage-be/internal/domain"
	"pharmacy-storage-be/internal/errs"
	"pharmacy-storage-be/internal/repository"
)

type InventoryMovementService struct {
	movementRepo     repository.InventoryMovementRepository
	productRepo      repository.ProductRepository
	productBatchRepo repository.ProductBatchRepository
}

func NewInventoryMovementService(
	movementRepo repository.InventoryMovementRepository,
	productRepo repository.ProductRepository,
	productBatchRepo repository.ProductBatchRepository,
) *InventoryMovementService {
	return &InventoryMovementService{
		movementRepo:     movementRepo,
		productRepo:      productRepo,
		productBatchRepo: productBatchRepo,
	}
}

func (s *InventoryMovementService) CreateInventoryMovement(ctx context.Context, input *domain.InventoryMovement) error {
	if input == nil {
		return fmt.Errorf("inventory movement input is required: %w", errs.ErrBadRequest)
	}

	if strings.TrimSpace(input.ProductID) == "" {
		return fmt.Errorf("product id is required: %w", errs.ErrBadRequest)
	}

	if strings.TrimSpace(input.ProductBatchID) == "" {
		return fmt.Errorf("product batch id is required: %w", errs.ErrBadRequest)
	}

	if strings.TrimSpace(input.MovementType) == "" {
		return fmt.Errorf("movement type is required: %w", errs.ErrBadRequest)
	}

	if input.Qty <= 0 {
		return fmt.Errorf("qty must be greater than zero: %w", errs.ErrBadRequest)
	}

	if input.MovementDate.IsZero() {
		return fmt.Errorf("movement date is required: %w", errs.ErrBadRequest)
	}

	productID := strings.TrimSpace(input.ProductID)
	productBatchID := strings.TrimSpace(input.ProductBatchID)
	movementType := strings.ToUpper(strings.TrimSpace(input.MovementType))

	if movementType != "IN" && movementType != "OUT" && movementType != "ADJUSTMENT" {
		return fmt.Errorf("movement type must be one of IN, OUT, or ADJUSTMENT: %w", errs.ErrBadRequest)
	}

	_, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return err
	}

	batch, err := s.productBatchRepo.FindByID(ctx, productBatchID)
	if err != nil {
		return err
	}

	if batch.ProductID != productID {
		return fmt.Errorf("product batch does not belong to the given product: %w", errs.ErrBadRequest)
	}

	newQtyAvailable := batch.QtyAvailable

	switch movementType {
	case "IN":
		newQtyAvailable = batch.QtyAvailable + input.Qty
	case "OUT":
		if input.Qty > batch.QtyAvailable {
			return fmt.Errorf("insufficient batch stock: %w", errs.ErrBadRequest)
		}
		newQtyAvailable = batch.QtyAvailable - input.Qty
	case "ADJUSTMENT":
		newQtyAvailable = batch.QtyAvailable + input.Qty
	}

	err = s.productBatchRepo.UpdateQtyAvailable(ctx, productBatchID, newQtyAvailable)
	if err != nil {
		return err
	}

	now := time.Now()

	if input.ID == "" {
		input.ID = generateMovementID()
	}

	input.ProductID = productID
	input.ProductBatchID = productBatchID
	input.MovementType = movementType
	input.ReferenceNo = strings.TrimSpace(input.ReferenceNo)
	input.Note = strings.TrimSpace(input.Note)
	input.CreatedBy = strings.TrimSpace(input.CreatedBy)
	input.CreatedAt = now
	input.UpdatedAt = now

	return s.movementRepo.Create(ctx, input)
}

func (s *InventoryMovementService) GetAllInventoryMovements(ctx context.Context, page int, limit int, search string) ([]domain.InventoryMovement, int, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	search = strings.TrimSpace(search)

	return s.movementRepo.FindAll(ctx, page, limit, search)
}

func (s *InventoryMovementService) GetInventoryMovementByID(ctx context.Context, id string) (*domain.InventoryMovement, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("inventory movement id is required: %w", errs.ErrBadRequest)
	}

	return s.movementRepo.FindByID(ctx, id)
}

func (s *InventoryMovementService) GetInventoryMovementsByBatchID(ctx context.Context, batchID string) ([]domain.InventoryMovement, error) {
	batchID = strings.TrimSpace(batchID)
	if batchID == "" {
		return nil, fmt.Errorf("product batch id is required: %w", errs.ErrBadRequest)
	}

	_, err := s.productBatchRepo.FindByID(ctx, batchID)
	if err != nil {
		return nil, err
	}

	return s.movementRepo.FindByBatchID(ctx, batchID)
}

func generateMovementID() string {
	return "MOV-" + time.Now().Format("20060102150405")
}
