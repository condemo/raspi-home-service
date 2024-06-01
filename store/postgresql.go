package store

import (
	"fmt"
	"os"

	"github.com/condemo/raspi-home-service/config"
	"github.com/condemo/raspi-home-service/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresStorage struct {
	db *gorm.DB
}

func NewPostgresStore() *PostgresStorage {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Madrid",
		config.Envs.DBHost, config.Envs.DBUser, config.Envs.DBPass,
		config.Envs.DBName, config.Envs.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("database error:", err)
		os.Exit(1)
	}

	return &PostgresStorage{db: db}
}

func (s *PostgresStorage) Init() (*gorm.DB, error) {
	// Load Tables
	s.db.AutoMigrate(&types.User{})

	return s.db, nil
}
