package repository

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"star-wms/app/base/dto/container"
	"star-wms/app/base/dto/store"
	"star-wms/app/base/models"
	warehouseModels "star-wms/app/warehouse/models"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/types"
	"star-wms/core/utils"
)

type ContainerRepository interface {
	GetAllRequiredApproval(plantID uint, stores []*store.Form, filter container.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Container, int64, error)
	GetAll(plantID uint, filter container.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Container, int64, error)
	Create(plantID uint, container *models.Container) error
	GetByID(plantID uint, id uint) (*models.Container, error)
	GetByCode(plantID uint, code string, needContents bool, needProduct bool, needStore bool, needLocation bool) (*models.Container, error)
	Update(plantID uint, container *models.Container) error
	Delete(plantID uint, id uint) error
	DeleteMulti(plantID uint, ids []uint) error
	ExistsByID(plantID uint, ID uint) bool
	ExistsByName(plantID uint, name string, ID uint) bool
	ExistsByCode(plantID uint, code string, ID uint) bool
	MarkedContainerFull(plantID uint, id uint) error
	Approve(plantID uint, id uint) error
	ApproveMulti(plantID uint, ids []uint) error
	Reject(plantID uint, id uint) error
	RejectMulti(plantID uint, ids []uint) error
}

type ContainerGormRepository struct {
	db *gorm.DB
}

func NewContainerGormRepository(database *gorm.DB) ContainerRepository {
	return &ContainerGormRepository{db: database}
}

func (p *ContainerGormRepository) GetAllRequiredApproval(plantID uint, storeForms []*store.Form, filter container.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Container, int64, error) {
	var containers []*models.Container
	var count int64
	storeIds := make([]uint, len(storeForms))
	for _, storeForm := range storeForms {
		storeIds = append(storeIds, storeForm.ID)
	}

	query := p.db.Model(&models.Container{})
	query.Where("plant_id = ?", plantID)
	query.Where("store_id in ?", storeIds)
	query.Where("product_id is not null")
	query.Where("approved = 3")
	query = utils.BuildQuery(query, filter)

	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count containers")
		return nil, 0, err
	}

	query = utils.ApplySorting(query, sorting)
	query = utils.ApplyPagination(query, pagination)

	if err := query.Preload("Store").Preload("Storelocation").Preload("Contents").Preload("Product").Find(&containers).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get all containers")
		return nil, 0, err
	}

	return containers, count, nil
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

	if err := query.Preload("Store").Preload("Storelocation").Preload("Product").Find(&containers).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get all containers")
		return nil, 0, err
	}

	return containers, count, nil
}

func (p *ContainerGormRepository) Create(plantID uint, containerModel *models.Container) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		containerModel.PlantID = plantID
		if err := tx.Create(&containerModel).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create container")
	}
	return err
}

func (p *ContainerGormRepository) GetByID(plantID uint, id uint) (*models.Container, error) {
	var containerModel *models.Container
	if err := p.db.Preload("Contents").Preload("Contents.Product").Preload("Product").Preload("Store").Preload("Storelocation").Where("plant_id = ?", plantID).First(&containerModel, id).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get container by ID")
		return nil, err
	}
	return containerModel, nil
}

func (p *ContainerGormRepository) GetByCode(plantID uint, code string, needContents bool, needProduct bool, needStore bool, needLocation bool) (*models.Container, error) {
	var containerModel *models.Container
	query := p.db
	if needContents {
		query = query.Preload("Contents")
	}
	if needProduct {
		query = query.Preload("Product")
	}
	if needStore {
		query = query.Preload("Store")
	}
	if needLocation {
		query = query.Preload("Storelocation")
	}
	if err := query.Where("plant_id = ?", plantID).Where("code", code).First(&containerModel).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get container by code")
		return nil, err
	}
	return containerModel, nil
}

func (p *ContainerGormRepository) Update(plantID uint, containerModel *models.Container) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		containerModel.PlantID = plantID
		if err := tx.Save(&containerModel).Error; err != nil {
			return err
		}
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

func (p *ContainerGormRepository) MarkedContainerFull(plantID uint, id uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Container{}).Where("plant_id = ?", plantID).Where("id = ?", id).Update("stock_level", "FULL").Error; err != nil {
			log.Error().Err(err).Msg("Failed to marked container full")
			return err
		}
		return nil
	})
}

func (p *ContainerGormRepository) Approve(plantID uint, id uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		var containerModel models.Container
		if err := tx.Where("plant_id = ?", plantID).First(&containerModel, id).Error; err != nil {
			log.Debug().Err(err).Msg("Failed to get container by ID")
			return err
		}

		if err := tx.Model(&models.Container{}).Where("plant_id = ?", plantID).Where("id = ?", id).Update("approved", 1).Error; err != nil {
			log.Error().Err(err).Msg("Failed to update the container field")
			return err
		}
		var containerContents []*models.ContainerContent
		query := p.db.Model(&models.ContainerContent{})
		query.Where("plant_id = ?", plantID)
		query.Where("container_id = ?", id)

		if err := query.Find(&containerContents).Error; err != nil {
			log.Error().Err(err).Msg("Failed to get all container contentss")
			return err
		}
		if containerContents != nil && len(containerContents) > 0 {
			for _, containerContent := range containerContents {
				if containerContent.RMBatchID != nil && *containerContent.RMBatchID > 0 {
					var rmBatchModel *warehouseModels.RMBatch
					if err := p.db.Where("plant_id = ?", plantID).First(&rmBatchModel, *containerContent.RMBatchID).Error; err != nil {
						log.Debug().Err(err).Msg("Failed to get raw material batch by ID")
						return err
					}
					rmBatchModel.Status = types.InventoryIn
					if err := tx.Omit(clause.Associations).Save(&rmBatchModel).Error; err != nil {
						return err
					}
					transactionHistory := rmBatchModel.NewTransactionHistory()

					if err := tx.Omit(clause.Associations).Create(&transactionHistory).Error; err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
}

func (p *ContainerGormRepository) ApproveMulti(plantID uint, ids []uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Container{}).Where("plant_id = ?", plantID).Where("id IN ?", ids).Update("approved", 1).Error; err != nil {
			log.Error().Err(err).Msg("Failed to update the container field")
			return err
		}
		var containerContents []*models.ContainerContent
		query := p.db.Model(&models.ContainerContent{})
		query.Where("plant_id = ?", plantID)
		query.Where("container_id in ?", ids)

		if err := query.Find(&containerContents).Error; err != nil {
			log.Error().Err(err).Msg("Failed to get all container contentss")
			return err
		}
		if containerContents != nil && len(containerContents) > 0 {
			for _, containerContent := range containerContents {
				if containerContent.RMBatchID != nil && *containerContent.RMBatchID > 0 {
					var rmBatchModel *warehouseModels.RMBatch
					if err := p.db.Where("plant_id = ?", plantID).First(&rmBatchModel, *containerContent.RMBatchID).Error; err != nil {
						log.Debug().Err(err).Msg("Failed to get raw material batch by ID")
						return err
					}
					rmBatchModel.Status = types.InventoryIn
					if err := tx.Omit(clause.Associations).Save(&rmBatchModel).Error; err != nil {
						return err
					}
					transactionHistory := rmBatchModel.NewTransactionHistory()

					if err := tx.Omit(clause.Associations).Create(&transactionHistory).Error; err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
}

func (p *ContainerGormRepository) Reject(plantID uint, id uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		var containerModel models.Container
		if err := tx.Where("plant_id = ?", plantID).First(&containerModel, id).Error; err != nil {
			log.Debug().Err(err).Msg("Failed to get container by ID")
			return err
		}

		if err := tx.Model(&models.Container{}).Where("plant_id = ?", plantID).Where("id = ?", id).Update("approved", 1).Error; err != nil {
			log.Error().Err(err).Msg("Failed to update the container field")
			return err
		}
		var containerContents []*models.ContainerContent
		query := p.db.Model(&models.ContainerContent{})
		query.Where("plant_id = ?", plantID)
		query.Where("container_id = ?", id)

		if err := query.Find(&containerContents).Error; err != nil {
			log.Error().Err(err).Msg("Failed to get all container contentss")
			return err
		}
		if containerContents != nil && len(containerContents) > 0 {
			for _, containerContent := range containerContents {
				if containerContent.RMBatchID != nil && *containerContent.RMBatchID > 0 {
					var rmBatchModel *warehouseModels.RMBatch
					if err := p.db.Where("plant_id = ?", plantID).First(&rmBatchModel, *containerContent.RMBatchID).Error; err != nil {
						log.Debug().Err(err).Msg("Failed to get raw material batch by ID")
						return err
					}
					rmBatchModel.Status = types.InventoryIn
					if err := tx.Omit(clause.Associations).Save(&rmBatchModel).Error; err != nil {
						return err
					}
					transactionHistory := rmBatchModel.NewTransactionHistory()

					if err := tx.Omit(clause.Associations).Create(&transactionHistory).Error; err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
}

func (p *ContainerGormRepository) RejectMulti(plantID uint, ids []uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Container{}).Where("plant_id = ?", plantID).Where("id IN ?", ids).Update("approved", 1).Error; err != nil {
			log.Error().Err(err).Msg("Failed to update the container field")
			return err
		}
		var containerContents []*models.ContainerContent
		query := p.db.Model(&models.ContainerContent{})
		query.Where("plant_id = ?", plantID)
		query.Where("container_id in ?", ids)

		if err := query.Find(&containerContents).Error; err != nil {
			log.Error().Err(err).Msg("Failed to get all container contentss")
			return err
		}
		if containerContents != nil && len(containerContents) > 0 {
			for _, containerContent := range containerContents {
				if containerContent.RMBatchID != nil && *containerContent.RMBatchID > 0 {
					var rmBatchModel *warehouseModels.RMBatch
					if err := p.db.Where("plant_id = ?", plantID).First(&rmBatchModel, *containerContent.RMBatchID).Error; err != nil {
						log.Debug().Err(err).Msg("Failed to get raw material batch by ID")
						return err
					}
					rmBatchModel.Status = types.InventoryIn
					if err := tx.Omit(clause.Associations).Save(&rmBatchModel).Error; err != nil {
						return err
					}
					transactionHistory := rmBatchModel.NewTransactionHistory()

					if err := tx.Omit(clause.Associations).Create(&transactionHistory).Error; err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
}
