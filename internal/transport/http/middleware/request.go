package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDKey = "X-Request-Id"

func RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.GetHeader(RequestIDKey)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx.Set(RequestIDKey, requestID)
		ctx.Writer.Header().Set(RequestIDKey, requestID)
		ctx.Next()
	}
}
