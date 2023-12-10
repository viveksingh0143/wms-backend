package repository

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"star-wms/app/base/dto/machine"
	"star-wms/app/base/models"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/utils"
)

type MachineRepository interface {
	GetAll(plantID uint, filter machine.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Machine, int64, error)
	Create(plantID uint, machine *models.Machine) error
	GetByID(plantID uint, id uint) (*models.Machine, error)
	GetByCode(plantID uint, code string) (*models.Machine, error)
	Update(plantID uint, machine *models.Machine) error
	Delete(plantID uint, id uint) error
	DeleteMulti(plantID uint, ids []uint) error
	ExistsByID(plantID uint, ID uint) bool
	ExistsByName(plantID uint, name string, ID uint) bool
	ExistsByCode(plantID uint, code string, ID uint) bool
}

type MachineGormRepository struct {
	db *gorm.DB
}

func NewMachineGormRepository(database *gorm.DB) MachineRepository {
	return &MachineGormRepository{db: database}
}

func (p *MachineGormRepository) GetAll(plantID uint, filter machine.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Machine, int64, error) {
	var machines []*models.Machine
	var count int64

	query := p.db.Model(&models.Machine{})
	query.Where("plant_id = ?", plantID)
	query = utils.BuildQuery(query, filter)

	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count machines")
		return nil, 0, err
	}

	query = utils.ApplySorting(query, sorting)
	query = utils.ApplyPagination(query, pagination)

	if err := query.Find(&machines).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get all machines")
		return nil, 0, err
	}

	return machines, count, nil
}

func (p *MachineGormRepository) Create(plantID uint, machineModel *models.Machine) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		machineModel.PlantID = plantID
		if err := tx.Create(&machineModel).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create machine")
	}
	return err
}

func (p *MachineGormRepository) GetByID(plantID uint, id uint) (*models.Machine, error) {
	var machineModel *models.Machine
	if err := p.db.Where("plant_id = ?", plantID).First(&machineModel, id).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get machine by ID")
		return nil, err
	}
	return machineModel, nil
}

func (p *MachineGormRepository) GetByCode(plantID uint, code string) (*models.Machine, error) {
	var machineModel *models.Machine
	if err := p.db.Where("plant_id = ?", plantID).Where("code = ?", code).First(&machineModel).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get machine by Code")
		return nil, err
	}
	return machineModel, nil
}

func (p *MachineGormRepository) Update(plantID uint, machineModel *models.Machine) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		machineModel.PlantID = plantID
		if err := tx.Save(&machineModel).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to update machine")
	}
	return err
}

func (p *MachineGormRepository) Delete(plantID uint, id uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		var machineModel models.Machine
		if err := tx.Where("plant_id = ?", plantID).First(&machineModel, id).Error; err != nil {
			log.Debug().Err(err).Msg("Failed to get machine by ID")
			return err
		}
		if err := tx.Where("plant_id = ?", plantID).Delete(&machineModel).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete machine")
			return err
		}
		return nil
	})
}

func (p *MachineGormRepository) DeleteMulti(plantID uint, ids []uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("plant_id = ?", plantID).Where("id IN ?", ids).Delete(&models.Machine{}).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete machines")
			return err
		}
		return nil
	})
}

func (p *MachineGormRepository) ExistsByID(plantID uint, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Machine{}).Where("plant_id = ?", plantID).Where("id = ?", ID)
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by id")
		return false
	}
	return count > 0
}

func (p *MachineGormRepository) ExistsByName(plantID uint, name string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Machine{}).Where("plant_id = ?", plantID).Where("name = ?", name)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by name")
		return false
	}
	return count > 0
}

func (p *MachineGormRepository) ExistsByCode(plantID uint, code string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Machine{}).Where("plant_id = ?", plantID).Where("code = ?", code)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by code")
		return false
	}
	return count > 0
}
