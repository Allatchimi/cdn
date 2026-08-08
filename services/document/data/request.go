<<<<<<< HEAD
package data

import "github.com/danielgtaylor/huma/v2"

type DocumentQuery struct {
}

type DocumentData struct {
	Document huma.FormFile `form:"file" required:"true" contentType:"application/octet-stream" doc:"File"`
}
=======
package data

import "github.com/danielgtaylor/huma/v2"

type DocumentQuery struct {
}

type DocumentData struct {
	Document huma.FormFile `form:"file" required:"true" contentType:"application/octet-stream" doc:"File"`
}
>>>>>>> 22022f0081c75477042da66cd81443ff4401ca37
