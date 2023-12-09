package models

import (
	adminModels "star-wms/app/admin/models"
	baseModels "star-wms/app/base/models"
	"star-wms/core/common/models"
	"star-wms/core/types"
)

type RMBatch struct {
	models.MyModel
	ProductID    uint                  `gorm:"index;not null;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Product      *baseModels.Product   `gorm:"foreignKey:ProductID;references:ID;"`
	BatchNumber  string                `gorm:"type:varchar(255)"`
	Quantity     float64               `gorm:"column:quantity"`
	Unit         string                `gorm:"type:varchar(255)"`
	ContainerID  *uint                 `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Container    *baseModels.Container `gorm:"foreignKey:ContainerID;references:ID;"`
	StoreID      *uint                 `gorm:"index;not null;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Store        *baseModels.Store     `gorm:"foreignKey:StoreID;references:ID;"`
	PlantID      uint                  `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Plant        adminModels.Plant     `gorm:"foreignKey:PlantID;references:ID;"`
	Status       types.InventoryStatus `gorm:"type:int;default:1"`
	Transactions []*RMBatchTransaction `gorm:"foreignKey:RMBatchID;constraint:OnDelete:CASCADE;"`
}

type RMBatchTransaction struct {
	models.MyModel
	RMBatchID       uint               `gorm:"not null;index"`
	ProductID       uint               `gorm:"not null;index;"`
	TransactionType string             `gorm:"not null;type:varchar(20);"`
	Quantity        float64            `gorm:"not null;"`
	Notes           string             `gorm:"type:text;"`
	RMBatch         RMBatch            `gorm:"foreignKey:RMBatchID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Product         baseModels.Product `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PlantID         uint               `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Plant           adminModels.Plant  `gorm:"foreignKey:PlantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (RMBatchTransaction) TableName() string {
	return "rmbatch_transactions"
}

func (r *RMBatch) NewTransactionHistory() *RMBatchTransaction {
	return &RMBatchTransaction{
		RMBatchID:       r.ID,
		ProductID:       r.ProductID,
		TransactionType: r.Status.String(),
		Quantity:        r.Quantity,
		Notes:           "",
		PlantID:         r.PlantID,
	}
}
