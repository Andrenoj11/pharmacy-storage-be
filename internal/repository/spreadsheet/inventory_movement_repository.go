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

type InventoryMovementRepository struct {
	sheetsService *sheets.Service
	spreadsheetID string
	sheetName     string
}

func NewInventoryMovementRepository(sheetsService *sheets.Service, spreadsheetID string, sheetName string) *InventoryMovementRepository {
	return &InventoryMovementRepository{
		sheetsService: sheetsService,
		spreadsheetID: spreadsheetID,
		sheetName:     sheetName,
	}
}

func (r *InventoryMovementRepository) Create(ctx context.Context, movement *domain.InventoryMovement) error {
	if movement == nil {
		return fmt.Errorf("inventory movement is required: %w", errs.ErrBadRequest)
	}

	row := []interface{}{
		movement.ID,
		movement.ProductID,
		movement.ProductBatchID,
		movement.MovementType,
		strconv.Itoa(movement.Qty),
		movement.MovementDate.Format(time.RFC3339),
		movement.ReferenceNo,
		movement.Note,
		movement.CreatedBy,
		movement.CreatedAt.Format(time.RFC3339),
		movement.UpdatedAt.Format(time.RFC3339),
	}

	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{row},
	}

	rangeName := fmt.Sprintf("%s!A:K", r.sheetName)

	_, err := r.sheetsService.Spreadsheets.Values.
		Append(r.spreadsheetID, rangeName, valueRange).
		ValueInputOption("RAW").
		Context(ctx).
		Do()
	if err != nil {
		return fmt.Errorf("failed to append inventory movement row to spreadsheet: %w", err)
	}

	return nil
}

func (r *InventoryMovementRepository) FindAll(ctx context.Context, page int, limit int, search string) ([]domain.InventoryMovement, int, error) {
	rangeName := fmt.Sprintf("%s!A2:K", r.sheetName)

	resp, err := r.sheetsService.Spreadsheets.Values.
		Get(r.spreadsheetID, rangeName).
		Context(ctx).
		Do()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get inventory movement rows from spreadsheet: %w", err)
	}

	var movements []domain.InventoryMovement
	search = strings.TrimSpace(strings.ToLower(search))

	for _, row := range resp.Values {
		movement := mapRowToInventoryMovement(row)

		if search == "" || containsInventoryMovementSearch(movement, search) {
			movements = append(movements, movement)
		}
	}

	total := len(movements)

	start := (page - 1) * limit
	end := start + limit

	if start >= total {
		return []domain.InventoryMovement{}, total, nil
	}

	if end > total {
		end = total
	}

	return movements[start:end], total, nil
}

func (r *InventoryMovementRepository) FindByID(ctx context.Context, id string) (*domain.InventoryMovement, error) {
	rangeName := fmt.Sprintf("%s!A2:K", r.sheetName)

	resp, err := r.sheetsService.Spreadsheets.Values.
		Get(r.spreadsheetID, rangeName).
		Context(ctx).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory movement rows from spreadsheet: %w", err)
	}

	for _, row := range resp.Values {
		movement := mapRowToInventoryMovement(row)

		if movement.ID == id {
			return &movement, nil
		}
	}

	return nil, fmt.Errorf("inventory movement with id %s not found: %w", id, errs.ErrNotFound)
}

func (r *InventoryMovementRepository) FindByBatchID(ctx context.Context, batchID string) ([]domain.InventoryMovement, error) {
	rangeName := fmt.Sprintf("%s!A2:K", r.sheetName)

	resp, err := r.sheetsService.Spreadsheets.Values.
		Get(r.spreadsheetID, rangeName).
		Context(ctx).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory movement rows from spreadsheet: %w", err)
	}

	var result []domain.InventoryMovement

	for _, row := range resp.Values {
		movement := mapRowToInventoryMovement(row)

		if movement.ProductBatchID == batchID {
			result = append(result, movement)
		}
	}

	return result, nil
}

func mapRowToInventoryMovement(row []interface{}) domain.InventoryMovement {
	movement := domain.InventoryMovement{}

	if len(row) > 0 {
		movement.ID = toString(row[0])
	}
	if len(row) > 1 {
		movement.ProductID = toString(row[1])
	}
	if len(row) > 2 {
		movement.ProductBatchID = toString(row[2])
	}
	if len(row) > 3 {
		movement.MovementType = toString(row[3])
	}
	if len(row) > 4 {
		movement.Qty = toInt(row[4])
	}
	if len(row) > 5 {
		movement.MovementDate = toTime(row[5])
	}
	if len(row) > 6 {
		movement.ReferenceNo = toString(row[6])
	}
	if len(row) > 7 {
		movement.Note = toString(row[7])
	}
	if len(row) > 8 {
		movement.CreatedBy = toString(row[8])
	}
	if len(row) > 9 {
		movement.CreatedAt = toTime(row[9])
	}
	if len(row) > 10 {
		movement.UpdatedAt = toTime(row[10])
	}

	return movement
}

func containsInventoryMovementSearch(movement domain.InventoryMovement, search string) bool {
	return strings.Contains(strings.ToLower(movement.ID), search) ||
		strings.Contains(strings.ToLower(movement.ProductID), search) ||
		strings.Contains(strings.ToLower(movement.ProductBatchID), search) ||
		strings.Contains(strings.ToLower(movement.MovementType), search) ||
		strings.Contains(strings.ToLower(movement.ReferenceNo), search) ||
		strings.Contains(strings.ToLower(movement.Note), search) ||
		strings.Contains(strings.ToLower(movement.CreatedBy), search)
}
