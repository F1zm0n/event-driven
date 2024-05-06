package usecases

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/F1zm0n/event-driven/auth/internal/dto"
	"github.com/F1zm0n/event-driven/auth/internal/entity"
	"github.com/F1zm0n/event-driven/auth/internal/repository"
)

type Customer struct {
	repo   repository.CustomerRepository
	logger *slog.Logger
	dur    time.Duration
}

type CustomerUsecases interface {
	LoginJWT(ctx context.Context, customer dto.CustomerDto) (token string, err error)
	// LoginOAuth()()

	Register(ctx context.Context, customer dto.CustomerDto) error
	GetByEmail(ctx context.Context, email string) (customer dto.CustomerDto, err error)
	GetByID(ctx context.Context, id uuid.UUID) (customer dto.CustomerDto, err error)

	// RegisterOAuth()error
}

func NewCustomerUsecases(
	repo repository.CustomerRepository,
	logger *slog.Logger,
	dur time.Duration,
) CustomerUsecases {
	return &Customer{
		repo,
		logger,
		dur,
	}
}

func (c Customer) Register(ctx context.Context, customer dto.CustomerDto) error {
	hashedPassword, err := entity.HashPassword(customer.Password)
	if err != nil {
		return err
	}

	customerEntity := entity.CustomerEntity{
		CustomerID: uuid.New(),
		Email:      customer.Email,
		Password:   hashedPassword,
	}

	if err := c.repo.Create(ctx, customerEntity); err != nil {
		return err
	}

	return nil
}

func (c Customer) LoginJWT(ctx context.Context, customer dto.CustomerDto) (string, error) {
	const op = "auth#usecases#customer.LoginJWT"
	l := c.logger.With(slog.String("op", op))
	customerRepo, err := c.repo.GetByEmail(ctx, customer.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("not found")
		}
		return "", err
	}

	if err = entity.ComparePassword(customerRepo.Password, customer.Password); err != nil {
		l.Error("error comparing password", slog.String("error", err.Error()))
		return "", err
	}

	tok, err := NewToken(os.Getenv("JWT_SECRET"), customer, c.dur)
	if err != nil {
		l.Error("error generating jwt token", slog.String("error", err.Error()))
		return "", err
	}

	return tok, nil
}

func (c Customer) GetByEmail(ctx context.Context, email string) (dto.CustomerDto, error) {
	customer, err := c.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.CustomerDto{}, fmt.Errorf("not found")
		}
		return dto.CustomerDto{}, err
	}
	customerDto := dto.CustomerDto{
		CustomerID: customer.CustomerID,
		Email:      email,
	}
	return customerDto, nil
}

func (c Customer) GetByID(ctx context.Context, id uuid.UUID) (dto.CustomerDto, error) {
	customer, err := c.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.CustomerDto{}, fmt.Errorf("not found")
		}
		return dto.CustomerDto{}, err
	}
	customerDto := dto.CustomerDto{
		CustomerID: customer.CustomerID,
		Email:      customer.Email,
	}
	return customerDto, nil
}
