package supplierproduct

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"

	"ecommerce-intelligence-platform/backend/pkg/httpx"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	sp := rg.Group("/supplier-products")
	{
		sp.POST("", h.Create)
		sp.GET("", h.List)
		sp.GET("/:id", h.GetByID)
		sp.GET("/product/:productId", h.GetByProductID)
		sp.GET("/cheapest/:productId", h.GetCheapest)
		sp.PUT("/:id", h.Update)
		sp.DELETE("/:id", h.Delete)
	}
}

func (h *Handler) Create(c *gin.Context) {
	var input CreateSupplierProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondValidationError(c, err.Error())
		return
	}

	sp := &SupplierProduct{
		ProductID:    input.ProductID,
		SupplierID:   input.SupplierID,
		SupplierSKU:  input.SupplierSKU,
		Cost:         input.Cost,
		ShippingCost: input.ShippingCost,
		LeadTimeDays: input.LeadTimeDays,
		MinOrderQty:  input.MinOrderQty,
		URL:          input.URL,
	}

	if err := h.repo.Create(c.Request.Context(), sp); err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "CREATE_FAILED", "failed to create supplier product")
		return
	}

	httpx.RespondOK(c, http.StatusCreated, sp)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		httpx.RespondValidationError(c, "invalid ID")
		return
	}

	sp, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			httpx.RespondError(c, http.StatusNotFound, "NOT_FOUND", "supplier product not found")
			return
		}
		httpx.RespondError(c, http.StatusInternalServerError, "GET_FAILED", "failed to get supplier product")
		return
	}

	httpx.RespondOK(c, http.StatusOK, sp)
}

func (h *Handler) GetByProductID(c *gin.Context) {
	productID, err := bson.ObjectIDFromHex(c.Param("productId"))
	if err != nil {
		httpx.RespondValidationError(c, "invalid product ID")
		return
	}

	results, err := h.repo.GetByProductID(c.Request.Context(), productID)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "LIST_FAILED", "failed to list supplier products")
		return
	}

	httpx.RespondOK(c, http.StatusOK, results)
}

func (h *Handler) GetCheapest(c *gin.Context) {
	productID, err := bson.ObjectIDFromHex(c.Param("productId"))
	if err != nil {
		httpx.RespondValidationError(c, "invalid product ID")
		return
	}

	sp, err := h.repo.GetCheapest(c.Request.Context(), productID)
	if err != nil {
		if errors.Is(err, ErrNoCheapestFound) {
			httpx.RespondError(c, http.StatusNotFound, "NOT_FOUND", "no active supplier found for this product")
			return
		}
		httpx.RespondError(c, http.StatusInternalServerError, "GET_FAILED", "failed to get cheapest supplier")
		return
	}

	httpx.RespondOK(c, http.StatusOK, sp)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		httpx.RespondValidationError(c, "invalid ID")
		return
	}

	sp, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		httpx.RespondError(c, http.StatusNotFound, "NOT_FOUND", "supplier product not found")
		return
	}

	var input UpdateSupplierProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondValidationError(c, err.Error())
		return
	}

	if input.SupplierSKU != nil {
		sp.SupplierSKU = *input.SupplierSKU
	}
	if input.Cost != nil {
		sp.Cost = *input.Cost
	}
	if input.ShippingCost != nil {
		sp.ShippingCost = *input.ShippingCost
	}
	if input.LeadTimeDays != nil {
		sp.LeadTimeDays = *input.LeadTimeDays
	}
	if input.MinOrderQty != nil {
		sp.MinOrderQty = *input.MinOrderQty
	}
	if input.URL != nil {
		sp.URL = *input.URL
	}
	if input.IsActive != nil {
		sp.IsActive = *input.IsActive
	}

	if err := h.repo.Update(c.Request.Context(), sp); err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "UPDATE_FAILED", "failed to update supplier product")
		return
	}

	httpx.RespondOK(c, http.StatusOK, sp)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		httpx.RespondValidationError(c, "invalid ID")
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			httpx.RespondError(c, http.StatusNotFound, "NOT_FOUND", "supplier product not found")
			return
		}
		httpx.RespondError(c, http.StatusInternalServerError, "DELETE_FAILED", "failed to delete supplier product")
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) List(c *gin.Context) {
	var filter SupplierProductFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		httpx.RespondValidationError(c, err.Error())
		return
	}

	results, err := h.repo.List(c.Request.Context(), filter)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "LIST_FAILED", "failed to list supplier products")
		return
	}

	httpx.RespondOK(c, http.StatusOK, results)
}
