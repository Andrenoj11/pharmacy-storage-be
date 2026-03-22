package spreadsheet

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

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
	if product == nil {
		return fmt.Errorf("product is required: %w", errs.ErrBadRequest)
	}

	row := []interface{}{
		product.ID,
		product.Code,
		product.Name,
		product.Category,
		product.Unit,
		strconv.Itoa(product.MinStock),
		product.StorageType,
		product.CreatedAt.Format(time.RFC3339),
		product.UpdatedAt.Format(time.RFC3339),
	}

	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{row},
	}

	rangeName := fmt.Sprintf("%s!A:I", r.sheetName)

	_, err := r.sheetsService.Spreadsheets.Values.
		Append(r.spreadsheetID, rangeName, valueRange).
		ValueInputOption("RAW").
		Context(ctx).
		Do()
	if err != nil {
		return fmt.Errorf("failed to append product row to spreadsheet: %w", err)
	}

	return nil
}

func (r *ProductRepository) FindAll(ctx context.Context, page int, limit int, search string) ([]domain.Product, int, error) {
	rangeName := fmt.Sprintf("%s!A2:I", r.sheetName)

	resp, err := r.sheetsService.Spreadsheets.Values.
		Get(r.spreadsheetID, rangeName).
		Context(ctx).
		Do()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get product rows from spreadsheet: %w", err)
	}

	var products []domain.Product
	search = strings.TrimSpace(strings.ToLower(search))

	for _, row := range resp.Values {
		product := mapRowToProduct(row)

		if search == "" || containsSearch(product, search) {
			products = append(products, product)
		}
	}

	total := len(products)

	start := (page - 1) * limit
	end := start + limit

	if start >= total {
		return []domain.Product{}, total, nil
	}

	if end > total {
		end = total
	}

	return products[start:end], total, nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	return nil, fmt.Errorf("spreadsheet find product by id is not implemented yet: %w", errs.ErrBadRequest)
}

func mapRowToProduct(row []interface{}) domain.Product {
	product := domain.Product{}

	if len(row) > 0 {
		product.ID = toString(row[0])
	}
	if len(row) > 1 {
		product.Code = toString(row[1])
	}
	if len(row) > 2 {
		product.Name = toString(row[2])
	}
	if len(row) > 3 {
		product.Category = toString(row[3])
	}
	if len(row) > 4 {
		product.Unit = toString(row[4])
	}
	if len(row) > 5 {
		product.MinStock = toInt(row[5])
	}
	if len(row) > 6 {
		product.StorageType = toString(row[6])
	}
	if len(row) > 7 {
		product.CreatedAt = toTime(row[7])
	}
	if len(row) > 8 {
		product.UpdatedAt = toTime(row[8])
	}

	return product
}

func containsSearch(product domain.Product, search string) bool {
	return strings.Contains(strings.ToLower(product.Code), search) ||
		strings.Contains(strings.ToLower(product.Name), search) ||
		strings.Contains(strings.ToLower(product.Category), search) ||
		strings.Contains(strings.ToLower(product.Unit), search) ||
		strings.Contains(strings.ToLower(product.StorageType), search)
}

func toString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

func toInt(value interface{}) int {
	strValue := strings.TrimSpace(fmt.Sprintf("%v", value))
	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		return 0
	}
	return intValue
}

func toTime(value interface{}) time.Time {
	strValue := strings.TrimSpace(fmt.Sprintf("%v", value))
	parsedTime, err := time.Parse(time.RFC3339, strValue)
	if err != nil {
		return time.Time{}
	}
	return parsedTime
}
