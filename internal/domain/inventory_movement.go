package domain

import "time"

type InventoryMovement struct {
	ID             string    `json:"id"`
	ProductID      string    `json:"product_id"`
	ProductBatchID string    `json:"product_batch_id"`
	MovementType   string    `json:"movement_type"`
	Qty            int       `json:"qty"`
	MovementDate   time.Time `json:"movement_date"`
	ReferenceNo    string    `json:"reference_no"`
	Note           string    `json:"note"`
	CreatedBy      string    `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
