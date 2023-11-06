package repository

import (
	"errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	baseModels "star-wms/app/base/models"
	"star-wms/app/warehouse/dto/inventory"
	"star-wms/app/warehouse/models"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/utils"
)

type InventoryRepository interface {
	GetAll(plantID uint, filter inventory.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Inventory, int64, error)
	Create(plantID uint, inventory *models.Inventory) error
	GetByID(plantID uint, id uint) (*models.Inventory, error)
	Update(plantID uint, inventory *models.Inventory) error
	Delete(plantID uint, id uint) error
	DeleteMulti(plantID uint, ids []uint) error
	ExistsByID(plantID uint, ID uint) bool
	CreateRawMaterialStockIn(plantID uint, storeModel *baseModels.Store, containerModel *baseModels.Container, contentModel *baseModels.ContainerContent) error
}

type InventoryGormRepository struct {
	db *gorm.DB
}

func NewInventoryGormRepository(database *gorm.DB) InventoryRepository {
	return &InventoryGormRepository{db: database}
}

func (p *InventoryGormRepository) GetAll(plantID uint, filter inventory.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Inventory, int64, error) {
	var inventories []*models.Inventory
	var count int64

	query := p.db.Model(&models.Inventory{})
	query.Where("plant_id = ?", plantID)
	query = utils.BuildQuery(query, filter)

	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count inventories")
		return nil, 0, err
	}

	query = utils.ApplySorting(query, sorting)
	query = utils.ApplyPagination(query, pagination)
	if err := query.Preload("Store").Preload("Product").Find(&inventories).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get all inventories")
		return nil, 0, err
	}

	return inventories, count, nil
}

func (p *InventoryGormRepository) CreateRawMaterialStockIn(plantID uint, storeModel *baseModels.Store, containerModel *baseModels.Container, contentModel *baseModels.ContainerContent) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		contentModel.ContainerID = containerModel.ID
		contentModel.PlantID = plantID
		if err := tx.Omit(clause.Associations).Create(&contentModel).Error; err != nil {
			return err
		}

		containerModel.StockLevel = baseModels.Full
		containerModel.Approved = false
		containerModel.ProductID = &contentModel.ProductID
		containerModel.StoreID = &storeModel.ID
		if err := tx.Model(&containerModel).Updates(map[string]interface{}{
			"stock_level": baseModels.Full,
			"approved":    false,
			"product_id":  &contentModel.ProductID,
			"store_id":    &storeModel.ID,
		}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create inventory")
	}
	return err
}

func (p *InventoryGormRepository) Create(plantID uint, inventoryModel *models.Inventory) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		var store *baseModels.Store
		if inventoryModel.Store.ID > 0 {
			if err := tx.First(&store, inventoryModel.Store.ID).Error; err != nil {
				log.Debug().Err(err).Msg("Failed to get store by ID")
				return err
			}
		} else {
			err := errors.New("failed to get store")
			log.Debug().Err(err).Msg("Failed to get store by ID")
			return err
		}
		inventoryModel.Store = store

		var product *baseModels.Product
		if inventoryModel.Product.ID > 0 {
			if err := tx.First(&product, inventoryModel.Product.ID).Error; err != nil {
				log.Debug().Err(err).Msg("Failed to get product by ID")
				return err
			}
		} else {
			err := errors.New("failed to get product")
			log.Debug().Err(err).Msg("Failed to get product by ID")
			return err
		}
		inventoryModel.Product = product

		inventoryModel.PlantID = plantID
		if err := tx.Create(&inventoryModel).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create inventory")
	}
	return err
}

func (p *InventoryGormRepository) GetByID(plantID uint, id uint) (*models.Inventory, error) {
	var inventoryModel *models.Inventory
	if err := p.db.Where("plant_id = ?", plantID).Preload("Store").Preload("Product").First(&inventoryModel, id).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get inventory by ID")
		return nil, err
	}
	return inventoryModel, nil
}

func (p *InventoryGormRepository) Update(plantID uint, inventoryModel *models.Inventory) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		var store *baseModels.Store
		if inventoryModel.Store.ID > 0 {
			if err := tx.First(&store, inventoryModel.Store.ID).Error; err != nil {
				log.Debug().Err(err).Msg("Failed to get store by ID")
				return err
			}
		} else {
			err := errors.New("failed to get store")
			log.Debug().Err(err).Msg("Failed to get store by ID")
			return err
		}
		inventoryModel.Store = store

		var product *baseModels.Product
		if inventoryModel.Product.ID > 0 {
			if err := tx.First(&product, inventoryModel.Product.ID).Error; err != nil {
				log.Debug().Err(err).Msg("Failed to get product by ID")
				return err
			}
		} else {
			err := errors.New("failed to get product")
			log.Debug().Err(err).Msg("Failed to get product by ID")
			return err
		}
		inventoryModel.Product = product

		inventoryModel.PlantID = plantID
		if err := tx.Save(&inventoryModel).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to update inventory")
	}
	return err
}

func (p *InventoryGormRepository) Delete(plantID uint, id uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		var inventoryModel models.Inventory
		if err := tx.Where("plant_id = ?", plantID).First(&inventoryModel, id).Error; err != nil {
			log.Debug().Err(err).Msg("Failed to get inventory by ID")
			return err
		}
		if err := tx.Where("plant_id = ?", plantID).Delete(&inventoryModel).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete inventory")
			return err
		}
		return nil
	})
}

func (p *InventoryGormRepository) DeleteMulti(plantID uint, ids []uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("plant_id = ?", plantID).Where("id IN ?", ids).Delete(&models.Inventory{}).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete inventories")
			return err
		}
		return nil
	})
}

func (p *InventoryGormRepository) ExistsByID(plantID uint, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Inventory{}).Where("plant_id = ?", plantID).Where("id = ?", ID)
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by id")
		return false
	}
	return count > 0
}
