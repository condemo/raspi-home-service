package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/condemo/raspi-home-service/config"
	"github.com/condemo/raspi-home-service/types"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type PostgresStorage struct {
	db *bun.DB
}

func NewPostgresStore() *PostgresStorage {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.Envs.DBUser, config.Envs.DBPass, config.Envs.DBHost,
		config.Envs.DBPort, config.Envs.DBName,
	)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	return &PostgresStorage{db: db}
}

func (s *PostgresStorage) Init() (*bun.DB, error) {
	// Load Tables
	_, err := s.db.NewCreateTable().Model((*types.User)(nil)).
		IfNotExists().Exec(context.Background())
	if err != nil {
		return nil, err
	}

	return s.db, nil
}
