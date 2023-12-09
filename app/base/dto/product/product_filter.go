package product

import (
	"star-wms/core/types"
)

type Filter struct {
	Query                  string       `form:"query" db:"name,product_type,category_path,code,cms_code" whereType:"like" binding:"omitempty,max=100"`
	IDNotEqual             uint         `form:"id_ne" db:"id" whereType:"ne" binding:"omitempty,gt=0"`
	ID                     uint         `form:"id" db:"id" binding:"omitempty,gt=0"`
	ProductType            string       `form:"product_type" db:"product_type" binding:"omitempty,max=100"`
	ProductTypes           []string     `form:"product_types" db:"product_type" whereType:"in" binding:"omitempty,max=100"`
	ProductTypeNotEqual    string       `form:"product_type_ne" db:"product_type" whereType:"ne" binding:"omitempty,max=100"`
	NameLike               string       `form:"name_like" db:"name" whereType:"like" binding:"omitempty,max=100"`
	Name                   string       `form:"name" db:"name" binding:"omitempty,max=100"`
	Slug                   string       `form:"slug" db:"slug" binding:"omitempty,max=100"`
	Code                   string       `form:"code" db:"code" binding:"omitempty,max=100"`
	CodeLike               string       `form:"code" db:"code" whereType:"like" binding:"omitempty,max=100"`
	CmsCode                string       `form:"cms_code" db:"cms_code" binding:"omitempty,max=100"`
	CategoryID             uint         `form:"category_id" db:"category_id" binding:"omitempty,gt=0"`
	CategoryPathLike       string       `form:"category_path_like" db:"category_path" whereType:"like" binding:"omitempty,gt=0"`
	CategoryPathStartsWith string       `form:"category_path_startswith" db:"category_path" whereType:"startswith" binding:"omitempty,gt=0"`
	Status                 types.Status `form:"status" db:"status" binding:"omitempty,gt=0"`
}
