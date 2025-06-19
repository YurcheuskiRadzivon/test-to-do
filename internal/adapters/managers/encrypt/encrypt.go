package encryptmanage

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

const (
	cost = 10

	ErrEncryptPassword = "FAILED_ENCRYPT_PASSWORD"
	ErrInvalidPassword = "INVALID_PASSWORD"
)

type EncryptManager interface {
	EncodePassword(password string) (string, error)
	CheckPassword(password, hashedPassword string) error
}

type Encrypter struct {
	cost int
}

func NewEncrypter() *Encrypter {
	return &Encrypter{
		cost: cost,
	}
}

func (e *Encrypter) EncodePassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		e.cost,
	)
	if err != nil {
		log.Printf("Failed while encrrypt password: %v", err)
		return "", errors.New(ErrEncryptPassword)
	}

	return string(hashedPassword), nil
}

func (e *Encrypter) CheckPassword(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
	if err != nil {
		log.Printf("Invalid password while compare and hash: %v", err)
		return errors.New(ErrInvalidPassword)
	}

	return nil
}
