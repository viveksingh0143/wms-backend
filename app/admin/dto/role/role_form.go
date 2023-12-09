package role

import (
	"star-wms/app/admin/dto/ability"
	"star-wms/core/types"
)

type Form struct {
	ID        uint            `json:"id" binding:"-"`
	Name      string          `json:"name" validate:"required,min=4,max=100"`
	Status    types.Status    `json:"status" validate:"required,gt=0"`
	Abilities []*ability.Form `json:"abilities" validate:"required"`
}
