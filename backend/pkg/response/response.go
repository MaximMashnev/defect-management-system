package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func Error(c *gin.Context, statusCode int, message string, err error) {
	resp := ErrorResponse{Message: message}
	if err != nil {
		resp.Error = err.Error()
	}
	c.JSON(statusCode, resp)
}

func BadRequest(c *gin.Context, message string, err error) {
	Error(c, http.StatusBadRequest, message, err)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message, nil)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message, nil)
}

func Internal(c *gin.Context, message string, err error) {
	Error(c, http.StatusInternalServerError, message, err)
}
