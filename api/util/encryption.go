package util

import "golang.org/x/crypto/bcrypt"

func EncryptPass(p string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(p), 14)
	if err != nil {
		return "", err
	}

	return string(hashedPass), nil
}

func VerifyPass(hashedPass, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(pass))

	return err == nil
}
