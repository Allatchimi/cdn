package data

import "github.com/danielgtaylor/huma/v2"

type DocumentQuery struct {
}

type DocumentData struct {
	Document huma.FormFile `form:"file" required:"true" contentType:"application/octet-stream" doc:"File"`
}
