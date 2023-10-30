package models

import (
	"star-wms/core/common/models"
	"star-wms/core/types"
)

type Product struct {
	models.MyModel
	ProductType ProductType  `gorm:"type:enum('RAW MATERIAL','FINISHED GOODS','SEMI FINISHED GOODS');not null;default:'RAW MATERIAL';column:product_type"`
	Name        string       `gorm:"type:varchar(255);uniqueIndex;not null;column:name"`
	Slug        string       `gorm:"type:varchar(255);uniqueIndex;not null;column:slug"`
	Code        string       `gorm:"type:varchar(255);uniqueIndex;not null;column:code"`
	CmsCode     string       `gorm:"type:varchar(255);index;column:cms_code"`
	Description string       `gorm:"type:text;column:description"`
	UnitType    UnitType     `gorm:"type:enum('WEIGHT','PIECE','LIQUID');not null;default:'WEIGHT';column:unit_type"`
	UnitWeight  float64      `gorm:"column:unit_weight"`
	UnitValue   UnitValue    `gorm:"type:enum('Kilogram','Gram','Liter','Milliliter','Piece');column:unit_weight_type;default:'Gram'"`
	Status      types.Status `gorm:"type:int;default:1"`
	CategoryID  *uint        `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Category    *Category    `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type ProductType string
type UnitType string
type UnitValue string

const (
	// Product Types
	RAW_MATERIAL_TYPE        ProductType = "RAW MATERIAL"
	FINISHED_GOODS_TYPE      ProductType = "FINISHED GOODS"
	SEMI_FINISHED_GOODS_TYPE ProductType = "SEMI FINISHED GOODS"

	// Unit Types
	UnitTypeWeight UnitType = "WEIGHT"
	UnitTypePiece  UnitType = "PIECE"
	UnitTypeLiquid UnitType = "LIQUID"

	// Unit Values for UnitTypeWeight
	UnitValueKilogram UnitValue = "Kilogram"
	UnitValueGram     UnitValue = "Gram"

	// Unit Values for UnitTypeLiquid
	UnitValueLiter      UnitValue = "Liter"
	UnitValueMilliliter UnitValue = "Milliliter"

	// Unit Values for UnitTypePiece
	UnitValuePiece UnitValue = "Piece"
)
