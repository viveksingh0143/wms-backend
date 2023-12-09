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

type JoborderRepository interface {
	GetAll(plantID uint, filter joborder.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Joborder, int64, error)
	Create(plantID uint, joborder *models.Joborder) error
	GetByID(plantID uint, id uint) (*models.Joborder, error)
	Update(plantID uint, joborder *models.Joborder) error
	Delete(plantID uint, id uint) error
	DeleteMulti(plantID uint, ids []uint) error
	ExistsByItemId(jororderID uint, ID uint) bool
	ExistsByID(plantID uint, ID uint) bool
	ExistsByOrderNo(plantID uint, orderNo string, ID uint) bool
}

type JoborderGormRepository struct {
	db *gorm.DB
}

func NewJoborderGormRepository(database *gorm.DB) JoborderRepository {
	return &JoborderGormRepository{db: database}
}

func (p *JoborderGormRepository) GetAll(plantID uint, filter joborder.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Joborder, int64, error) {
	var joborders []*models.Joborder
	var count int64

	query := p.db.Model(&models.Joborder{})
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

func (p *JoborderGormRepository) Create(plantID uint, joborderModel *models.Joborder) error {
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
			item.JoborderID = joborderModel.ID
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

func (p *JoborderGormRepository) GetByID(plantID uint, id uint) (*models.Joborder, error) {
	var joborderModel *models.Joborder
	if err := p.db.Where("plant_id = ?", plantID).Preload("Customer").Preload("Items.Product").First(&joborderModel, id).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get joborder by ID")
		return nil, err
	}
	return joborderModel, nil
}

func (p *JoborderGormRepository) Update(plantID uint, joborderModel *models.Joborder) error {
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
		if err := tx.Where("job_order_id = ?", joborderModel.ID).Delete(&models.JoborderItem{}).Error; err != nil {
			return err
		}
		if err := tx.Omit("Items").Save(&joborderModel).Error; err != nil {
			return err
		}
		for _, item := range joborderModel.Items {
			item.JoborderID = joborderModel.ID
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

func (p *JoborderGormRepository) Delete(plantID uint, id uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		var joborderModel models.Joborder
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

func (p *JoborderGormRepository) DeleteMulti(plantID uint, ids []uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("plant_id = ?", plantID).Where("id IN ?", ids).Delete(&models.Joborder{}).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete joborders")
			return err
		}
		return nil
	})
}

func (p *JoborderGormRepository) ExistsByItemId(jororderID uint, ID uint) bool {
	var count int64
	query := p.db.Model(&models.JoborderItem{}).Where("joborder_id = ?", jororderID).Where("id = ?", ID)
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by id")
		return false
	}
	return count > 0
}

func (p *JoborderGormRepository) ExistsByID(plantID uint, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Joborder{}).Where("plant_id = ?", plantID).Where("id = ?", ID)
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by id")
		return false
	}
	return count > 0
}

func (p *JoborderGormRepository) ExistsByOrderNo(plantID uint, orderNo string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Joborder{}).Where("plant_id = ?", plantID).Where("order_no = ?", orderNo)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by order number")
		return false
	}
	return count > 0
}
