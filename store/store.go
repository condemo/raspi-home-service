package store

import (
	"context"

	"github.com/condemo/raspi-home-service/types"
	"github.com/uptrace/bun"
)

type Store interface {
	CreateUser(*types.User) error
	GetUserByUsername(string) (*types.User, error)
}

type Storage struct {
	db *bun.DB
}

func NewStorage(db *bun.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) CreateUser(u *types.User) error {
	_, err := s.db.NewInsert().Model(u).
		Returning("username", "id").Exec(context.Background())

	return err
}

func (s *Storage) GetUserByUsername(us string) (*types.User, error) {
	user := new(types.User)
	err := s.db.NewSelect().Model(user).
		Where("username = ?", us).Limit(1).Scan(context.Background())

	return user, err
}
