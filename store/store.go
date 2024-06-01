package store

import (
	"github.com/condemo/raspi-home-service/types"
	"gorm.io/gorm"
)

type Store interface {
	CreateUser(*types.User) error
	GetUserByUsername(string) (*types.User, error)
}

type Storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) CreateUser(u *types.User) error {
	res := s.db.Create(u)

	return res.Error
}

func (s *Storage) GetUserByUsername(us string) (*types.User, error) {
	user := new(types.User)
	res := s.db.First(user, "username = ?", us)

	return user, res.Error
}
