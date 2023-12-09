package store

import (
	"star-wms/core/types"
)

type Filter struct {
	Query                  string       `form:"query" db:"name,code,category_path" whereType:"like" binding:"omitempty,max=100"`
	ID                     uint         `form:"id" db:"id" binding:"omitempty,gt=0"`
	Name                   string       `form:"name" db:"name" binding:"omitempty,max=100"`
	Code                   string       `form:"code" db:"code" binding:"omitempty,max=100"`
	Status                 types.Status `form:"status" db:"status" binding:"omitempty,gt=0"`
	CategoryPathLike       string       `form:"category_path_like" db:"category_path" whereType:"like" binding:"omitempty,gt=0"`
	CategoryPathStartsWith string       `form:"category_path_startswith" db:"category_path" whereType:"startswith" binding:"omitempty,gt=0"`
	CategoryID             uint         `form:"category_id" db:"category_id" binding:"omitempty,gt=0"`
	PlantID                uint
}
