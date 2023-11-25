package models

import (
	"star-wms/core/common/models"
	"star-wms/core/types"
)

type Product struct {
	models.MyModel
	ProductType  ProductType  `gorm:"type:enum('RAW MATERIAL','FINISHED GOODS','SEMI FINISHED GOODS');not null;default:'RAW MATERIAL';column:product_type"`
	Name         string       `gorm:"type:varchar(255);not null;column:name"`
	Slug         string       `gorm:"type:varchar(255);uniqueIndex;not null;column:slug"`
	Code         string       `gorm:"type:varchar(255);uniqueIndex;not null;column:code"`
	CmsCode      string       `gorm:"type:varchar(255);index;column:cms_code"`
	Description  string       `gorm:"type:text;column:description"`
	UnitType     UnitType     `gorm:"type:enum('WEIGHT','PIECE','LIQUID','LENGTH');not null;default:'WEIGHT';column:unit_type"`
	UnitWeight   float64      `gorm:"column:unit_weight"`
	UnitValue    UnitValue    `gorm:"type:enum('PC','GM','KG','MT','LT','YD','SM');column:unit_weight_type;default:'GM'"`
	Status       types.Status `gorm:"type:int;default:1"`
	CategoryPath string       `gorm:"index;"`
	CategoryID   *uint        `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Category     *Category    `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Ingredients  []*ProductIngredient
}

type ProductType string
type UnitType string
type UnitValue string

//const (
//	TypeRawMaterial       ProductType = "RAW MATERIAL"
//	TypeFinishedGoods     ProductType = "FINISHED GOODS"
//	TypeSemiFinishedGoods ProductType = "SEMI FINISHED GOODS"
//
//	UnitTypeWeight UnitType = "WEIGHT"
//	UnitTypePiece  UnitType = "PIECE"
//	UnitTypeLiquid UnitType = "LIQUID"
//	UnitTypeLength UnitType = "LENGTH"
//
//	UnitValueKilogram UnitValue = "Kilogram"
//	UnitValueGram     UnitValue = "Gram"
//
//	UnitValueLiter      UnitValue = "Liter"
//	UnitValueMilliliter UnitValue = "Milliliter"
//
//	UnitValuePiece UnitValue = "Piece"
//)
