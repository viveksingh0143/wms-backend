package repository

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"star-wms/app/base/dto/storelocation"
	"star-wms/app/base/models"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/utils"
)

type StorelocationRepository interface {
	GetAll(plantID uint, storeID uint, filter storelocation.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Storelocation, int64, error)
	Create(plantID uint, storeID uint, storelocation *models.Storelocation) error
	GetByID(plantID uint, storeID uint, id uint) (*models.Storelocation, error)
	GetByCode(plantID uint, code string) (*models.Storelocation, error)
	Update(plantID uint, storeID uint, storelocation *models.Storelocation) error
	Delete(plantID uint, storeID uint, id uint) error
	DeleteMulti(plantID uint, storeID uint, ids []uint) error
	ExistsByID(plantID uint, storeID uint, ID uint) bool
	ExistsByCode(plantID uint, storeID uint, code string, ID uint) bool
	ExistsByOnlyCode(plantID uint, code string) bool
}

type StorelocationGormRepository struct {
	db *gorm.DB
}

func NewStorelocationGormRepository(database *gorm.DB) StorelocationRepository {
	return &StorelocationGormRepository{db: database}
}

func (p *StorelocationGormRepository) GetAll(plantID uint, storeID uint, filter storelocation.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Storelocation, int64, error) {
	var storelocations []*models.Storelocation
	var count int64

	query := p.db.Model(&models.Storelocation{})
	query.Where("plant_id = ?", plantID).Where("store_id = ?", storeID)
	query = utils.BuildQuery(query, filter)

	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count store locations")
		return nil, 0, err
	}

	query = utils.ApplySorting(query, sorting)
	query = utils.ApplyPagination(query, pagination)

	if err := query.Find(&storelocations).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get all store locations")
		return nil, 0, err
	}

	return storelocations, count, nil
}

func (p *StorelocationGormRepository) Create(plantID uint, storeID uint, storelocationModel *models.Storelocation) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		storelocationModel.PlantID = plantID
		storelocationModel.StoreID = storeID
		if err := tx.Create(&storelocationModel).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create store location")
	}
	return err
}

func (p *StorelocationGormRepository) GetByID(plantID uint, storeID uint, id uint) (*models.Storelocation, error) {
	var storelocationModel *models.Storelocation
	if err := p.db.Preload("Store").Where("plant_id = ?", plantID).Where("store_id = ?", storeID).First(&storelocationModel, id).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get storelocation by ID")
		return nil, err
	}
	return storelocationModel, nil
}

func (p *StorelocationGormRepository) GetByCode(plantID uint, code string) (*models.Storelocation, error) {
	var storelocationModel *models.Storelocation
	if err := p.db.Preload("Store").Where("plant_id = ?", plantID).Where("code = ?", code).First(&storelocationModel).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get storelocation by Code")
		return nil, err
	}
	return storelocationModel, nil
}

func (p *StorelocationGormRepository) Update(plantID uint, storeID uint, storelocationModel *models.Storelocation) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		storelocationModel.PlantID = plantID
		storelocationModel.StoreID = storeID
		if err := tx.Save(&storelocationModel).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to update storelocation")
	}
	return err
}

func (p *StorelocationGormRepository) Delete(plantID uint, storeID uint, id uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		var storelocationModel models.Storelocation
		if err := tx.Where("plant_id = ?", plantID).Where("store_id = ?", storeID).First(&storelocationModel, id).Error; err != nil {
			log.Debug().Err(err).Msg("Failed to get storelocation by ID")
			return err
		}
		if err := tx.Where("plant_id = ?", plantID).Where("store_id = ?", storeID).Delete(&storelocationModel).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete storelocation")
			return err
		}
		return nil
	})
}

func (p *StorelocationGormRepository) DeleteMulti(plantID uint, storeID uint, ids []uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("plant_id = ?", plantID).Where("store_id = ?", storeID).Where("id IN ?", ids).Delete(&models.Storelocation{}).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete storelocations")
			return err
		}
		return nil
	})
}

func (p *StorelocationGormRepository) ExistsByID(plantID uint, storeID uint, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Storelocation{}).Where("plant_id = ?", plantID).Where("store_id = ?", storeID).Where("id = ?", ID)
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by id")
		return false
	}
	return count > 0
}

func (p *StorelocationGormRepository) ExistsByCode(plantID uint, storeID uint, code string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Storelocation{}).Where("plant_id = ?", plantID).Where("store_id = ?", storeID).Where("code = ?", code)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by code")
		return false
	}
	return count > 0
}

func (p *StorelocationGormRepository) ExistsByOnlyCode(plantID uint, code string) bool {
	var count int64
	query := p.db.Model(&models.Storelocation{}).Where("plant_id = ?", plantID).Where("code = ?", code)
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by code")
		return false
	}
	return count > 0
}
