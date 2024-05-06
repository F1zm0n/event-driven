package repository

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/F1zm0n/event-driven/inventory/internal/entity"
)

type Inventory struct {
	logger *slog.Logger
	db     *sqlx.DB
}

type InventoryRepository interface {
	Create(ctx context.Context, Inventory entity.InventoryEntity) (entity.InventoryEntity, error)
	GetByID(ctx context.Context, id uuid.UUID) (entity.InventoryEntity, error)
	UpdateCountByID(
		ctx context.Context,
		productID uuid.UUID,
		count int,
	) error
}

func NewInventoryRepository(logger *slog.Logger, db *sqlx.DB) InventoryRepository {
	return &Inventory{
		logger: logger,
		db:     db,
	}
}

func (i Inventory) Create(
	ctx context.Context,
	inventory entity.InventoryEntity,
) (entity.InventoryEntity, error) {
	const op = "inventory#repository#inventory.Create"
	l := i.logger.With(slog.String("op", op))
	l.InfoContext(ctx, "creating inventory item")
	_, err := i.db.ExecContext(
		ctx,
		`INSERT INTO inventory (inventory_id, product_name, product_count, base_price, sale_price)
		VALUES ($1, $2, $3, $4, $5)
	`,
		inventory.InventoryID,
		inventory.ProductName,
		inventory.ProductCount,
		inventory.BasePrice,
		inventory.SalePrice,
	)
	if err != nil {
		l.InfoContext(ctx, "error creating inventory item", slog.String("error", err.Error()))
		return entity.InventoryEntity{}, err
	}
	l.InfoContext(ctx, "successfully created inventory item")
	return inventory, nil
}

func (i Inventory) GetByID(ctx context.Context, id uuid.UUID) (entity.InventoryEntity, error) {
	const op = "inventory#repository#inventory.GetByID"
	l := i.logger.With(slog.String("op", op))

	var inventory entity.InventoryEntity
	l.InfoContext(ctx, "getting inventory item by id")
	err := i.db.GetContext(ctx, &inventory, "SELECT * FROM inventory WHERE product_id=$1", id)
	if err != nil {
		l.ErrorContext(ctx, "error getting inventory item by id", slog.String("error", err.Error()))
		return entity.InventoryEntity{}, err
	}

	l.InfoContext(ctx, "successfully got inventory item by id")
	return inventory, nil
}

func (i Inventory) UpdateCountByID(
	ctx context.Context,
	productID uuid.UUID,
	count int,
) error {
	const op = "inventory#repository#inventory.UpdateCount"
	l := i.logger.With(slog.String("op", op))

	l.Info("updating product_count by product_id")

	_, err := i.db.ExecContext(
		ctx,
		"UPDATE inventory SET product_count=$1 WHERE product_id=$2",
		count,
	)
	if err != nil {
		l.Error("error updating product_count by product_id", slog.String("error", err.Error()))
		return err
	}

	l.Info("successfully updated product_count by product_id")

	return nil
}
