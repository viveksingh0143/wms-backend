package service

import (
	"github.com/rs/zerolog/log"
	"star-wms/app/admin/dto/role"
	"star-wms/app/admin/dto/user"
	"star-wms/app/admin/models"
	"star-wms/app/admin/repository"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/common/responses"
	"star-wms/core/utils"
)

type UserService interface {
	GetAllUsers(filter user.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*user.Form, int64, error)
	CreateUser(userForm *user.Form) error
	GetUserByID(id uint) (*user.Form, error)
	UpdateUser(id uint, userForm *user.Form) error
	DeleteUser(id uint) error
	DeleteUsers(ids []uint) error
	ExistsByStaffID(staffid string, ID uint) bool
	ExistsByUsername(username string, ID uint) bool
	ExistsByEMail(email string, ID uint) bool
	ToModel(userForm *user.Form) *models.User
	FormToModel(userForm *user.Form, userModel *models.User)
	ToForm(userModel *models.User) *user.Form
	ToFormSlice(userModels []*models.User) []*user.Form
	ToModelSlice(userForms []*user.Form) []*models.User
}

type DefaultUserService struct {
	repo         repository.UserRepository
	roleService  RoleService
	plantService PlantService
}

func NewUserService(repo repository.UserRepository, roleService RoleService, plantService PlantService) UserService {
	return &DefaultUserService{repo: repo, roleService: roleService, plantService: plantService}
}

func (s *DefaultUserService) GetAllUsers(filter user.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*user.Form, int64, error) {
	data, count, err := s.repo.GetAll(filter, pagination, sorting)
	if err != nil {
		return nil, count, err
	}
	return s.ToFormSlice(data), count, err
}

func (s *DefaultUserService) CreateUser(userForm *user.Form) error {
	if s.ExistsByStaffID(userForm.StaffID, 0) {
		return responses.NewInputError("staff_id", "already exists", userForm.Name)
	}
	if s.ExistsByUsername(userForm.Username, 0) {
		return responses.NewInputError("username", "already exists", userForm.Name)
	}
	if s.ExistsByEMail(userForm.EMail, 0) {
		return responses.NewInputError("email", "already exists", userForm.Name)
	}
	if userForm.Plant != nil {
		if !s.plantService.ExistsById(userForm.Plant.ID) {
			return responses.NewInputError("plant.id", "plant not exists", userForm.Plant.ID)
		}
	}
	userModel := s.ToModel(userForm)
	hashedPassword, err := utils.GenerateFromPassword(userModel.Password)
	if err != nil {
		log.Error().Err(err).Msg("failed to encrypt password")
		return responses.NewInputError("password", "failed to encrypt", userForm.Password)
	}
	userModel.Password = hashedPassword
	return s.repo.Create(userModel)
}

func (s *DefaultUserService) GetUserByID(id uint) (*user.Form, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.ToForm(data), nil
}

func (s *DefaultUserService) UpdateUser(id uint, userForm *user.Form) error {
	if s.ExistsByStaffID(userForm.StaffID, id) {
		return responses.NewInputError("staff_id", "already exists", userForm.Name)
	}
	if s.ExistsByUsername(userForm.Username, id) {
		return responses.NewInputError("username", "already exists", userForm.Name)
	}
	if s.ExistsByEMail(userForm.EMail, id) {
		return responses.NewInputError("email", "already exists", userForm.Name)
	}
	if userForm.Plant != nil {
		if !s.plantService.ExistsById(userForm.Plant.ID) {
			return responses.NewInputError("plant.id", "plant not exists", userForm.Plant.ID)
		}
	}
	userModel, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	s.FormToModel(userForm, userModel)
	if userForm.Password != "" {
		hashedPassword, err := utils.GenerateFromPassword(userForm.Password)
		if err != nil {
			log.Error().Err(err).Msg("failed to encrypt password")
			return responses.NewInputError("password", "failed to encrypt", userModel.Password)
		}
		userModel.Password = hashedPassword
	}
	return s.repo.Update(userModel)
}

func (s *DefaultUserService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}

func (s *DefaultUserService) DeleteUsers(ids []uint) error {
	return s.repo.DeleteMulti(ids)
}

func (s *DefaultUserService) ExistsByStaffID(staffid string, ID uint) bool {
	return s.repo.ExistsByStaffID(staffid, ID)
}

func (s *DefaultUserService) ExistsByUsername(username string, ID uint) bool {
	return s.repo.ExistsByUsername(username, ID)
}

func (s *DefaultUserService) ExistsByEMail(email string, ID uint) bool {
	return s.repo.ExistsByEMail(email, ID)
}

func (s *DefaultUserService) ToModel(userForm *user.Form) *models.User {
	userModel := &models.User{
		Name:     userForm.Name,
		StaffID:  userForm.StaffID,
		Username: userForm.Username,
		EMail:    userForm.EMail,
		Password: userForm.Password,
		Status:   userForm.Status,
	}
	userModel.ID = userForm.ID

	if userForm.Roles != nil {
		roles := make([]*models.Role, 0)
		if len(userForm.Roles) > 0 {
			roles = s.roleService.ToModelSlice(userForm.Roles)
		}
		userModel.Roles = roles
	}
	if userForm.Plant != nil {
		userModel.Plant = s.plantService.ToModel(userForm.Plant)
	}
	return userModel
}

func (s *DefaultUserService) FormToModel(userForm *user.Form, userModel *models.User) {
	userModel.Name = userForm.Name
	userModel.StaffID = userForm.StaffID
	userModel.Username = userForm.Username
	userModel.EMail = userForm.EMail
	userModel.Status = userForm.Status

	if userForm.Roles != nil {
		roles := make([]*models.Role, 0)
		if len(userForm.Roles) > 0 {
			roles = s.roleService.ToModelSlice(userForm.Roles)
		}
		userModel.Roles = roles
	}
	if userForm.Plant != nil {
		userModel.Plant = s.plantService.ToModel(userForm.Plant)
	} else {
		userModel.Plant = nil
		userModel.PlantID = nil
	}
}

func (s *DefaultUserService) ToForm(userModel *models.User) *user.Form {
	userForm := &user.Form{
		ID:       userModel.ID,
		Name:     userModel.Name,
		StaffID:  userModel.StaffID,
		Username: userModel.Username,
		EMail:    userModel.EMail,
		Password: "",
		Status:   userModel.Status,
	}
	if userModel.Roles != nil {
		roles := make([]*role.Form, 0)
		if len(userModel.Roles) > 0 {
			roles = s.roleService.ToFormSlice(userModel.Roles)
		}
		userForm.Roles = roles
	}
	if userModel.Plant != nil {
		userForm.Plant = s.plantService.ToForm(userModel.Plant)
	}
	return userForm
}

func (s *DefaultUserService) ToFormSlice(userModels []*models.User) []*user.Form {
	data := make([]*user.Form, 0)
	for _, userModel := range userModels {
		data = append(data, s.ToForm(userModel))
	}
	return data
}

func (s *DefaultUserService) ToModelSlice(userForms []*user.Form) []*models.User {
	data := make([]*models.User, 0)
	for _, userForm := range userForms {
		data = append(data, s.ToModel(userForm))
	}
	return data
}
