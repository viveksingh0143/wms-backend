package rmbatch

import (
	"star-wms/app/base/dto/container"
	"star-wms/app/base/dto/product"
	"star-wms/app/base/dto/store"
	"star-wms/core/types"
	"time"
)

type Form struct {
	ID           uint                  `json:"id" binding:"-"`
	BatchNumber  string                `json:"batch_number"`
	Quantity     float64               `json:"quantity"`
	Unit         string                `json:"unit"`
	Container    *container.Form       `json:"container"`
	Store        *store.Form           `json:"store"`
	Product      *product.Form         `json:"product"`
	Status       types.InventoryStatus `json:"status"`
	Transactions []*Transaction        `json:"transactions"`
}

type Transaction struct {
	ID              uint      `json:"id" binding:"-"`
	TransactionType string    `json:"transaction_type"`
	Quantity        float64   `json:"quantity"`
	Notes           string    `json:"notes"`
	RMBatchID       uint      `json:"rm_batch_id"`
	ProductID       uint      `json:"product_id"`
	CreatedAt       time.Time `json:"created_at"`
}
