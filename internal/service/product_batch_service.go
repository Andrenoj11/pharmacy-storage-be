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

type ProductBatchService struct {
	productBatchRepo repository.ProductBatchRepository
	productRepo      repository.ProductRepository
}

func NewProductBatchService(
	productBatchRepo repository.ProductBatchRepository,
	productRepo repository.ProductRepository,
) *ProductBatchService {
	return &ProductBatchService{
		productBatchRepo: productBatchRepo,
		productRepo:      productRepo,
	}
}

func (s *ProductBatchService) CreateProductBatch(ctx context.Context, input *domain.ProductBatch) error {
	if input == nil {
		return fmt.Errorf("product batch input is required: %w", errs.ErrBadRequest)
	}

	if strings.TrimSpace(input.ProductID) == "" {
		return fmt.Errorf("product id is required: %w", errs.ErrBadRequest)
	}

	if strings.TrimSpace(input.BatchNumber) == "" {
		return fmt.Errorf("batch number is required: %w", errs.ErrBadRequest)
	}

	if input.ExpiryDate.IsZero() {
		return fmt.Errorf("expiry date is required: %w", errs.ErrBadRequest)
	}

	if input.ReceivedDate.IsZero() {
		return fmt.Errorf("received date is required: %w", errs.ErrBadRequest)
	}

	if input.QtyAvailable < 0 {
		return fmt.Errorf("qty available cannot be negative: %w", errs.ErrBadRequest)
	}

	productID := strings.TrimSpace(input.ProductID)

	_, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return err
	}

	now := time.Now()

	if input.ID == "" {
		input.ID = generateBatchID()
	}

	input.ProductID = productID
	input.BatchNumber = strings.TrimSpace(input.BatchNumber)
	input.SupplierName = strings.TrimSpace(input.SupplierName)
	input.CreatedAt = now
	input.UpdatedAt = now

	return s.productBatchRepo.Create(ctx, input)
}

func (s *ProductBatchService) GetAllProductBatches(ctx context.Context, page int, limit int, search string) ([]domain.ProductBatch, int, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	search = strings.TrimSpace(search)

	return s.productBatchRepo.FindAll(ctx, page, limit, search)
}

func (s *ProductBatchService) GetProductBatchByID(ctx context.Context, id string) (*domain.ProductBatch, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("product batch id is required: %w", errs.ErrBadRequest)
	}

	return s.productBatchRepo.FindByID(ctx, id)
}

func (s *ProductBatchService) GetProductBatchesByProductID(ctx context.Context, productID string) ([]domain.ProductBatch, error) {
	productID = strings.TrimSpace(productID)
	if productID == "" {
		return nil, fmt.Errorf("product id is required: %w", errs.ErrBadRequest)
	}

	_, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	return s.productBatchRepo.FindByProductID(ctx, productID)
}

func (s *ProductBatchService) GetFEFOBatchByProductID(ctx context.Context, productID string) (*domain.ProductBatch, error) {
	productID = strings.TrimSpace(productID)
	if productID == "" {
		return nil, fmt.Errorf("product id is required: %w", errs.ErrBadRequest)
	}

	_, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	return s.productBatchRepo.FindFEFOBatchByProductID(ctx, productID)
}

func generateBatchID() string {
	return "BAT-" + time.Now().Format("20060102150405")
}
