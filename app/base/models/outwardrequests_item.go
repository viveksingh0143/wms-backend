package models

import (
	"star-wms/core/common/models"
)

type OutwardrequestItem struct {
	models.MyModel
	OutwardrequestID uint     `gorm:"not null;index;constraint:OnDelete:CASCADE;"`
	ProductID        uint     `gorm:"not null;index;"`
	Product          *Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	Quantity         float64  `gorm:"not null"`
}
