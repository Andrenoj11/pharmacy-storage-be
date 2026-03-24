package handler

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"

	"pharmacy-storage-be/internal/domain"
	"pharmacy-storage-be/internal/errs"
	"pharmacy-storage-be/internal/service"

	"github.com/gin-gonic/gin"
)

type InventoryMovementHandler struct {
	inventoryMovementService *service.InventoryMovementService
}

func NewInventoryMovementHandler(inventoryMovementService *service.InventoryMovementService) *InventoryMovementHandler {
	return &InventoryMovementHandler{
		inventoryMovementService: inventoryMovementService,
	}
}

func (h *InventoryMovementHandler) CreateInventoryMovement(c *gin.Context) {
	var input domain.InventoryMovement

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.inventoryMovementService.CreateInventoryMovement(c.Request.Context(), &input)
	if err != nil {
		statusCode := http.StatusInternalServerError

		switch {
		case errors.Is(err, errs.ErrBadRequest):
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrNotFound):
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "inventory movement created successfully",
		"data":    input,
	})
}

func (h *InventoryMovementHandler) GetAllInventoryMovements(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := strings.TrimSpace(c.Query("search"))

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	movements, total, err := h.inventoryMovementService.GetAllInventoryMovements(c.Request.Context(), page, limit, search)
	if err != nil {
		statusCode := http.StatusInternalServerError

		if errors.Is(err, errs.ErrBadRequest) {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "inventory movements fetched successfully",
		"data":    movements,
		"meta": gin.H{
			"page":        page,
			"limit":       limit,
			"search":      search,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

func (h *InventoryMovementHandler) GetInventoryMovementByID(c *gin.Context) {
	id := c.Param("id")

	movement, err := h.inventoryMovementService.GetInventoryMovementByID(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError

		switch {
		case errors.Is(err, errs.ErrBadRequest):
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrNotFound):
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "inventory movement fetched successfully",
		"data":    movement,
	})
}

func (h *InventoryMovementHandler) GetInventoryMovementsByBatchID(c *gin.Context) {
	batchID := c.Param("id")

	movements, err := h.inventoryMovementService.GetInventoryMovementsByBatchID(c.Request.Context(), batchID)
	if err != nil {
		statusCode := http.StatusInternalServerError

		switch {
		case errors.Is(err, errs.ErrBadRequest):
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrNotFound):
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "inventory movements fetched successfully",
		"data":    movements,
	})
}
