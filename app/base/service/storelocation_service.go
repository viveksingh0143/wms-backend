package service

import (
	"star-wms/app/base/dto/storelocation"
	"star-wms/app/base/models"
	"star-wms/app/base/repository"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/common/responses"
)

type StorelocationService interface {
	GetAllStorelocations(plantID uint, storeID uint, filter storelocation.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*storelocation.Form, int64, error)
	CreateStorelocation(plantID uint, storeID uint, storelocationForm *storelocation.Form) error
	GetStorelocationByID(plantID uint, storeID uint, id uint) (*storelocation.Form, error)
	GetStorelocationByCode(plantID uint, code string) (*storelocation.Form, error)
	UpdateStorelocation(plantID uint, storeID uint, id uint, storelocationForm *storelocation.Form) error
	DeleteStorelocation(plantID uint, storeID uint, id uint) error
	DeleteStorelocations(plantID uint, storeID uint, ids []uint) error
	ExistsById(plantID uint, storeID uint, ID uint) bool
	ExistsByCode(plantID uint, storeID uint, code string, ID uint) bool
	ExistsByOnlyCode(plantID uint, code string) bool
	ToModel(plantID uint, storeID uint, storelocationForm *storelocation.Form) *models.Storelocation
	FormToModel(plantID uint, storeID uint, storelocationForm *storelocation.Form, storelocationModel *models.Storelocation)
	ToForm(plantID uint, storeID uint, storelocationModel *models.Storelocation) *storelocation.Form
	ToFormSlice(plantID uint, storeID uint, storelocationModels []*models.Storelocation) []*storelocation.Form
	ToModelSlice(plantID uint, storeID uint, storelocationForms []*storelocation.Form) []*models.Storelocation
}

type DefaultStorelocationService struct {
	repo         repository.StorelocationRepository
	storeService StoreService
}

func NewStorelocationService(repo repository.StorelocationRepository, storeService StoreService) StorelocationService {
	return &DefaultStorelocationService{repo: repo, storeService: storeService}
}

func (s *DefaultStorelocationService) GetAllStorelocations(plantID uint, storeID uint, filter storelocation.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*storelocation.Form, int64, error) {
	data, count, err := s.repo.GetAll(plantID, storeID, filter, pagination, sorting)
	if err != nil {
		return nil, count, err
	}
	return s.ToFormSlice(plantID, storeID, data), count, err
}

func (s *DefaultStorelocationService) CreateStorelocation(plantID uint, storeID uint, storelocationForm *storelocation.Form) error {
	if s.ExistsByCode(plantID, storeID, storelocationForm.Code, 0) {
		return responses.NewInputError("code", "already exists", storelocationForm.Code)
	}
	storelocationModel := s.ToModel(plantID, storeID, storelocationForm)
	return s.repo.Create(plantID, storeID, storelocationModel)
}

func (s *DefaultStorelocationService) GetStorelocationByID(plantID uint, storeID uint, id uint) (*storelocation.Form, error) {
	data, err := s.repo.GetByID(plantID, storeID, id)
	if err != nil {
		return nil, err
	}
	return s.ToForm(plantID, storeID, data), nil
}

func (s *DefaultStorelocationService) GetStorelocationByCode(plantID uint, code string) (*storelocation.Form, error) {
	data, err := s.repo.GetByCode(plantID, code)
	if err != nil {
		return nil, err
	}
	return s.ToForm(plantID, data.StoreID, data), nil
}

func (s *DefaultStorelocationService) UpdateStorelocation(plantID uint, storeID uint, id uint, storelocationForm *storelocation.Form) error {
	if s.ExistsByCode(plantID, storeID, storelocationForm.Code, id) {
		return responses.NewInputError("code", "already exists", storelocationForm.Code)
	}
	storelocationModel, err := s.repo.GetByID(plantID, storeID, id)
	if err != nil {
		return err
	}
	s.FormToModel(plantID, storeID, storelocationForm, storelocationModel)
	return s.repo.Update(plantID, storeID, storelocationModel)
}

func (s *DefaultStorelocationService) DeleteStorelocation(plantID uint, storeID uint, id uint) error {
	return s.repo.Delete(plantID, storeID, id)
}

func (s *DefaultStorelocationService) DeleteStorelocations(plantID uint, storeID uint, ids []uint) error {
	return s.repo.DeleteMulti(plantID, storeID, ids)
}

func (s *DefaultStorelocationService) ExistsById(plantID uint, storeID uint, ID uint) bool {
	return s.repo.ExistsByID(plantID, storeID, ID)
}

func (s *DefaultStorelocationService) ExistsByCode(plantID uint, storeID uint, code string, ID uint) bool {
	return s.repo.ExistsByCode(plantID, storeID, code, ID)
}

func (s *DefaultStorelocationService) ExistsByOnlyCode(plantID uint, code string) bool {
	return s.repo.ExistsByOnlyCode(plantID, code)
}

func (s *DefaultStorelocationService) ToModel(plantID uint, storeID uint, storelocationForm *storelocation.Form) *models.Storelocation {
	storelocationModel := &models.Storelocation{
		Code:        storelocationForm.Code,
		ZoneName:    storelocationForm.ZoneName,
		AisleNumber: storelocationForm.AisleNumber,
		RackNumber:  storelocationForm.RackNumber,
		ShelfNumber: storelocationForm.ShelfNumber,
		Description: storelocationForm.Description,
		Status:      storelocationForm.Status,
	}
	storelocationModel.ID = storelocationForm.ID
	storelocationModel.PlantID = plantID
	storelocationModel.StoreID = storeID
	return storelocationModel
}

func (s *DefaultStorelocationService) FormToModel(plantID uint, storeID uint, storelocationForm *storelocation.Form, storelocationModel *models.Storelocation) {
	storelocationModel.Code = storelocationForm.Code
	storelocationModel.ZoneName = storelocationForm.ZoneName
	storelocationModel.AisleNumber = storelocationForm.AisleNumber
	storelocationModel.RackNumber = storelocationForm.RackNumber
	storelocationModel.ShelfNumber = storelocationForm.ShelfNumber
	storelocationModel.Description = storelocationForm.Description
	storelocationModel.Status = storelocationForm.Status
}

func (s *DefaultStorelocationService) ToForm(plantID uint, storeID uint, storelocationModel *models.Storelocation) *storelocation.Form {
	storelocationForm := &storelocation.Form{
		ID:          storelocationModel.ID,
		Code:        storelocationModel.Code,
		ZoneName:    storelocationModel.ZoneName,
		AisleNumber: storelocationModel.AisleNumber,
		RackNumber:  storelocationModel.RackNumber,
		ShelfNumber: storelocationModel.ShelfNumber,
		Description: storelocationModel.Description,
		Status:      storelocationModel.Status,
	}
	storelocationForm.PlantID = storelocationModel.PlantID
	storelocationForm.StoreID = storelocationModel.StoreID
	if storelocationModel.Store != nil {
		storelocationForm.Store = s.storeService.ToForm(plantID, storelocationModel.Store)
	}
	return storelocationForm
}

func (s *DefaultStorelocationService) ToFormSlice(plantID uint, storeID uint, storelocationModels []*models.Storelocation) []*storelocation.Form {
	data := make([]*storelocation.Form, 0)
	for _, storelocationModel := range storelocationModels {
		data = append(data, s.ToForm(plantID, storeID, storelocationModel))
	}
	return data
}

func (s *DefaultStorelocationService) ToModelSlice(plantID uint, storeID uint, storelocationForms []*storelocation.Form) []*models.Storelocation {
	data := make([]*models.Storelocation, 0)
	for _, storelocationForm := range storelocationForms {
		data = append(data, s.ToModel(plantID, storeID, storelocationForm))
	}
	return data
}
