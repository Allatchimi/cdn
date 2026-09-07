package middlewares

import (
	"context"

	"github.com/gin-gonic/gin"
)

type ctxKey string

const GIN_CONTEXT_KEY ctxKey = "gin-context"

func GinContextRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Inject *gin.Context to standard context
		ctx := context.WithValue(c.Request.Context(), GIN_CONTEXT_KEY, c)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
