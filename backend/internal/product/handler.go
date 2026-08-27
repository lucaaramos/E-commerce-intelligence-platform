package product

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"

	"ecommerce-intelligence-platform/backend/pkg/httpx"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	products := rg.Group("/products")
	{
		products.POST("", h.Create)
		products.GET("", h.List)
		products.GET("/:id", h.GetByID)
		products.GET("/sku/:sku", h.GetBySKU)
		products.PUT("/:id", h.Update)
		products.DELETE("/:id", h.Delete)
	}
}

func (h *Handler) Create(c *gin.Context) {
	var input CreateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondValidationError(c, err.Error())
		return
	}

	p, err := h.svc.Create(c.Request.Context(), input)
	if err != nil {
		switch err {
		case ErrSKUExists:
			httpx.RespondError(c, http.StatusConflict, "SKU_EXISTS", "product with this SKU already exists")
		case ErrInvalidInput:
			httpx.RespondValidationError(c, "invalid input: SKU, name, and selling price are required")
		default:
			httpx.RespondError(c, http.StatusInternalServerError, "CREATE_FAILED", "failed to create product")
		}
		return
	}

	httpx.RespondOK(c, http.StatusCreated, p)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := parseObjectID(c.Param("id"))
	if err != nil {
		httpx.RespondValidationError(c, "invalid product ID")
		return
	}

	p, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			httpx.RespondError(c, http.StatusNotFound, "NOT_FOUND", "product not found")
			return
		}
		httpx.RespondError(c, http.StatusInternalServerError, "GET_FAILED", "failed to get product")
		return
	}

	httpx.RespondOK(c, http.StatusOK, p)
}

func (h *Handler) GetBySKU(c *gin.Context) {
	sku := c.Param("sku")
	if sku == "" {
		httpx.RespondValidationError(c, "SKU is required")
		return
	}

	p, err := h.svc.GetBySKU(c.Request.Context(), sku)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			httpx.RespondError(c, http.StatusNotFound, "NOT_FOUND", "product not found")
			return
		}
		httpx.RespondError(c, http.StatusInternalServerError, "GET_FAILED", "failed to get product")
		return
	}

	httpx.RespondOK(c, http.StatusOK, p)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := parseObjectID(c.Param("id"))
	if err != nil {
		httpx.RespondValidationError(c, "invalid product ID")
		return
	}

	var input UpdateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondValidationError(c, err.Error())
		return
	}

	p, err := h.svc.Update(c.Request.Context(), id, input)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			httpx.RespondError(c, http.StatusNotFound, "NOT_FOUND", "product not found")
			return
		}
		if errors.Is(err, ErrInvalidInput) {
			httpx.RespondValidationError(c, "selling price must be greater than 0")
			return
		}
		httpx.RespondError(c, http.StatusInternalServerError, "UPDATE_FAILED", "failed to update product")
		return
	}

	httpx.RespondOK(c, http.StatusOK, p)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := parseObjectID(c.Param("id"))
	if err != nil {
		httpx.RespondValidationError(c, "invalid product ID")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			httpx.RespondError(c, http.StatusNotFound, "NOT_FOUND", "product not found")
			return
		}
		httpx.RespondError(c, http.StatusInternalServerError, "DELETE_FAILED", "failed to delete product")
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) List(c *gin.Context) {
	var filter ProductFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		httpx.RespondValidationError(c, err.Error())
		return
	}

	products, total, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "LIST_FAILED", "failed to list products")
		return
	}

	httpx.RespondOK(c, http.StatusOK, gin.H{
		"products": products,
		"total":    total,
		"limit":    filter.Limit,
		"offset":   filter.Offset,
	})
}

func parseObjectID(s string) (bson.ObjectID, error) {
	return bson.ObjectIDFromHex(s)
}