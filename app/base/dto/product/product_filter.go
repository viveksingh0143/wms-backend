package product

import (
	"star-wms/core/types"
)

type Filter struct {
	Query       string       `form:"query" db:"name,product_type,code,cms_code" whereType:"like" binding:"omitempty,max=100"` // Filters by Name, ProductType, Code, CmsCode, max length 100
	ID          uint         `form:"id" db:"id" binding:"omitempty,gt=0"`                                                     // Filters by ID, should be greater than 0
	ProductType string       `form:"product_type" db:"product_type" binding:"omitempty,max=100"`                              // Filters by ProductType, max length 100
	Name        string       `form:"name" db:"name" binding:"omitempty,max=100"`                                              // Filters by Name, max length 100
	Slug        string       `form:"slug" db:"slug" binding:"omitempty,max=100"`                                              // Filters by Slug, max length 100
	Code        string       `form:"code" db:"code" binding:"omitempty,max=100"`                                              // Filters by Code, max length 100
	CmsCode     string       `form:"cms_code" db:"cms_code" binding:"omitempty,max=100"`                                      // Filters by CmsCode, max length 100
	CategoryID  uint         `form:"category_id" db:"catgory_id" binding:"omitempty,gt=0"`                                    // Filters by CategoryID
	Status      types.Status `form:"status" db:"status" binding:"omitempty,gt=0"`                                             // Filters by Status, should be greater than 0
}
