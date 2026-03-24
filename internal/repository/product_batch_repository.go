package repository

import (
	"context"

	"pharmacy-storage-be/internal/domain"
)

type ProductBatchRepository interface {
	Create(ctx context.Context, batch *domain.ProductBatch) error
	FindAll(ctx context.Context, page int, limit int, search string) ([]domain.ProductBatch, int, error)
	FindByID(ctx context.Context, id string) (*domain.ProductBatch, error)
	FindByProductID(ctx context.Context, productID string) ([]domain.ProductBatch, error)
	FindFEFOBatchByProductID(ctx context.Context, productID string) (*domain.ProductBatch, error)
	UpdateQtyAvailable(ctx context.Context, batchID string, qtyAvailable int) error
}
