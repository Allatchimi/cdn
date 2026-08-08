<<<<<<< HEAD
package di

import (
	"cdn/cmd/api"
	"cdn/services/document"
	"cdn/services/image"
)

// Inject all dependencies
func InjectDependencies() {
	// Image
	api.Controllers.ImageController = image.NewController(
		image.NewService(),
	)
	// Document
	api.Controllers.DocumentController = document.NewController(
		document.NewService(),
	)
}
=======
package di

import (
	"cdn/cmd/api"
	"cdn/services/document"
	"cdn/services/image"
)

// Inject all dependencies
func InjectDependencies() {
	// Image
	api.Controllers.ImageController = image.NewController(
		image.NewService(),
	)
	// Document
	api.Controllers.DocumentController = document.NewController(
		document.NewService(),
	)
}
>>>>>>> 22022f0081c75477042da66cd81443ff4401ca37
