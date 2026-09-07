package data

type UploadImageResponse struct {
	Url    string `json:"url" required:"false" doc:"Asset URL" example:"https://cdn.application.com/images/aKMiF_an5diG2j"`
	Path   string `json:"path" required:"false" doc:"Asset path from /images" example:"aKMiF_an5diG2j"`
	Width  int    `json:"width" required:"false" doc:"Default width in pixels" example:"1280"`
	Height int    `json:"height" required:"false" doc:"Default height in pixels" example:"720"`
}
