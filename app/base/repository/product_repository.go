package repository

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"star-wms/app/base/dto/product"
	"star-wms/app/base/models"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/utils"
)

type ProductRepository interface {
	GetAll(filter product.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Product, int64, error)
	Create(product *models.Product) error
	GetByID(id uint) (*models.Product, error)
	GetByCode(code string) (*models.Product, error)
	Update(product *models.Product) error
	Delete(id uint) error
	DeleteMulti(ids []uint) error
	ExistsByID(ID uint) bool
	ExistsBySlug(slug string, ID uint) bool
	ExistsByName(name string, ID uint) bool
	ExistsByCode(code string, ID uint) bool
	ExistsByCmsCode(cmsCode string, ID uint) bool
}

type ProductGormRepository struct {
	db *gorm.DB
}

func NewProductGormRepository(database *gorm.DB) ProductRepository {
	return &ProductGormRepository{db: database}
}

func (p *ProductGormRepository) GetAll(filter product.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Product, int64, error) {
	var products []*models.Product
	var count int64

	query := p.db.Model(&models.Product{})
	query = utils.BuildQuery(query, filter)

	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count products")
		return nil, 0, err
	}

	query = utils.ApplySorting(query, sorting)
	query = utils.ApplyPagination(query, pagination)
	if err := query.Preload("Ingredients.Ingredient").Preload("Category").Find(&products).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get all products")
		return nil, 0, err
	}

	return products, count, nil
}

func (p *ProductGormRepository) Create(productModel *models.Product) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		if productModel.Category != nil {
			var category *models.Category
			if productModel.Category.ID > 0 {
				if err := tx.First(&category, productModel.Category.ID).Error; err != nil {
					log.Debug().Err(err).Msg("Failed to get category by ID")
					return err
				}
			}
			productModel.Category = category
		}
		if err := tx.Omit("Ingredients").Create(&productModel).Error; err != nil {
			return err
		}
		if productModel.Ingredients != nil && len(productModel.Ingredients) > 0 {
			for _, ingredient := range productModel.Ingredients {
				ingredient.ProductID = productModel.ID
			}
			if err := tx.Omit(clause.Associations).Create(&productModel.Ingredients).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create product")
	}
	return err
}

func (p *ProductGormRepository) GetByID(id uint) (*models.Product, error) {
	var productModel *models.Product
	if err := p.db.Preload("Ingredients.Ingredient").Preload("Category").First(&productModel, id).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get product by ID")
		return nil, err
	}
	return productModel, nil
}

func (p *ProductGormRepository) GetByCode(code string) (*models.Product, error) {
	var productModel *models.Product
	if err := p.db.Preload("Ingredients.Ingredient").Preload("Category").Where("code = ?", code).First(&productModel).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get product by Code")
		return nil, err
	}
	return productModel, nil
}

func (p *ProductGormRepository) Update(productModel *models.Product) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		if productModel.Category != nil {
			var category *models.Category
			if productModel.Category.ID > 0 {
				if err := tx.First(&category, productModel.Category.ID).Error; err != nil {
					log.Debug().Err(err).Msg("Failed to get category by ID")
					return err
				}
			}
			productModel.Category = category
		}
		if err := tx.Where("product_id = ?", productModel.ID).Delete(&models.ProductIngredient{}).Error; err != nil {
			return err
		}
		if err := tx.Omit("Ingredients").Save(&productModel).Error; err != nil {
			return err
		}
		if productModel.Ingredients != nil && len(productModel.Ingredients) > 0 {
			for _, ingredient := range productModel.Ingredients {
				ingredient.ProductID = productModel.ID
			}
			if err := tx.Omit(clause.Associations).Create(&productModel.Ingredients).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to update product")
	}
	return err
}

func (p *ProductGormRepository) Delete(id uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		var productModel models.Product
		if err := tx.First(&productModel, id).Error; err != nil {
			log.Debug().Err(err).Msg("Failed to get product by ID")
			return err
		}
		if err := tx.Delete(&productModel).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete product")
			return err
		}
		return nil
	})
}

func (p *ProductGormRepository) DeleteMulti(ids []uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id IN ?", ids).Delete(&models.Product{}).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete products")
			return err
		}
		return nil
	})
}

func (p *ProductGormRepository) ExistsByID(ID uint) bool {
	var count int64
	query := p.db.Model(&models.Product{}).Where("id = ?", ID)
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by id")
		return false
	}
	return count > 0
}

func (p *ProductGormRepository) ExistsByName(name string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Product{}).Where("name = ?", name)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by name")
		return false
	}
	return count > 0
}

func (p *ProductGormRepository) ExistsBySlug(slug string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Product{}).Where("slug = ?", slug)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by slug")
		return false
	}
	return count > 0
}

func (p *ProductGormRepository) ExistsByCode(code string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Product{}).Where("code = ?", code)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by code")
		return false
	}
	return count > 0
}

func (p *ProductGormRepository) ExistsByCmsCode(cmsCode string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Product{}).Where("cms_code = ?", cmsCode)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by cms_code")
		return false
	}
	return count > 0
}
