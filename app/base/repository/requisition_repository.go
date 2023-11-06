package repository

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"star-wms/app/base/dto/requisition"
	"star-wms/app/base/dto/store"
	"star-wms/app/base/models"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/utils"
)

type RequisitionRepository interface {
	GetAllRequiredApproval(plantID uint, stores []*store.Form, filter requisition.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Requisition, int64, error)
	GetAll(plantID uint, filter requisition.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Requisition, int64, error)
	Create(plantID uint, requisition *models.Requisition) error
	GetByID(plantID uint, id uint) (*models.Requisition, error)
	Update(plantID uint, requisition *models.Requisition) error
	Delete(plantID uint, id uint) error
	DeleteMulti(plantID uint, ids []uint) error
	ExistsByItemId(requisitionID uint, ID uint) bool
	ExistsByID(plantID uint, ID uint) bool
	ExistsByOrderNo(plantID uint, orderNo string, ID uint) bool
	Approve(plantID uint, id uint) error
	ApproveMulti(plantID uint, ids []uint) error
}

type RequisitionGormRepository struct {
	db *gorm.DB
}

func NewRequisitionGormRepository(database *gorm.DB) RequisitionRepository {
	return &RequisitionGormRepository{db: database}
}

func (p *RequisitionGormRepository) GetAllRequiredApproval(plantID uint, storeForms []*store.Form, filter requisition.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Requisition, int64, error) {
	var requisitions []*models.Requisition
	var count int64
	storeIds := make([]uint, len(storeForms))
	for _, storeForm := range storeForms {
		storeIds = append(storeIds, storeForm.ID)
	}

	query := p.db.Model(&models.Requisition{})
	query.Where("plant_id = ?", plantID)
	query.Where("store_id in ?", storeIds)
	query.Where("approved = 0")
	query = utils.BuildQuery(query, filter)

	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count requisitions")
		return nil, 0, err
	}

	query = utils.ApplySorting(query, sorting)
	query = utils.ApplyPagination(query, pagination)

	if err := query.Preload("Store").Preload("Items").Preload("Items.Product").Find(&requisitions).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get all requisitions")
		return nil, 0, err
	}

	return requisitions, count, nil
}

func (p *RequisitionGormRepository) GetAll(plantID uint, filter requisition.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Requisition, int64, error) {
	var requisitions []*models.Requisition
	var count int64

	query := p.db.Model(&models.Requisition{})
	query.Where("plant_id = ?", plantID)
	query = utils.BuildQuery(query, filter)

	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count requisitions")
		return nil, 0, err
	}

	query = utils.ApplySorting(query, sorting)
	query = utils.ApplyPagination(query, pagination)
	if err := query.Preload("Store").Preload("Items.Product").Find(&requisitions).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get all requisitions")
		return nil, 0, err
	}

	return requisitions, count, nil
}

func (p *RequisitionGormRepository) Create(plantID uint, requisitionModel *models.Requisition) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		if requisitionModel.Store != nil {
			var store *models.Store
			if requisitionModel.Store.ID > 0 {
				if err := tx.First(&store, requisitionModel.Store.ID).Error; err != nil {
					log.Debug().Err(err).Msg("Failed to get store by ID")
					return err
				}
			}
			requisitionModel.Store = store
		}
		requisitionModel.PlantID = plantID
		if err := tx.Omit("Items").Create(&requisitionModel).Error; err != nil {
			return err
		}
		for _, item := range requisitionModel.Items {
			item.RequisitionID = requisitionModel.ID
		}
		if err := tx.Omit(clause.Associations).Create(&requisitionModel.Items).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create requisition")
	}
	return err
}

func (p *RequisitionGormRepository) GetByID(plantID uint, id uint) (*models.Requisition, error) {
	var requisitionModel *models.Requisition
	if err := p.db.Where("plant_id = ?", plantID).Preload("Store").Preload("Items").Preload("Items.Product").First(&requisitionModel, id).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get requisition by ID")
		return nil, err
	}
	return requisitionModel, nil
}

func (p *RequisitionGormRepository) Update(plantID uint, requisitionModel *models.Requisition) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		if requisitionModel.Store != nil {
			var store *models.Store
			if requisitionModel.Store.ID > 0 {
				if err := tx.First(&store, requisitionModel.Store.ID).Error; err != nil {
					log.Debug().Err(err).Msg("Failed to get store by ID")
					return err
				}
			}
			requisitionModel.Store = store
		}
		requisitionModel.PlantID = plantID
		if err := tx.Where("requisition_id = ?", requisitionModel.ID).Delete(&models.RequisitionItem{}).Error; err != nil {
			return err
		}
		if err := tx.Omit("Items").Save(&requisitionModel).Error; err != nil {
			return err
		}
		for _, item := range requisitionModel.Items {
			item.RequisitionID = requisitionModel.ID
		}
		if err := tx.Omit(clause.Associations).Create(&requisitionModel.Items).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to update requisition")
	}
	return err
}

func (p *RequisitionGormRepository) Approve(plantID uint, id uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		var requisitionModel models.Requisition
		if err := tx.Where("plant_id = ?", plantID).First(&requisitionModel, id).Error; err != nil {
			log.Debug().Err(err).Msg("Failed to get requisition by ID")
			return err
		}

		if err := tx.Model(&models.Requisition{}).Where("plant_id = ?", plantID).Where("id = ?", id).Update("approved", 1).Error; err != nil {
			log.Error().Err(err).Msg("Failed to update the requisition field")
			return err
		}
		return nil
	})
}

func (p *RequisitionGormRepository) ApproveMulti(plantID uint, ids []uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Requisition{}).Where("plant_id = ?", plantID).Where("id IN ?", ids).Update("approved", 1).Error; err != nil {
			log.Error().Err(err).Msg("Failed to update the requisition field")
			return err
		}
		return nil
	})
}

func (p *RequisitionGormRepository) Delete(plantID uint, id uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		var requisitionModel models.Requisition
		if err := tx.Where("plant_id = ?", plantID).First(&requisitionModel, id).Error; err != nil {
			log.Debug().Err(err).Msg("Failed to get requisition by ID")
			return err
		}
		if err := tx.Where("plant_id = ?", plantID).Delete(&requisitionModel).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete requisition")
			return err
		}
		return nil
	})
}

func (p *RequisitionGormRepository) DeleteMulti(plantID uint, ids []uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("plant_id = ?", plantID).Where("id IN ?", ids).Delete(&models.Requisition{}).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete requisitions")
			return err
		}
		return nil
	})
}

func (p *RequisitionGormRepository) ExistsByItemId(requisitionID uint, ID uint) bool {
	var count int64
	query := p.db.Model(&models.RequisitionItem{}).Where("requisition_id = ?", requisitionID).Where("id = ?", ID)
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by id")
		return false
	}
	return count > 0
}

func (p *RequisitionGormRepository) ExistsByID(plantID uint, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Requisition{}).Where("plant_id = ?", plantID).Where("id = ?", ID)
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by id")
		return false
	}
	return count > 0
}

func (p *RequisitionGormRepository) ExistsByOrderNo(plantID uint, orderNo string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Requisition{}).Where("plant_id = ?", plantID).Where("order_no = ?", orderNo)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by order number")
		return false
	}
	return count > 0
}
