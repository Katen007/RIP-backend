package repository

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

func (r *Repository) UploadComponentImg(ctx context.Context, file io.Reader, fileName string, fileSize int64, fileType string) (string, error) {
	var filePath string
	filePath, err := CreateNewFilePath(fileName)
	if err != nil {
		return "", err
	}

	putOptions := minio.PutObjectOptions{
		ContentType: fileType,
	}
	_, err = r.minio.Client.PutObject(ctx, r.minio.Config.MinioBucket, filePath, file, fileSize, putOptions)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return fmt.Sprintf("http://%v:%v/%v/%v", r.minio.Config.MinioHost, r.minio.Config.MinioPort, r.minio.Config.MinioBucket, filePath), nil
}

func (r *Repository) DeleteComponentImg(ctx context.Context, url *string) error {
	filePath, err := r.GetUrlComponentImg(url)
	if err != nil {
		return err
	}
	err = r.minio.Client.RemoveObject(context.Background(), r.minio.Config.MinioBucket, filePath, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) GetUrlComponentImg(urlPtr *string) (string, error) {
	url := *urlPtr
	if url == "" {
		return "", nil
	}
	lenBucket := len(r.minio.Config.MinioBucket)
	indexBucket := strings.Index(url, r.minio.Config.MinioBucket)
	if indexBucket == -1 {
		return "", fmt.Errorf("not found bucket %v in url-path %v", r.minio.Config.MinioBucket, url)
	}
	filePath := url[indexBucket+lenBucket+1:]
	logrus.Info(filePath)
	return filePath, nil
}
