package main

import (
	"log"

	"pharmacy-storage-be/internal/app"
	"pharmacy-storage-be/internal/config"
	"pharmacy-storage-be/internal/handler"
	"pharmacy-storage-be/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	productRepo, err := app.NewProductRepository(cfg)
	if err != nil {
		log.Fatal(err)
	}

	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	productBatchRepo, err := app.NewProductBatchRepository(cfg)
	if err != nil {
		log.Fatal(err)
	}

	productBatchService := service.NewProductBatchService(productBatchRepo, productRepo)
	productBatchHandler := handler.NewProductBatchHandler(productBatchService)

	inventoryMovementRepo, err := app.NewInventoryMovementRepository(cfg)
	if err != nil {
		log.Fatal(err)
	}

	inventoryMovementService := service.NewInventoryMovementService(
		inventoryMovementRepo,
		productRepo,
		productBatchRepo,
	)
	inventoryMovementHandler := handler.NewInventoryMovementHandler(inventoryMovementService)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":                    "server is running",
			"storage_provider":           cfg.StorageProvider,
			"product_service":            productService != nil,
			"product_batch_service":      productBatchService != nil,
			"inventory_movement_service": inventoryMovementService != nil,
		})
	})

	r.POST("/products", productHandler.CreateProduct)
	r.GET("/products", productHandler.GetAllProducts)
	r.GET("/products/:id", productHandler.GetProductByID)
	r.GET("/products/:id/batches", productBatchHandler.GetProductBatchesByProductID)
	r.GET("/products/:id/fefo-batch", productBatchHandler.GetFEFOBatchByProductID)

	r.POST("/product-batches", productBatchHandler.CreateProductBatch)
	r.GET("/product-batches", productBatchHandler.GetAllProductBatches)
	r.GET("/product-batches/:id", productBatchHandler.GetProductBatchByID)

	r.POST("/stock-in", inventoryMovementHandler.CreateStockIn)
	r.POST("/stock-out", inventoryMovementHandler.CreateStockOut)

	r.POST("/inventory-movements", inventoryMovementHandler.CreateInventoryMovement)
	r.GET("/inventory-movements", inventoryMovementHandler.GetAllInventoryMovements)
	r.GET("/inventory-movements/:id", inventoryMovementHandler.GetInventoryMovementByID)
	r.GET("/product-batches/:id/movements", inventoryMovementHandler.GetInventoryMovementsByBatchID)

	log.Printf("server running on port %s with provider %s", cfg.AppPort, cfg.StorageProvider)

	err = r.Run(":" + cfg.AppPort)
	if err != nil {
		log.Fatal(err)
	}
}
