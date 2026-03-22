package domain

import "time"

type Product struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Unit        string    `json:"unit"`
	MinStock    int       `json:"min_stock"`
	StorageType string    `json:"storage_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
