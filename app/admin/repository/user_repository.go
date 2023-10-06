package repository

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"star-wms/app/admin/dto/user"
	"star-wms/app/admin/models"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/utils"
)

type UserRepository interface {
	GetAll(filter user.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.User, int64, error)
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	DeleteMulti(ids []uint) error
	ExistsByStaffID(staffid string, ID uint) bool
	ExistsByUsername(username string, ID uint) bool
	ExistsByEMail(email string, ID uint) bool
}

type UserGormRepository struct {
	db *gorm.DB
}

func NewUserGormRepository(database *gorm.DB) UserRepository {
	return &UserGormRepository{db: database}
}

func (p *UserGormRepository) GetAll(filter user.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*models.User, int64, error) {
	var users []*models.User
	var count int64

	query := p.db.Model(&models.User{})
	query = utils.BuildQuery(query, filter)

	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count users")
		return nil, 0, err
	}

	query = utils.ApplySorting(query, sorting)
	query = utils.ApplyPagination(query, pagination)

	if err := query.Preload("Plant").Preload("Roles").Find(&users).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get all users")
		return nil, 0, err
	}

	return users, count, nil
}

func (p *UserGormRepository) Create(userModel *models.User) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		if userModel.Roles != nil {
			var existingRoles []*models.Role
			for _, roleModel := range userModel.Roles {
				var existingRole *models.Role
				if roleModel.ID > 0 {
					if err := tx.First(&existingRole, roleModel.ID).Error; err != nil {
						log.Debug().Err(err).Msg("Failed to get role by ID")
						return err
					}
				} else {
					continue
				}
				existingRoles = append(existingRoles, existingRole)
			}
			userModel.Roles = existingRoles
		}

		if err := tx.Create(&userModel).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user")
	}
	return err
}

func (p *UserGormRepository) GetByID(id uint) (*models.User, error) {
	var userModel *models.User
	if err := p.db.Preload("Plant").Preload("Roles").First(&userModel, id).Error; err != nil {
		log.Debug().Err(err).Msg("Failed to get user by ID")
		return nil, err
	}
	return userModel, nil
}

func (p *UserGormRepository) Update(userModel *models.User) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		if userModel.Roles != nil {
			var existingRoles []*models.Role
			for _, roleModel := range userModel.Roles {
				var existingRole *models.Role
				if roleModel.ID > 0 {
					if err := tx.First(&existingRole, roleModel.ID).Error; err != nil {
						log.Debug().Err(err).Msg("Failed to get role by ID")
						return err
					}
				} else {
					continue
				}
				existingRoles = append(existingRoles, existingRole)
			}
			userModel.Roles = existingRoles
		}
		if userModel.Plant != nil {
			var existingPlant *models.Plant
			if userModel.Plant.ID > 0 {
				if err := tx.First(&existingPlant, userModel.Plant.ID).Error; err != nil {
					log.Debug().Err(err).Msg("Failed to get plant by ID")
					return err
				}
			}
			userModel.Plant = existingPlant
		}
		if err := tx.Save(&userModel).Error; err != nil {
			return err
		}
		if userModel.Roles != nil {
			if err := tx.Model(&userModel).Association("Roles").Replace(userModel.Roles); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user")
	}
	return err
}

func (p *UserGormRepository) Delete(id uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		var userModel models.User
		if err := tx.First(&userModel, id).Error; err != nil {
			log.Debug().Err(err).Msg("Failed to get user by ID")
			return err
		}
		if err := tx.Delete(&userModel).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete user")
			return err
		}
		return nil
	})
}

func (p *UserGormRepository) DeleteMulti(ids []uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id IN ?", ids).Delete(&models.User{}).Error; err != nil {
			log.Error().Err(err).Msg("Failed to delete users")
			return err
		}
		return nil
	})
}

func (p *UserGormRepository) ExistsByStaffID(staffid string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.User{}).Where("staff_id = ?", staffid)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by staff id")
		return false
	}
	return count > 0
}

func (p *UserGormRepository) ExistsByUsername(username string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.User{}).Where("username = ?", username)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by username")
		return false
	}
	return count > 0
}

func (p *UserGormRepository) ExistsByEMail(email string, ID uint) bool {
	var count int64
	query := p.db.Model(&models.User{}).Where("email = ?", email)
	if ID > 0 {
		query = query.Where("ID <> ?", ID)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count by email")
		return false
	}
	return count > 0
}
