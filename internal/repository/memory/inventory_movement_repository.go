package memory

import (
	"context"
	"fmt"
	"strings"

	"pharmacy-storage-be/internal/domain"
	"pharmacy-storage-be/internal/errs"
)

type InventoryMovementRepository struct {
	movements []domain.InventoryMovement
}

func NewInventoryMovementRepository() *InventoryMovementRepository {
	return &InventoryMovementRepository{
		movements: []domain.InventoryMovement{},
	}
}

func (r *InventoryMovementRepository) Create(ctx context.Context, movement *domain.InventoryMovement) error {
	if movement == nil {
		return fmt.Errorf("inventory movement is required: %w", errs.ErrBadRequest)
	}

	r.movements = append(r.movements, *movement)
	return nil
}

func (r *InventoryMovementRepository) FindAll(ctx context.Context, page int, limit int, search string) ([]domain.InventoryMovement, int, error) {
	var filtered []domain.InventoryMovement

	search = strings.TrimSpace(strings.ToLower(search))

	if search == "" {
		filtered = r.movements
	} else {
		for _, movement := range r.movements {
			if containsMovementSearch(movement, search) {
				filtered = append(filtered, movement)
			}
		}
	}

	total := len(filtered)

	start := (page - 1) * limit
	end := start + limit

	if start >= total {
		return []domain.InventoryMovement{}, total, nil
	}

	if end > total {
		end = total
	}

	return filtered[start:end], total, nil
}

func (r *InventoryMovementRepository) FindByID(ctx context.Context, id string) (*domain.InventoryMovement, error) {
	for _, movement := range r.movements {
		if movement.ID == id {
			result := movement
			return &result, nil
		}
	}

	return nil, fmt.Errorf("inventory movement with id %s not found: %w", id, errs.ErrNotFound)
}

func (r *InventoryMovementRepository) FindByBatchID(ctx context.Context, batchID string) ([]domain.InventoryMovement, error) {
	var result []domain.InventoryMovement

	for _, movement := range r.movements {
		if movement.ProductBatchID == batchID {
			result = append(result, movement)
		}
	}

	return result, nil
}

func containsMovementSearch(movement domain.InventoryMovement, search string) bool {
	return strings.Contains(strings.ToLower(movement.ID), search) ||
		strings.Contains(strings.ToLower(movement.ProductID), search) ||
		strings.Contains(strings.ToLower(movement.ProductBatchID), search) ||
		strings.Contains(strings.ToLower(movement.MovementType), search) ||
		strings.Contains(strings.ToLower(movement.ReferenceNo), search) ||
		strings.Contains(strings.ToLower(movement.Note), search) ||
		strings.Contains(strings.ToLower(movement.CreatedBy), search)
}
