package service

import (
	"fmt"
	"star-wms/app/base/dto/container"
	baseModels "star-wms/app/base/models"
	baseService "star-wms/app/base/service"
	"star-wms/app/warehouse/dto/batchlabel"
	"star-wms/app/warehouse/dto/inventory"
	"star-wms/app/warehouse/models"
	"star-wms/app/warehouse/repository"
	"star-wms/core/app"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/common/responses"
	"star-wms/core/types"
)

type InventoryService interface {
	CreateRawMaterialStockIn(plantID uint, inventoryForm *inventory.RawMaterialStockInForm) error
	CreateFinishedGoodsStockIn(plantID uint, inventoryForm *inventory.FinishedGoodsStockInForm) error
	CreateFinishedGoodStockIn(plantID uint, inventoryForm *inventory.FinishedGoodStockInForm) (*batchlabel.StickerForm, error)
	AttachContainerToLocation(plantID uint, attachForm *inventory.AttachContainerForm) error
	GetAllInventorys(plantID uint, filter inventory.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*inventory.Form, int64, error)
	GetInventoryByID(plantID uint, id uint) (*inventory.Form, error)
	ToModel(plantID uint, inventoryForm *inventory.Form) *models.Inventory
	FormToModel(plantID uint, inventoryForm *inventory.Form, inventoryModel *models.Inventory)
	ToForm(plantID uint, inventoryModel *models.Inventory) *inventory.Form
	ToFormSlice(plantID uint, inventoryModels []*models.Inventory) []*inventory.Form
	ToModelSlice(plantID uint, inventoryForms []*inventory.Form) []*models.Inventory

	CreateInventory(plantID uint, inventoryForm *inventory.Form) error
	UpdateInventory(plantID uint, id uint, inventoryForm *inventory.Form) error
	DeleteInventory(plantID uint, id uint) error
	DeleteInventorys(plantID uint, ids []uint) error
	ExistsById(plantID uint, ID uint) bool
}

type DefaultInventoryService struct {
	repo              repository.InventoryRepository
	productService    baseService.ProductService
	storeService      baseService.StoreService
	locationService   baseService.StorelocationService
	containerService  baseService.ContainerService
	batchlabelService BatchlabelService
}

func NewInventoryService(repo repository.InventoryRepository, productService baseService.ProductService, storeService baseService.StoreService, locationService baseService.StorelocationService, containerService baseService.ContainerService, batchlabelService BatchlabelService) InventoryService {
	return &DefaultInventoryService{repo: repo, productService: productService, storeService: storeService, locationService: locationService, containerService: containerService, batchlabelService: batchlabelService}
}

func (s *DefaultInventoryService) GetAllInventorys(plantID uint, filter inventory.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*inventory.Form, int64, error) {
	data, count, err := s.repo.GetAll(plantID, filter, pagination, sorting)
	if err != nil {
		return nil, count, err
	}
	return s.ToFormSlice(plantID, data), count, err
}

func (s *DefaultInventoryService) GetInventoryByID(plantID uint, id uint) (*inventory.Form, error) {
	data, err := s.repo.GetByID(plantID, id)
	if err != nil {
		return nil, err
	}
	return s.ToForm(plantID, data), nil
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

func (s *DefaultInventoryService) CreateRawMaterialStockIn(plantID uint, inventoryForm *inventory.RawMaterialStockInForm) error {
	if !s.storeService.ExistsById(plantID, inventoryForm.Store.ID) {
		return responses.NewInputError("store.id", "store not exists", inventoryForm.Store.ID)
	}
	if !s.productService.ExistsById(inventoryForm.Product.ID) {
		return responses.NewInputError("product.id", "product not exists", inventoryForm.Product.ID)
	}
	if inventoryForm.Container.Code == "" {
		return responses.NewInputError("container.code", "code is required", inventoryForm.Container.Code)
	}
	var containerForm *container.Form
	if !s.containerService.ExistsByCode(plantID, inventoryForm.Container.Code, 0) {
		containerForm = &container.Form{
			PlantID:       plantID,
			ContainerType: string(baseModels.Pallet),
			Name:          inventoryForm.Container.Code,
			Code:          inventoryForm.Container.Code,
		}
		err := s.containerService.CreateContainer(plantID, containerForm)
		if err != nil {
			return err
		}
	}
	containerForm, err := s.containerService.GetContainerByCode(plantID, inventoryForm.Container.Code, false, false, false, false)
	if err != nil {
		return err
	}
	if containerForm.StockLevel != baseModels.Empty {
		return responses.NewInputError("container.code", "is not empty", inventoryForm.Container.Code)
	}
	storeModel := s.storeService.ToModel(plantID, inventoryForm.Store)
	containerModel := s.containerService.ToModel(plantID, containerForm)

	rmBatchModel := &models.RMBatch{
		ProductID:   inventoryForm.Product.ID,
		Product:     s.productService.ToModel(inventoryForm.Product),
		Quantity:    inventoryForm.Quantity,
		BatchNumber: inventoryForm.BatchNo,
	}
	return s.repo.CreateRawMaterialStockIn(plantID, storeModel, containerModel, rmBatchModel)
}

func (s *DefaultInventoryService) CreateFinishedGoodsStockIn(plantID uint, inventoryForm *inventory.FinishedGoodsStockInForm) error {
	if !s.storeService.ExistsById(plantID, inventoryForm.Store.ID) {
		return responses.NewInputError("store.id", "store not exists", inventoryForm.Store.ID)
	}
	if len(inventoryForm.Barcodes) <= 0 {
		return responses.NewInputError("barcodes", "no barcode found", inventoryForm.Barcodes)
	}

	containerForm, err := s.containerService.GetContainerByCode(plantID, inventoryForm.ContainerCode, false, false, false, false)
	if err != nil {
		return responses.NewInputError("container_code", fmt.Sprintf("givem (%s) container not exists", inventoryForm.ContainerCode), inventoryForm.ContainerCode)
	} else if containerForm.StockLevel == baseModels.Full {
		return responses.NewInputError("container_code", fmt.Sprintf("givem (%s) container is already full", inventoryForm.ContainerCode), inventoryForm.ContainerCode)
	}

	stickerForms := make([]*batchlabel.StickerForm, 0)
	for _, barcode := range inventoryForm.Barcodes {
		sticker, err := s.batchlabelService.GetStickerByBarcodePlantwise(plantID, barcode)
		if err != nil {
			return responses.NewInputError("barcodes", fmt.Sprintf("given (%s) barcode not exists", barcode), inventoryForm.Barcodes)
		} else if sticker.IsUsed {
			return responses.NewInputError("barcodes", fmt.Sprintf("givem (%s) barcode already used", barcode), inventoryForm.Barcodes)
		} else if containerForm.ProductID != nil && sticker.ProductID != *containerForm.ProductID {
			return responses.NewInputError("barcodes", fmt.Sprintf("barcode (%s) having different product than container", barcode), inventoryForm.Barcodes)
		}
		stickerForms = append(stickerForms, sticker)
	}
	storeModel := s.storeService.ToModel(plantID, inventoryForm.Store)
	containerModel := s.containerService.ToModel(plantID, containerForm)
	stickerModels := s.batchlabelService.ToStickerModelFormSlice(plantID, stickerForms)
	return s.repo.CreateFinishedGoodsStockIn(plantID, storeModel, containerModel, stickerModels)
}

func (s *DefaultInventoryService) CreateFinishedGoodStockIn(plantID uint, inventoryForm *inventory.FinishedGoodStockInForm) (*batchlabel.StickerForm, error) {
	store, err := s.storeService.GetStoreByCode(plantID, app.APP_FG_STORE_CODE)

	if err != nil {
		return nil, responses.NewInputError("barcode", "No finished goods store found", app.APP_FG_STORE_CODE)
	}
	containerForm, err := s.containerService.GetContainerByCode(plantID, inventoryForm.ContainerCode, false, false, false, false)
	if err != nil {
		return nil, responses.NewInputError("container_code", fmt.Sprintf("givem (%s) container not exists", inventoryForm.ContainerCode), inventoryForm.ContainerCode)
	} else if containerForm.StockLevel == baseModels.Full {
		return nil, responses.NewInputError("container_code", fmt.Sprintf("givem (%s) container is already full", inventoryForm.ContainerCode), inventoryForm.ContainerCode)
	}

	stickerForms := make([]*batchlabel.StickerForm, 0)
	sticker, err := s.batchlabelService.GetStickerByBarcodePlantwise(plantID, inventoryForm.Barcode)
	if err != nil {
		return nil, responses.NewInputError("barcodes", fmt.Sprintf("given (%s) barcode not exists", inventoryForm.Barcode), inventoryForm.Barcode)
	} else if sticker.IsUsed {
		return nil, responses.NewInputError("barcodes", fmt.Sprintf("givem (%s) barcode already used", inventoryForm.Barcode), inventoryForm.Barcode)
	} else if containerForm.ProductID != nil && sticker.ProductID != *containerForm.ProductID {
		return nil, responses.NewInputError("barcodes", fmt.Sprintf("barcode (%s) having different product than container", inventoryForm.Barcode), inventoryForm.Barcode)
	}
	stickerForms = append(stickerForms, sticker)
	storeModel := s.storeService.ToModel(plantID, store)
	containerModel := s.containerService.ToModel(plantID, containerForm)
	stickerModels := s.batchlabelService.ToStickerModelFormSlice(plantID, stickerForms)
	return sticker, s.repo.CreateFinishedGoodsStockIn(plantID, storeModel, containerModel, stickerModels)
}

func (s *DefaultInventoryService) AttachContainerToLocation(plantID uint, attachForm *inventory.AttachContainerForm) error {
	if !s.containerService.ExistsByCode(plantID, attachForm.ContainerCode, 0) {
		return responses.NewInputError("container_code", "container not exists", attachForm.ContainerCode)
	}
	if !s.locationService.ExistsByOnlyCode(plantID, attachForm.LocationCode) {
		return responses.NewInputError("location_code", "location not exists", attachForm.LocationCode)
	}
	locationForm, err := s.locationService.GetStorelocationByCode(plantID, attachForm.LocationCode)
	if err != nil {
		return err
	}
	if locationForm.Status != types.FillStatusEmpty {
		return responses.NewInputError("location_code", "is not empty", attachForm.LocationCode)
	}

	containerForm, err := s.containerService.GetContainerByCode(plantID, attachForm.ContainerCode, false, false, false, true)
	if err != nil {
		return err
	}
	if containerForm.StockLevel == baseModels.Empty {
		return responses.NewInputError("container_code", "is empty", attachForm.ContainerCode)
	} else if containerForm.Storelocation != nil {
		return responses.NewInputError("container_code", "is already attached", attachForm.ContainerCode)
	} else if containerForm.Approved != types.ApprovalYes {
		return responses.NewInputError("container_code", "is not approved yet", attachForm.ContainerCode)
	} else if containerForm.StoreID == nil || *containerForm.StoreID != locationForm.StoreID {
		return responses.NewInputError("container_code", "is not approved for selected store", attachForm.ContainerCode)
	}
	containerModel := s.containerService.ToModel(plantID, containerForm)
	locationModel := s.locationService.ToModel(plantID, locationForm.StoreID, locationForm)
	return s.repo.AttachContainerToLocation(plantID, containerModel, locationModel)
}
