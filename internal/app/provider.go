package app

import (
	"context"
	"fmt"

	"pharmacy-storage-be/internal/config"
	googleclient "pharmacy-storage-be/internal/google"
	"pharmacy-storage-be/internal/repository"
	memoryrepo "pharmacy-storage-be/internal/repository/memory"
	spreadsheetrepo "pharmacy-storage-be/internal/repository/spreadsheet"
)

func NewProductRepository(cfg *config.Config) (repository.ProductRepository, error) {
	switch cfg.StorageProvider {
	case "memory":
		return memoryrepo.NewProductRepository(), nil

	case "spreadsheet":
		sheetsService, err := googleclient.NewSheetsService(context.Background(), cfg.GoogleCredentialsFile)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize google sheets service: %w", err)
		}

		if cfg.GoogleSpreadsheetID == "" {
			return nil, fmt.Errorf("google spreadsheet id is required for spreadsheet provider")
		}

		return spreadsheetrepo.NewProductRepository(
			sheetsService,
			cfg.GoogleSpreadsheetID,
			cfg.GoogleProductsSheet,
		), nil

	case "postgres":
		return nil, fmt.Errorf("postgres repository is not implemented yet")

	default:
		return nil, fmt.Errorf("unsupported storage provider: %s", cfg.StorageProvider)
	}
}
