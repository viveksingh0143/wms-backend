package models

type ProductIngredient struct {
	ProductID    uint     `gorm:"primaryKey;autoIncrement:false"`
	IngredientID uint     `gorm:"primaryKey;autoIncrement:false"`
	Product      *Product `gorm:"foreignKey:ProductID"`
	Ingredient   *Product `gorm:"foreignKey:IngredientID"`
	Quantity     float64
}
