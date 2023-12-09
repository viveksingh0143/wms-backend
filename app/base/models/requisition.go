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
	Approved   types.Approval     `gorm:"type:int;default:3"`
	StoreID    uint               `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Store      *Store             `gorm:"foreignKey:StoreID;references:ID;"`
	PlantID    uint               `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Plant      adminModels.Plant  `gorm:"foreignKey:PlantID;references:ID;"`
	Items      []*RequisitionItem `gorm:"foreignKey:RequisitionID;constraint:OnDelete:CASCADE;"`
}
