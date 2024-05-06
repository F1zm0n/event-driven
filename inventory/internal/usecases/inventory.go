package usecases

import (
	"log/slog"

	"github.com/F1zm0n/event-driven/inventory/internal/repository"
)

type InventoryUsecases interface{}

type Inventory struct {
	repo   repository.InventoryRepository
	logger *slog.Logger
}
