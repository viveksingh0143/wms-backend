package repository

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"star-wms/app/base/dto/joborder"
	"star-wms/app/base/models"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/utils"
)

type JobOrderRepository interface {
	GetAll(plantID uint, filter joborder.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.JobOrder, int64, error)
	Create(plantID uint, joborder *models.JobOrder) error
	GetByID(plantID uint, id uint) (*models.JobOrder, error)
	Update(plantID uint, joborder *models.JobOrder) error
	Delete(plantID uint, id uint) error
	DeleteMulti(plantID uint, ids []uint) error
	ExistsByID(plantID uint, ID uint) bool
	ExistsByOrderNo(plantID uint, orderNo string, ID uint) bool
}

type JobOrderGormRepository struct {
	db *gorm.DB
}

func NewJobOrderGormRepository(database *gorm.DB) JobOrderRepository {
	return &JobOrderGormRepository{db: database}
}

func (p *JobOrderGormRepository) GetAll(plantID uint, filter joborder.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.JobOrder, int64, error) {
	var joborders []*models.JobOrder
	var count int64

	query := p.db.Model(&models.JobOrder{})
	query.Where("plant_id = ?", plantID)
	query = utils.BuildQuery(query, filter)

	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count joborders")
		return nil, 0, err
	}

	query = utils.ApplySorting(query, sorting)
	query = utils.ApplyPagination(query, pagination)
	if err := query.Preload("Customer").Preload("Items.Product").Find(&joborders).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get all joborders")
		return nil, 0, err
	}

	return joborders, count, nil
}

func (p *JobOrderGormRepository) Create(plantID uint, joborderModel *models.JobOrder) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		if joborderModel.Customer != nil {
			var customer *models.Customer
			if joborderModel.Customer.ID > 0 {
				if err := tx.First(&customer, joborderModel.Customer.ID).Error; err != nil {
					log.Debug().Err(err).Msg("Failed to get customer by ID")
					return err
				}
			}
			joborderModel.Customer = customer
		}
		joborderModel.PlantID = plantID
		if err := tx.Omit("Items").Create(&joborderModel).Error; err != nil {
			return err
		}
		for _, item := range joborderModel.Items {
			item.JobOrderID = joborderModel.ID
		}
		if err := tx.Omit(clause.Associations).Create(&joborderModel.Items).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create joborder")
	}
	return err
}

func (p *JobOrderGormRepository) GetByID(plantID uint, id uint) (*models.JobOrder, error) {
	var joborderModel *models.JobOrder
	if err := p.db.Where("plant_id = ?", plantID).Preload("Customer").Preload("Items.Product").First(&joborderModel, id).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get joborder by ID")
		return nil, err
	}
	return joborderModel, nil
}

func (p *JobOrderGormRepository) Update(plantID uint, joborderModel *models.JobOrder) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		if joborderModel.Customer != nil {
			var customer *models.Customer
			if joborderModel.Customer.ID > 0 {
				if err := tx.First(&customer, joborderModel.Customer.ID).Error; err != nil {
					log.Debug().Err(err).Msg("Failed to get customer by ID")
					return err
				}
			}
			joborderModel.Customer = customer
		}
		joborderModel.PlantID = plantID
		if err := tx.Where("job_order_id = ?", joborderModel.ID).Delete(&models.JobOrderItem{}).Error; err != nil {
			return err
		}
		if err := tx.Omit("Items").Save(&joborderModel).Error; err != nil {
			return err
		}
		for _, item := range joborderModel.Items {
			item.JobOrderID = joborderModel.ID
		}
		if err := tx.Omit(clause.Associations).Create(&joborderModel.Items).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to update joborder")
	}
	return err
}

func (p *JobOrderGormRepository) Delete(plantID uint, id uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		var joborderModel models.JobOrder
		if err := tx.Where("plant_id = ?", plantID).First(&joborderModel, id).Error; err != nil {
			log.Debug().Err(err).Msg("Failed to get joborder by ID")
			return err
		}
		if err := tx.Where("plant_id = ?", plantID).Delete(&joborderModel).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete joborder")
			return err
		}
		return nil
	})
}

func (p *JobOrderGormRepository) DeleteMulti(plantID uint, ids []uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("plant_id = ?", plantID).Where("id IN ?", ids).Delete(&models.JobOrder{}).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete joborders")
			return err
		}
		return nil
	})
}

func (p *JobOrderGormRepository) ExistsByID(plantID uint, ID uint) bool {
	var count int64
	query := p.db.Model(&models.JobOrder{}).Where("plant_id = ?", plantID).Where("id = ?", ID)
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by id")
		return false
	}
	return count > 0
}

func (p *JobOrderGormRepository) ExistsByOrderNo(plantID uint, orderNo string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.JobOrder{}).Where("plant_id = ?", plantID).Where("order_no = ?", orderNo)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by order number")
		return false
	}
	return count > 0
}
