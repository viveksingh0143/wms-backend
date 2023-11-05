package repository

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"star-wms/app/warehouse/dto/batchlabel"
	"star-wms/app/warehouse/models"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/utils"
	"time"
)

type StickerRepository interface {
	GetCount(plantID uint, batchlabelID uint) (int64, error)
	GetCountForShift(plantID uint, batchlabelID uint, shift string, createdAt time.Time) (int64, error)
	GetCountForBatchlabel(plantID uint, batchlabelID uint) (int64, error)
	GetAll(plantID uint, batchlabelID uint, filter batchlabel.StickerFilter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Sticker, int64, error)
	Create(plantID uint, batchlabelID uint, sticker *models.Sticker) error
	CreateAll(plantID uint, batchlabelID uint, stickers []*models.Sticker) error
	GetByID(plantID uint, batchlabelID uint, id uint) (*models.Sticker, error)
	ExistsByID(plantID uint, batchlabelID uint, ID uint) bool
	ExistsByBarcode(plantID uint, batchlabelID uint, barcode string, ID uint) bool
}

type StickerGormRepository struct {
	db *gorm.DB
}

func NewStickerGormRepository(database *gorm.DB) StickerRepository {
	return &StickerGormRepository{db: database}
}

func (p *StickerGormRepository) GetCount(plantID uint, batchlabelID uint) (int64, error) {
	var count int64
	query := p.db.Model(&models.Sticker{})
	query.Where("plant_id = ?", plantID).Where("batchlabel_id = ?", batchlabelID)
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count stickers")
		return 0, err
	}
	return count, nil
}

func (p *StickerGormRepository) GetCountForBatchlabel(plantID uint, batchlabelID uint) (int64, error) {
	var count int64
	query := p.db.Model(&models.Sticker{})
	query.
		Where("plant_id = ?", plantID).
		Where("batchlabel_id = ?", batchlabelID)

	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count stickers")
		return 0, err
	}
	return count, nil
}

func (p *StickerGormRepository) GetCountForShift(plantID uint, batchlabelID uint, shift string, createdAt time.Time) (int64, error) {
	var count int64
	query := p.db.Model(&models.Sticker{})
	query.
		Where("plant_id = ?", plantID).
		Where("batchlabel_id = ?", batchlabelID).
		Where("shift", shift).
		Where("DATE(created_at) = DATE(?)", createdAt)

	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count stickers")
		return 0, err
	}
	return count, nil
}

func (p *StickerGormRepository) GetAll(plantID uint, batchlabelID uint, filter batchlabel.StickerFilter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.Sticker, int64, error) {
	var stickers []*models.Sticker
	var count int64

	query := p.db.Model(&models.Sticker{})
	query.Where("plant_id = ?", plantID).Where("batchlabel_id = ?", batchlabelID)
	query = utils.BuildQuery(query, filter)

	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count stickers")
		return nil, 0, err
	}

	query = utils.ApplySorting(query, sorting)
	query = utils.ApplyPagination(query, pagination)
	if err := query.Preload("Product").Find(&stickers).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get all stickers")
		return nil, 0, err
	}

	return stickers, count, nil
}

func (p *StickerGormRepository) Create(plantID uint, batchlabelID uint, stickerModel *models.Sticker) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		stickerModel.BatchlabelID = batchlabelID
		stickerModel.PlantID = plantID
		if err := tx.Omit(clause.Associations).Create(&stickerModel).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create sticker")
	}
	return err
}

func (p *StickerGormRepository) CreateAll(plantID uint, batchlabelID uint, stickers []*models.Sticker) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		for _, sticker := range stickers {
			sticker.BatchlabelID = batchlabelID
			sticker.PlantID = plantID
		}
		if err := tx.CreateInBatches(&stickers, len(stickers)).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create sticker")
	}
	return err
}

func (p *StickerGormRepository) GetByID(plantID uint, batchlabelID uint, id uint) (*models.Sticker, error) {
	var stickerModel *models.Sticker
	if err := p.db.Where("plant_id = ?", plantID).Where("batchlabel_id = ?", batchlabelID).
		Preload("Product").
		Preload("Batchlabel").
		Preload("StickerItems.Product").
		First(&stickerModel, id).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get sticker by ID")
		return nil, err
	}
	return stickerModel, nil
}

func (p *StickerGormRepository) ExistsByID(plantID uint, batchlabelID uint, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Sticker{}).Where("plant_id = ?", plantID).Where("batchlabel_id = ?", batchlabelID).Where("id = ?", ID)
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by id")
		return false
	}
	return count > 0
}

func (p *StickerGormRepository) ExistsByBarcode(plantID uint, batchlabelID uint, barcode string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.Sticker{}).Where("plant_id = ?", plantID).Where("batchlabel_id = ?", batchlabelID).Where("barcode = ?", barcode)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by order number")
		return false
	}
	return count > 0
}
