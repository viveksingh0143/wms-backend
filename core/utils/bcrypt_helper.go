package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func GenerateFromPassword(text string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CompareHashAndPassword(password string, storedHashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(password))
}
