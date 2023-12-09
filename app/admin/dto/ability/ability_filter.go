package ability

type Filter struct {
	Query  string `form:"query" db:"name" whereType:"like" binding:"omitempty,max=100"` // Filters by Group, ModuleName, max length 100
	ID     uint   `form:"id" db:"id" binding:"omitempty,gt=0"`                          // Filters by ID, should be greater than 0
	Name   string `form:"name" db:"name" binding:"omitempty,max=100"`                   // Filters by Name, max length 100
	Module string `form:"module" db:"module" binding:"omitempty,max=100"`               // Filters by Module, max length 100
}
