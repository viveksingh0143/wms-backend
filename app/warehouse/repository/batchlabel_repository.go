package repository

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	baseModels "star-wms/app/base/models"
	"star-wms/app/warehouse/dto/batchlabel"
	"star-wms/app/warehouse/models"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/utils"
)

type BatchlabelRepository interface {
	GetAll(plantID uint, filter batchlabel.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Batchlabel, int64, error)
	Create(plantID uint, batchlabel *models.Batchlabel) error
	GetByID(plantID uint, id uint) (*models.Batchlabel, error)
	Update(plantID uint, batchlabel *models.Batchlabel) error
	Delete(plantID uint, id uint) error
	DeleteMulti(plantID uint, ids []uint) error
	ExistsByID(plantID uint, ID uint) bool
	ExistsByBatchNo(plantID uint, batchNo string, ID uint) bool
}

type BatchlabelGormRepository struct {
	db *gorm.DB
}

func NewBatchlabelGormRepository(database *gorm.DB) BatchlabelRepository {
	return &BatchlabelGormRepository{db: database}
}

func (p *BatchlabelGormRepository) GetAll(plantID uint, filter batchlabel.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Batchlabel, int64, error) {
	var batchlabels []*models.Batchlabel
	var count int64

	query := p.db.Model(&models.Batchlabel{})
	query.Where("plant_id = ?", plantID)
	query = utils.BuildQuery(query, filter)

	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count batchlabels")
		return nil, 0, err
	}

	query = utils.ApplySorting(query, sorting)
	query = utils.ApplyPagination(query, pagination)
	if err := query.Preload("Joborder").Preload("Customer").Preload("Product").Find(&batchlabels).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get all batchlabels")
		return nil, 0, err
	}

	return batchlabels, count, nil
}

func (p *BatchlabelGormRepository) Create(plantID uint, batchlabelModel *models.Batchlabel) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		var customer *baseModels.Customer
		if batchlabelModel.Customer.ID > 0 {
			if err := tx.First(&customer, batchlabelModel.Customer.ID).Error; err != nil {
				log.Debug().Err(err).Msg("Failed to get customer by ID")
				return err
			}
		}
		batchlabelModel.Customer = customer

		var joborder *baseModels.Joborder
		if batchlabelModel.Joborder.ID > 0 {
			if err := tx.First(&joborder, batchlabelModel.Joborder.ID).Error; err != nil {
				log.Debug().Err(err).Msg("Failed to get job order by ID")
				return err
			}
		}
		batchlabelModel.Joborder = joborder

		var joborderItem *baseModels.JoborderItem
		if batchlabelModel.JoborderItem != nil && batchlabelModel.JoborderItem.ID > 0 {
			if err := tx.First(&joborderItem, batchlabelModel.JoborderItem.ID).Error; err != nil {
				log.Debug().Err(err).Msg("Failed to get job order item by ID")
				return err
			}
		}
		batchlabelModel.JoborderItem = joborderItem

		batchlabelModel.PlantID = plantID
		if err := tx.Omit("Stickers").Create(&batchlabelModel).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create batchlabel")
	}
	return err
}

func (p *BatchlabelGormRepository) GetByID(plantID uint, id uint) (*models.Batchlabel, error) {
	var batchlabelModel *models.Batchlabel
	if err := p.db.Where("plant_id = ?", plantID).Preload("Joborder.Items").Preload("JoborderItem").Preload("Customer").Preload("Product").Preload("Machine").Preload("Stickers").First(&batchlabelModel, id).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get batchlabel by ID")
		return nil, err
	}
	return batchlabelModel, nil
}

func (p *BatchlabelGormRepository) Update(plantID uint, batchlabelModel *models.Batchlabel) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		var customer *baseModels.Customer
		if batchlabelModel.Customer.ID > 0 {
			if err := tx.First(&customer, batchlabelModel.Customer.ID).Error; err != nil {
				log.Debug().Err(err).Msg("Failed to get customer by ID")
				return err
			}
		}
		batchlabelModel.Customer = customer

		var joborder *baseModels.Joborder
		if batchlabelModel.Joborder.ID > 0 {
			if err := tx.First(&joborder, batchlabelModel.Joborder.ID).Error; err != nil {
				log.Debug().Err(err).Msg("Failed to get job order by ID")
				return err
			}
		}
		batchlabelModel.Joborder = joborder

		var joborderItem *baseModels.JoborderItem
		if batchlabelModel.JoborderItem != nil && batchlabelModel.JoborderItem.ID > 0 {
			if err := tx.First(&joborderItem, batchlabelModel.JoborderItem.ID).Error; err != nil {
				log.Debug().Err(err).Msg("Failed to get job order item by ID")
				return err
			}
		}
		batchlabelModel.JoborderItem = joborderItem

		batchlabelModel.PlantID = plantID
		if err := tx.Omit("Stickers").Save(&batchlabelModel).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to update batchlabel")
	}
	return err
}

func (p *BatchlabelGormRepository) Delete(plantID uint, id uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		var batchlabelModel models.Batchlabel
		if err := tx.Where("plant_id = ?", plantID).First(&batchlabelModel, id).Error; err != nil {
			log.Debug().Err(err).Msg("Failed to get batchlabel by ID")
			return err
		}
		if err := tx.Where("plant_id = ?", plantID).Delete(&batchlabelModel).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete batchlabel")
			return err
		}
		return nil
	})
}

func (p *BatchlabelGormRepository) DeleteMulti(plantID uint, ids []uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("plant_id = ?", plantID).Where("id IN ?", ids).Delete(&models.Batchlabel{}).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete batchlabels")
			return err
		}
		return nil
	})
}

func (p *BatchlabelGormRepository) ExistsByID(plantID uint, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Batchlabel{}).Where("plant_id = ?", plantID).Where("id = ?", ID)
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by id")
		return false
	}
	return count > 0
}

func (p *BatchlabelGormRepository) ExistsByBatchNo(plantID uint, batchNo string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Batchlabel{}).Where("plant_id = ?", plantID).Where("batch_no = ?", batchNo)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by order number")
		return false
	}
	return count > 0
}
