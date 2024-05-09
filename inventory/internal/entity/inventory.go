package entity

import "github.com/google/uuid"

type InventoryEntity struct {
	InventoryID  uuid.UUID `db:"inventory_id"`
	ProductName  string    `db:"product_name"`
	ProductCount int32     `db:"product_count"`
	BasePrice    int       `db:"base_price"`
	SalePrice    int       `db:"sale_price"`
}

func PriceToFloat(price int) float32 {
	return float32(price) / 100
}

func PriceToInt(price float32) int {
	return int(price) * 100
}
