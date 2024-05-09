package usecases

import (
	"context"

	"github.com/google/uuid"

	"github.com/F1zm0n/event-driven/inventory/internal/dto"
	"github.com/F1zm0n/event-driven/inventory/internal/entity"
	"github.com/F1zm0n/event-driven/inventory/internal/repository"
)

type InventoryUsecases interface {
	Create(
		ctx context.Context,
		inventory dto.InventoryDto,
	) (dto.InventoryDto, error)
	GetByID(ctx context.Context, inventoryID uuid.UUID) (dto.InventoryDto, error)
	UpdateCount(
		ctx context.Context,
		inventoryID uuid.UUID,
		price float32,
	) error
}

type Inventory struct {
	repo repository.InventoryRepository
}

func NewInventoryUsecases(repo repository.InventoryRepository) InventoryUsecases {
	return &Inventory{
		repo: repo,
	}
}

func (i Inventory) Create(
	ctx context.Context,
	inventory dto.InventoryDto,
) (dto.InventoryDto, error) {
	inventoryRepo := entity.InventoryEntity{
		InventoryID:  uuid.New(),
		ProductName:  inventory.ProductName,
		ProductCount: inventory.ProductCount,
		BasePrice:    entity.PriceToInt(inventory.BasePrice),
		SalePrice:    entity.PriceToInt(inventory.SalePrice),
	}
	_, err := i.repo.Create(ctx, inventoryRepo)
	if err != nil {
		return dto.InventoryDto{}, err
	}
	return dto.InventoryDto{
		InventoryID:  inventoryRepo.InventoryID,
		ProductName:  inventoryRepo.ProductName,
		ProductCount: inventoryRepo.ProductCount,
		BasePrice:    entity.PriceToFloat(inventoryRepo.BasePrice),
		SalePrice:    entity.PriceToFloat(inventoryRepo.SalePrice),
	}, err
}

func (i Inventory) GetByID(ctx context.Context, inventoryID uuid.UUID) (dto.InventoryDto, error) {
	inventoryRepo, err := i.repo.GetByID(ctx, inventoryID)
	if err != nil {
		return dto.InventoryDto{}, err
	}

	return dto.InventoryDto{
		InventoryID:  inventoryRepo.InventoryID,
		ProductName:  inventoryRepo.ProductName,
		ProductCount: inventoryRepo.ProductCount,
		BasePrice:    entity.PriceToFloat(inventoryRepo.BasePrice),
		SalePrice:    entity.PriceToFloat(inventoryRepo.SalePrice),
	}, err
}

func (i Inventory) UpdateCount(
	ctx context.Context,
	inventoryID uuid.UUID,
	price float32,
) error {
	err := i.repo.UpdateCountByID(ctx, inventoryID, entity.PriceToInt(price))
	if err != nil {
		return err
	}
	return nil
}
