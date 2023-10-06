package plant

import (
	"star-wms/core/types"
)

type Form struct {
	ID     uint         `json:"id" binding:"-"`
	Code   string       `json:"code" validate:"required,len=10"`
	Name   string       `json:"name" validate:"required,min=4,max=100"`
	Status types.Status `json:"status" validate:"required,gt=0"`
}
