package user

import (
	"star-wms/core/types"
)

type Filter struct {
	Query    string       `form:"query" db:"name,staff_id,username,email" whereType:"like" binding:"omitempty,max=100"` // Filters by Name, StaffID, Username, EMail, max length 100
	ID       uint         `form:"id" db:"id" binding:"omitempty,gt=0"`                                                  // Filters by ID, should be greater than 0
	Name     string       `form:"name" db:"name" binding:"omitempty,max=100"`                                           // Filters by Name, max length 100
	StaffID  string       `form:"staff_id" db:"staff_id" binding:"omitempty,max=100"`                                   // Filters by StaffID, max length 100
	Username string       `form:"username" db:"username" binding:"omitempty,max=100"`                                   // Filters by Username, max length 100
	EMail    string       `form:"email" db:"email" binding:"omitempty,max=100"`                                         // Filters by EMail, max length 100
	PlantID  *uint        `form:"plant_id" db:"plant_id" binding:"omitempty"`                                           // Filters by PlantID
	Status   types.Status `form:"status" db:"status" binding:"omitempty,gt=0"`                                          // Filters by Status, should be greater than 0
}
