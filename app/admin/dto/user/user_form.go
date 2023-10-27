package user

import (
	"star-wms/app/admin/dto/plant"
	"star-wms/app/admin/dto/role"
	"star-wms/core/types"
)

type Form struct {
	ID       uint         `json:"id" binding:"-"`
	Name     string       `json:"name" validate:"required,min=4,max=100"`
	StaffID  string       `json:"staff_id" validate:"required,min=4,max=100"`
	Username string       `json:"username" validate:"required,min=4,max=100"`
	EMail    string       `json:"email" validate:"required,email,min=4,max=100"`
	Password string       `json:"password" validate:"required,min=4,max=100"`
	Status   types.Status `json:"status" validate:"required,gt=0"`
	Roles    []*role.Form `json:"roles"`
	Plant    *plant.Form  `json:"plant" validationTag:"plant.id" validate:"omitempty,validRelationID,structonly"`
}
