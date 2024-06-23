package types

import (
	"fmt"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type User struct {
	Username string `bun:"username,unique,notnull" json:"username"`
	Password string `json:"-"`
	ID       uint8  `bun:",pk,autoincrement" json:"id"`
}

// PERF: Mejorar usando un mejor validator
func (u User) Validate() bool {
	if u.Password == "" || u.Username == "" {
		fmt.Println("empty crendentials")
		return false
	}

	for _, c := range u.Username {
		if unicode.IsSpace(c) {
			return false
		}
	}
	for _, c := range u.Password {
		if unicode.IsSpace(c) {
			return false
		}
	}

	return true
}
