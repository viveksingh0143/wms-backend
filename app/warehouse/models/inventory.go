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
	Store     *baseModels.Store   `gorm:"foreignKey:StoreID;references:ID;"`
	ProductID uint                `gorm:"index;not null;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Product   *baseModels.Product `gorm:"foreignKey:ProductID;references:ID;"`
	Quantity  float64             `gorm:"column:quantity"`
	PlantID   uint                `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Plant     adminModels.Plant   `gorm:"foreignKey:PlantID;references:ID;"`
}

type StockMovements struct {
	models.MyModel
	StoreID      *uint                 `gorm:"index;not null;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Store        *baseModels.Store     `gorm:"foreignKey:StoreID;references:ID;"`
	ProductID    uint                  `gorm:"index;not null;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Product      *baseModels.Product   `gorm:"foreignKey:ProductID;references:ID;"`
	Quantity     float64               `gorm:"column:quantity"`
	JoborderID   *uint                 `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Joborder     *baseModels.Joborder  `gorm:"foreignKey:JoborderID;references:ID;"`
	ContainerID  *uint                 `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Container    *baseModels.Container `gorm:"foreignKey:ContainerID;references:ID;"`
	BatchlabelID *uint                 `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Batchlabel   *Batchlabel           `gorm:"foreignKey:BatchlabelID;references:ID;"`
	StickerID    *uint                 `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Sticker      *Sticker              `gorm:"foreignKey:StickerID;references:ID;"`
	Status       types.InventoryStatus `gorm:"type:int;default:1"`
	PlantID      uint                  `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Plant        adminModels.Plant     `gorm:"foreignKey:PlantID;references:ID;"`
}
