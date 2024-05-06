package repository

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/F1zm0n/event-driven/auth/internal/entity"
)

type Customer struct {
	logger *slog.Logger
	db     *sqlx.DB
}

type CustomerRepository interface {
	Create(ctx context.Context, customer entity.CustomerEntity) error
	GetByID(ctx context.Context, customerID uuid.UUID) (customer entity.CustomerEntity, err error)
	GetByEmail(ctx context.Context, email string) (customer entity.CustomerEntity, err error)
}

func NewCustomerRepository(logger *slog.Logger, db *sqlx.DB) CustomerRepository {
	return &Customer{
		logger,
		db,
	}
}

func (r Customer) Create(ctx context.Context, customer entity.CustomerEntity) error {
	const op = "auth#repository.Create"
	l := r.logger.With(slog.String("op", op))
	l.Info("creating customer")
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO customer (customer_id, email, password) VALUES($1, $2, $3)`,
		customer.CustomerID,
		customer.Email,
		customer.Password,
	)
	if err != nil {
		l.Error("error creating customer", slog.String("error", err.Error()))
		return err
	}

	l.Info("created customer successfully")

	return nil
}

func (r Customer) GetByID(
	ctx context.Context,
	customerID uuid.UUID,
) (entity.CustomerEntity, error) {
	const op = "auth#repository.GetByID"
	l := r.logger.With(slog.String("op", op))
	l.Info("getting customer by id")
	var customer entity.CustomerEntity
	if err := r.db.GetContext(ctx, &customer, "SELECT * FROM customer WHERE customer_id = $1", customerID); err != nil {

		l.Error("error getting customer by id", slog.String("error", err.Error()))
		return entity.CustomerEntity{}, err
	}

	l.Info("successfully got customer by id")

	return customer, nil
}

func (r Customer) GetByEmail(ctx context.Context, email string) (entity.CustomerEntity, error) {
	const op = "auth#repository.GetByEmail"
	l := r.logger.With(slog.String("op", op))
	l.Info("getting customer by email")
	var customer entity.CustomerEntity
	if err := r.db.GetContext(ctx, &customer, "SELECT * FROM customer WHERE email = $1", email); err != nil {
		l.Error("error getting customer by email", slog.String("error", err.Error()))
		return entity.CustomerEntity{}, err
	}

	l.Info("successfully got customer by email")

	return customer, nil
}
