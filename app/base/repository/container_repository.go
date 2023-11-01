package repository

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"star-wms/app/base/dto/container"
	"star-wms/app/base/models"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/utils"
)

type ContainerRepository interface {
	GetAll(plantID uint, filter container.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Container, int64, error)
	Create(plantID uint, container *models.Container) error
	GetByID(plantID uint, id uint) (*models.Container, error)
	Update(plantID uint, container *models.Container) error
	Delete(plantID uint, id uint) error
	DeleteMulti(plantID uint, ids []uint) error
	ExistsByID(plantID uint, ID uint) bool
	ExistsByName(plantID uint, name string, ID uint) bool
	ExistsByCode(plantID uint, code string, ID uint) bool
}

type ContainerGormRepository struct {
	db *gorm.DB
}

func NewContainerGormRepository(database *gorm.DB) ContainerRepository {
	return &ContainerGormRepository{db: database}
}

func (p *ContainerGormRepository) GetAll(plantID uint, filter container.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Container, int64, error) {
	var containers []*models.Container
	var count int64

	query := p.db.Model(&models.Container{})
	query.Where("plant_id = ?", plantID)
	query = utils.BuildQuery(query, filter)

	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count containers")
		return nil, 0, err
	}

	query = utils.ApplySorting(query, sorting)
	query = utils.ApplyPagination(query, pagination)

	if err := query.Preload("Store").Preload("Product").Find(&containers).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get all containers")
		return nil, 0, err
	}

	return containers, count, nil
}

func (p *ContainerGormRepository) Create(plantID uint, containerModel *models.Container) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		//if containerModel.Category != nil {
		//	var category *models.Category
		//	if containerModel.Category.ID > 0 {
		//		if err := tx.First(&category, containerModel.Category.ID).Error; err != nil {
		//			log.Debug().Err(err).Msg("Failed to get category by ID")
		//			return err
		//		}
		//	}
		//	containerModel.Category = category
		//}
		//if containerModel.Approvers != nil {
		//	var existingApprovers []*adminModels.User
		//	for _, approverModel := range containerModel.Approvers {
		//		var existingApprover *adminModels.User
		//		if approverModel.ID > 0 {
		//			if err := tx.First(&existingApprover, approverModel.ID).Error; err != nil {
		//				log.Debug().Err(err).Msg("Failed to get approver by ID")
		//				return err
		//			}
		//		} else {
		//			continue
		//		}
		//		existingApprovers = append(existingApprovers, existingApprover)
		//	}
		//	containerModel.Approvers = existingApprovers
		//}
		containerModel.PlantID = plantID
		if err := tx.Create(&containerModel).Error; err != nil {
			return err
		}
		//if containerModel.Approvers != nil {
		//	if err := tx.Model(&containerModel).Association("Approvers").Replace(containerModel.Approvers); err != nil {
		//		return err
		//	}
		//}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create container")
	}
	return err
}

func (p *ContainerGormRepository) GetByID(plantID uint, id uint) (*models.Container, error) {
	var containerModel *models.Container
	if err := p.db.Preload("Product").Preload("Store").Where("plant_id = ?", plantID).First(&containerModel, id).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get container by ID")
		return nil, err
	}
	return containerModel, nil
}

func (p *ContainerGormRepository) Update(plantID uint, containerModel *models.Container) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		//if containerModel.Category != nil {
		//	var category *models.Category
		//	if containerModel.Category.ID > 0 {
		//		if err := tx.First(&category, containerModel.Category.ID).Error; err != nil {
		//			log.Debug().Err(err).Msg("Failed to get category by ID")
		//			return err
		//		}
		//	}
		//	containerModel.Category = category
		//}
		//if containerModel.Approvers != nil {
		//	var existingApprovers []*adminModels.User
		//	for _, approverModel := range containerModel.Approvers {
		//		var existingApprover *adminModels.User
		//		if approverModel.ID > 0 {
		//			if err := tx.First(&existingApprover, approverModel.ID).Error; err != nil {
		//				log.Debug().Err(err).Msg("Failed to get approver by ID")
		//				return err
		//			}
		//		} else {
		//			continue
		//		}
		//		existingApprovers = append(existingApprovers, existingApprover)
		//	}
		//	containerModel.Approvers = existingApprovers
		//}
		containerModel.PlantID = plantID
		if err := tx.Save(&containerModel).Error; err != nil {
			return err
		}
		//if containerModel.Approvers != nil {
		//	if err := tx.Model(&containerModel).Association("Approvers").Replace(containerModel.Approvers); err != nil {
		//		return err
		//	}
		//}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to update container")
	}
	return err
}

func (p *ContainerGormRepository) Delete(plantID uint, id uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		var containerModel models.Container
		if err := tx.Where("plant_id = ?", plantID).First(&containerModel, id).Error; err != nil {
			log.Debug().Err(err).Msg("Failed to get container by ID")
			return err
		}
		if err := tx.Where("plant_id = ?", plantID).Delete(&containerModel).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete container")
			return err
		}
		return nil
	})
}

func (p *ContainerGormRepository) DeleteMulti(plantID uint, ids []uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("plant_id = ?", plantID).Where("id IN ?", ids).Delete(&models.Container{}).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete containers")
			return err
		}
		return nil
	})
}

func (p *ContainerGormRepository) ExistsByID(plantID uint, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Container{}).Where("plant_id = ?", plantID).Where("id = ?", ID)
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by id")
		return false
	}
	return count > 0
}

func (p *ContainerGormRepository) ExistsByName(plantID uint, name string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Container{}).Where("plant_id = ?", plantID).Where("name = ?", name)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by name")
		return false
	}
	return count > 0
}

func (p *ContainerGormRepository) ExistsByCode(plantID uint, code string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Container{}).Where("plant_id = ?", plantID).Where("code = ?", code)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by code")
		return false
	}
	return count > 0
}
