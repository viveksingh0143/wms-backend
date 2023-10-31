package store

import (
	"star-wms/app/admin/dto/user"
	"star-wms/app/base/dto/category"
	"star-wms/core/types"
)

type Form struct {
	PlantID   uint           `json:"plant_id" binding:"-"`
	ID        uint           `json:"id" binding:"-"`
	Name      string         `json:"name" validate:"required,min=4,max=100"`
	Code      string         `json:"code" validate:"required,min=4,max=100"`
	Address   string         `json:"address" validate:"omitempty,min=4,max=100"`
	Status    types.Status   `json:"status" validate:"required,gt=0"`
	Approvers []*user.Form   `json:"approvers"`
	Category  *category.Form `json:"category" validationTag:"category.id" validate:"omitempty,validRelationID,structonly"`
}
