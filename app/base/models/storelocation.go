package models

import (
	adminModels "star-wms/app/admin/models"
	"star-wms/core/common/models"
	"star-wms/core/types"
)

type Storelocation struct {
	models.MyModel
	Code        string            `gorm:"type:varchar(255);uniqueIndex;"`
	ZoneName    string            `gorm:"type:varchar(255);"`
	AisleNumber string            `gorm:"type:varchar(255);"`
	RackNumber  string            `gorm:"type:varchar(255);"`
	ShelfNumber string            `gorm:"type:varchar(255);"`
	Description string            `gorm:"type:varchar(255);"`
	Status      types.FillStatus  `gorm:"type:int;default:1"`
	StoreID     uint              `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Store       *Store            `gorm:"foreignKey:StoreID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PlantID     uint              `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Plant       adminModels.Plant `gorm:"foreignKey:PlantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
