package transport

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/F1zm0n/event-driven/inventory/internal/usecases"
	event_pb_inventoryv1 "github.com/F1zm0n/event-driven/inventory/protos/gen/inventory/v1"
)

type Inventory struct {
	event_pb_inventoryv1.UnimplementedInventoryServiceServer
	usecase usecases.InventoryUsecases
}

func (i Inventory) GetByID(
	ctx context.Context,
	req *event_pb_inventoryv1.GetByIDRequest,
) (*event_pb_inventoryv1.GetByIDResponse, error) {
	if err := verifyGetByID(req); err != nil {
		return nil, err
	}

	id, err := uuid.Parse(req.InventoryId)
	if err != nil {
		return nil, err
	}

	item, err := i.usecase.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &event_pb_inventoryv1.GetByIDResponse{
		InventoryId:  req.InventoryId,
		ProductName:  item.ProductName,
		ProductCount: item.ProductCount,
		BasePrice:    item.BasePrice,
		SalePrice:    item.SalePrice,
	}, nil
}

func verifyGetByID(req *event_pb_inventoryv1.GetByIDRequest) error {
	if req.InventoryId == "" {
		return fmt.Errorf("inventory_id shouldn't be empty")
	}
	return nil
}

func RegisterServer(gRPC *grpc.Server, service usecases.InventoryUsecases) {
	event_pb_inventoryv1.RegisterInventoryServiceServer(gRPC, &Inventory{
		usecase: service,
	})
}
