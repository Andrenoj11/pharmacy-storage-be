package repository

import (
	"context"

	"pharmacy-storage-be/internal/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	FindAll(ctx context.Context, page int, limit int, search string) ([]domain.Product, int, error)
	FindByID(ctx context.Context, id string) (*domain.Product, error)
}
