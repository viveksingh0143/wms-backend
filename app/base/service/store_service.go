package service

import (
	"star-wms/app/admin/dto/user"
	adminModels "star-wms/app/admin/models"
	"star-wms/app/admin/service"
	"star-wms/app/base/dto/store"
	"star-wms/app/base/models"
	"star-wms/app/base/repository"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/common/responses"
)

type StoreService interface {
	GetAllStores(plantID uint, filter store.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*store.Form, int64, error)
	CreateStore(plantID uint, storeForm *store.Form) error
	GetStoreByID(plantID uint, id uint) (*store.Form, error)
	UpdateStore(plantID uint, id uint, storeForm *store.Form) error
	DeleteStore(plantID uint, id uint) error
	DeleteStores(plantID uint, ids []uint) error
	ExistsById(plantID uint, ID uint) bool
	ExistsByName(plantID uint, name string, ID uint) bool
	ExistsByCode(plantID uint, code string, ID uint) bool
	ToModel(plantID uint, storeForm *store.Form) *models.Store
	FormToModel(plantID uint, storeForm *store.Form, storeModel *models.Store)
	ToForm(plantID uint, storeModel *models.Store) *store.Form
	ToFormSlice(plantID uint, storeModels []*models.Store) []*store.Form
	ToModelSlice(plantID uint, storeForms []*store.Form) []*models.Store
}

type DefaultStoreService struct {
	repo            repository.StoreRepository
	categoryService CategoryService
	userService     service.UserService
}

func NewStoreService(repo repository.StoreRepository, categoryService CategoryService, userService service.UserService) StoreService {
	return &DefaultStoreService{repo: repo, categoryService: categoryService, userService: userService}
}

func (s *DefaultStoreService) GetAllStores(plantID uint, filter store.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*store.Form, int64, error) {
	data, count, err := s.repo.GetAll(plantID, filter, pagination, sorting)
	if err != nil {
		return nil, count, err
	}
	return s.ToFormSlice(plantID, data), count, err
}

func (s *DefaultStoreService) CreateStore(plantID uint, storeForm *store.Form) error {
	if s.ExistsByName(plantID, storeForm.Name, 0) {
		return responses.NewInputError("name", "already exists", storeForm.Name)
	}
	if s.ExistsByCode(plantID, storeForm.Code, 0) {
		return responses.NewInputError("code", "already exists", storeForm.Code)
	}
	if storeForm.Category != nil {
		if !s.categoryService.ExistsById(storeForm.Category.ID) {
			return responses.NewInputError("category.id", "category not exists", storeForm.Category.ID)
		}
	}
	storeModel := s.ToModel(plantID, storeForm)
	if storeForm.Category != nil {
		categoryForm, _ := s.categoryService.GetCategoryShortInfoByID(storeForm.Category.ID)
		storeModel.CategoryPath = categoryForm.FullPath
	}
	return s.repo.Create(plantID, storeModel)
}

func (s *DefaultStoreService) GetStoreByID(plantID uint, id uint) (*store.Form, error) {
	data, err := s.repo.GetByID(plantID, id)
	if err != nil {
		return nil, err
	}
	return s.ToForm(plantID, data), nil
}

func (s *DefaultStoreService) UpdateStore(plantID uint, id uint, storeForm *store.Form) error {
	if s.ExistsByName(plantID, storeForm.Name, id) {
		return responses.NewInputError("name", "already exists", storeForm.Name)
	}
	if s.ExistsByCode(plantID, storeForm.Code, id) {
		return responses.NewInputError("code", "already exists", storeForm.Code)
	}
	if storeForm.Category != nil {
		if !s.categoryService.ExistsById(storeForm.Category.ID) {
			return responses.NewInputError("category.id", "category not exists", storeForm.Category.ID)
		}
	}
	storeModel, err := s.repo.GetByID(plantID, id)
	if err != nil {
		return err
	}
	s.FormToModel(plantID, storeForm, storeModel)
	if storeForm.Category != nil {
		categoryForm, _ := s.categoryService.GetCategoryShortInfoByID(storeForm.Category.ID)
		storeModel.CategoryPath = categoryForm.FullPath
	}
	return s.repo.Update(plantID, storeModel)
}

func (s *DefaultStoreService) DeleteStore(plantID uint, id uint) error {
	return s.repo.Delete(plantID, id)
}

func (s *DefaultStoreService) DeleteStores(plantID uint, ids []uint) error {
	return s.repo.DeleteMulti(plantID, ids)
}

func (s *DefaultStoreService) ExistsById(plantID uint, ID uint) bool {
	return s.repo.ExistsByID(plantID, ID)
}

func (s *DefaultStoreService) ExistsByName(plantID uint, name string, ID uint) bool {
	return s.repo.ExistsByName(plantID, name, ID)
}

func (s *DefaultStoreService) ExistsByCode(plantID uint, code string, ID uint) bool {
	return s.repo.ExistsByCode(plantID, code, ID)
}

func (s *DefaultStoreService) ToModel(plantID uint, storeForm *store.Form) *models.Store {
	storeModel := &models.Store{
		Name:    storeForm.Name,
		Code:    storeForm.Code,
		Address: storeForm.Address,
		Status:  storeForm.Status,
	}
	storeModel.ID = storeForm.ID
	if storeForm.Category != nil {
		storeModel.Category = s.categoryService.ToModel(storeForm.Category)
		storeModel.CategoryPath = storeModel.Category.FullPath
	} else {
		storeModel.CategoryPath = ""
	}
	storeModel.PlantID = plantID
	if storeForm.Approvers != nil {
		approvers := make([]*adminModels.User, 0)
		if len(storeForm.Approvers) > 0 {
			approvers = s.userService.ToModelSlice(storeForm.Approvers)
		}
		storeModel.Approvers = approvers
	}
	return storeModel
}

func (s *DefaultStoreService) FormToModel(plantID uint, storeForm *store.Form, storeModel *models.Store) {
	storeModel.Name = storeForm.Name
	storeModel.Code = storeForm.Code
	storeModel.Address = storeForm.Address
	storeModel.Status = storeForm.Status

	if storeForm.Category != nil {
		storeModel.Category = s.categoryService.ToModel(storeForm.Category)
		storeModel.CategoryPath = storeModel.Category.FullPath
	} else {
		storeModel.Category = nil
		storeModel.CategoryID = nil
		storeModel.CategoryPath = ""
	}

	if storeForm.Approvers != nil {
		approvers := make([]*adminModels.User, 0)
		if len(storeForm.Approvers) > 0 {
			approvers = s.userService.ToModelSlice(storeForm.Approvers)
		}
		storeModel.Approvers = approvers
	}
}

func (s *DefaultStoreService) ToForm(plantID uint, storeModel *models.Store) *store.Form {
	storeForm := &store.Form{
		ID:           storeModel.ID,
		Name:         storeModel.Name,
		Code:         storeModel.Code,
		Address:      storeModel.Address,
		Status:       storeModel.Status,
		CategoryPath: storeModel.CategoryPath,
	}
	storeForm.PlantID = storeModel.PlantID
	if storeModel.Category != nil {
		storeForm.Category = s.categoryService.ToForm(storeModel.Category)
	}
	if storeModel.Approvers != nil {
		approvers := make([]*user.Form, 0)
		if len(storeModel.Approvers) > 0 {
			approvers = s.userService.ToFormSlice(storeModel.Approvers)
		}
		storeForm.Approvers = approvers
	}
	return storeForm
}

func (s *DefaultStoreService) ToFormSlice(plantID uint, storeModels []*models.Store) []*store.Form {
	data := make([]*store.Form, 0)
	for _, storeModel := range storeModels {
		data = append(data, s.ToForm(plantID, storeModel))
	}
	return data
}

func (s *DefaultStoreService) ToModelSlice(plantID uint, storeForms []*store.Form) []*models.Store {
	data := make([]*models.Store, 0)
	for _, storeForm := range storeForms {
		data = append(data, s.ToModel(plantID, storeForm))
	}
	return data
}
