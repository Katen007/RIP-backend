package repository

import (
	"lab1_rip/internal/app/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db    *gorm.DB
	minio *Minio
}

func NewRepository(dsn string, config *config.Config) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) // подключаемся к БД
	if err != nil {
		return nil, err
	}
	minio, err := NewMinio(config)
	if err != nil {
		return nil, err
	}

	// Возвращаем объект Repository с подключенной базой данных
	return &Repository{
		db:    db,
		minio: minio,
	}, nil
}
