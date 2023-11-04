package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"star-wms/core/common/models"
	"star-wms/core/types"
)

type Category struct {
	models.MyModel
	Name     string       `gorm:"type:varchar(255);uniqueIndex;not null"`
	Slug     string       `gorm:"type:varchar(255);uniqueIndex;not null"`
	FullPath string       `gorm:"type:varchar(255);"`
	ParentID *uint        `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Parent   *Category    `gorm:"foreignKey:ParentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Children []*Category  `gorm:"foreignKey:ParentID"`
	Status   types.Status `gorm:"type:int;default:1"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {
	if c.Parent != nil {
		var parent *Category
		if c.Parent.Slug == "" {
			if err := tx.First(&parent, c.Parent.ID).Error; err != nil {
				return err
			}
		} else {
			parent = c.Parent
		}
		c.FullPath = parent.FullPath + c.Slug + "/"
	} else {
		c.FullPath = "/" + c.Slug + "/"
	}
	return nil
}

func updateChildPaths(tx *gorm.DB, parent *Category) error {
	var children []Category
	if err := tx.Where("parent_id = ?", parent.ID).Find(&children).Error; err != nil {
		return err
	}

	if len(children) > 0 {
		for _, child := range children {
			child.FullPath = parent.FullPath + child.Slug + "/"
			if err := tx.Model(&Category{}).Where("id = ?", child.ID).Omit(clause.Associations).Updates(map[string]interface{}{"FullPath": child.FullPath}).Error; err != nil {
				return err
			}
			// Update this record's FullPath to Products
			if err := tx.Model(&Product{}).Where("category_id = ?", child.ID).Omit(clause.Associations).Updates(map[string]interface{}{"CategoryPath": child.FullPath}).Error; err != nil {
				return err
			}
			// Update this record's FullPath to Stores
			if err := tx.Model(&Store{}).Where("category_id = ?", child.ID).Omit(clause.Associations).Updates(map[string]interface{}{"CategoryPath": child.FullPath}).Error; err != nil {
				return err
			}
			err := updateChildPaths(tx, &child)
			if err != nil {
				return err
			} // Recursive call
		}
	}

	return nil
}

func (c *Category) AfterUpdate(tx *gorm.DB) (err error) {
	if c.ID == 0 {
		return nil
	}
	if c.Parent != nil {
		var parent *Category
		if c.Parent.Slug == "" {
			if err := tx.First(&parent, c.Parent.ID).Error; err != nil {
				return err
			}
		} else {
			parent = c.Parent
		}
		c.FullPath = parent.FullPath + c.Slug + "/"
	} else {
		c.FullPath = "/" + c.Slug + "/"
	}
	// Update this record's FullPath
	if err := tx.Model(&Category{}).Where("id = ?", c.ID).Omit(clause.Associations).Updates(map[string]interface{}{"FullPath": c.FullPath}).Error; err != nil {
		return err
	}
	// Update this record's FullPath to Products
	if err := tx.Model(&Product{}).Where("category_id = ?", c.ID).Omit(clause.Associations).Updates(map[string]interface{}{"CategoryPath": c.FullPath}).Error; err != nil {
		return err
	}
	// Update this record's FullPath to Stores
	if err := tx.Model(&Store{}).Where("category_id = ?", c.ID).Omit(clause.Associations).Updates(map[string]interface{}{"CategoryPath": c.FullPath}).Error; err != nil {
		return err
	}
	// Update the FullPath for all children
	return updateChildPaths(tx, c)
}

func (c *Category) BeforeDelete(tx *gorm.DB) (err error) {
	return deleteChildCategories(tx, c.ID)
}

func deleteChildCategories(tx *gorm.DB, parentID uint) error {
	var children []Category
	if err := tx.Where("parent_id = ?", parentID).Find(&children).Error; err != nil {
		return err
	}

	for _, child := range children {
		if err := tx.Delete(&Category{}, child.ID).Error; err != nil {
			return err
		}
		err := deleteChildCategories(tx, child.ID)
		if err != nil {
			return err
		} // Recursive call
	}

	return nil
}
