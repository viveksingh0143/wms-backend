package models

import (
	adminModels "star-wms/app/admin/models"
	baseModels "star-wms/app/base/models"
	"star-wms/core/common/models"
	"star-wms/core/types"
)

type Inventory struct {
	models.MyModel
	StoreID   *uint               `gorm:"index;not null;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Store     *baseModels.Store   `gorm:"foreignKey:StoreID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	ProductID uint                `gorm:"index;not null;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Product   *baseModels.Product `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Quantity  float64             `gorm:"column:quantity"`
	PlantID   uint                `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Plant     adminModels.Plant   `gorm:"foreignKey:PlantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type InventoryTransaction struct {
	models.MyModel
	StoreID      *uint                 `gorm:"index;not null;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Store        *baseModels.Store     `gorm:"foreignKey:StoreID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	ProductID    uint                  `gorm:"index;not null;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Product      *baseModels.Product   `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Quantity     float64               `gorm:"column:quantity"`
	JoborderID   *uint                 `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Joborder     *baseModels.Joborder  `gorm:"foreignKey:JoborderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ContainerID  *uint                 `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Container    *baseModels.Container `gorm:"foreignKey:ContainerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	BatchlabelID *uint                 `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Batchlabel   *Batchlabel           `gorm:"foreignKey:BatchlabelID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	StickerID    *uint                 `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Sticker      *Sticker              `gorm:"foreignKey:StickerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Status       types.InventoryStatus `gorm:"type:int;default:1"`
	PlantID      uint                  `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Plant        adminModels.Plant     `gorm:"foreignKey:PlantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
