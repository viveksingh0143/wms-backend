package service

import (
	"star-wms/app/base/dto/container"
	"star-wms/app/base/models"
	"star-wms/app/base/repository"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/common/responses"
)

type ContainerService interface {
	GetAllContainers(plantID uint, filter container.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*container.Form, int64, error)
	CreateContainer(plantID uint, containerForm *container.Form) error
	GetContainerByID(plantID uint, id uint) (*container.Form, error)
	UpdateContainer(plantID uint, id uint, containerForm *container.Form) error
	DeleteContainer(plantID uint, id uint) error
	DeleteContainers(plantID uint, ids []uint) error
	ExistsById(plantID uint, ID uint) bool
	ExistsByName(plantID uint, name string, ID uint) bool
	ExistsByCode(plantID uint, code string, ID uint) bool
	ToModel(plantID uint, containerForm *container.Form) *models.Container
	FormToModel(plantID uint, containerForm *container.Form, containerModel *models.Container)
	ToForm(plantID uint, containerModel *models.Container) *container.Form
	ToFormSlice(plantID uint, containerModels []*models.Container) []*container.Form
	ToModelSlice(plantID uint, containerForms []*container.Form) []*models.Container
}

type DefaultContainerService struct {
	repo           repository.ContainerRepository
	storeService   StoreService
	productService ProductService
}

func NewContainerService(repo repository.ContainerRepository, storeService StoreService, productService ProductService) ContainerService {
	return &DefaultContainerService{repo: repo, storeService: storeService, productService: productService}
}

func (s *DefaultContainerService) GetAllContainers(plantID uint, filter container.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*container.Form, int64, error) {
	data, count, err := s.repo.GetAll(plantID, filter, pagination, sorting)
	if err != nil {
		return nil, count, err
	}
	return s.ToFormSlice(plantID, data), count, err
}

func (s *DefaultContainerService) CreateContainer(plantID uint, containerForm *container.Form) error {
	if s.ExistsByName(plantID, containerForm.Name, 0) {
		return responses.NewInputError("name", "already exists", containerForm.Name)
	}
	if s.ExistsByCode(plantID, containerForm.Code, 0) {
		return responses.NewInputError("code", "already exists", containerForm.Code)
	}
	containerModel := s.ToModel(plantID, containerForm)
	return s.repo.Create(plantID, containerModel)
}

func (s *DefaultContainerService) GetContainerByID(plantID uint, id uint) (*container.Form, error) {
	data, err := s.repo.GetByID(plantID, id)
	if err != nil {
		return nil, err
	}
	return s.ToForm(plantID, data), nil
}

func (s *DefaultContainerService) UpdateContainer(plantID uint, id uint, containerForm *container.Form) error {
	if s.ExistsByName(plantID, containerForm.Name, id) {
		return responses.NewInputError("name", "already exists", containerForm.Name)
	}
	if s.ExistsByCode(plantID, containerForm.Code, id) {
		return responses.NewInputError("code", "already exists", containerForm.Code)
	}
	containerModel, err := s.repo.GetByID(plantID, id)
	if err != nil {
		return err
	}
	s.FormToModel(plantID, containerForm, containerModel)
	return s.repo.Update(plantID, containerModel)
}

func (s *DefaultContainerService) DeleteContainer(plantID uint, id uint) error {
	return s.repo.Delete(plantID, id)
}

func (s *DefaultContainerService) DeleteContainers(plantID uint, ids []uint) error {
	return s.repo.DeleteMulti(plantID, ids)
}

func (s *DefaultContainerService) ExistsById(plantID uint, ID uint) bool {
	return s.repo.ExistsByID(plantID, ID)
}

func (s *DefaultContainerService) ExistsByName(plantID uint, name string, ID uint) bool {
	return s.repo.ExistsByName(plantID, name, ID)
}

func (s *DefaultContainerService) ExistsByCode(plantID uint, code string, ID uint) bool {
	return s.repo.ExistsByCode(plantID, code, ID)
}

func (s *DefaultContainerService) ToModel(plantID uint, containerForm *container.Form) *models.Container {
	containerModel := &models.Container{
		ContainerType: models.ContainerType(containerForm.ContainerType),
		Name:          containerForm.Name,
		Code:          containerForm.Code,
		Address:       containerForm.Address,
		Status:        containerForm.Status,
	}
	containerModel.ID = containerForm.ID
	if containerForm.Approved {
		containerModel.Approved = containerForm.Approved
	}
	if containerForm.StockLevel != "" {
		containerModel.StockLevel = containerForm.StockLevel
	}
	if containerForm.Store != nil {
		containerModel.Store = s.storeService.ToModel(plantID, containerForm.Store)
	}
	if containerForm.Product != nil {
		containerModel.Product = s.productService.ToModel(containerForm.Product)
	}
	containerModel.PlantID = plantID
	return containerModel
}

func (s *DefaultContainerService) FormToModel(plantID uint, containerForm *container.Form, containerModel *models.Container) {
	containerModel.ContainerType = models.ContainerType(containerForm.ContainerType)
	containerModel.Name = containerForm.Name
	containerModel.Code = containerForm.Code
	containerModel.Address = containerForm.Address
	containerModel.Status = containerForm.Status
	if containerForm.Approved {
		containerModel.Approved = containerForm.Approved
	}
	if containerForm.StockLevel != "" {
		containerModel.StockLevel = containerForm.StockLevel
	}
	if containerForm.Store != nil {
		containerModel.Store = s.storeService.ToModel(plantID, containerForm.Store)
	} else {
		containerModel.Store = nil
		containerModel.StoreID = nil
	}

	if containerForm.Product != nil {
		containerModel.Product = s.productService.ToModel(containerForm.Product)
	} else {
		containerModel.Product = nil
		containerModel.ProductID = nil
	}
}

func (s *DefaultContainerService) ToForm(plantID uint, containerModel *models.Container) *container.Form {
	containerForm := &container.Form{
		ID:            containerModel.ID,
		ContainerType: string(containerModel.ContainerType),
		Name:          containerModel.Name,
		Code:          containerModel.Code,
		Address:       containerModel.Address,
		Status:        containerModel.Status,
		StockLevel:    containerModel.StockLevel,
		Approved:      containerModel.Approved,
	}
	containerForm.PlantID = containerModel.PlantID
	if containerModel.Store != nil {
		containerForm.Store = s.storeService.ToForm(plantID, containerModel.Store)
	}
	if containerModel.Product != nil {
		containerForm.Product = s.productService.ToForm(containerModel.Product)
	}
	return containerForm
}

func (s *DefaultContainerService) ToFormSlice(plantID uint, containerModels []*models.Container) []*container.Form {
	data := make([]*container.Form, 0)
	for _, containerModel := range containerModels {
		data = append(data, s.ToForm(plantID, containerModel))
	}
	return data
}

func (s *DefaultContainerService) ToModelSlice(plantID uint, containerForms []*container.Form) []*models.Container {
	data := make([]*models.Container, 0)
	for _, containerForm := range containerForms {
		data = append(data, s.ToModel(plantID, containerForm))
	}
	return data
}
