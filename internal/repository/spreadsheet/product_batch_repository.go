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

type ProductBatchRepository struct {
	sheetsService *sheets.Service
	spreadsheetID string
	sheetName     string
}

func NewProductBatchRepository(sheetsService *sheets.Service, spreadsheetID string, sheetName string) *ProductBatchRepository {
	return &ProductBatchRepository{
		sheetsService: sheetsService,
		spreadsheetID: spreadsheetID,
		sheetName:     sheetName,
	}
}

func (r *ProductBatchRepository) Create(ctx context.Context, batch *domain.ProductBatch) error {
	if batch == nil {
		return fmt.Errorf("product batch is required: %w", errs.ErrBadRequest)
	}

	row := []interface{}{
		batch.ID,
		batch.ProductID,
		batch.BatchNumber,
		batch.ExpiryDate.Format(time.RFC3339),
		strconv.Itoa(batch.QtyAvailable),
		batch.ReceivedDate.Format(time.RFC3339),
		batch.SupplierName,
		batch.CreatedAt.Format(time.RFC3339),
		batch.UpdatedAt.Format(time.RFC3339),
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
		return fmt.Errorf("failed to append product batch row to spreadsheet: %w", err)
	}

	return nil
}

func (r *ProductBatchRepository) FindAll(ctx context.Context, page int, limit int, search string) ([]domain.ProductBatch, int, error) {
	rangeName := fmt.Sprintf("%s!A2:I", r.sheetName)

	resp, err := r.sheetsService.Spreadsheets.Values.
		Get(r.spreadsheetID, rangeName).
		Context(ctx).
		Do()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get product batch rows from spreadsheet: %w", err)
	}

	var batches []domain.ProductBatch
	search = strings.TrimSpace(strings.ToLower(search))

	for _, row := range resp.Values {
		batch := mapRowToProductBatch(row)

		if search == "" || containsProductBatchSearch(batch, search) {
			batches = append(batches, batch)
		}
	}

	total := len(batches)

	start := (page - 1) * limit
	end := start + limit

	if start >= total {
		return []domain.ProductBatch{}, total, nil
	}

	if end > total {
		end = total
	}

	return batches[start:end], total, nil
}

func (r *ProductBatchRepository) FindByID(ctx context.Context, id string) (*domain.ProductBatch, error) {
	rangeName := fmt.Sprintf("%s!A2:I", r.sheetName)

	resp, err := r.sheetsService.Spreadsheets.Values.
		Get(r.spreadsheetID, rangeName).
		Context(ctx).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get product batch rows from spreadsheet: %w", err)
	}

	for _, row := range resp.Values {
		batch := mapRowToProductBatch(row)

		if batch.ID == id {
			return &batch, nil
		}
	}

	return nil, fmt.Errorf("product batch with id %s not found: %w", id, errs.ErrNotFound)
}

func (r *ProductBatchRepository) FindByProductID(ctx context.Context, productID string) ([]domain.ProductBatch, error) {
	rangeName := fmt.Sprintf("%s!A2:I", r.sheetName)

	resp, err := r.sheetsService.Spreadsheets.Values.
		Get(r.spreadsheetID, rangeName).
		Context(ctx).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get product batch rows from spreadsheet: %w", err)
	}

	var result []domain.ProductBatch

	for _, row := range resp.Values {
		batch := mapRowToProductBatch(row)

		if batch.ProductID == productID {
			result = append(result, batch)
		}
	}

	return result, nil
}

func (r *ProductBatchRepository) FindFEFOBatchByProductID(ctx context.Context, productID string) (*domain.ProductBatch, error) {
	rangeName := fmt.Sprintf("%s!A2:I", r.sheetName)

	resp, err := r.sheetsService.Spreadsheets.Values.
		Get(r.spreadsheetID, rangeName).
		Context(ctx).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get product batch rows from spreadsheet: %w", err)
	}

	var selected *domain.ProductBatch

	for _, row := range resp.Values {
		batch := mapRowToProductBatch(row)

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
	rangeName := fmt.Sprintf("%s!A2:I", r.sheetName)

	resp, err := r.sheetsService.Spreadsheets.Values.
		Get(r.spreadsheetID, rangeName).
		Context(ctx).
		Do()
	if err != nil {
		return fmt.Errorf("failed to get product batch rows from spreadsheet: %w", err)
	}

	targetRowIndex := -1

	for i, row := range resp.Values {
		batch := mapRowToProductBatch(row)

		if batch.ID == batchID {
			targetRowIndex = i + 2
			break
		}
	}

	if targetRowIndex == -1 {
		return fmt.Errorf("product batch with id %s not found: %w", batchID, errs.ErrNotFound)
	}

	updateRange := fmt.Sprintf("%s!E%d", r.sheetName, targetRowIndex)

	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{
			{strconv.Itoa(qtyAvailable)},
		},
	}

	_, err = r.sheetsService.Spreadsheets.Values.
		Update(r.spreadsheetID, updateRange, valueRange).
		ValueInputOption("RAW").
		Context(ctx).
		Do()
	if err != nil {
		return fmt.Errorf("failed to update qty available in spreadsheet: %w", err)
	}

	return nil
}

func mapRowToProductBatch(row []interface{}) domain.ProductBatch {
	batch := domain.ProductBatch{}

	if len(row) > 0 {
		batch.ID = toString(row[0])
	}
	if len(row) > 1 {
		batch.ProductID = toString(row[1])
	}
	if len(row) > 2 {
		batch.BatchNumber = toString(row[2])
	}
	if len(row) > 3 {
		batch.ExpiryDate = toTime(row[3])
	}
	if len(row) > 4 {
		batch.QtyAvailable = toInt(row[4])
	}
	if len(row) > 5 {
		batch.ReceivedDate = toTime(row[5])
	}
	if len(row) > 6 {
		batch.SupplierName = toString(row[6])
	}
	if len(row) > 7 {
		batch.CreatedAt = toTime(row[7])
	}
	if len(row) > 8 {
		batch.UpdatedAt = toTime(row[8])
	}

	return batch
}

func containsProductBatchSearch(batch domain.ProductBatch, search string) bool {
	return strings.Contains(strings.ToLower(batch.ID), search) ||
		strings.Contains(strings.ToLower(batch.ProductID), search) ||
		strings.Contains(strings.ToLower(batch.BatchNumber), search) ||
		strings.Contains(strings.ToLower(batch.SupplierName), search)
}
