package models

import (
	adminModels "star-wms/app/admin/models"
	"star-wms/core/common/models"
	"star-wms/core/types"
)

type Store struct {
	models.MyModel
	Name         string              `gorm:"type:varchar(255);uniqueIndex;not null;column:name"`
	Code         string              `gorm:"type:varchar(255);uniqueIndex;not null;column:code"`
	Address      string              `gorm:"type:varchar(255);column:address"`
	Status       types.Status        `gorm:"type:int;default:1"`
	CategoryPath string              `gorm:"index;"`
	CategoryID   *uint               `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Category     *Category           `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Approvers    []*adminModels.User `gorm:"many2many:store_approvers;constraint:OnDelete:CASCADE;"`
	PlantID      uint                `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Plant        adminModels.Plant   `gorm:"foreignKey:PlantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
