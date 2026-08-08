<<<<<<< HEAD
package helpers

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

// Enable the logger to print beautiful log messages.
func EnableLogger() {
	Logger, _ = zap.NewDevelopment()
	defer func(Logger *zap.Logger) {
		err := Logger.Sync()
		if err != nil {
			return
		}
	}(Logger)
}
=======
package helpers

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

// Enable the logger to print beautiful log messages.
func EnableLogger() {
	Logger, _ = zap.NewDevelopment()
	defer func(Logger *zap.Logger) {
		err := Logger.Sync()
		if err != nil {
			return
		}
	}(Logger)
}
>>>>>>> 22022f0081c75477042da66cd81443ff4401ca37
