package service

import (
	baseService "star-wms/app/base/service"
	"star-wms/app/warehouse/dto/inventory"
	"star-wms/app/warehouse/models"
	"star-wms/app/warehouse/repository"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/common/responses"
)

type InventoryService interface {
	GetAllInventorys(plantID uint, filter inventory.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*inventory.Form, int64, error)
	CreateInventory(plantID uint, inventoryForm *inventory.Form) error
	GetInventoryByID(plantID uint, id uint) (*inventory.Form, error)
	UpdateInventory(plantID uint, id uint, inventoryForm *inventory.Form) error
	DeleteInventory(plantID uint, id uint) error
	DeleteInventorys(plantID uint, ids []uint) error
	ExistsById(plantID uint, ID uint) bool
	ToModel(plantID uint, inventoryForm *inventory.Form) *models.Inventory
	FormToModel(plantID uint, inventoryForm *inventory.Form, inventoryModel *models.Inventory)
	ToForm(plantID uint, inventoryModel *models.Inventory) *inventory.Form
	ToFormSlice(plantID uint, inventoryModels []*models.Inventory) []*inventory.Form
	ToModelSlice(plantID uint, inventoryForms []*inventory.Form) []*models.Inventory
}

type DefaultInventoryService struct {
	repo           repository.InventoryRepository
	productService baseService.ProductService
	storeService   baseService.StoreService
}

func NewInventoryService(repo repository.InventoryRepository, productService baseService.ProductService, storeService baseService.StoreService) InventoryService {
	return &DefaultInventoryService{repo: repo, productService: productService, storeService: storeService}
}

func (s *DefaultInventoryService) GetAllInventorys(plantID uint, filter inventory.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*inventory.Form, int64, error) {
	data, count, err := s.repo.GetAll(plantID, filter, pagination, sorting)
	if err != nil {
		return nil, count, err
	}
	return s.ToFormSlice(plantID, data), count, err
}

func (s *DefaultInventoryService) CreateInventory(plantID uint, inventoryForm *inventory.Form) error {
	if !s.storeService.ExistsById(plantID, inventoryForm.Store.ID) {
		return responses.NewInputError("store.id", "store not exists", inventoryForm.Store.ID)
	}
	if !s.productService.ExistsById(inventoryForm.Product.ID) {
		return responses.NewInputError("product.id", "product not exists", inventoryForm.Product.ID)
	}
	inventoryModel := s.ToModel(plantID, inventoryForm)
	return s.repo.Create(plantID, inventoryModel)
}

func (s *DefaultInventoryService) GetInventoryByID(plantID uint, id uint) (*inventory.Form, error) {
	data, err := s.repo.GetByID(plantID, id)
	if err != nil {
		return nil, err
	}
	return s.ToForm(plantID, data), nil
}

func (s *DefaultInventoryService) UpdateInventory(plantID uint, id uint, inventoryForm *inventory.Form) error {
	if !s.storeService.ExistsById(plantID, inventoryForm.Store.ID) {
		return responses.NewInputError("store.id", "store not exists", inventoryForm.Store.ID)
	}
	if !s.productService.ExistsById(inventoryForm.Product.ID) {
		return responses.NewInputError("product.id", "product not exists", inventoryForm.Product.ID)
	}
	inventoryModel, err := s.repo.GetByID(plantID, id)
	if err != nil {
		return err
	}
	s.FormToModel(plantID, inventoryForm, inventoryModel)
	return s.repo.Update(plantID, inventoryModel)
}

func (s *DefaultInventoryService) DeleteInventory(plantID uint, id uint) error {
	return s.repo.Delete(plantID, id)
}

func (s *DefaultInventoryService) DeleteInventorys(plantID uint, ids []uint) error {
	return s.repo.DeleteMulti(plantID, ids)
}

func (s *DefaultInventoryService) ExistsById(plantID uint, ID uint) bool {
	return s.repo.ExistsByID(plantID, ID)
}

func (s *DefaultInventoryService) ToModel(plantID uint, inventoryForm *inventory.Form) *models.Inventory {
	inventoryModel := &models.Inventory{
		Quantity: inventoryForm.Quantity,
	}
	inventoryModel.ID = inventoryForm.ID

	if inventoryForm.Product != nil {
		inventoryModel.Product = s.productService.ToModel(inventoryForm.Product)
	}
	if inventoryForm.Store != nil {
		inventoryModel.Store = s.storeService.ToModel(plantID, inventoryForm.Store)
	}
	return inventoryModel
}

func (s *DefaultInventoryService) FormToModel(plantID uint, inventoryForm *inventory.Form, inventoryModel *models.Inventory) {
	inventoryModel.Quantity = inventoryForm.Quantity

	if inventoryForm.Product != nil {
		inventoryModel.Product = s.productService.ToModel(inventoryForm.Product)
	}
	if inventoryForm.Store != nil {
		inventoryModel.Store = s.storeService.ToModel(plantID, inventoryForm.Store)
	}
}

func (s *DefaultInventoryService) ToForm(plantID uint, inventoryModel *models.Inventory) *inventory.Form {
	inventoryForm := &inventory.Form{
		ID:       inventoryModel.ID,
		Quantity: inventoryModel.Quantity,
	}

	if inventoryModel.Product != nil {
		inventoryForm.Product = s.productService.ToForm(inventoryModel.Product)
	}
	if inventoryModel.Store != nil {
		inventoryForm.Store = s.storeService.ToForm(plantID, inventoryModel.Store)
	}
	return inventoryForm
}

func (s *DefaultInventoryService) ToFormSlice(plantID uint, inventoryModels []*models.Inventory) []*inventory.Form {
	data := make([]*inventory.Form, 0)
	for _, inventoryModel := range inventoryModels {
		data = append(data, s.ToForm(plantID, inventoryModel))
	}
	return data
}

func (s *DefaultInventoryService) ToModelSlice(plantID uint, inventoryForms []*inventory.Form) []*models.Inventory {
	data := make([]*models.Inventory, 0)
	for _, inventoryForm := range inventoryForms {
		data = append(data, s.ToModel(plantID, inventoryForm))
	}
	return data
}
