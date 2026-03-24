package memory

import (
	"context"
	"fmt"
	"strings"

	"pharmacy-storage-be/internal/domain"
	"pharmacy-storage-be/internal/errs"
)

type ProductBatchRepository struct {
	batches []domain.ProductBatch
}

func NewProductBatchRepository() *ProductBatchRepository {
	return &ProductBatchRepository{
		batches: []domain.ProductBatch{},
	}
}

func (r *ProductBatchRepository) Create(ctx context.Context, batch *domain.ProductBatch) error {
	if batch == nil {
		return fmt.Errorf("product batch is required: %w", errs.ErrBadRequest)
	}

	r.batches = append(r.batches, *batch)
	return nil
}

func (r *ProductBatchRepository) FindAll(ctx context.Context, page int, limit int, search string) ([]domain.ProductBatch, int, error) {
	var filtered []domain.ProductBatch

	search = strings.TrimSpace(strings.ToLower(search))

	if search == "" {
		filtered = r.batches
	} else {
		for _, batch := range r.batches {
			if containsBatchSearch(batch, search) {
				filtered = append(filtered, batch)
			}
		}
	}

	total := len(filtered)

	start := (page - 1) * limit
	end := start + limit

	if start >= total {
		return []domain.ProductBatch{}, total, nil
	}

	if end > total {
		end = total
	}

	return filtered[start:end], total, nil
}

func (r *ProductBatchRepository) FindByID(ctx context.Context, id string) (*domain.ProductBatch, error) {
	for _, batch := range r.batches {
		if batch.ID == id {
			result := batch
			return &result, nil
		}
	}

	return nil, fmt.Errorf("product batch with id %s not found: %w", id, errs.ErrNotFound)
}

func (r *ProductBatchRepository) FindByProductID(ctx context.Context, productID string) ([]domain.ProductBatch, error) {
	var result []domain.ProductBatch

	for _, batch := range r.batches {
		if batch.ProductID == productID {
			result = append(result, batch)
		}
	}

	return result, nil
}

func (r *ProductBatchRepository) FindFEFOBatchByProductID(ctx context.Context, productID string) (*domain.ProductBatch, error) {
	var selected *domain.ProductBatch

	for _, batch := range r.batches {
		if batch.ProductID != productID {
			continue
		}

		if batch.QtyAvailable <= 0 {
			continue
		}

		if batch.ExpiryDate.IsZero() {
			continue
		}

		current := batch

		if selected == nil {
			selected = &current
			continue
		}

		if current.ExpiryDate.Before(selected.ExpiryDate) {
			selected = &current
		}
	}

	if selected == nil {
		return nil, fmt.Errorf("no FEFO batch available for product id %s: %w", productID, errs.ErrNotFound)
	}

	return selected, nil
}

func (r *ProductBatchRepository) UpdateQtyAvailable(ctx context.Context, batchID string, qtyAvailable int) error {
	for i, batch := range r.batches {
		if batch.ID == batchID {
			r.batches[i].QtyAvailable = qtyAvailable
			return nil
		}
	}

	return fmt.Errorf("product batch with id %s not found: %w", batchID, errs.ErrNotFound)
}

func containsBatchSearch(batch domain.ProductBatch, search string) bool {
	return strings.Contains(strings.ToLower(batch.ID), search) ||
		strings.Contains(strings.ToLower(batch.ProductID), search) ||
		strings.Contains(strings.ToLower(batch.BatchNumber), search) ||
		strings.Contains(strings.ToLower(batch.SupplierName), search)
}
