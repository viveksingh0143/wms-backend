package models

import (
	adminModels "star-wms/app/admin/models"
	"star-wms/core/common/models"
	"star-wms/core/types"
)

type Container struct {
	models.MyModel
	ContainerType   ContainerType       `gorm:"type:enum('PALLET','BIN');not null;default:'PALLET';column:container_type"`
	Code            string              `gorm:"type:varchar(255);uniqueIndex;column:code"`
	Name            string              `gorm:"type:varchar(255);not null;column:name"`
	Address         string              `gorm:"type:varchar(255);column:address"`
	Status          types.Status        `gorm:"type:int;default:1"`
	StockLevel      StockLevel          `gorm:"type:enum('EMPTY','PARTIAL','FULL');not null;default:'EMPTY';column:stock_level"`
	Approved        types.Approval      `gorm:"type:int;default:3"`
	ProductID       *uint               `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Product         *Product            `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	StorelocationID *uint               `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Storelocation   *Storelocation      `gorm:"foreignKey:StorelocationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	StoreID         *uint               `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Store           *Store              `gorm:"foreignKey:StoreID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PlantID         uint                `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Plant           adminModels.Plant   `gorm:"foreignKey:PlantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Contents        []*ContainerContent `gorm:"foreignKey:ContainerID;constraint:OnDelete:CASCADE;"`
}

type ContainerType string
type StockLevel string

const (
	Pallet ContainerType = "PALLET"
	Bin    ContainerType = "BIN"
)

const (
	Empty   StockLevel = "EMPTY"
	Partial StockLevel = "PARTIAL"
	Full    StockLevel = "FULL"
)
