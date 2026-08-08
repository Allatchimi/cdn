<<<<<<< HEAD
package helpers

import (
	"github.com/danielgtaylor/huma/v2"
)

// ExtractApiKeyHeader Retrieves the api key from the current request context.
func ExtractApiKeyHeader(ctx *huma.Context) string {
	return (*ctx).Header("X-API-Key")
}
=======
package helpers

import (
	"github.com/danielgtaylor/huma/v2"
)

// ExtractApiKeyHeader Retrieves the api key from the current request context.
func ExtractApiKeyHeader(ctx *huma.Context) string {
	return (*ctx).Header("X-API-Key")
}
>>>>>>> 22022f0081c75477042da66cd81443ff4401ca37
