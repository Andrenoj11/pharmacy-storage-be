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

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var input domain.Product

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.productService.CreateProduct(c.Request.Context(), &input)
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

	c.JSON(http.StatusCreated, gin.H{
		"message": "product created successfully",
		"data":    input,
	})
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := strings.TrimSpace(c.Query("search"))

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	products, total, err := h.productService.GetAllProducts(c.Request.Context(), page, limit, search)
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
		"message": "products fetched successfully",
		"data":    products,
		"meta": gin.H{
			"page":        page,
			"limit":       limit,
			"search":      search,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")

	product, err := h.productService.GetProductByID(c.Request.Context(), id)
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
		"message": "product fetched successfully",
		"data":    product,
	})
}
