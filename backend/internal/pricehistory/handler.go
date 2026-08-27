package pricehistory

import (
	"net/http"
	"time"

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
	history := rg.Group("/price-history")
	{
		history.GET("/product/:productId", h.GetByProductID)
		history.GET("/supplier-product/:supplierProductId", h.GetBySupplierProduct)
		history.GET("/latest/:supplierProductId", h.GetLatest)
		history.GET("/change/:supplierProductId", h.GetPriceChange)
	}
}

func (h *Handler) GetByProductID(c *gin.Context) {
	productID, err := bson.ObjectIDFromHex(c.Param("productId"))
	if err != nil {
		httpx.RespondValidationError(c, "invalid product ID")
		return
	}

	limit := 50
	history, err := h.repo.List(c.Request.Context(), PriceHistoryFilter{
		ProductID: productID,
		Limit:     limit,
	})
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "LIST_FAILED", "failed to get price history")
		return
	}

	httpx.RespondOK(c, http.StatusOK, history)
}

func (h *Handler) GetBySupplierProduct(c *gin.Context) {
	spID, err := bson.ObjectIDFromHex(c.Param("supplierProductId"))
	if err != nil {
		httpx.RespondValidationError(c, "invalid supplier product ID")
		return
	}

	var since *time.Time
	if d := c.Query("since"); d != "" {
		t, err := time.Parse(time.RFC3339, d)
		if err == nil {
			since = &t
		}
	}

	history, err := h.repo.List(c.Request.Context(), PriceHistoryFilter{
		SupplierProductID: spID,
		Since:             since,
		Limit:             100,
	})
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "LIST_FAILED", "failed to get price history")
		return
	}

	httpx.RespondOK(c, http.StatusOK, history)
}

func (h *Handler) GetLatest(c *gin.Context) {
	spID, err := bson.ObjectIDFromHex(c.Param("supplierProductId"))
	if err != nil {
		httpx.RespondValidationError(c, "invalid supplier product ID")
		return
	}

	latest, err := h.repo.GetLatest(c.Request.Context(), spID)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "GET_FAILED", "failed to get latest price")
		return
	}

	httpx.RespondOK(c, http.StatusOK, latest)
}

func (h *Handler) GetPriceChange(c *gin.Context) {
	spID, err := bson.ObjectIDFromHex(c.Param("supplierProductId"))
	if err != nil {
		httpx.RespondValidationError(c, "invalid supplier product ID")
		return
	}

	change, err := h.repo.GetPriceChange(c.Request.Context(), spID)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "GET_FAILED", "failed to get price change")
		return
	}

	httpx.RespondOK(c, http.StatusOK, change)
}
