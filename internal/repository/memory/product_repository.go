package memory

import (
	"context"
	"fmt"
	"strings"

	"pharmacy-storage-be/internal/domain"
	"pharmacy-storage-be/internal/errs"
)

type ProductRepository struct {
	products []domain.Product
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		products: []domain.Product{},
	}
}

func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) error {
	if product == nil {
		return fmt.Errorf("product is required: %w", errs.ErrBadRequest)
	}

	r.products = append(r.products, *product)
	return nil
}

func (r *ProductRepository) FindAll(ctx context.Context, page int, limit int, search string) ([]domain.Product, int, error) {
	var filtered []domain.Product

	search = strings.TrimSpace(strings.ToLower(search))

	if search == "" {
		filtered = r.products
	} else {
		for _, product := range r.products {
			if containsSearch(product, search) {
				filtered = append(filtered, product)
			}
		}
	}

	total := len(filtered)

	start := (page - 1) * limit
	end := start + limit

	if start >= total {
		return []domain.Product{}, total, nil
	}

	if end > total {
		end = total
	}

	return filtered[start:end], total, nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	for _, product := range r.products {
		if product.ID == id {
			result := product
			return &result, nil
		}
	}

	return nil, fmt.Errorf("product with id %s not found: %w", id, errs.ErrNotFound)
}

func containsSearch(product domain.Product, search string) bool {
	return strings.Contains(strings.ToLower(product.Code), search) ||
		strings.Contains(strings.ToLower(product.Name), search) ||
		strings.Contains(strings.ToLower(product.Category), search) ||
		strings.Contains(strings.ToLower(product.Unit), search) ||
		strings.Contains(strings.ToLower(product.StorageType), search)
}
