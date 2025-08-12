package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ru *ResponseUtil) ParsingRequest(ctx *gin.Context, data any) bool {
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ru.Log.Errorf("Request parameter parsing error, Error: %v", err)
		ru.BadRequest(ctx, "Request parameter parsing error")
		return false
	}
	return true
}

func (ru *ResponseUtil) Build(ctx *gin.Context, code int, message string, data any) {
	ctx.JSON(code, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

func (ru *ResponseUtil) Unauthorized(ctx *gin.Context, message string) {
	ru.Build(ctx, http.StatusUnauthorized, fmt.Sprintf("Error: %s", message), nil)
}

func (ru *ResponseUtil) BadRequest(ctx *gin.Context, message string) {
	ru.Build(ctx, http.StatusBadRequest, fmt.Sprintf("Error: %s", message), nil)
}

func (ru *ResponseUtil) Success(ctx *gin.Context, message any, data any) {
	ru.Build(ctx, http.StatusOK, fmt.Sprintf("Success: %v", message), data)
}

func (ru *ResponseUtil) Error(ctx *gin.Context, message any) {
	ru.Build(ctx, http.StatusInternalServerError, fmt.Sprintf("Error: %v", message), nil)
}
