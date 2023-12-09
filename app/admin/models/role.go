package models

import (
	"star-wms/core/common/models"
	"star-wms/core/types"
)

type Role struct {
	models.MyModel
	Name      string       `gorm:"type:varchar(100);uniqueIndex;not null"`
	Status    types.Status `gorm:"type:int;default:1"`
	Abilities []*Ability   `gorm:"many2many:role_abilities;constraint:OnDelete:CASCADE;"`
}
