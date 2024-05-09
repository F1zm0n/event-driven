package dto

import "github.com/google/uuid"

type InventoryDto struct {
	InventoryID  uuid.UUID `json:"inventory_id,omitempty"`
	ProductName  string    `json:"product_name"`
	ProductCount int32     `json:"product_count"`
	BasePrice    float32   `json:"base_price"`
	SalePrice    float32   `json:"sale_price,omitempty"`
}
