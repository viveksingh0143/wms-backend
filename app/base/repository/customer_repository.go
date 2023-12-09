package repository

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"star-wms/app/base/dto/customer"
	"star-wms/app/base/models"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/utils"
)

type CustomerRepository interface {
	GetAll(plantID uint, filter customer.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Customer, int64, error)
	Create(plantID uint, customer *models.Customer) error
	GetByID(plantID uint, id uint) (*models.Customer, error)
	Update(plantID uint, customer *models.Customer) error
	Delete(plantID uint, id uint) error
	DeleteMulti(plantID uint, ids []uint) error
	ExistsByID(plantID uint, ID uint) bool
	ExistsByName(plantID uint, name string, ID uint) bool
	ExistsByCode(plantID uint, code string, ID uint) bool
}

type CustomerGormRepository struct {
	db *gorm.DB
}

func NewCustomerGormRepository(database *gorm.DB) CustomerRepository {
	return &CustomerGormRepository{db: database}
}

func (p *CustomerGormRepository) GetAll(plantID uint, filter customer.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Customer, int64, error) {
	var customers []*models.Customer
	var count int64

	query := p.db.Model(&models.Customer{})
	query.Where("plant_id = ?", plantID)
	query = utils.BuildQuery(query, filter)

	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count customers")
		return nil, 0, err
	}

	query = utils.ApplySorting(query, sorting)
	query = utils.ApplyPagination(query, pagination)

	if err := query.Find(&customers).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get all customers")
		return nil, 0, err
	}

	return customers, count, nil
}

func (p *CustomerGormRepository) Create(plantID uint, customerModel *models.Customer) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		customerModel.PlantID = plantID
		if err := tx.Create(&customerModel).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create customer")
	}
	return err
}

func (p *CustomerGormRepository) GetByID(plantID uint, id uint) (*models.Customer, error) {
	var customerModel *models.Customer
	if err := p.db.Where("plant_id = ?", plantID).First(&customerModel, id).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get customer by ID")
		return nil, err
	}
	return customerModel, nil
}

func (p *CustomerGormRepository) Update(plantID uint, customerModel *models.Customer) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		customerModel.PlantID = plantID
		if err := tx.Save(&customerModel).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to update customer")
	}
	return err
}

func (p *CustomerGormRepository) Delete(plantID uint, id uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		var customerModel models.Customer
		if err := tx.Where("plant_id = ?", plantID).First(&customerModel, id).Error; err != nil {
			log.Debug().Err(err).Msg("Failed to get customer by ID")
			return err
		}
		if err := tx.Where("plant_id = ?", plantID).Delete(&customerModel).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete customer")
			return err
		}
		return nil
	})
}

func (p *CustomerGormRepository) DeleteMulti(plantID uint, ids []uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("plant_id = ?", plantID).Where("id IN ?", ids).Delete(&models.Customer{}).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete customers")
			return err
		}
		return nil
	})
}

func (p *CustomerGormRepository) ExistsByID(plantID uint, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Customer{}).Where("plant_id = ?", plantID).Where("id = ?", ID)
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by id")
		return false
	}
	return count > 0
}

func (p *CustomerGormRepository) ExistsByName(plantID uint, name string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Customer{}).Where("plant_id = ?", plantID).Where("name = ?", name)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by name")
		return false
	}
	return count > 0
}

func (p *CustomerGormRepository) ExistsByCode(plantID uint, code string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Customer{}).Where("plant_id = ?", plantID).Where("code = ?", code)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by code")
		return false
	}
	return count > 0
}
