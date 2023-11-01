package machine

import (
	"star-wms/core/types"
)

type Form struct {
	PlantID uint         `json:"plant_id" binding:"-"`
	ID      uint         `json:"id" binding:"-"`
	Name    string       `json:"name" validate:"required,min=4,max=100"`
	Code    string       `json:"code" validate:"required,min=4,max=100"`
	Status  types.Status `json:"status" validate:"required,gt=0"`
}
