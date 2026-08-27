package supplier

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
	suppliers := rg.Group("/suppliers")
	{
		suppliers.POST("", h.Create)
		suppliers.GET("", h.List)
		suppliers.GET("/:id", h.GetByID)
		suppliers.GET("/name/:name", h.GetByName)
		suppliers.PUT("/:id", h.Update)
		suppliers.DELETE("/:id", h.Delete)
	}
}

func (h *Handler) Create(c *gin.Context) {
	var input CreateSupplierInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondValidationError(c, err.Error())
		return
	}

	sup, err := h.svc.Create(c.Request.Context(), input)
	if err != nil {
		switch err {
		case ErrNameExists:
			httpx.RespondError(c, http.StatusConflict, "NAME_EXISTS", "supplier with this name already exists")
		case ErrInvalidInput:
			httpx.RespondValidationError(c, "name is required")
		default:
			httpx.RespondError(c, http.StatusInternalServerError, "CREATE_FAILED", "failed to create supplier")
		}
		return
	}

	httpx.RespondOK(c, http.StatusCreated, sup)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := parseObjectID(c.Param("id"))
	if err != nil {
		httpx.RespondValidationError(c, "invalid supplier ID")
		return
	}

	sup, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			httpx.RespondError(c, http.StatusNotFound, "NOT_FOUND", "supplier not found")
			return
		}
		httpx.RespondError(c, http.StatusInternalServerError, "GET_FAILED", "failed to get supplier")
		return
	}

	httpx.RespondOK(c, http.StatusOK, sup)
}

func (h *Handler) GetByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		httpx.RespondValidationError(c, "name is required")
		return
	}

	sup, err := h.svc.GetByName(c.Request.Context(), name)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			httpx.RespondError(c, http.StatusNotFound, "NOT_FOUND", "supplier not found")
			return
		}
		httpx.RespondError(c, http.StatusInternalServerError, "GET_FAILED", "failed to get supplier")
		return
	}

	httpx.RespondOK(c, http.StatusOK, sup)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := parseObjectID(c.Param("id"))
	if err != nil {
		httpx.RespondValidationError(c, "invalid supplier ID")
		return
	}

	var input UpdateSupplierInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondValidationError(c, err.Error())
		return
	}

	sup, err := h.svc.Update(c.Request.Context(), id, input)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			httpx.RespondError(c, http.StatusNotFound, "NOT_FOUND", "supplier not found")
			return
		}
		if errors.Is(err, ErrNameExists) {
			httpx.RespondError(c, http.StatusConflict, "NAME_EXISTS", "supplier with this name already exists")
			return
		}
		httpx.RespondError(c, http.StatusInternalServerError, "UPDATE_FAILED", "failed to update supplier")
		return
	}

	httpx.RespondOK(c, http.StatusOK, sup)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := parseObjectID(c.Param("id"))
	if err != nil {
		httpx.RespondValidationError(c, "invalid supplier ID")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			httpx.RespondError(c, http.StatusNotFound, "NOT_FOUND", "supplier not found")
			return
		}
		httpx.RespondError(c, http.StatusInternalServerError, "DELETE_FAILED", "failed to delete supplier")
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) List(c *gin.Context) {
	var filter SupplierFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		httpx.RespondValidationError(c, err.Error())
		return
	}

	suppliers, total, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "LIST_FAILED", "failed to list suppliers")
		return
	}

	httpx.RespondOK(c, http.StatusOK, gin.H{
		"suppliers": suppliers,
		"total":     total,
		"limit":     filter.Limit,
		"offset":    filter.Offset,
	})
}

func parseObjectID(s string) (bson.ObjectID, error) {
	return bson.ObjectIDFromHex(s)
}
