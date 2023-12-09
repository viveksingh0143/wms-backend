package models

import (
	"star-wms/core/common/models"
)

type Ability struct {
	models.MyModel
	Name   string `gorm:"uniqueIndex:idx_name_module;type:varchar(255);not null" json:"name"`
	Module string `gorm:"uniqueIndex:idx_name_module;type:varchar(255);not null" json:"module"`
}
