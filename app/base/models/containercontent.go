package models

import (
	adminModels "star-wms/app/admin/models"
	"star-wms/core/common/models"
)

type ContainerContent struct {
	models.MyModel
	Barcode     string            `gorm:"index;"`
	RMBatchID   *uint             `gorm:"index;"`
	ContainerID uint              `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Container   *Container        `gorm:"foreignKey:ContainerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProductID   uint              `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Product     *Product          `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Quantity    float64           `gorm:"column:quantity"`
	PlantID     uint              `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Plant       adminModels.Plant `gorm:"foreignKey:PlantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
