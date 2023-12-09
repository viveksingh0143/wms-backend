package models

import (
	"gorm.io/gorm"
)

type MyModel struct {
	gorm.Model
	UpdatedBy string `gorm:"column:updated_by;type:varchar(100)"`
}

//func (m *MyModel) BeforeCreate(tx *gorm.DB) error {
//	if tx.Statement.Context == nil {
//		return nil
//	}
//
//	ginContext, ok := tx.Statement.Context.(*gin.Context)
//	if !ok {
//		return nil
//	}
//
//	contextValue, exists := ginContext.Get(gin.AuthUserKey)
//	if !exists {
//		return nil
//	}
//
//	// Assuming userForm is a struct and Identifier is a string field
//	userForm := contextValue.(user.Form)
//	m.UpdatedBy = fmt.Sprintf("%s (%s)", userForm.Name, userForm.StaffID)
//	return nil
//}
//
//func (m *MyModel) BeforeUpdate(tx *gorm.DB) error {
//	if tx.Statement.Context == nil {
//		return nil
//	}
//
//	ginContext, ok := tx.Statement.Context.(*gin.Context)
//	if !ok {
//		return nil
//	}
//
//	contextValue, exists := ginContext.Get(gin.AuthUserKey)
//	if !exists {
//		return nil
//	}
//
//	// Assuming userForm is a struct and Identifier is a string field
//	userForm := contextValue.(user.Form)
//	m.UpdatedBy = fmt.Sprintf("%s (%s)", userForm.Name, userForm.StaffID)
//	return nil
//}
