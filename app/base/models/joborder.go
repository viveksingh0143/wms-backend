package models

import (
	adminModels "star-wms/app/admin/models"
	"star-wms/core/common/models"
	"star-wms/core/types"
	"time"
)

type JobOrder struct {
	models.MyModel
	IssuedDate time.Time         `gorm:"not null;"`
	OrderNo    string            `gorm:"type:varchar(255);uniqueIndex;not null;column:order_no"`
	POCategory POCategory        `gorm:"type:enum('PRODUCTION','TRAILS','NPD','SAMPLES');not null;default:'PRODUCTION';column:po_category"`
	Status     types.Status      `gorm:"type:int;default:1"`
	CustomerID *uint             `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;column:customer_id"`
	Customer   *Customer         `gorm:"foreignKey:CustomerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Items      []*JobOrderItem   `gorm:"foreignKey:JobOrderID;constraint:OnDelete:CASCADE;"`
	PlantID    uint              `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Plant      adminModels.Plant `gorm:"foreignKey:PlantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
type POCategory string

const (
	PRODUCTION POCategory = "PRODUCTION"
	TRAILS     POCategory = "TRAILS"
	NPD        POCategory = "NPD"
	SAMPLES    POCategory = "SAMPLES"
)
