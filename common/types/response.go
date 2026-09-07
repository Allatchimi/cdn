package types

type DeletedResponse struct {
	Deleted bool `json:"deleted" required:"false" doc:"Deleted ?" example:"true"`
}
