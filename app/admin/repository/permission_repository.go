package repository

import (
	"errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"star-wms/app/admin/dto/permission"
	"star-wms/app/admin/models"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/utils"
)

type PermissionRepository interface {
	GetAll(filter permission.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Permission, int64, error)
	Create(permission *models.Permission) error
	GetByID(id uint) (*models.Permission, error)
	Update(permission *models.Permission) error
	Delete(id uint) error
	DeleteMulti(ids []uint) error
	ExistsByModuleName(moduleName string, ID uint) bool
}

type PermissionGormRepository struct {
	db *gorm.DB
}

func NewPermissionGormRepository(database *gorm.DB) PermissionRepository {
	return &PermissionGormRepository{db: database}
}

func (p *PermissionGormRepository) GetAll(filter permission.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Permission, int64, error) {
	var permissions = make([]*models.Permission, 0)
	var count int64

	query := p.db.Model(&models.Permission{})
	query = utils.BuildQuery(query, filter)

	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count permissions")
		return nil, 0, err
	}

	query = utils.ApplySorting(query, sorting)
	query = utils.ApplyPagination(query, pagination)

	if err := query.Find(&permissions).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get all permissions")
		return nil, 0, err
	}

	return permissions, count, nil
}

func (p *PermissionGormRepository) Create(permission *models.Permission) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(permission).Error; err != nil {
			log.Error().Err(err).Msg("Failed to create permission")
			return err
		}

		abilities, err := p.createAbilities(permission)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create abilities")
			return err
		}

		if err := tx.Create(abilities).Error; err != nil {
			log.Error().Err(err).Msg("Failed to create abilities in DB")
			return err
		}

		return nil
	})
}

func (p *PermissionGormRepository) GetByID(id uint) (*models.Permission, error) {
	var permissionModel *models.Permission
	if err := p.db.First(&permissionModel, id).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get permission by ID")
		return nil, err
	}
	return permissionModel, nil
}

func (p *PermissionGormRepository) Update(permission *models.Permission) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(permission).Error; err != nil {
			log.Error().Err(err).Msg("Failed to update permission")
			return err
		}
		abilitiesToDelete, abilitiesToInsert := p.computeAbilities(permission)
		if err := p.deleteAbilities(tx, abilitiesToDelete); err != nil {
			return err
		}
		if err := p.insertOrUpdateAbilities(tx, abilitiesToInsert); err != nil {
			return err
		}
		return nil
	})
}

func (p *PermissionGormRepository) Delete(id uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		return p.deleteSinglePermission(tx, id)
	})
}

func (p *PermissionGormRepository) DeleteMulti(ids []uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		for _, id := range ids {
			if err := p.deleteSinglePermission(tx, id); err != nil {
				return err
			}
		}
		return nil
	})
}

func (p *PermissionGormRepository) ExistsByModuleName(moduleName string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Permission{}).Where("module_name = ?", moduleName)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by module name")
		return false
	}
	return count > 0
}

func (p *PermissionGormRepository) createAbilities(permission *models.Permission) ([]*models.Ability, error) {
	var abilities []*models.Ability

	if permission.ReadPerm {
		abilities = append(abilities, &models.Ability{Module: permission.ModuleName, Name: permission.ReadPermName()})
	}
	if permission.CreatePerm {
		abilities = append(abilities, &models.Ability{Module: permission.ModuleName, Name: permission.CreatePermName()})
	}
	if permission.UpdatePerm {
		abilities = append(abilities, &models.Ability{Module: permission.ModuleName, Name: permission.UpdatePermName()})
	}
	if permission.DeletePerm {
		abilities = append(abilities, &models.Ability{Module: permission.ModuleName, Name: permission.DeletePermName()})
	}
	if permission.ImportPerm {
		abilities = append(abilities, &models.Ability{Module: permission.ModuleName, Name: permission.ImportPermName()})
	}
	if permission.ExportPerm {
		abilities = append(abilities, &models.Ability{Module: permission.ModuleName, Name: permission.ExportPermName()})
	}

	return abilities, nil
}

func (p *PermissionGormRepository) computeAbilities(permission *models.Permission) ([]*models.Ability, []*models.Ability) {
	abilitiesToDelete := make([]*models.Ability, 0)
	abilitiesToInsert := make([]*models.Ability, 0)
	if permission.ReadPerm {
		abilitiesToInsert = append(abilitiesToInsert, &models.Ability{Name: permission.ReadPermName(), Module: permission.ModuleName})
	} else {
		abilitiesToDelete = append(abilitiesToDelete, &models.Ability{Name: permission.ReadPermName(), Module: permission.ModuleName})
	}
	if permission.CreatePerm {
		abilitiesToInsert = append(abilitiesToInsert, &models.Ability{Name: permission.CreatePermName(), Module: permission.ModuleName})
	} else {
		abilitiesToDelete = append(abilitiesToDelete, &models.Ability{Name: permission.CreatePermName(), Module: permission.ModuleName})
	}
	if permission.UpdatePerm {
		abilitiesToInsert = append(abilitiesToInsert, &models.Ability{Name: permission.UpdatePermName(), Module: permission.ModuleName})
	} else {
		abilitiesToDelete = append(abilitiesToDelete, &models.Ability{Name: permission.UpdatePermName(), Module: permission.ModuleName})
	}
	if permission.DeletePerm {
		abilitiesToInsert = append(abilitiesToInsert, &models.Ability{Name: permission.DeletePermName(), Module: permission.ModuleName})
	} else {
		abilitiesToDelete = append(abilitiesToDelete, &models.Ability{Name: permission.DeletePermName(), Module: permission.ModuleName})
	}
	if permission.ImportPerm {
		abilitiesToInsert = append(abilitiesToInsert, &models.Ability{Name: permission.ImportPermName(), Module: permission.ModuleName})
	} else {
		abilitiesToDelete = append(abilitiesToDelete, &models.Ability{Name: permission.ImportPermName(), Module: permission.ModuleName})
	}
	if permission.ExportPerm {
		abilitiesToInsert = append(abilitiesToInsert, &models.Ability{Name: permission.ExportPermName(), Module: permission.ModuleName})
	} else {
		abilitiesToDelete = append(abilitiesToDelete, &models.Ability{Name: permission.ExportPermName(), Module: permission.ModuleName})
	}
	return abilitiesToDelete, abilitiesToInsert
}

func (p *PermissionGormRepository) deleteAbilities(tx *gorm.DB, abilities []*models.Ability) error {
	for _, abilityModel := range abilities {
		if err := tx.Where("name = ?", abilityModel.Name).Where("module = ?", abilityModel.Module).Delete(&models.Ability{}).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete ability")
			return err
		}
	}
	return nil
}

func (p *PermissionGormRepository) insertOrUpdateAbilities(tx *gorm.DB, abilities []*models.Ability) error {
	for _, abilityModel := range abilities {
		err := tx.Unscoped().Where(&abilityModel).First(&abilityModel).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := tx.Create(&abilityModel).Error; err != nil {
				log.Error().Err(err).Msg("Failed to create abilityModel")
				return err
			}
		} else if err != nil {
			log.Error().Err(err).Msg("Failed to get abilityModel")
			return err
		} else if abilityModel.DeletedAt.Valid {
			if err := tx.Model(&abilityModel).Unscoped().Update("deleted_at", nil).Error; err != nil {
				log.Error().Err(err).Msg("Failed to restore ability")
				return err
			}
		}
	}
	return nil
}

func (p *PermissionGormRepository) deleteSinglePermission(tx *gorm.DB, id uint) error {
	var permissionModel models.Permission

	if err := tx.First(&permissionModel, id).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get permission by ID")
		return err
	}

	if err := tx.Delete(&permissionModel).Error; err != nil {
		log.Error().Err(err).Msg("Failed to delete permission")
		return err
	}

	if err := p.deleteAssociatedAbilities(tx, permissionModel.ModuleName); err != nil {
		return err
	}
	return nil
}

func (p *PermissionGormRepository) deleteAssociatedAbilities(tx *gorm.DB, moduleName string) error {
	if err := tx.Where("module = ?", moduleName).Delete(&models.Ability{}).Error; err != nil {
		log.Error().Err(err).Msg("Failed to delete abilities associated with permission")
		return err
	}
	return nil
}
