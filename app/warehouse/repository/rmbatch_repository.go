package repository

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"star-wms/app/warehouse/dto/rmbatch"
	"star-wms/app/warehouse/models"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/utils"
)

type RMBatchRepository interface {
	GetAll(plantID uint, filter rmbatch.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.RMBatch, int64, error)
	GetByID(plantID uint, id uint) (*models.RMBatch, error)
}

type RMBatchGormRepository struct {
	db *gorm.DB
}

func NewRMBatchGormRepository(database *gorm.DB) RMBatchRepository {
	return &RMBatchGormRepository{db: database}
}

func (p *RMBatchGormRepository) GetAll(plantID uint, filter rmbatch.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.RMBatch, int64, error) {
	var inventories []*models.RMBatch
	var count int64

	query := p.db.Model(&models.RMBatch{})
	query.Where("plant_id = ?", plantID)
	query = utils.BuildQuery(query, filter)

	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count inventories")
		return nil, 0, err
	}

	query = utils.ApplySorting(query, sorting)
	query = utils.ApplyPagination(query, pagination)
	if err := query.Preload("Container").Preload("Store").Preload("Product").Find(&inventories).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get all inventories")
		return nil, 0, err
	}

	return inventories, count, nil
}

func (p *RMBatchGormRepository) GetByID(plantID uint, id uint) (*models.RMBatch, error) {
	var rmbatchModel *models.RMBatch
	if err := p.db.Where("plant_id = ?", plantID).Preload("Transactions").Preload("Container").Preload("Store").Preload("Product").First(&rmbatchModel, id).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get rmbatch by ID")
		return nil, err
	}
	return rmbatchModel, nil
}
