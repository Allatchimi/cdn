package document

import (
	"cdn/common/types"
	"cdn/services/document/data"
	"context"

	"github.com/danielgtaylor/huma/v2"
)

type Controller struct {
	Service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{Service: service}
}

func (controller *Controller) Create(
	ctx *context.Context,
	input *struct {
		data.DocumentQuery
		RawBody huma.MultipartFormFiles[data.DocumentData]
	},
) (result *data.UploadDocumentResponse, errCode int, err error) {
	result, errCode, err = controller.Service.Create(ctx, &input.DocumentQuery, input.RawBody.Data())
	return
}

func (controller *Controller) Update(
	ctx *context.Context,
	input *struct {
		types.FilePath
		data.DocumentQuery
		RawBody huma.MultipartFormFiles[data.DocumentData]
	},
) (result *data.UploadDocumentResponse, errCode int, err error) {

	result, errCode, err = controller.Service.Update(ctx, input.FilePath.Path, &input.DocumentQuery, input.RawBody.Data())
	return
}

func (controller *Controller) Delete(
	ctx *context.Context,
	input *struct {
		types.FilePath
	},
) (result bool, errCode int, err error) {
	result, errCode, err = controller.Service.Delete(ctx, input.FilePath.Path)
	return
}

func (controller *Controller) Get(
	ctx *context.Context,
	input *struct {
		types.FilePath
		data.DocumentQuery
	},
) (result string, errCode int, err error) {
	result, errCode, err = controller.Service.Get(ctx, input.FilePath.Path, input.DocumentQuery)
	return
}
