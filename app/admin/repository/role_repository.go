package repository

import (
	"errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"star-wms/app/admin/dto/role"
	"star-wms/app/admin/models"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/utils"
)

type RoleRepository interface {
	GetAll(filter role.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Role, int64, error)
	Create(role *models.Role) error
	GetByID(id uint) (*models.Role, error)
	Update(role *models.Role) error
	Delete(id uint) error
	DeleteMulti(ids []uint) error
	ExistsByName(name string, ID uint) bool
}

type RoleGormRepository struct {
	db *gorm.DB
}

func NewRoleGormRepository(database *gorm.DB) RoleRepository {
	return &RoleGormRepository{db: database}
}

func (p *RoleGormRepository) GetAll(filter role.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Role, int64, error) {
	var roles []*models.Role
	var count int64

	query := p.db.Model(&models.Role{})
	query = utils.BuildQuery(query, filter)

	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count roles")
		return nil, 0, err
	}

	query = utils.ApplySorting(query, sorting)
	query = utils.ApplyPagination(query, pagination)

	if err := query.Find(&roles).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get all roles")
		return nil, 0, err
	}

	return roles, count, nil
}

func (p *RoleGormRepository) Create(roleModel *models.Role) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		for i, abilityModel := range roleModel.Abilities {
			var existingAbility models.Ability
			result := tx.Where("name = ? AND module = ?", abilityModel.Name, abilityModel.Module).First(&existingAbility)

			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				// Record does not exist, proceed to insert
				if err := tx.Create(&abilityModel).Error; err != nil {
					return err
				}
			} else if result.Error != nil {
				// An error occurred during the check
				return result.Error
			} else {
				// Record exists, update the ID of abilityModel to match the existing record
				abilityModel.ID = existingAbility.ID
			}

			// Update roleModel.Abilities with either the new or existing abilityModel
			roleModel.Abilities[i] = abilityModel
		}
		if err := tx.Create(&roleModel).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create role")
	}
	return err
}

func (p *RoleGormRepository) GetByID(id uint) (*models.Role, error) {
	var roleModel *models.Role
	if err := p.db.Preload("Abilities").First(&roleModel, id).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get role by ID")
		return nil, err
	}
	return roleModel, nil
}

func (p *RoleGormRepository) Update(roleModel *models.Role) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		for i, abilityModel := range roleModel.Abilities {
			var existingAbility models.Ability
			result := tx.Where("name = ? AND module = ?", abilityModel.Name, abilityModel.Module).First(&existingAbility)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				if err := tx.Create(&abilityModel).Error; err != nil {
					return err
				}
			} else if result.Error != nil {
				return result.Error
			} else {
				// Record exists, you can choose to update it or do nothing
				// Example: Update existing record with new data
				// tx.Model(&existingAbility).Updates(abilityModel)
			}
			if abilityModel.ID == 0 {
				abilityModel.ID = existingAbility.ID
			}
			roleModel.Abilities[i] = abilityModel
		}
		if err := tx.Save(&roleModel).Error; err != nil {
			return err
		}
		if err := tx.Model(&roleModel).Association("Abilities").Replace(roleModel.Abilities); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create role")
	}
	return err
}

func (p *RoleGormRepository) Delete(id uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		var roleModel models.Role
		if err := tx.First(&roleModel, id).Error; err != nil {
			log.Debug().Err(err).Msg("Failed to get role by ID")
			return err
		}
		if err := tx.Delete(&roleModel).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete role")
			return err
		}
		return nil
	})
}

func (p *RoleGormRepository) DeleteMulti(ids []uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id IN ?", ids).Delete(&models.Role{}).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete roles")
			return err
		}
		return nil
	})
}

func (p *RoleGormRepository) ExistsByName(name string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Role{}).Where("name = ?", name)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by name")
		return false
	}
	return count > 0
}
