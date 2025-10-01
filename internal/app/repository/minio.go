package repository

import (
	"fmt"
	"lab1_rip/internal/app/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

type Minio struct {
	Client *minio.Client
	Config *config.Config
}

func NewMinio(config *config.Config) (*Minio, error) {
	endpoint := fmt.Sprintf("%s:%d", config.MinioHost, config.MinioPort)
	fmt.Println(endpoint)
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinioAccessKey, config.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	return &Minio{Client: client, Config: config}, nil

}
