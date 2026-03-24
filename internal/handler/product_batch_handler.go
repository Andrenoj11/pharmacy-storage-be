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

type ProductBatchHandler struct {
	productBatchService *service.ProductBatchService
}

func NewProductBatchHandler(productBatchService *service.ProductBatchService) *ProductBatchHandler {
	return &ProductBatchHandler{
		productBatchService: productBatchService,
	}
}

func (h *ProductBatchHandler) CreateProductBatch(c *gin.Context) {
	var input domain.ProductBatch

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.productBatchService.CreateProductBatch(c.Request.Context(), &input)
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
		"message": "product batch created successfully",
		"data":    input,
	})
}

func (h *ProductBatchHandler) GetAllProductBatches(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := strings.TrimSpace(c.Query("search"))

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	batches, total, err := h.productBatchService.GetAllProductBatches(c.Request.Context(), page, limit, search)
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
		"message": "product batches fetched successfully",
		"data":    batches,
		"meta": gin.H{
			"page":        page,
			"limit":       limit,
			"search":      search,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

func (h *ProductBatchHandler) GetProductBatchByID(c *gin.Context) {
	id := c.Param("id")

	batch, err := h.productBatchService.GetProductBatchByID(c.Request.Context(), id)
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
		"message": "product batch fetched successfully",
		"data":    batch,
	})
}

func (h *ProductBatchHandler) GetProductBatchesByProductID(c *gin.Context) {
	productID := c.Param("id")

	batches, err := h.productBatchService.GetProductBatchesByProductID(c.Request.Context(), productID)
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
		"message": "product batches fetched successfully",
		"data":    batches,
	})
}

func (h *ProductBatchHandler) GetFEFOBatchByProductID(c *gin.Context) {
	productID := c.Param("id")

	batch, err := h.productBatchService.GetFEFOBatchByProductID(c.Request.Context(), productID)
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
		"message": "FEFO batch fetched successfully",
		"data":    batch,
	})
}
