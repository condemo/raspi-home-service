package types

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type User struct {
	Username string `bun:"username,unique,notnull" validate:"gt=0,lt=10,alphanum" json:"username"`
	Password string `bun:"password" validate:"gt=0,lt=15,alphanum" json:"-"`
	ID       uint8  `bun:",pk,autoincrement" json:"id"`
}

// PERF: Mejorar el sistema
func (u User) Validate() error {
	err := validate.Struct(u)
	if err != nil {
		valErr := err.(validator.ValidationErrors)
		return valErr
	}
	return nil
}
