package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	PlainText string
	Hash      []byte
}

func (p *Password) Set(passwordPlainText string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwordPlainText), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.PlainText = passwordPlainText
	p.Hash = hash
	return nil
}

func (p *Password) Validate(plainTextPassword string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plainTextPassword)); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
