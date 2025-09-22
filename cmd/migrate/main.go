package main

import (
	"lab1_rip/internal/app/ds"
	"lab1_rip/internal/app/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&ds.Text{},
		&ds.User{},
		&ds.ReadIndxsToText{},
		&ds.ReadIndxs{},
	)
	if err != nil {
		panic("cant migrate db")
	}
}
