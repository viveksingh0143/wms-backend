package service

import (
	"star-wms/app/admin/dto/plant"
	"star-wms/app/admin/models"
	"star-wms/app/admin/repository"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/common/responses"
)

type PlantService interface {
	GetAllPlants(filter plant.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*plant.Form, int64, error)
	CreatePlant(plantForm *plant.Form) error
	GetPlantByID(id uint) (*plant.Form, error)
	GetPlantByCode(code string) (*plant.Form, error)
	UpdatePlant(id uint, plantForm *plant.Form) error
	DeletePlant(id uint) error
	DeletePlants(ids []uint) error
	ExistsById(ID uint) bool
	ExistsByCode(code string, ID uint) bool
	ExistsByName(name string, ID uint) bool
	ToModel(plantForm *plant.Form) *models.Plant
	ToForm(plantModel *models.Plant) *plant.Form
	FormToModel(plantForm *plant.Form, plantModel *models.Plant)
	ToFormSlice(plantModels []*models.Plant) []*plant.Form
	ToModelSlice(plantForms []*plant.Form) []*models.Plant
}

type DefaultPlantService struct {
	repo repository.PlantRepository
}

func NewPlantService(repo repository.PlantRepository) PlantService {
	return &DefaultPlantService{repo: repo}
}

func (s *DefaultPlantService) GetAllPlants(filter plant.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*plant.Form, int64, error) {
	data, count, err := s.repo.GetAll(filter, pagination, sorting)
	if err != nil {
		return nil, count, err
	}
	return s.ToFormSlice(data), count, err
}

func (s *DefaultPlantService) CreatePlant(plantForm *plant.Form) error {
	if s.ExistsByCode(plantForm.Code, 0) {
		return responses.NewInputError("code", "already exists", plantForm.Code)
	}
	if s.ExistsByName(plantForm.Name, 0) {
		return responses.NewInputError("name", "already exists", plantForm.Name)
	}
	resultModel := s.ToModel(plantForm)
	return s.repo.Create(resultModel)
}

func (s *DefaultPlantService) GetPlantByID(id uint) (*plant.Form, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.ToForm(data), nil
}

func (s *DefaultPlantService) GetPlantByCode(code string) (*plant.Form, error) {
	data, err := s.repo.GetByCode(code)
	if err != nil {
		return nil, err
	}
	return s.ToForm(data), nil
}

func (s *DefaultPlantService) UpdatePlant(id uint, plantForm *plant.Form) error {
	if s.ExistsByCode(plantForm.Code, id) {
		return responses.NewInputError("code", "already exists", plantForm.Code)
	}
	if s.ExistsByName(plantForm.Name, id) {
		return responses.NewInputError("name", "already exists", plantForm.Name)
	}
	plantModel, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	s.FormToModel(plantForm, plantModel)
	return s.repo.Update(plantModel)
}

func (s *DefaultPlantService) DeletePlant(id uint) error {
	return s.repo.Delete(id)
}

func (s *DefaultPlantService) DeletePlants(ids []uint) error {
	return s.repo.DeleteMulti(ids)
}

func (s *DefaultPlantService) ExistsById(ID uint) bool {
	return s.repo.ExistsByID(ID)
}

func (s *DefaultPlantService) ExistsByCode(code string, ID uint) bool {
	return s.repo.ExistsByCode(code, ID)
}

func (s *DefaultPlantService) ExistsByName(name string, ID uint) bool {
	return s.repo.ExistsByName(name, ID)
}

func (s *DefaultPlantService) ToModel(plantForm *plant.Form) *models.Plant {
	plantModel := &models.Plant{
		Code:   plantForm.Code,
		Name:   plantForm.Name,
		Status: plantForm.Status,
	}
	plantModel.ID = plantForm.ID
	return plantModel
}

func (s *DefaultPlantService) FormToModel(plantForm *plant.Form, plantModel *models.Plant) {
	plantModel.Code = plantForm.Code
	plantModel.Name = plantForm.Name
	plantModel.Status = plantForm.Status
}

func (s *DefaultPlantService) ToForm(plantModel *models.Plant) *plant.Form {
	plantForm := &plant.Form{
		ID:     plantModel.ID,
		Code:   plantModel.Code,
		Name:   plantModel.Name,
		Status: plantModel.Status,
	}
	return plantForm
}

func (s *DefaultPlantService) ToFormSlice(plantModels []*models.Plant) []*plant.Form {
	data := make([]*plant.Form, 0)
	for _, plantModel := range plantModels {
		data = append(data, s.ToForm(plantModel))
	}
	return data
}

func (s *DefaultPlantService) ToModelSlice(plantForms []*plant.Form) []*models.Plant {
	data := make([]*models.Plant, 0)
	for _, plantForm := range plantForms {
		data = append(data, s.ToModel(plantForm))
	}
	return data
}
