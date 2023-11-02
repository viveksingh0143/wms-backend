package models

import (
	"star-wms/core/common/models"
)

type JobOrderItem struct {
	models.MyModel
	JobOrderID uint     `gorm:"not null;index;constraint:OnDelete:CASCADE;"`
	ProductID  uint     `gorm:"not null;index;"`
	Product    *Product `gorm:"foreignKey:ProductID"`
	Quantity   float64  `gorm:"not null"`
}
