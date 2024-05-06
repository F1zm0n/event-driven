package usecases

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/F1zm0n/event-driven/auth/internal/dto"
)

func NewToken(secret string, user dto.CustomerDto, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.CustomerID.String()
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string, secret string) (dto.CustomerDto, error) {
	var customer dto.CustomerDto

	// Define the JWT parsing key and options
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, keyFunc)
	if err != nil {
		return customer, err
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		// Extract customer information from claims
		customerID, okUID := (*claims)["uid"].(string)
		email, okEmail := (*claims)["email"].(string)

		if !okUID || !okEmail {
			return dto.CustomerDto{}, fmt.Errorf("invalid claims data")
		}
		customerUUID, err := uuid.Parse(customerID)
		if err != nil {
			return dto.CustomerDto{}, fmt.Errorf("invalid uuid format")
		}
		customer.CustomerID = customerUUID
		customer.Email = email
		return customer, nil
	}

	return customer, errors.New("invalid token")
}
