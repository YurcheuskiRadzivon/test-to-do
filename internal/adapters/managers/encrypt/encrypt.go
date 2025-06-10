package encryptmanage

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	cost = 10
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
		return "", err
	}

	return string(hashedPassword), nil
}

func (e *Encrypter) CheckPassword(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
	if err != nil {
		return err
	}

	return nil
}
