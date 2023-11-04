package storelocation

import (
	"star-wms/app/base/dto/store"
	"star-wms/core/types"
)

type Form struct {
	PlantID     uint             `json:"plant_id" binding:"-"`
	StoreID     uint             `json:"store_id" binding:"-"`
	Store       *store.Form      `json:"store" binding:"-"`
	ID          uint             `json:"id" binding:"-"`
	Code        string           `json:"code" validate:"required,min=4,max=100"`
	ZoneName    string           `json:"zone_name" validate:"omitempty,min=1,max=100"`
	AisleNumber string           `json:"aisle_number" validate:"omitempty,min=1,max=100"`
	RackNumber  string           `json:"rack_number" validate:"omitempty,min=1,max=100"`
	ShelfNumber string           `json:"shelf_number" validate:"omitempty,min=1,max=100"`
	Description string           `json:"description" validate:"omitempty,min=4,max=100"`
	Status      types.FillStatus `json:"status" binding:"-"`
}
