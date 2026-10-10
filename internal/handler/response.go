package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

type statusResponse struct {
	Status string `json:"status"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	slog.Error("request failed", "status", statusCode, "error", message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
