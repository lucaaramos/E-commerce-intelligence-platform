package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorBody struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type DataResponse struct {
	Data any `json:"data"`
}

// RespondOK writes a 2xx response with the payload under a "data" key so the
// API has a stable envelope for all success responses.
func RespondOK(c *gin.Context, status int, data any) {
	c.JSON(status, DataResponse{Data: data})
}

// RespondError writes an error envelope with a stable machine-readable code.
func RespondError(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, ErrorResponse{
		Error: ErrorBody{
			Code:      code,
			Message:   message,
			RequestID: GetRequestID(c),
		},
	})
}

// RespondValidationError is a shortcut for validation failures.
func RespondValidationError(c *gin.Context, message string) {
	RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", message)
}
