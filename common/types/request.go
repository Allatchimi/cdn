package types

type FilePath struct {
	Path string `json:"id" path:"id" required:"true" minLength:"3" doc:"Path" example:""`
}
