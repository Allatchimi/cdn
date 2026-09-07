package config

import (
	"cdn/common/helpers"

	"github.com/maypok86/otter"
	"go.uber.org/zap"
)

var OtterCache otter.CacheWithVariableTTL[string, string]

// Setup otter cache.
func SetupOtterCache() error {
	var err error
	OtterCache, err = otter.MustBuilder[string, string](10000).
		CollectStats().
		Cost(func(key string, value string) uint32 {
			return 1
		}).
		WithVariableTTL().
		Build()
	if err != nil {
		helpers.Logger.Warn(
			"Failed to initialize otter cache!",
			zap.String("Error: ", err.Error()),
		)
	}

	return err
}
