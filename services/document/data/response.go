package data

type UploadDocumentResponse struct {
	Url  string `json:"url" required:"false" doc:"Asset URL" example:"https://cdn.application.com/documents/aKMiF_an5diG2j"`
	Path string `json:"path" required:"false" doc:"Asset path from /documents" example:"aKMiF_an5diG2j"`
}
