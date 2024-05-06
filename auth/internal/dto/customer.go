package dto

import (
	"fmt"

	"github.com/google/uuid"
)

type CustomerDto struct {
	CustomerID uuid.UUID `json:"customer_id,omitempty"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
}

func ParseUUID(id string) (uuid.UUID, error) {
	customerID, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error parsing uuid customerId")
	}
	return customerID, nil
}
