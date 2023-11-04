package category

import (
	"star-wms/core/types"
)

type Filter struct {
	Query        string       `form:"query" db:"name,full_path" whereType:"like" binding:"omitempty,max=100"`
	ID           uint         `form:"id" db:"id" binding:"omitempty,gt=0"`
	NameLike     string       `form:"name_like" db:"name" whereType:"like" binding:"omitempty,max=100"`
	ParentID     *uint        `form:"parent_id" db:"parent_id" binding:"omitempty"`
	Status       types.Status `form:"status" db:"status" binding:"omitempty,gt=0"`
	CreatedAtGte string       `form:"createdAt_gte" db:"created_at" whereType:"gte" binding:"omitempty,max=100"`
	CreatedAtLte string       `form:"createdAt_lte" db:"created_at" whereType:"lte" binding:"omitempty,max=100"`
}
