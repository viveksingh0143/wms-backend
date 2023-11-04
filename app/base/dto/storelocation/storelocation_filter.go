package storelocation

import (
	"star-wms/core/types"
)

type Filter struct {
	Query           string           `form:"query" db:"code,zone_name,aisle_number,rack_number,shelf_number" whereType:"like" binding:"omitempty,max=100"`
	ID              uint             `form:"id" db:"id" binding:"omitempty,gt=0"`
	StoreID         uint             `form:"store_id" db:"store_id" binding:"omitempty,max=100"`
	CodeLike        string           `form:"code_like" db:"code" whereType:"like" binding:"omitempty,max=100"`
	ZoneNameLike    string           `form:"zone_name_like" db:"zone_name" whereType:"like" binding:"omitempty,max=100"`
	AisleNumberLike string           `form:"aisle_number_like" db:"aisle_number" whereType:"like" binding:"omitempty,max=100"`
	RackNumberLike  string           `form:"rack_number_like" db:"rack_number" whereType:"like" binding:"omitempty,max=100"`
	ShelfNumberLike string           `form:"shelf_number_like" db:"shelf_number" whereType:"like" binding:"omitempty,max=100"`
	Status          types.FillStatus `form:"status" db:"status" binding:"omitempty,max=100"`
}
