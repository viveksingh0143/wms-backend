package models

import (
	"star-wms/core/common/models"
)

type Permission struct {
	models.MyModel
	Group      string `gorm:"type:varchar(100);not null" json:"group"`
	ModuleName string `gorm:"type:varchar(100);uniqueIndex;not null" json:"module_name"`
	ReadPerm   bool   `gorm:"default:false" json:"read_perm"`
	CreatePerm bool   `gorm:"default:false" json:"create_perm"`
	UpdatePerm bool   `gorm:"default:false" json:"update_perm"`
	DeletePerm bool   `gorm:"default:false" json:"delete_perm"`
	ImportPerm bool   `gorm:"default:false" json:"import_perm"`
	ExportPerm bool   `gorm:"default:false" json:"export_perm"`
}

func (receiver *Permission) ReadPermName() string {
	return "READ"
}
func (receiver *Permission) CreatePermName() string {
	return "CREATE"
}
func (receiver *Permission) UpdatePermName() string {
	return "UPDATE"
}
func (receiver *Permission) DeletePermName() string {
	return "DELETE"
}
func (receiver *Permission) ImportPermName() string {
	return "IMPORT"
}
func (receiver *Permission) ExportPermName() string {
	return "EXPORT"
}
