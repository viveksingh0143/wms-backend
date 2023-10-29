package category

import (
	"star-wms/core/types"
)

type Filter struct {
	Query    string       `form:"query" db:"name,full_path" whereType:"like" binding:"omitempty,max=100"` // Filters by Name, FullPath, max length 100
	ID       uint         `form:"id" db:"id" binding:"omitempty,gt=0"`                                    // Filters by ID, should be greater than 0
	Name     string       `form:"name" db:"name" binding:"omitempty,max=100"`                             // Filters by Name, max length 100
	ParentID *uint        `form:"parent_id" db:"parent_id" binding:"omitempty"`                           // Filters by ParentID
	Status   types.Status `form:"status" db:"status" binding:"omitempty,gt=0"`                            // Filters by Status, should be greater than 0
}
