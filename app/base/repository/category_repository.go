package repository

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"star-wms/app/base/dto/category"
	"star-wms/app/base/models"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/utils"
)

type CategoryRepository interface {
	GetAll(filter category.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting, withParent bool, withChildren bool) ([]*models.Category, int64, error)
	Create(category *models.Category) error
	GetByID(id uint, withParent bool, withChildren bool) (*models.Category, error)
	GetBySlug(slug string, withParent bool, withChildren bool) (*models.Category, error)
	Update(category *models.Category) error
	Delete(id uint) error
	DeleteMulti(ids []uint) error
	ExistsByID(ID uint) bool
	ExistsBySlug(slug string, ID uint) bool
	ExistsByName(name string, ID uint) bool
}

type CategoryGormRepository struct {
	db *gorm.DB
}

func NewCategoryGormRepository(database *gorm.DB) CategoryRepository {
	return &CategoryGormRepository{db: database}
}

func (p *CategoryGormRepository) GetAll(filter category.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting, withParent bool, withChildren bool) ([]*models.Category, int64, error) {
	var categories []*models.Category
	var count int64

	query := p.db.Model(&models.Category{})
	query = utils.BuildQuery(query, filter)

	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count categories")
		return nil, 0, err
	}

	query = utils.ApplySorting(query, sorting)
	query = utils.ApplyPagination(query, pagination)

	if withParent {
		query = query.Preload("Parent")
	}
	if withChildren {
		query = query.Preload("Children")
	}

	if err := query.Find(&categories).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get all categories")
		return nil, 0, err
	}

	return categories, count, nil
}

func (p *CategoryGormRepository) Create(categoryModel *models.Category) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&categoryModel).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create category")
	}
	return err
}

func (p *CategoryGormRepository) GetByID(id uint, withParent bool, withChildren bool) (*models.Category, error) {
	var categoryModel *models.Category
	query := p.db
	if withParent {
		query = query.Preload("Parent")
	}
	if withChildren {
		query = query.Preload("Children")
	}
	if err := query.First(&categoryModel, id).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get category by ID")
		return nil, err
	}
	return categoryModel, nil
}

func (p *CategoryGormRepository) GetBySlug(slug string, withParent bool, withChildren bool) (*models.Category, error) {
	var categoryModel *models.Category
	query := p.db
	if withParent {
		query = query.Preload("Parent")
	}
	if withChildren {
		query = query.Preload("Children")
	}
	if err := query.Where("slug = ?", slug).First(&categoryModel).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get category by Slug")
		return nil, err
	}
	return categoryModel, nil
}

func (p *CategoryGormRepository) Update(categoryModel *models.Category) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		if categoryModel.Parent != nil {
			var existingParent *models.Category
			if categoryModel.Parent.ID > 0 {
				if err := tx.First(&existingParent, categoryModel.Parent.ID).Error; err != nil {
					log.Debug().Err(err).Msg("Failed to get parent by ID")
					return err
				}
			}
			categoryModel.Parent = existingParent
		}
		if err := tx.Save(&categoryModel).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to update category")
	}
	return err
}

func (p *CategoryGormRepository) Delete(id uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		var categoryModel models.Category
		if err := tx.First(&categoryModel, id).Error; err != nil {
			log.Debug().Err(err).Msg("Failed to get category by ID")
			return err
		}
		if err := tx.Delete(&categoryModel).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete category")
			return err
		}
		return nil
	})
}

func (p *CategoryGormRepository) DeleteMulti(ids []uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id IN ?", ids).Delete(&models.Category{}).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete categories")
			return err
		}
		return nil
	})
}

func (p *CategoryGormRepository) ExistsByID(ID uint) bool {
	var count int64
	query := p.db.Model(&models.Category{}).Where("id = ?", ID)
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by id")
		return false
	}
	return count > 0
}

func (p *CategoryGormRepository) ExistsByName(name string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Category{}).Where("name = ?", name)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by name")
		return false
	}
	return count > 0
}

func (p *CategoryGormRepository) ExistsBySlug(slug string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Category{}).Where("slug = ?", slug)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by slug")
		return false
	}
	return count > 0
}
