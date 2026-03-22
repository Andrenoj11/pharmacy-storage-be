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

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":               "server is running",
			"storage_provider":      cfg.StorageProvider,
			"product_service":       productService != nil,
			"google_products_sheet": cfg.GoogleProductsSheet,
		})
	})

	r.POST("/products", productHandler.CreateProduct)
	r.GET("/products", productHandler.GetAllProducts)
	r.GET("/products/:id", productHandler.GetProductByID)

	log.Printf("server running on port %s with provider %s", cfg.AppPort, cfg.StorageProvider)

	err = r.Run(":" + cfg.AppPort)
	if err != nil {
		log.Fatal(err)
	}
}
