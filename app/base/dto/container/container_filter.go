package container

import (
	"star-wms/core/types"
)

type Filter struct {
	Query         string       `form:"query" db:"name,code" whereType:"like" binding:"omitempty,max=100"`
	ID            uint         `form:"id" db:"id" binding:"omitempty,gt=0"`
	ContainerType string       `form:"container_type" db:"container_type" binding:"omitempty,max=100"`
	Name          string       `form:"name" db:"name" binding:"omitempty,max=100"`
	Code          string       `form:"code" db:"code" binding:"omitempty,max=100"`
	CategoryID    uint         `form:"category_id" db:"category_id" binding:"omitempty,gt=0"`
	Status        types.Status `form:"status" db:"status" binding:"omitempty,gt=0"`
	StockLevel    string       `form:"stock_level" db:"stock_level" binding:"omitempty,gt=0"`
	Approved      bool         `form:"approved" db:"approved" binding:"omitempty"`
	ProductID     uint         `form:"product_id" db:"product_id" binding:"omitempty,gt=0"`
	StoreID       uint         `form:"store_id" db:"store_id" binding:"omitempty,gt=0"`
	StoreIDsIn    []uint       `form:"store_ids_in" db:"store_id" whereType:"IN" binding:"omitempty,gt=0"`
}
