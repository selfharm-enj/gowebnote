package utils

import (
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct{}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

func (h BcryptHasher) HashPassword(password string) (string, error) {
	const op = "infrastructure.crypto.bcrypt.Hash"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", fmt.Errorf("%s:%w", op, err)
	}
	return string(hashedPassword), nil
}

func (h BcryptHasher) Compare(hash, password string) error {
	const op = "infrastructure.crypto.bcrypt.Compare"
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		slog.Debug("passwords are not the same", "op", op)
		return err
	}
	return nil
}
