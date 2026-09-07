package middlewares

import (
	"fmt"
	"github.com/danielgtaylor/huma/v2"
)

// HeadersMiddleware Sets common HTTP headers for responses.
func HeadersMiddleware(api huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {

		ctx.SetHeader("X-Frame-Options", "DENY")
		ctx.SetHeader("X-Content-Type-Options", "nosniff")
		ctx.SetHeader("X-Xss-Protection", "1; mode=block")
		ctx.SetHeader("Content-Security-Policy", "default-src 'self'")
		ctx.SetHeader("Referrer-Policy", "strict-origin-when-cross-origin")
		ctx.SetHeader("Strict-Transport-Security", fmt.Sprintf("max-age=%d; %s", 31536000, "includeSubDomains"))

		next(ctx)
	}
}
