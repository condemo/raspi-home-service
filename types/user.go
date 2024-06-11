package types

import (
	"unicode"
)

type User struct {
	Username string `gorm:"primaryKey,unique" json:"username"`
	Password string `json:"-"`
	ID       uint8  `json:"id"`
}

// PERF: Mejorar usando un mejor validator
func (u User) Validate() bool {
	if u.Password == "" || u.Username == "" {
		return false
	}

	for _, c := range u.Username {
		if unicode.IsSpace(c) || unicode.IsSymbol(c) {
			return false
		}
	}
	for _, c := range u.Password {
		if unicode.IsSpace(c) || unicode.IsSymbol(c) {
			return false
		}
	}

	return true
}
