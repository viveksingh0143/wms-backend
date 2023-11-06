package models

import (
	adminModels "star-wms/app/admin/models"
	"star-wms/core/common/models"
	"star-wms/core/types"
	"time"
)

type Outwardrequest struct {
	models.MyModel
	IssuedDate time.Time             `gorm:"not null;"`
	OrderNo    string                `gorm:"type:varchar(255);uniqueIndex;not null;column:order_no"`
	Status     types.Status          `gorm:"type:int;default:1"`
	CustomerID uint                  `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Customer   *Customer             `gorm:"foreignKey:CustomerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PlantID    uint                  `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Plant      adminModels.Plant     `gorm:"foreignKey:PlantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Items      []*OutwardrequestItem `gorm:"foreignKey:OutwardrequestID;constraint:OnDelete:CASCADE;"`
}
