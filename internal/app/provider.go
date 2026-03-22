package app

import (
	"fmt"

	"pharmacy-storage-be/internal/config"
	"pharmacy-storage-be/internal/repository"
	memoryrepo "pharmacy-storage-be/internal/repository/memory"
)

func NewProductRepository(cfg *config.Config) (repository.ProductRepository, error) {
	switch cfg.StorageProvider {
	case "memory":
		return memoryrepo.NewProductRepository(), nil
	case "spreadsheet":
		return nil, fmt.Errorf("spreadsheet repository is not implemented yet")
	case "postgres":
		return nil, fmt.Errorf("postgres repository is not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported storage provider: %s", cfg.StorageProvider)
	}
}
