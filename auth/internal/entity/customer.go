package entity

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type CustomerEntity struct {
	CustomerID uuid.UUID `db:"customer_id"`
	Email      string    `db:"email"`
	Password   []byte    `db:"password"`
}

func HashPassword(password string) ([]byte, error) {
	pssw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return pssw, nil
}

func ComparePassword(hashed []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashed, []byte(password))
}
