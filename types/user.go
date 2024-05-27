package types

type User struct {
	Username string `gorm:"primaryKey,unique" json:"username"`
	Password string `json:"-"`
	ID       uint8  `json:"id"`
}

// TODO: Implementar
func (u User) Validate() bool { return true }
