package repository

import (
	"context"

	"pharmacy-storage-be/internal/domain"
)

type InventoryMovementRepository interface {
	Create(ctx context.Context, movement *domain.InventoryMovement) error
	FindAll(ctx context.Context, page int, limit int, search string) ([]domain.InventoryMovement, int, error)
	FindByID(ctx context.Context, id string) (*domain.InventoryMovement, error)
	FindByBatchID(ctx context.Context, batchID string) ([]domain.InventoryMovement, error)
}
