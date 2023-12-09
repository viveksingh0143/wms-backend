package models

import (
	adminModels "star-wms/app/admin/models"
	baseModels "star-wms/app/base/models"
	"star-wms/core/common/models"
	"star-wms/core/types"
	"time"
)

type Batchlabel struct {
	models.MyModel
	BatchDate time.Time `gorm:"not null;"`
	BatchNo   string    `gorm:"type:varchar(255);uniqueIndex;not null;column:batch_no"`
	SoNumber  string    `gorm:"type:varchar(255)"`
	//POCategory      baseModels.POCategory    `gorm:"type:enum('PRODUCTION','TRAILS','NPD','SAMPLES');not null;default:'PRODUCTION';column:po_category"`
	//UnitType        baseModels.UnitType      `gorm:"type:enum('WEIGHT','PIECE','LIQUID','LENGTH');not null;default:'WEIGHT';column:unit_type"`
	//UnitValue       baseModels.UnitValue     `gorm:"type:enum('PC','GM','KG','MT','LT','YD','SM');column:unit_weight_type;default:'GM'"`
	POCategory      baseModels.POCategory    `gorm:"type:varchar(255);not null;default:'PRODUCTION';column:po_category"`
	UnitType        baseModels.UnitType      `gorm:"type:varchar(255);not null;default:'WEIGHT';column:unit_type"`
	UnitValue       baseModels.UnitValue     `gorm:"type:varchar(255);column:unit_weight_type;default:'GM'"`
	UnitWeight      float64                  `gorm:"column:unit_weight"`
	TargetQuantity  float64                  `gorm:"not null"`
	PackageQuantity float64                  `gorm:"not null"`
	Status          types.Status             `gorm:"type:int;default:1"`
	ProcessStatus   types.ProcessStatus      `gorm:"type:int;default:1"`
	CustomerID      uint                     `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;column:customer_id"`
	Customer        *baseModels.Customer     `gorm:"foreignKey:CustomerID;references:ID;"`
	ProductID       uint                     `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;column:product_id"`
	Product         *baseModels.Product      `gorm:"foreignKey:ProductID;references:ID;"`
	MachineID       uint                     `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;column:machine_id"`
	Machine         *baseModels.Machine      `gorm:"foreignKey:MachineID;references:ID;"`
	JoborderID      *uint                    `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Joborder        *baseModels.Joborder     `gorm:"foreignKey:JoborderID;references:ID;"`
	JoborderItemID  *uint                    `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	JoborderItem    *baseModels.JoborderItem `gorm:"foreignKey:JoborderItemID;references:ID;"`
	PlantID         uint                     `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Plant           adminModels.Plant        `gorm:"foreignKey:PlantID;references:ID;"`
	Stickers        []*Sticker               `gorm:"foreignKey:BatchlabelID;constraint:OnDelete:CASCADE;"`
}

type Sticker struct {
	models.MyModel
	Barcode        string              `gorm:"type:varchar(255);uniqueIndex;not null;"`
	PacketNo       string              `gorm:"type:varchar(255);not null;"`
	PrintCount     int32               `gorm:"type:int;default:0"`
	Shift          string              `gorm:"type:varchar(255);not null;"`
	ProductLine    string              `gorm:"type:varchar(255);not null;"`
	BatchNo        string              `gorm:"type:varchar(255);not null;"`
	MachineNo      string              `gorm:"type:varchar(255);not null;"`
	IsUsed         bool                `gorm:"default:false"`
	UnitWeightLine string              `gorm:"type:varchar(255);not null;"`
	QuantityLine   string              `gorm:"type:varchar(255);not null;"`
	Quantity       float64             `gorm:"not null"`
	Supervisor     string              `gorm:"type:varchar(255);not null;"`
	ProductID      uint                `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Product        *baseModels.Product `gorm:"foreignKey:ProductID;references:ID;"`
	BatchlabelID   uint                `gorm:"not null;index;constraint:OnDelete:CASCADE;"`
	Batchlabel     *Batchlabel         `gorm:"foreignKey:BatchlabelID;"`
	StickerItems   []*StickerItem      `gorm:"foreignKey:StickerID;constraint:OnDelete:CASCADE;"`
	PlantID        uint                `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Plant          adminModels.Plant   `gorm:"foreignKey:PlantID;references:ID;"`
}

type StickerItem struct {
	models.MyModel
	ProductID    uint                `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Product      *baseModels.Product `gorm:"foreignKey:ProductID;references:ID;"`
	Quantity     float64             `gorm:"not null"`
	BatchNo      string              `gorm:"type:varchar(255);not null;"`
	StickerID    uint                `gorm:"not null;index;constraint:OnDelete:CASCADE;"`
	Sticker      *Sticker            `gorm:"foreignKey:StickerID;"`
	BatchlabelID uint                `gorm:"index;constraint:OnDelete:CASCADE;"`
	Batchlabel   *Batchlabel         `gorm:"foreignKey:BatchlabelID;"`
}
