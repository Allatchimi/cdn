package helpers

import (
	"github.com/danielgtaylor/huma/v2"
)

// ExtractApiKeyHeader Retrieves the api key from the current request context.
func ExtractApiKeyHeader(ctx *huma.Context) string {
	return (*ctx).Header("X-API-Key")
}
