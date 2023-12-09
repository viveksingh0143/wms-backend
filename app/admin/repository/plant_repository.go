package repository

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"star-wms/app/admin/dto/plant"
	"star-wms/app/admin/models"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/utils"
)

type PlantRepository interface {
	GetAll(filter plant.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Plant, int64, error)
	Create(plant *models.Plant) error
	GetByID(id uint) (*models.Plant, error)
	Update(plant *models.Plant) error
	Delete(id uint) error
	DeleteMulti(ids []uint) error
	ExistsByID(ID uint) bool
	ExistsByCode(code string, ID uint) bool
	ExistsByName(name string, ID uint) bool
}

type PlantGormRepository struct {
	db *gorm.DB
}

func NewPlantGormRepository(database *gorm.DB) PlantRepository {
	return &PlantGormRepository{db: database}
}

func (p *PlantGormRepository) GetAll(filter plant.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Plant, int64, error) {
	var plants []*models.Plant
	var count int64

	query := p.db.Model(&models.Plant{})
	query = utils.BuildQuery(query, filter)

	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count plants")
		return nil, 0, err
	}

	query = utils.ApplySorting(query, sorting)
	query = utils.ApplyPagination(query, pagination)

	if err := query.Find(&plants).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get all plants")
		return nil, 0, err
	}

	return plants, count, nil
}

func (p *PlantGormRepository) Create(plantModel *models.Plant) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(plantModel).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create plant")
	}
	return err
}

func (p *PlantGormRepository) GetByID(id uint) (*models.Plant, error) {
	var plantModel *models.Plant
	if err := p.db.First(&plantModel, id).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get plant by ID")
		return nil, err
	}
	return plantModel, nil
}

func (p *PlantGormRepository) Update(plantModel *models.Plant) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(plantModel).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create plant")
	}
	return err
}

func (p *PlantGormRepository) Delete(id uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		var plantModel models.Plant
		if err := tx.First(&plantModel, id).Error; err != nil {
			log.Debug().Err(err).Msg("Failed to get plant by ID")
			return err
		}
		if err := tx.Delete(&plantModel).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete plant")
			return err
		}
		return nil
	})
}

func (p *PlantGormRepository) DeleteMulti(ids []uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id IN ?", ids).Delete(&models.Plant{}).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete plants")
			return err
		}
		return nil
	})
}

func (p *PlantGormRepository) ExistsByID(ID uint) bool {
	var count int64
	query := p.db.Model(&models.Plant{}).Where("id = ?", ID)
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by id")
		return false
	}
	return count > 0
}

func (p *PlantGormRepository) ExistsByCode(code string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Plant{}).Where("code = ?", code)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by code")
		return false
	}
	return count > 0
}

func (p *PlantGormRepository) ExistsByName(name string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Plant{}).Where("name = ?", name)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by name")
		return false
	}
	return count > 0
}
