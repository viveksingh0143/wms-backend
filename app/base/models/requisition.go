package models

import (
	adminModels "star-wms/app/admin/models"
	"star-wms/core/common/models"
	"star-wms/core/types"
	"time"
)

type Requisition struct {
	models.MyModel
	IssuedDate time.Time          `gorm:"not null;"`
	OrderNo    string             `gorm:"type:varchar(255);uniqueIndex;not null;column:order_no"`
	Department string             `gorm:"type:varchar(255);"`
	Status     types.Status       `gorm:"type:int;default:1"`
	Approved   bool               `gorm:"type:tinyint;not null;default:0;column:approved"`
	StoreID    uint               `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Store      *Store             `gorm:"foreignKey:StoreID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PlantID    uint               `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Plant      adminModels.Plant  `gorm:"foreignKey:PlantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Items      []*RequisitionItem `gorm:"foreignKey:RequisitionID;constraint:OnDelete:CASCADE;"`
}
