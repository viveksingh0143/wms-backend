package models

import (
	adminModels "star-wms/app/admin/models"
	"star-wms/core/common/models"
	"star-wms/core/types"
)

type Machine struct {
	models.MyModel
	Name    string            `gorm:"type:varchar(255);not null;column:name"`
	Code    string            `gorm:"type:varchar(255);uniqueIndex;not null;column:code"`
	Status  types.Status      `gorm:"type:int;default:1"`
	PlantID uint              `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Plant   adminModels.Plant `gorm:"foreignKey:PlantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
