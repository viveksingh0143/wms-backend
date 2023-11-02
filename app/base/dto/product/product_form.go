package product

import (
	"star-wms/app/base/dto/category"
	"star-wms/core/types"
)

type Form struct {
	ID          uint              `json:"id" binding:"-"`
	ProductType string            `json:"product_type" validate:"required,oneof='RAW MATERIAL' 'FINISHED GOODS' 'SEMI FINISHED GOODS'"`
	Name        string            `json:"name" validate:"required,min=4,max=100"`
	Slug        string            `json:"slug"`
	Code        string            `json:"code" validate:"required,min=4,max=100"`
	CmsCode     string            `json:"cms_code" validate:"omitempty,min=4,max=100"`
	Description string            `json:"description" validate:"required,min=4,max=400"`
	UnitType    string            `json:"unit_type" validate:"required,oneof='WEIGHT' 'PIECE' 'LIQUID'"`
	UnitWeight  float64           `json:"unit_weight" validate:"omitempty,lte=10000"`
	UnitValue   string            `json:"unit_weight_type" validate:"required,oneof='Kilogram' 'Gram' 'Liter' 'Milliliter' 'Piece'"`
	Status      types.Status      `json:"status" validate:"required,gt=0"`
	Category    *category.Form    `json:"category" validationTag:"category.id" validate:"omitempty,validRelationID,structonly"`
	Ingredients []*IngredientForm `json:"ingredients"`
}

type IngredientForm struct {
	IngredientID uint    `json:"ingredient_id" binding:"required"`
	Ingredient   *Form   `json:"ingredient"`
	Quantity     float64 `json:"quantity"`
}
