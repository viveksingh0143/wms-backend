package models

import (
	"star-wms/core/common/models"
	"star-wms/core/types"
)

type User struct {
	models.MyModel
	Name     string       `gorm:"type:varchar(100);not null"`
	StaffID  string       `gorm:"type:varchar(100);uniqueIndex;not null"`
	Username string       `gorm:"type:varchar(100);uniqueIndex;not null"`
	EMail    string       `gorm:"column:email;type:varchar(100);uniqueIndex;not null"`
	Password string       `gorm:"type:varchar(100);not null"`
	Status   types.Status `gorm:"type:int;default:1"`
	PlantID  *uint        `gorm:"foreignKey:PlantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Roles    []*Role      `gorm:"many2many:user_roles;"`
	Plant    *Plant       `gorm:"foreignKey:PlantID"`
}
