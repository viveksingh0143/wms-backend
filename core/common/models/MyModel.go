package models

import "gorm.io/gorm"

type MyModel struct {
	gorm.Model
	UpdatedBy string `gorm:"column:updated_by;type:varchar(100)"`
}
