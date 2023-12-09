package models

type ProductIngredient struct {
	ProductID    uint     `gorm:"primaryKey;autoIncrement:false"`
	IngredientID uint     `gorm:"primaryKey;autoIncrement:false"`
	Product      *Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	Ingredient   *Product `gorm:"foreignKey:IngredientID;constraint:OnDelete:CASCADE;"`
	Quantity     float64
}
