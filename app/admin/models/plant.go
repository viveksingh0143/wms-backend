package models

import (
	"star-wms/core/common/models"
	"star-wms/core/types"
)

type Plant struct {
	models.MyModel
	Code   string       `gorm:"type:varchar(10);uniqueIndex;not null"`
	Name   string       `gorm:"type:varchar(100);uniqueIndex;not null"`
	Status types.Status `gorm:"type:int;default:1"`
	Users  []*User      `gorm:"foreignKey:PlantID"`
}
