package container

import (
	"star-wms/app/base/dto/product"
	"star-wms/app/base/dto/store"
	"star-wms/app/base/models"
	"star-wms/core/types"
)

type Form struct {
	PlantID       uint              `json:"plant_id" binding:"-"`
	ID            uint              `json:"id" binding:"-"`
	ContainerType string            `json:"container_type" validate:"required,oneof='PALLET' 'BIN'"`
	Name          string            `json:"name" validate:"required,min=4,max=100"`
	Code          string            `json:"code" validate:"required,min=4,max=100"`
	Address       string            `json:"address" validate:"omitempty,min=4,max=400"`
	Status        types.Status      `json:"status" validate:"required,gt=0"`
	StockLevel    models.StockLevel `json:"stock_level" binding:"-"`
	Approved      bool              `json:"approved" binding:"-"`
	Product       *product.Form     `json:"product" binding:"-"`
	Store         *store.Form       `json:"store" binding:"-"`
}
