package repository

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"star-wms/app/base/dto/outwardrequest"
	"star-wms/app/base/models"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/utils"
)

type OutwardrequestRepository interface {
	GetAll(plantID uint, filter outwardrequest.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Outwardrequest, int64, error)
	Create(plantID uint, outwardrequest *models.Outwardrequest) error
	GetByID(plantID uint, id uint) (*models.Outwardrequest, error)
	Update(plantID uint, outwardrequest *models.Outwardrequest) error
	Delete(plantID uint, id uint) error
	DeleteMulti(plantID uint, ids []uint) error
	ExistsByItemId(outwardrequestID uint, ID uint) bool
	ExistsByID(plantID uint, ID uint) bool
	ExistsByOrderNo(plantID uint, orderNo string, ID uint) bool
}

type OutwardrequestGormRepository struct {
	db *gorm.DB
}

func NewOutwardrequestGormRepository(database *gorm.DB) OutwardrequestRepository {
	return &OutwardrequestGormRepository{db: database}
}

func (p *OutwardrequestGormRepository) GetAll(plantID uint, filter outwardrequest.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Outwardrequest, int64, error) {
	var outwardrequests []*models.Outwardrequest
	var count int64

	query := p.db.Model(&models.Outwardrequest{})
	query.Where("plant_id = ?", plantID)
	query = utils.BuildQuery(query, filter)

	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count outwardrequests")
		return nil, 0, err
	}

	query = utils.ApplySorting(query, sorting)
	query = utils.ApplyPagination(query, pagination)
	if err := query.Preload("Customer").Preload("Items.Product").Find(&outwardrequests).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get all outwardrequests")
		return nil, 0, err
	}

	return outwardrequests, count, nil
}

func (p *OutwardrequestGormRepository) Create(plantID uint, outwardrequestModel *models.Outwardrequest) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		if outwardrequestModel.Customer != nil {
			var customer *models.Customer
			if outwardrequestModel.Customer.ID > 0 {
				if err := tx.First(&customer, outwardrequestModel.Customer.ID).Error; err != nil {
					log.Debug().Err(err).Msg("Failed to get customer by ID")
					return err
				}
			}
			outwardrequestModel.Customer = customer
		}
		outwardrequestModel.PlantID = plantID
		if err := tx.Omit("Items").Create(&outwardrequestModel).Error; err != nil {
			return err
		}
		for _, item := range outwardrequestModel.Items {
			item.OutwardrequestID = outwardrequestModel.ID
		}
		if err := tx.Omit(clause.Associations).Create(&outwardrequestModel.Items).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create outwardrequest")
	}
	return err
}

func (p *OutwardrequestGormRepository) GetByID(plantID uint, id uint) (*models.Outwardrequest, error) {
	var outwardrequestModel *models.Outwardrequest
	if err := p.db.Where("plant_id = ?", plantID).Preload("Customer").Preload("Items").Preload("Items.Product").First(&outwardrequestModel, id).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get outwardrequest by ID")
		return nil, err
	}
	return outwardrequestModel, nil
}

func (p *OutwardrequestGormRepository) Update(plantID uint, outwardrequestModel *models.Outwardrequest) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		if outwardrequestModel.Customer != nil {
			var customer *models.Customer
			if outwardrequestModel.Customer.ID > 0 {
				if err := tx.First(&customer, outwardrequestModel.Customer.ID).Error; err != nil {
					log.Debug().Err(err).Msg("Failed to get customer by ID")
					return err
				}
			}
			outwardrequestModel.Customer = customer
		}
		outwardrequestModel.PlantID = plantID
		if err := tx.Where("outwardrequest_id = ?", outwardrequestModel.ID).Delete(&models.OutwardrequestItem{}).Error; err != nil {
			return err
		}
		if err := tx.Omit("Items").Save(&outwardrequestModel).Error; err != nil {
			return err
		}
		for _, item := range outwardrequestModel.Items {
			item.OutwardrequestID = outwardrequestModel.ID
		}
		if err := tx.Omit(clause.Associations).Create(&outwardrequestModel.Items).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to update outwardrequest")
	}
	return err
}

func (p *OutwardrequestGormRepository) Delete(plantID uint, id uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		var outwardrequestModel models.Outwardrequest
		if err := tx.Where("plant_id = ?", plantID).First(&outwardrequestModel, id).Error; err != nil {
			log.Debug().Err(err).Msg("Failed to get outwardrequest by ID")
			return err
		}
		if err := tx.Where("plant_id = ?", plantID).Delete(&outwardrequestModel).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete outwardrequest")
			return err
		}
		return nil
	})
}

func (p *OutwardrequestGormRepository) DeleteMulti(plantID uint, ids []uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("plant_id = ?", plantID).Where("id IN ?", ids).Delete(&models.Outwardrequest{}).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete outwardrequests")
			return err
		}
		return nil
	})
}

func (p *OutwardrequestGormRepository) ExistsByItemId(outwardrequestID uint, ID uint) bool {
	var count int64
	query := p.db.Model(&models.OutwardrequestItem{}).Where("outwardrequest_id = ?", outwardrequestID).Where("id = ?", ID)
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by id")
		return false
	}
	return count > 0
}

func (p *OutwardrequestGormRepository) ExistsByID(plantID uint, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Outwardrequest{}).Where("plant_id = ?", plantID).Where("id = ?", ID)
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by id")
		return false
	}
	return count > 0
}

func (p *OutwardrequestGormRepository) ExistsByOrderNo(plantID uint, orderNo string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Outwardrequest{}).Where("plant_id = ?", plantID).Where("order_no = ?", orderNo)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by order number")
		return false
	}
	return count > 0
}
