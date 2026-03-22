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

type ProductService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, input *domain.Product) error {
	if input == nil {
		return fmt.Errorf("product input is required: %w", errs.ErrBadRequest)
	}

	if strings.TrimSpace(input.Code) == "" {
		return fmt.Errorf("product code is required: %w", errs.ErrBadRequest)
	}

	if strings.TrimSpace(input.Name) == "" {
		return fmt.Errorf("product name is required: %w", errs.ErrBadRequest)
	}

	if strings.TrimSpace(input.Unit) == "" {
		return fmt.Errorf("product unit is required: %w", errs.ErrBadRequest)
	}

	now := time.Now()

	if input.ID == "" {
		input.ID = generateID()
	}

	input.Code = strings.TrimSpace(input.Code)
	input.Name = strings.TrimSpace(input.Name)
	input.Category = strings.TrimSpace(input.Category)
	input.Unit = strings.TrimSpace(input.Unit)
	input.StorageType = strings.TrimSpace(input.StorageType)
	input.CreatedAt = now
	input.UpdatedAt = now

	return s.productRepo.Create(ctx, input)
}

func (s *ProductService) GetAllProducts(ctx context.Context, page int, limit int, search string) ([]domain.Product, int, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	search = strings.TrimSpace(search)

	return s.productRepo.FindAll(ctx, page, limit, search)
}

func (s *ProductService) GetProductByID(ctx context.Context, id string) (*domain.Product, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("product id is required: %w", errs.ErrBadRequest)
	}

	return s.productRepo.FindByID(ctx, id)
}

func generateID() string {
	return "PRD-" + time.Now().Format("20060102150405")
}
