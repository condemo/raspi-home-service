package util

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	jwt.RegisteredClaims
	UserID uint8 `json:"userID"`
}

func CreateJWT(id uint8) (string, time.Time, error) {
	jwtKey := []byte(os.Getenv("JWT_KEY"))

	expire, _ := strconv.Atoi(os.Getenv("JWT_EXP_DAYS"))
	expDays := time.Now().Add(time.Hour * 24 * time.Duration(expire))

	claims := UserClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expDays),
		},
		id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtKey)

	return signedToken, expDays, err
}

func ValidateJWT(t string) (*UserClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(
		t, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_KEY")), nil
		})
	return parsedToken.Claims.(*UserClaims), err
}
