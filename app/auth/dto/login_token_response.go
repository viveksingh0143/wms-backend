package dto

import (
	"star-wms/app/admin/dto/ability"
	"star-wms/app/admin/dto/plant"
	"star-wms/app/admin/dto/role"
)

type LoginTokenResponse struct {
	AccessToken  string          `json:"access_token"`
	RefreshToken string          `json:"refresh_token"`
	Name         string          `json:"name"`
	StaffID      string          `json:"staff_id"`
	Roles        []*role.Form    `json:"roles"`
	Abilities    []*ability.Form `json:"permissions"`
	Plant        *plant.Form     `json:"plant"`
}
