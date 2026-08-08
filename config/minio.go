<<<<<<< HEAD
package config

import (
	"context"
	tls "crypto/tls"
	"io"
	http "net/http"
	"net/url"
	"time"

	"cdn/common/constants"
	"cdn/common/helpers"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"go.uber.org/zap"
)

var MinioClient *minio.Client

// Establishes a connection to Minio server.
func ConnectMinio() error {
	var err error

	// Connect to minio endpoint
	MinioClient, err = minio.New(Env.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(Env.MinioAccessKeyID, Env.MinioAccessKeySecret, ""),
		Secure: Env.MinioUseSSL,
		Region: constants.MINIO_LOCATION_DEFAULT,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: Env.InsecureSkipVerify, // Ignore les certificats invalides
			},
		},
	})
	if err != nil || MinioClient == nil {
		if err != nil {
			helpers.Logger.Warn(
				"Failed to connect to minio server!",
				zap.Error(err),
			)
		} else {
			helpers.Logger.Warn("Failed to connect to minio server: no client was returned")
		}
		return fmt.Errorf("failed to initialize minio client")
	}

	_, err = MinioClient.GetCreds()
	if err != nil {
		helpers.Logger.Warn(
			"Failed to get creds!",
			zap.Error(err),
		)
		return err
	}

	// Create the buckets
	errBucketImg := CreateBucket(constants.MINIO_BUCKET_IMAGES, constants.MINIO_LOCATION_DEFAULT)
	if errBucketImg != nil {
		helpers.Logger.Warn(
			"Bucket for images error!",
			zap.Error(errBucketImg),
		)
		return errBucketImg
	}
	errBucketDoc := CreateBucket(constants.MINIO_BUCKET_DOCUMENTS, constants.MINIO_LOCATION_DEFAULT)
	if errBucketDoc != nil {
		helpers.Logger.Warn(
			"Bucket for documents error!",
			zap.Error(errBucketDoc),
		)
		return errBucketDoc
	}
	return nil
}

// Create a bucket
func CreateBucket(bucketName string, location string) error {
	if MinioClient == nil {
		helpers.Logger.Warn(
			"No available minio client!",
		)
		return fmt.Errorf("No available minio client!")
	}
	err := MinioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket
		exists, errBucketExists := MinioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists != nil {
			helpers.Logger.Warn(
				"Failed to check if the bucket exists!",
				zap.String("Error: ", errBucketExists.Error()),
			)
			return errBucketExists
		}
		if exists {
			helpers.Logger.Info(
				"This bucket already exists!",
				zap.String("Bucket name: ", bucketName),
			)
			return nil
		}
		helpers.Logger.Warn(
			"Failed to add new bucket!",
			zap.String("Error: ", err.Error()),
		)
		return err

	}
	helpers.Logger.Info(
		"Minio bucked created!",
		zap.String("Bucket name: ", bucketName),
	)

	return nil
}

// UploadObjectToMinio Uploads object to minio server
func UploadFObjectToMinio(bucketName string, path string, objectName string) (*minio.UploadInfo, error) {
	if MinioClient == nil {
		helpers.Logger.Warn(
			"No available minio client!",
		)
		return nil, fmt.Errorf("No available minio client!")
	}
	contentType := "application/octet-stream"
	info, err := MinioClient.FPutObject(
		context.Background(),
		bucketName,
		objectName,
		path,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		helpers.Logger.Warn(
			"Minio upload error!",
			zap.String("Error: ", err.Error()),
		)
		return nil, err
	}
	return &info, nil
}

// UploadObjectToMinio Uploads object to minio server
func UploadObjectToMinio(bucketName string, objectName string, reader io.Reader, objectSize int64) (*minio.UploadInfo, error) {
	if MinioClient == nil {
		helpers.Logger.Warn(
			"No available minio client!",
		)
		return nil, fmt.Errorf("No available minio client!")
	}
	contentType := "application/octet-stream"
	info, err := MinioClient.PutObject(
		context.Background(),
		bucketName,
		objectName,
		reader,
		objectSize,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		helpers.Logger.Warn(
			"Minio upload error!",
			zap.String("Error: ", err.Error()),
		)
		return nil, err
	}
	return &info, nil
}

// GetObjectFromMinio Gets object from minio server
func GetPresignedObjectFromMinio(bucketName string, objectName string, expiry time.Duration) (*url.URL, error) {
	if MinioClient == nil {
		helpers.Logger.Warn(
			"No available minio client!",
		)
		return nil, fmt.Errorf("No available minio client!")
	}

	// Generate a signed url
	reqParams := make(url.Values)
	presignedUrl, err := MinioClient.PresignedGetObject(
		context.Background(),
		bucketName,
		objectName,
		expiry,
		reqParams,
	)
	if err != nil || presignedUrl == nil {
		helpers.Logger.Warn(
			"Failed to get public object!",
		)
		return nil, fmt.Errorf(err.Error())
	}
	helpers.Logger.Info(
		"Presigned object",
		zap.String("String", presignedUrl.String()),
		zap.String("Host", presignedUrl.Host),
		zap.String("Path", presignedUrl.Path),
		zap.String("RawPath", presignedUrl.RawPath),
		zap.String("RawQuery", presignedUrl.RawQuery),
		zap.String("Scheme", presignedUrl.Scheme),
		zap.String("Redacted", presignedUrl.Redacted()),
	)
	return presignedUrl, nil
}

// DeleteObjectFromMinio Deletes object from minio server
func DeleteObjectFromMinio(bucketName string, objectName string) error {
	if MinioClient == nil {
		helpers.Logger.Warn(
			"No available minio client!",
		)
		return fmt.Errorf("No available minio client!")
	}
	err := MinioClient.RemoveObject(
		context.Background(),
		bucketName,
		objectName,
		minio.RemoveObjectOptions{},
	)
	if err != nil {
		helpers.Logger.Warn(
			"Minio delete object error!",
			zap.String("Error: ", err.Error()),
		)
		return err
	}
	return nil
}
=======
package config

import (
	"context"
	tls "crypto/tls"
	"io"
	http "net/http"
	"net/url"
	"time"

	"cdn/common/constants"
	"cdn/common/helpers"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"go.uber.org/zap"
)

var MinioClient *minio.Client

// Establishes a connection to Minio server.
func ConnectMinio() error {
	var err error

	// Connect to minio endpoint
	MinioClient, err = minio.New(Env.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(Env.MinioAccessKeyID, Env.MinioAccessKeySecret, ""),
		Secure: Env.MinioUseSSL,
		Region: constants.MINIO_LOCATION_DEFAULT,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: Env.InsecureSkipVerify, // Ignore les certificats invalides
			},
		},
	})
	if err != nil || MinioClient == nil {
		helpers.Logger.Warn(
			"Failed to connect to minio server!",
			zap.String("Error: ", err.Error()),
		)
	}
	_, err = MinioClient.GetCreds()
	if err != nil {
		helpers.Logger.Warn(
			"Failed to get creds!",
			zap.String("Error: ", err.Error()),
		)
		return err
	}

	// Create the buckets
	errBucketImg := CreateBucket(constants.MINIO_BUCKET_IMAGES, constants.MINIO_LOCATION_DEFAULT)
	if errBucketImg != nil {
		helpers.Logger.Warn(
			"Bucket for images error!",
			zap.String("Error: ", err.Error()),
		)
		return errBucketImg
	}
	errBucketDoc := CreateBucket(constants.MINIO_BUCKET_DOCUMENTS, constants.MINIO_LOCATION_DEFAULT)
	if errBucketDoc != nil {
		helpers.Logger.Warn(
			"Bucket for documents error!",
			zap.String("Error: ", err.Error()),
		)
		return errBucketDoc
	}
	return err
}

// Create a bucket
func CreateBucket(bucketName string, location string) error {
	if MinioClient == nil {
		helpers.Logger.Warn(
			"No available minio client!",
		)
		return fmt.Errorf("No available minio client!")
	}
	err := MinioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket
		exists, errBucketExists := MinioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists != nil {
			helpers.Logger.Warn(
				"Failed to check if the bucket exists!",
				zap.String("Error: ", errBucketExists.Error()),
			)
			return errBucketExists
		}
		if exists {
			helpers.Logger.Info(
				"This bucket already exists!",
				zap.String("Bucket name: ", bucketName),
			)
			return nil
		}
		helpers.Logger.Warn(
			"Failed to add new bucket!",
			zap.String("Error: ", err.Error()),
		)
		return err

	}
	helpers.Logger.Info(
		"Minio bucked created!",
		zap.String("Bucket name: ", bucketName),
	)

	return nil
}

// UploadObjectToMinio Uploads object to minio server
func UploadFObjectToMinio(bucketName string, path string, objectName string) (*minio.UploadInfo, error) {
	if MinioClient == nil {
		helpers.Logger.Warn(
			"No available minio client!",
		)
		return nil, fmt.Errorf("No available minio client!")
	}
	contentType := "application/octet-stream"
	info, err := MinioClient.FPutObject(
		context.Background(),
		bucketName,
		objectName,
		path,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		helpers.Logger.Warn(
			"Minio upload error!",
			zap.String("Error: ", err.Error()),
		)
		return nil, err
	}
	return &info, nil
}

// UploadObjectToMinio Uploads object to minio server
func UploadObjectToMinio(bucketName string, objectName string, reader io.Reader, objectSize int64) (*minio.UploadInfo, error) {
	if MinioClient == nil {
		helpers.Logger.Warn(
			"No available minio client!",
		)
		return nil, fmt.Errorf("No available minio client!")
	}
	contentType := "application/octet-stream"
	info, err := MinioClient.PutObject(
		context.Background(),
		bucketName,
		objectName,
		reader,
		objectSize,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		helpers.Logger.Warn(
			"Minio upload error!",
			zap.String("Error: ", err.Error()),
		)
		return nil, err
	}
	return &info, nil
}

// GetObjectFromMinio Gets object from minio server
func GetPresignedObjectFromMinio(bucketName string, objectName string, expiry time.Duration) (*url.URL, error) {
	if MinioClient == nil {
		helpers.Logger.Warn(
			"No available minio client!",
		)
		return nil, fmt.Errorf("No available minio client!")
	}

	// Generate a signed url
	reqParams := make(url.Values)
	presignedUrl, err := MinioClient.PresignedGetObject(
		context.Background(),
		bucketName,
		objectName,
		expiry,
		reqParams,
	)
	if err != nil || presignedUrl == nil {
		helpers.Logger.Warn(
			"Failed to get public object!",
		)
		return nil, fmt.Errorf(err.Error())
	}
	helpers.Logger.Info(
		"Presigned object",
		zap.String("String", presignedUrl.String()),
		zap.String("Host", presignedUrl.Host),
		zap.String("Path", presignedUrl.Path),
		zap.String("RawPath", presignedUrl.RawPath),
		zap.String("RawQuery", presignedUrl.RawQuery),
		zap.String("Scheme", presignedUrl.Scheme),
		zap.String("Redacted", presignedUrl.Redacted()),
	)
	return presignedUrl, nil
}

// DeleteObjectFromMinio Deletes object from minio server
func DeleteObjectFromMinio(bucketName string, objectName string) error {
	if MinioClient == nil {
		helpers.Logger.Warn(
			"No available minio client!",
		)
		return fmt.Errorf("No available minio client!")
	}
	err := MinioClient.RemoveObject(
		context.Background(),
		bucketName,
		objectName,
		minio.RemoveObjectOptions{},
	)
	if err != nil {
		helpers.Logger.Warn(
			"Minio delete object error!",
			zap.String("Error: ", err.Error()),
		)
		return err
	}
	return nil
}
>>>>>>> 22022f0081c75477042da66cd81443ff4401ca37
