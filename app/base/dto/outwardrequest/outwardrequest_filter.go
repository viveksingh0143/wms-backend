package outwardrequest

import (
	"star-wms/core/types"
)

type Filter struct {
	Query      string       `form:"query" db:"order_no" whereType:"like" binding:"omitempty,max=100"`
	ID         uint         `form:"id" db:"id" binding:"omitempty,gt=0"`
	CustomerID uint         `form:"customer_id" db:"customer_id" binding:"omitempty,gt=0"`
	Status     types.Status `form:"status" db:"status" binding:"omitempty,gt=0"`
}
