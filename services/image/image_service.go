package image

import (
	"cdn/common/constants"
	"cdn/common/helpers"
	"cdn/common/utils"
	"cdn/config"
	"cdn/services/image/data"
	"context"
	"net/http"
	"time"
)

type Service struct {
}

const subDir = "/images"

func NewService() *Service {
	return &Service{}
}

// Create new image
func (service *Service) Create(
	ctx *context.Context,
	option *data.ImageQuery,
	imageData *data.ImageData,
) (result *data.UploadImageResponse, errCode int, err error) {
	// Read buffer
	buffer, err := utils.ReadMultipartFile(imageData.Image.File)
	if err != nil {
		errCode = http.StatusInternalServerError
		err = constants.HTTP_500_ERROR_MESSAGE("read image")
		return
	}

	// Resize and compress the image.
	bufferResized, size, err := utils.ResizeImage(
		buffer, option.Width, option.Height, option.Quality,
	)
	if err != nil {
		bufferResized = buffer
		errCode = http.StatusBadRequest
		err = constants.HTTP_400_ERROR_MESSAGE("unsupported file type")
		return
	}

	// Create the file
	fileAbsPath, fileName, err := utils.SaveFile(bufferResized, constants.ASSET_UPLOADS_PATH+subDir)
	if err != nil {
		errCode = http.StatusInternalServerError
		err = constants.HTTP_500_ERROR_MESSAGE("save image")
		return
	}

	// Upload to minio
	info, err := config.UploadFObjectToMinio(constants.MINIO_BUCKET_IMAGES, fileAbsPath, fileName)
	// Delete the file
	_, _ = utils.DeleteFile(fileAbsPath)
	// Check errors
	if err != nil || info == nil {
		errCode = http.StatusInternalServerError
		err = constants.HTTP_500_ERROR_MESSAGE("save image")
		return
	}

	result = &data.UploadImageResponse{
		Url:    config.Env.Hostname + config.Env.ApiGroup + "/images/" + fileName,
		Path:   fileName,
		Width:  size.Width,
		Height: size.Height,
	}
	return
}

// Update existing image
func (service *Service) Update(
	ctx *context.Context,
	objectName string,
	option *data.ImageQuery,
	imageData *data.ImageData,
) (result *data.UploadImageResponse, errCode int, err error) {
	// Delete object
	err = config.DeleteObjectFromMinio(constants.MINIO_BUCKET_IMAGES, objectName)
	if err != nil {
		errCode = http.StatusNotFound
		err = constants.HTTP_404_ERROR_MESSAGE("resource")
		return
	}

	// Upload the new one
	result, errCode, err = service.Create(ctx, option, imageData)
	if err != nil {
		errCode = http.StatusInternalServerError
		err = constants.HTTP_500_ERROR_MESSAGE("upload")
		return
	}
	return
}

// Delete image with matching id and return affected rows
func (service *Service) Delete(
	ctx *context.Context,
	objectName string,
) (result bool, errCode int, err error) {
	err = config.DeleteObjectFromMinio(constants.MINIO_BUCKET_IMAGES, objectName)
	if err != nil {
		errCode = http.StatusNotFound
		err = constants.HTTP_404_ERROR_MESSAGE("resource")
	}
	result = true
	return
}

// Get image with matching id
func (service *Service) Get(
	ctx *context.Context,
	objectName string,
	option data.ImageQuery,
) (url string, errCode int, err error) {
	// Get cached presigned url
	cachedUrl, ok := config.OtterCache.Get(objectName)
	if ok && len(cachedUrl) > 0 {
		url = cachedUrl
		helpers.Logger.Info("Returned cached url!")
		return
	}

	// Get new presigned url
	presignedUrl, err := config.GetPresignedObjectFromMinio(constants.MINIO_BUCKET_IMAGES, objectName, time.Minute*15)
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
