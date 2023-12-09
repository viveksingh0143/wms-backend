package permission

type Filter struct {
	Query      string `form:"query" db:"group,module_name" whereType:"like" binding:"omitempty,max=100"` // Filters by Group, ModuleName, max length 100
	ID         uint   `form:"id" db:"id" binding:"omitempty,gt=0"`                                       // Filters by ID, should be greater than 0
	GroupName  string `form:"group_name" db:"group_name" binding:"omitempty,max=100"`                    // Filters by Group, max length 100
	ModuleName string `form:"module_name" db:"module_name" binding:"omitempty,max=100"`                  // Filters by Module Name, max length 100
}
