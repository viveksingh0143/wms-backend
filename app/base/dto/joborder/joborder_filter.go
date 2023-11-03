package joborder

import (
	"star-wms/core/types"
)

type Filter struct {
	Query         string              `form:"query" db:"order_no,po_category" whereType:"like" binding:"omitempty,max=100"` // Filters by Name, ProductType, Code, CmsCode, max length 100
	ID            uint                `form:"id" db:"id" binding:"omitempty,gt=0"`                                          // Filters by ID, should be greater than 0
	POCategory    string              `form:"po_category" db:"po_category" binding:"omitempty,max=100"`                     // Filters by ProductType, max length 100
	CustomerID    uint                `form:"customer_id" db:"customer_id" binding:"omitempty,gt=0"`                        // Filters by CategoryID
	Status        types.Status        `form:"status" db:"status" binding:"omitempty,gt=0"`                                  // Filters by Status, should be greater than 0
	ProcessStatus types.ProcessStatus `form:"process_status" db:"process_status" validate:"omitempty,gt=0"`
}
