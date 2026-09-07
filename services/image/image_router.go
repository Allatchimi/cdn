package image

import (
	"cdn/common/constants"
	"cdn/middlewares"
	"context"
	"fmt"
	"net/http"

	"cdn/common/types"
	"cdn/services/image/data"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
)

func RegisterEndpoints(
	humaApi *huma.API,
	controller *Controller,
) {
	var endpointConfig = types.APIEndpointConfig{
		Group: "/images",
		Tag:   []string{"Images"},
	}

	// Upload image
	huma.Register(
		*humaApi,
		huma.Operation{
			OperationID: "post-image",
			Summary:     "Upload image",
			Description: "Upload new image.",
			Method:      http.MethodPost,
			Path:        endpointConfig.Group,
			Tags:        endpointConfig.Tag,
			Security: []map[string][]string{
				{constants.SECURITY_AUTH_NAME: {}}, // Used to require authentication
			},
			MaxBodyBytes:  (1024 * 1000) * 10, // (1 KiB * 1000) * 10 = 10MiB
			DefaultStatus: http.StatusOK,
			Errors:        []int{http.StatusInternalServerError, http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden},
		},
		func(
			ctx context.Context,
			input *struct {
				data.ImageQuery
				RawBody huma.MultipartFormFiles[data.ImageData]
			},
		) (*struct{ Body *data.UploadImageResponse }, error) {
			result, errCode, err := controller.Create(&ctx, input)
			if err != nil {
				return nil, huma.NewError(errCode, err.Error())
			}
			return &struct{ Body *data.UploadImageResponse }{Body: result}, nil
		},
	)

	// Update image
	huma.Register(
		*humaApi,
		huma.Operation{
			OperationID: "update-image",
			Summary:     "Update image",
			Description: "Update existing image.",
			Method:      http.MethodPut,
			Path:        fmt.Sprintf("%s/{id}", endpointConfig.Group),
			Tags:        endpointConfig.Tag,
			Security: []map[string][]string{
				{constants.SECURITY_AUTH_NAME: {}}, // Used to require authentication
			},
			MaxBodyBytes:  (1024 * 1000) * 10, // (1 KiB * 1000) * 10 = 10MiB
			DefaultStatus: http.StatusOK,
			Errors:        []int{http.StatusInternalServerError, http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound},
		},
		func(
			ctx context.Context,
			input *struct {
				types.FilePath
				data.ImageQuery
				RawBody huma.MultipartFormFiles[data.ImageData]
			},
		) (*struct{ Body *data.UploadImageResponse }, error) {
			result, errCode, err := controller.Update(&ctx, input)
			if err != nil {
				return nil, huma.NewError(errCode, err.Error())
			}
			return &struct{ Body *data.UploadImageResponse }{Body: result}, nil
		},
	)

	// Delete image
	huma.Register(
		*humaApi,
		huma.Operation{
			OperationID: "delete-image",
			Summary:     "Delete image",
			Description: "Delete existing image.",
			Method:      http.MethodDelete,
			Path:        fmt.Sprintf("%s/{id}", endpointConfig.Group),
			Tags:        endpointConfig.Tag,
			Security: []map[string][]string{
				{constants.SECURITY_AUTH_NAME: {}}, // Used to require authentication
			},
			MaxBodyBytes:  1024, // 1 KiB
			DefaultStatus: http.StatusOK,
			Errors:        []int{http.StatusInternalServerError, http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound},
		},
		func(
			ctx context.Context,
			input *struct {
				types.FilePath
			},
		) (*struct{ Body *types.DeletedResponse }, error) {
			result, errCode, err := controller.Delete(&ctx, input)
			if err != nil {
				return nil, huma.NewError(errCode, err.Error())
			}
			return &struct{ Body *types.DeletedResponse }{Body: &types.DeletedResponse{Deleted: result}}, nil
		},
	)

	// Get image
	huma.Register(
		*humaApi,
		huma.Operation{
			OperationID:   "get-image",
			Summary:       "Get image",
			Description:   "Return existing image data",
			Method:        http.MethodGet,
			Path:          fmt.Sprintf("%s/{id}", endpointConfig.Group),
			Tags:          endpointConfig.Tag,
			MaxBodyBytes:  1024, // 1 KiB
			DefaultStatus: http.StatusOK,
			Errors:        []int{http.StatusInternalServerError, http.StatusBadRequest, http.StatusForbidden, http.StatusNotFound},
		},
		func(
			ctx context.Context,
			input *struct {
				types.FilePath
				data.ImageQuery
			},
		) (
			*huma.StreamResponse,
			error,
		) {
			result, errCode, err := controller.Get(&ctx, input)
			if err != nil {
				return nil, huma.NewError(errCode, err.Error())
			}
			if len(result) < 1 {
				return nil, huma.NewError(http.StatusNotFound, "Invalid url!")
			}

			ginCtx := ctx.Value(middlewares.GIN_CONTEXT_KEY).(*gin.Context)
			if ginCtx != nil {
				ginCtx.Redirect(http.StatusTemporaryRedirect, result)
			}
			return nil, nil
		},
	)
}
