package spreadsheet

import (
	"context"
	"fmt"

	"pharmacy-storage-be/internal/domain"
	"pharmacy-storage-be/internal/errs"

	"google.golang.org/api/sheets/v4"
)

type ProductRepository struct {
	sheetsService *sheets.Service
	spreadsheetID string
	sheetName     string
}

func NewProductRepository(sheetsService *sheets.Service, spreadsheetID string, sheetName string) *ProductRepository {
	return &ProductRepository{
		sheetsService: sheetsService,
		spreadsheetID: spreadsheetID,
		sheetName:     sheetName,
	}
}

func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) error {
	return fmt.Errorf("spreadsheet create product is not implemented yet: %w", errs.ErrBadRequest)
}

func (r *ProductRepository) FindAll(ctx context.Context, page int, limit int, search string) ([]domain.Product, int, error) {
	return nil, 0, fmt.Errorf("spreadsheet find all products is not implemented yet: %w", errs.ErrBadRequest)
}

func (r *ProductRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	return nil, fmt.Errorf("spreadsheet find product by id is not implemented yet: %w", errs.ErrBadRequest)
}
