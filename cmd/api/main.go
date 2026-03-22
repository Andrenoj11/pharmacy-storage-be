package main

import (
	"log"

	"pharmacy-storage-be/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":          "server is running well!",
			"storage_provider": cfg.StorageProvider,
		})
	})

	log.Printf("server running on port %s with provider %s", cfg.AppPort, cfg.StorageProvider)

	err = r.Run(":" + cfg.AppPort)
	if err != nil {
		log.Fatal(err)
	}
}
