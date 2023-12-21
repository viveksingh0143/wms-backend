package repository

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	adminModels "star-wms/app/admin/models"
	"star-wms/app/base/dto/store"
	"star-wms/app/base/models"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/utils"
)

type StoreRepository interface {
	GetAllByApprover(plantID uint, userID uint) ([]*models.Store, error)
	GetAll(plantID uint, filter store.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Store, int64, error)
	Create(plantID uint, store *models.Store) error
	GetByID(plantID uint, id uint) (*models.Store, error)
	GetByCode(plantID uint, code string) (*models.Store, error)
	Update(plantID uint, store *models.Store) error
	Delete(plantID uint, id uint) error
	DeleteMulti(plantID uint, ids []uint) error
	ExistsByID(plantID uint, ID uint) bool
	ExistsByName(plantID uint, name string, ID uint) bool
	ExistsByCode(plantID uint, code string, ID uint) bool
}

type StoreGormRepository struct {
	db *gorm.DB
}

func NewStoreGormRepository(database *gorm.DB) StoreRepository {
	return &StoreGormRepository{db: database}
}

func (p *StoreGormRepository) GetAllByApprover(plantID uint, userID uint) ([]*models.Store, error) {
	var stores []*models.Store

	if err := p.db.Model(&models.Store{}).
		Joins("JOIN store_approvers on store_approvers.store_id = stores.id").
		Where("store_approvers.user_id = ?", userID).
		Where("stores.plant_id = ?", plantID).
		Find(&stores).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get all stores with user approval authority")
		return nil, err
	}
	return stores, nil
}

func (p *StoreGormRepository) GetAll(plantID uint, filter store.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Store, int64, error) {
	var stores []*models.Store
	var count int64

	query := p.db.Model(&models.Store{})
	query.Where("plant_id = ?", plantID)
	query = utils.BuildQuery(query, filter)

	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count stores")
		return nil, 0, err
	}

	query = utils.ApplySorting(query, sorting)
	query = utils.ApplyPagination(query, pagination)

	if err := query.Preload("Approvers").Preload("Category").Find(&stores).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get all stores")
		return nil, 0, err
	}

	return stores, count, nil
}

func (p *StoreGormRepository) Create(plantID uint, storeModel *models.Store) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		if storeModel.Category != nil {
			var category *models.Category
			if storeModel.Category.ID > 0 {
				if err := tx.First(&category, storeModel.Category.ID).Error; err != nil {
					log.Debug().Err(err).Msg("Failed to get category by ID")
					return err
				}
			}
			storeModel.Category = category
		}
		storeModel.PlantID = plantID
		if storeModel.Approvers != nil {
			var existingApprovers []*adminModels.User
			for _, approverModel := range storeModel.Approvers {
				var existingApprover *adminModels.User
				if approverModel.ID > 0 {
					if err := tx.First(&existingApprover, approverModel.ID).Error; err != nil {
						log.Debug().Err(err).Msg("Failed to get approver by ID")
						return err
					}
				} else {
					continue
				}
				existingApprovers = append(existingApprovers, existingApprover)
			}
			storeModel.Approvers = existingApprovers
		}
		if err := tx.Create(&storeModel).Error; err != nil {
			return err
		}
		if storeModel.Approvers != nil {
			if err := tx.Model(&storeModel).Association("Approvers").Replace(storeModel.Approvers); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create store")
	}
	return err
}

func (p *StoreGormRepository) GetByID(plantID uint, id uint) (*models.Store, error) {
	var storeModel *models.Store
	if err := p.db.Preload("Approvers").Preload("Category").Where("plant_id = ?", plantID).First(&storeModel, id).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get store by ID")
		return nil, err
	}
	return storeModel, nil
}

func (p *StoreGormRepository) GetByCode(plantID uint, code string) (*models.Store, error) {
	var storeModel *models.Store
	if err := p.db.Where("plant_id = ?", plantID).Where("code = ?", code).First(&storeModel).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get store by Code")
		return nil, err
	}
	return storeModel, nil
}

func (p *StoreGormRepository) Update(plantID uint, storeModel *models.Store) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		if storeModel.Category != nil {
			var category *models.Category
			if storeModel.Category.ID > 0 {
				if err := tx.First(&category, storeModel.Category.ID).Error; err != nil {
					log.Debug().Err(err).Msg("Failed to get category by ID")
					return err
				}
			}
			storeModel.Category = category
		}
		storeModel.PlantID = plantID
		if storeModel.Approvers != nil {
			var existingApprovers []*adminModels.User
			for _, approverModel := range storeModel.Approvers {
				var existingApprover *adminModels.User
				if approverModel.ID > 0 {
					if err := tx.First(&existingApprover, approverModel.ID).Error; err != nil {
						log.Debug().Err(err).Msg("Failed to get approver by ID")
						return err
					}
				} else {
					continue
				}
				existingApprovers = append(existingApprovers, existingApprover)
			}
			storeModel.Approvers = existingApprovers
		}
		if err := tx.Save(&storeModel).Error; err != nil {
			return err
		}
		if storeModel.Approvers != nil {
			if err := tx.Model(&storeModel).Association("Approvers").Replace(storeModel.Approvers); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to update store")
	}
	return err
}

func (p *StoreGormRepository) Delete(plantID uint, id uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		var storeModel models.Store
		if err := tx.Where("plant_id = ?", plantID).First(&storeModel, id).Error; err != nil {
			log.Debug().Err(err).Msg("Failed to get store by ID")
			return err
		}
		if err := tx.Where("plant_id = ?", plantID).Delete(&storeModel).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete store")
			return err
		}
		return nil
	})
}

func (p *StoreGormRepository) DeleteMulti(plantID uint, ids []uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("plant_id = ?", plantID).Where("id IN ?", ids).Delete(&models.Store{}).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete stores")
			return err
		}
		return nil
	})
}

func (p *StoreGormRepository) ExistsByID(plantID uint, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Store{}).Where("plant_id = ?", plantID).Where("id = ?", ID)
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by id")
		return false
	}
	return count > 0
}

func (p *StoreGormRepository) ExistsByName(plantID uint, name string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Store{}).Where("plant_id = ?", plantID).Where("name = ?", name)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by name")
		return false
	}
	return count > 0
}

func (p *StoreGormRepository) ExistsByCode(plantID uint, code string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Store{}).Where("plant_id = ?", plantID).Where("code = ?", code)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by code")
		return false
	}
	return count > 0
}
