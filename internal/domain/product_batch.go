package domain

import "time"

type ProductBatch struct {
	ID           string    `json:"id"`
	ProductID    string    `json:"product_id"`
	BatchNumber  string    `json:"batch_number"`
	ExpiryDate   time.Time `json:"expiry_date"`
	QtyAvailable int       `json:"qty_available"`
	ReceivedDate time.Time `json:"received_date"`
	SupplierName string    `json:"supplier_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
