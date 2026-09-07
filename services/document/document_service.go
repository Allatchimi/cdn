package document

import (
	"cdn/common/constants"
	"cdn/common/helpers"
	"cdn/config"
	"cdn/services/document/data"
	"context"
	"fmt"
	"net/http"
	"time"
)

type Service struct {
}

const subDir = "/documents"

func NewService() *Service {
	return &Service{}
}

// Create new document
func (service *Service) Create(
	ctx *context.Context,
	option *data.DocumentQuery,
	documentData *data.DocumentData,
) (result *data.UploadDocumentResponse, errCode int, err error) {
	// Encode file name
	initialFileName := documentData.Document.Filename
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), initialFileName)

	// Upload to minio
	info, err := config.UploadObjectToMinio(constants.MINIO_BUCKET_DOCUMENTS, fileName, documentData.Document.File, documentData.Document.Size)
	if err != nil || info == nil {
		errCode = http.StatusInternalServerError
		err = constants.HTTP_500_ERROR_MESSAGE("save document")
		return
	}

	result = &data.UploadDocumentResponse{
		Url:  config.Env.Hostname + config.Env.ApiGroup + "/documents/" + fileName,
		Path: fileName,
	}
	return
}

// Update existing document
func (service *Service) Update(
	ctx *context.Context,
	objectName string,
	option *data.DocumentQuery,
	documentData *data.DocumentData,
) (result *data.UploadDocumentResponse, errCode int, err error) {
	// Delete object
	err = config.DeleteObjectFromMinio(constants.MINIO_BUCKET_DOCUMENTS, objectName)
	if err != nil {
		errCode = http.StatusNotFound
		err = constants.HTTP_404_ERROR_MESSAGE("resource")
		return
	}

	// Upload the new one
	result, errCode, err = service.Create(ctx, option, documentData)
	if err != nil {
		errCode = http.StatusInternalServerError
		err = constants.HTTP_500_ERROR_MESSAGE("upload")
		return
	}
	return
}

// Delete document with matching id and return affected rows
func (service *Service) Delete(
	ctx *context.Context,
	objectName string,
) (result bool, errCode int, err error) {
	err = config.DeleteObjectFromMinio(constants.MINIO_BUCKET_DOCUMENTS, objectName)
	if err != nil {
		errCode = http.StatusNotFound
		err = constants.HTTP_404_ERROR_MESSAGE("resource")
	}
	result = true
	return
}

// Get document with matching id
func (service *Service) Get(
	ctx *context.Context,
	objectName string,
	option data.DocumentQuery,
) (url string, errCode int, err error) {
	// Get cached presigned url
	cachedUrl, ok := config.OtterCache.Get(objectName)
	if ok && len(cachedUrl) > 0 {
		url = cachedUrl
		helpers.Logger.Info("Returned cached url!")
		return
	}

	// Get new presigned url
	presignedUrl, err := config.GetPresignedObjectFromMinio(constants.MINIO_BUCKET_DOCUMENTS, objectName, time.Minute*15)
	if err != nil || presignedUrl == nil {
		errCode = http.StatusNotFound
		err = constants.HTTP_404_ERROR_MESSAGE("resource")
	}
	url = presignedUrl.String()

	// Cache the new presigned url
	config.OtterCache.Set(objectName, presignedUrl.String(), time.Minute*14)
	helpers.Logger.Info("New url cached!")
	return
}
