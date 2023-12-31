package service

import (
	"star-wms/app/base/dto/container"
	"star-wms/app/base/models"
	"star-wms/app/base/repository"
	"star-wms/core/common/dto"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/common/responses"
)

type ContainerService interface {
	GetAllContainersRequiredApproval(plantID uint, filter container.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*container.Form, int64, error)
	GetContainersReports(plantID uint, filter container.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*container.Form, int64, error)
	GetAllContainers(plantID uint, filter container.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*container.Form, int64, error)
	GetStatistics(plantID uint, filter container.Filter) *container.Statistics
	CreateContainer(plantID uint, containerForm *container.Form) error
	GetContainerByID(plantID uint, id uint) (*container.Form, error)
	GetContainerByCode(plantID uint, code string, needContents bool, needProduct bool, needStore bool, needLocation bool) (*container.Form, error)
	UpdateContainer(plantID uint, id uint, containerForm *container.Form) error
	MarkedContainerFull(plantID uint, containerForm *container.Form) error
	DeleteContainer(plantID uint, id uint) error
	DeleteContainers(plantID uint, ids []uint) error
	ExistsById(plantID uint, ID uint) bool
	ExistsByName(plantID uint, name string, ID uint) bool
	ExistsByCode(plantID uint, code string, ID uint) bool
	ToModel(plantID uint, containerForm *container.Form) *models.Container
	FormToModel(plantID uint, containerForm *container.Form, containerModel *models.Container)
	ToModelSlice(plantID uint, containerForms []*container.Form) []*models.Container
	ToForm(plantID uint, containerModel *models.Container) *container.Form
	ToFormSlice(plantID uint, containerModels []*models.Container) []*container.Form
	ToContentForm(plantID uint, containerModel *models.ContainerContent) *container.ContentForm
	ToContentFormSlice(plantID uint, containerModels []*models.ContainerContent) []*container.ContentForm
	ApproveContainer(plantID uint, id uint) error
	ApproveContainers(plantID uint, ids []uint) error

	ReportStockLevels(plantID uint, filter container.Filter) []*dto.ReportDto
	ReportApprovals(plantID uint, filter container.Filter) []*dto.ReportDto
}

type DefaultContainerService struct {
	repo                 repository.ContainerRepository
	storeService         StoreService
	storelocationService StorelocationService
	productService       ProductService
}

func NewContainerService(repo repository.ContainerRepository, storeService StoreService, storelocationService StorelocationService, productService ProductService) ContainerService {
	return &DefaultContainerService{repo: repo, storeService: storeService, storelocationService: storelocationService, productService: productService}
}

func (s *DefaultContainerService) GetAllContainersRequiredApproval(plantID uint, filter container.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*container.Form, int64, error) {
	data, count, err := s.repo.GetAllRequiredApproval(plantID, filter, pagination, sorting)
	if err != nil {
		return nil, count, err
	}
	return s.ToFormSlice(plantID, data), count, err
}

func (s *DefaultContainerService) GetContainersReports(plantID uint, filter container.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*container.Form, int64, error) {
	data, count, err := s.repo.GetAll(plantID, filter, pagination, sorting, true, true, true)
	if err != nil {
		return nil, count, err
	}
	return s.ToFormSlice(plantID, data), count, err
}

func (s *DefaultContainerService) GetAllContainers(plantID uint, filter container.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*container.Form, int64, error) {
	data, count, err := s.repo.GetAll(plantID, filter, pagination, sorting, false, false, false)
	if err != nil {
		return nil, count, err
	}
	return s.ToFormSlice(plantID, data), count, err
}

func (s *DefaultContainerService) GetStatistics(plantID uint, filter container.Filter) *container.Statistics {
	emptyCount, partialCount, fullCount, waitingCount := s.repo.GetStatistics(plantID, filter)
	return &container.Statistics{
		Empty:              emptyCount,
		Partial:            partialCount,
		Full:               fullCount,
		WaitingForApproval: waitingCount,
	}
}

func (s *DefaultContainerService) CreateContainer(plantID uint, containerForm *container.Form) error {
	//if s.ExistsByName(plantID, containerForm.Name, 0) {
	//	return responses.NewInputError("name", "already exists", containerForm.Name)
	//}
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

func (s *DefaultContainerService) GetContainerByCode(plantID uint, code string, needContents bool, needProduct bool, needStore bool, needLocation bool) (*container.Form, error) {
	data, err := s.repo.GetByCode(plantID, code, needContents, needProduct, needStore, needLocation)
	if err != nil {
		return nil, err
	}
	return s.ToForm(plantID, data), nil
}

func (s *DefaultContainerService) UpdateContainer(plantID uint, id uint, containerForm *container.Form) error {
	//if s.ExistsByName(plantID, containerForm.Name, id) {
	//	return responses.NewInputError("name", "already exists", containerForm.Name)
	//}
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

func (s *DefaultContainerService) MarkedContainerFull(plantID uint, containerForm *container.Form) error {
	return s.repo.MarkedContainerFull(plantID, containerForm.ID)
}

func (s *DefaultContainerService) ApproveContainer(plantID uint, id uint) error {
	return s.repo.Approve(plantID, id)
}

func (s *DefaultContainerService) ApproveContainers(plantID uint, ids []uint) error {
	return s.repo.ApproveMulti(plantID, ids)
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
	containerModel.Approved = containerForm.Approved
	if containerForm.StockLevel != "" {
		containerModel.StockLevel = containerForm.StockLevel
	}
	if containerForm.Store != nil {
		containerModel.Store = s.storeService.ToModel(plantID, containerForm.Store)
	}
	if containerForm.Storelocation != nil {
		containerModel.Storelocation = s.storelocationService.ToModel(plantID, containerForm.Storelocation.StoreID, containerForm.Storelocation)
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
	containerModel.Approved = containerForm.Approved
	if containerForm.StockLevel != "" {
		containerModel.StockLevel = containerForm.StockLevel
	}

	if containerForm.Store != nil {
		containerModel.Store = s.storeService.ToModel(plantID, containerForm.Store)
	} else {
		containerModel.Store = nil
		containerModel.StoreID = nil
	}

	if containerForm.Storelocation != nil {
		containerModel.Storelocation = s.storelocationService.ToModel(plantID, containerForm.Storelocation.StoreID, containerForm.Storelocation)
	} else {
		containerModel.Storelocation = nil
		containerModel.StorelocationID = nil
	}

	if containerForm.Product != nil {
		containerModel.Product = s.productService.ToModel(containerForm.Product)
	} else {
		containerModel.Product = nil
		containerModel.ProductID = nil
	}
}

func (s *DefaultContainerService) ToModelSlice(plantID uint, containerForms []*container.Form) []*models.Container {
	data := make([]*models.Container, 0)
	for _, containerForm := range containerForms {
		data = append(data, s.ToModel(plantID, containerForm))
	}
	return data
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
		StoreID:       containerModel.StoreID,
	}
	containerForm.PlantID = containerModel.PlantID
	containerForm.ProductID = containerModel.ProductID
	if containerModel.Storelocation != nil {
		containerForm.Storelocation = s.storelocationService.ToForm(plantID, containerModel.Storelocation.StoreID, containerModel.Storelocation)
	}
	if containerModel.Store != nil {
		containerForm.Store = s.storeService.ToForm(plantID, containerModel.Store)
	}
	if containerModel.Product != nil {
		containerForm.Product = s.productService.ToForm(containerModel.Product)
	}
	if containerModel.Contents != nil {
		containerForm.Contents = s.ToContentFormSlice(plantID, containerModel.Contents)
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

func (s *DefaultContainerService) ToContentForm(plantID uint, contentModel *models.ContainerContent) *container.ContentForm {
	contentForm := &container.ContentForm{
		ID:       contentModel.ID,
		Quantity: contentModel.Quantity,
		Barcode:  contentModel.Barcode,
	}
	contentForm.PlantID = contentModel.PlantID
	if contentModel.Product != nil {
		contentForm.Product = s.productService.ToForm(contentModel.Product)
	}
	return contentForm
}

func (s *DefaultContainerService) ToContentFormSlice(plantID uint, contentModels []*models.ContainerContent) []*container.ContentForm {
	data := make([]*container.ContentForm, 0)
	for _, contentModel := range contentModels {
		data = append(data, s.ToContentForm(plantID, contentModel))
	}
	return data
}

func (s *DefaultContainerService) ReportStockLevels(plantID uint, filter container.Filter) []*dto.ReportDto {
	return s.repo.GetReportStockLevels(plantID, filter)
}

func (s *DefaultContainerService) ReportApprovals(plantID uint, filter container.Filter) []*dto.ReportDto {
	return s.repo.GetReportApprovals(plantID, filter)
}
