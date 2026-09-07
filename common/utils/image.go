package utils

import (
	"cdn/common/helpers"
	"fmt"

	"github.com/gabriel-vasile/mimetype"
	"github.com/h2non/bimg"
	"go.uber.org/zap"
)

const defaultImageSize = 1000

func ResizeImage(inputBuffer []byte, width int, height int, quality int) (outputBuffer []byte, outputSize bimg.ImageSize, err error) {
	// Check for supported types
	inputType := bimg.NewImage(inputBuffer).Type()
	if inputType == "unknown" {
		tmpMimeType := mimetype.Detect(inputBuffer)
		if tmpMimeType != nil {
			inputType = tmpMimeType.Extension()
		}
	}
	if !bimg.IsTypeNameSupported(inputType) {
		if inputType != "ico" {
			err = fmt.Errorf("%s %s", "Unsupported image type:", inputType)
			helpers.Logger.Warn(
				"Unsupported image type",
				zap.String("Value: ", inputType),
			)
			return
		}
	}
	// Check for resizable types
	if inputType == "pdf" || inputType == "svg" || inputType == "magick" ||
		inputType == "heif" || inputType == "avif" || inputType == "ico" {
		outputSize, err = bimg.NewImage(inputBuffer).Size()
		outputBuffer = inputBuffer
		return
	}

	// Apply default type for maximum performance
	outputType := bimg.JPEG
	if inputType == "gif" {
		outputType = bimg.GIF
	}
	if inputType == "png" || inputType == "webp" || inputType == "tiff" {
		outputType = bimg.WEBP
	}

	// Fix width and height
	fixedWidth := width
	fixedHeight := height
	if fixedWidth <= 0 {
		fixedWidth = 1
	}
	if fixedHeight <= 0 {
		fixedHeight = 1
	}
	inputSize, err := bimg.NewImage(inputBuffer).Size()
	if fixedWidth > inputSize.Width {
		fixedWidth = inputSize.Width
	}
	if fixedHeight > inputSize.Height {
		fixedHeight = inputSize.Height
	}

	// Fix quality
	fixedQuality := quality
	if fixedQuality <= 0 {
		fixedQuality = bimg.Quality
	}

	// Apply image processing
	outputBuffer, err = bimg.NewImage(inputBuffer).Resize(fixedWidth, fixedHeight)
	if err != nil {
		helpers.Logger.Warn(
			"Failed to resize image!",
			zap.String("Error", err.Error()),
		)
	}
	outputBuffer, err = bimg.NewImage(outputBuffer).Process(bimg.Options{
		Quality:  quality,
		Lossless: false,
		Gravity:  bimg.GravityCentre,
		Type:     outputType,
	})
	if err != nil {
		helpers.Logger.Warn(
			"Failed to process image!",
			zap.String("Error", err.Error()),
		)
	}
	outputSize, err = bimg.NewImage(outputBuffer).Size()
	return
}
