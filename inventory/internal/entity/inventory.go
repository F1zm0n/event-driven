package entity

import "github.com/google/uuid"

type InventoryEntity struct {
	InventoryID  uuid.UUID `db:"inventory_id"`
	ProductName  string    `db:"product_name"`
	ProductCount int       `db:"product_count"`
	BasePrice    int       `db:"base_price"`
	SalePrice    int       `db:"sale_price"`
}
