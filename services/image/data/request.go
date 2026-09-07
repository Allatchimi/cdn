package data

import "github.com/danielgtaylor/huma/v2"

type ImageQuery struct {
	Width   int `query:"w" required:"false" minimum:"1" maximum:"5000" default:"900" doc:"Width in pixels"`
	Height  int `query:"h" required:"false" minimum:"1" maximum:"5000" default:"900" doc:"Height in pixels"`
	Quality int `query:"q" required:"false" minimum:"1" maximum:"100" default:"75" doc:"Quality in percent"`
}

type ImageData struct {
	Image huma.FormFile `form:"file" required:"true" contentType:"application/octet-stream" doc:"Select you image"`
}
