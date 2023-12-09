package service

import (
	"star-wms/app/base/dto/requisition"
	"star-wms/app/base/dto/store"
	"star-wms/app/base/models"
	"star-wms/app/base/repository"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/common/responses"
	"star-wms/core/types"
)

type RequisitionService interface {
	GetAllRequisitionsRequiredApproval(plantID uint, stores []*store.Form, filter requisition.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*requisition.Form, int64, error)
	GetAllRequisitions(plantID uint, filter requisition.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*requisition.Form, int64, error)
	CreateRequisition(plantID uint, requisitionForm *requisition.Form) error
	GetRequisitionByID(plantID uint, id uint) (*requisition.Form, error)
	UpdateRequisition(plantID uint, id uint, requisitionForm *requisition.Form) error
	DeleteRequisition(plantID uint, id uint) error
	DeleteRequisitions(plantID uint, ids []uint) error
	ApproveRequisition(plantID uint, id uint) error
	ApproveRequisitions(plantID uint, ids []uint) error
	ExistsByItemId(requisitionID uint, ID uint) bool
	ExistsById(plantID uint, ID uint) bool
	ExistsByOrderNo(plantID uint, orderNo string, ID uint) bool
	ToModel(plantID uint, requisitionForm *requisition.Form) *models.Requisition
	FormToModel(plantID uint, requisitionForm *requisition.Form, requisitionModel *models.Requisition)
	ToForm(plantID uint, requisitionModel *models.Requisition) *requisition.Form
	ToFormSlice(plantID uint, requisitionModels []*models.Requisition) []*requisition.Form
	ToModelSlice(plantID uint, requisitionForms []*requisition.Form) []*models.Requisition
	ToItemFormSlice(plantID uint, itemModels []*models.RequisitionItem) []*requisition.ItemsForm
	ToItemForm(plantID uint, itemModel *models.RequisitionItem) *requisition.ItemsForm
	ToItemModel(plantID uint, itemForm *requisition.ItemsForm) *models.RequisitionItem
	ToItemModelSlice(plantID uint, itemForms []*requisition.ItemsForm) []*models.RequisitionItem
}

type DefaultRequisitionService struct {
	repo           repository.RequisitionRepository
	storeService   StoreService
	productService ProductService
}

func NewRequisitionService(repo repository.RequisitionRepository, storeService StoreService, productService ProductService) RequisitionService {
	return &DefaultRequisitionService{repo: repo, storeService: storeService, productService: productService}
}

func (s *DefaultRequisitionService) GetAllRequisitionsRequiredApproval(plantID uint, stores []*store.Form, filter requisition.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*requisition.Form, int64, error) {
	data, count, err := s.repo.GetAllRequiredApproval(plantID, stores, filter, pagination, sorting)
	if err != nil {
		return nil, count, err
	}
	return s.ToFormSlice(plantID, data), count, err
}

func (s *DefaultRequisitionService) GetAllRequisitions(plantID uint, filter requisition.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*requisition.Form, int64, error) {
	data, count, err := s.repo.GetAll(plantID, filter, pagination, sorting)
	if err != nil {
		return nil, count, err
	}
	return s.ToFormSlice(plantID, data), count, err
}

func (s *DefaultRequisitionService) CreateRequisition(plantID uint, requisitionForm *requisition.Form) error {
	if s.ExistsByOrderNo(plantID, requisitionForm.OrderNo, 0) {
		return responses.NewInputError("order_no", "already exists", requisitionForm.OrderNo)
	}
	if requisitionForm.Store != nil {
		if !s.storeService.ExistsById(plantID, requisitionForm.Store.ID) {
			return responses.NewInputError("store.id", "store not exists", requisitionForm.Store.ID)
		}
	}
	requisitionForm.Approved = types.ApprovalWait
	requisitionModel := s.ToModel(plantID, requisitionForm)
	return s.repo.Create(plantID, requisitionModel)
}

func (s *DefaultRequisitionService) GetRequisitionByID(plantID uint, id uint) (*requisition.Form, error) {
	data, err := s.repo.GetByID(plantID, id)
	if err != nil {
		return nil, err
	}
	return s.ToForm(plantID, data), nil
}

func (s *DefaultRequisitionService) UpdateRequisition(plantID uint, id uint, requisitionForm *requisition.Form) error {
	if s.ExistsByOrderNo(plantID, requisitionForm.OrderNo, id) {
		return responses.NewInputError("order_no", "already exists", requisitionForm.OrderNo)
	}
	if requisitionForm.Store != nil {
		if !s.storeService.ExistsById(plantID, requisitionForm.Store.ID) {
			return responses.NewInputError("store.id", "store not exists", requisitionForm.Store.ID)
		}
	}
	requisitionModel, err := s.repo.GetByID(plantID, id)
	if err != nil {
		return err
	}
	s.FormToModel(plantID, requisitionForm, requisitionModel)
	return s.repo.Update(plantID, requisitionModel)
}

func (s *DefaultRequisitionService) DeleteRequisition(plantID uint, id uint) error {
	return s.repo.Delete(plantID, id)
}

func (s *DefaultRequisitionService) DeleteRequisitions(plantID uint, ids []uint) error {
	return s.repo.DeleteMulti(plantID, ids)
}

func (s *DefaultRequisitionService) ApproveRequisition(plantID uint, id uint) error {
	return s.repo.Approve(plantID, id)
}

func (s *DefaultRequisitionService) ApproveRequisitions(plantID uint, ids []uint) error {
	return s.repo.ApproveMulti(plantID, ids)
}

func (s *DefaultRequisitionService) ExistsByItemId(requisitionID uint, ID uint) bool {
	return s.repo.ExistsByItemId(requisitionID, ID)
}

func (s *DefaultRequisitionService) ExistsById(plantID uint, ID uint) bool {
	return s.repo.ExistsByID(plantID, ID)
}

func (s *DefaultRequisitionService) ExistsByOrderNo(plantID uint, orderNo string, ID uint) bool {
	return s.repo.ExistsByOrderNo(plantID, orderNo, ID)
}

func (s *DefaultRequisitionService) ToModel(plantID uint, requisitionForm *requisition.Form) *models.Requisition {
	requisitionModel := &models.Requisition{
		IssuedDate: requisitionForm.IssuedDate,
		OrderNo:    requisitionForm.OrderNo,
		Department: requisitionForm.Department,
		Status:     requisitionForm.Status,
	}
	requisitionModel.ID = requisitionForm.ID

	if requisitionForm.Store != nil {
		requisitionModel.Store = s.storeService.ToModel(plantID, requisitionForm.Store)
	}

	if requisitionForm.Items != nil {
		requisitionModel.Items = s.ToItemModelSlice(plantID, requisitionForm.Items)
	}
	return requisitionModel
}

func (s *DefaultRequisitionService) FormToModel(plantID uint, requisitionForm *requisition.Form, requisitionModel *models.Requisition) {
	requisitionModel.IssuedDate = requisitionForm.IssuedDate
	requisitionModel.OrderNo = requisitionForm.OrderNo
	requisitionModel.Department = requisitionForm.Department
	requisitionModel.Status = requisitionForm.Status

	if requisitionForm.Store != nil {
		requisitionModel.Store = s.storeService.ToModel(plantID, requisitionForm.Store)
	}
	if requisitionForm.Items != nil {
		requisitionModel.Items = s.ToItemModelSlice(plantID, requisitionForm.Items)
	} else {
		requisitionModel.Items = make([]*models.RequisitionItem, 0)
	}
}

func (s *DefaultRequisitionService) ToForm(plantID uint, requisitionModel *models.Requisition) *requisition.Form {
	requisitionForm := &requisition.Form{
		ID:         requisitionModel.ID,
		IssuedDate: requisitionModel.IssuedDate,
		OrderNo:    requisitionModel.OrderNo,
		Department: requisitionModel.Department,
		Status:     requisitionModel.Status,
		Approved:   requisitionModel.Approved,
	}
	if requisitionModel.Store != nil {
		requisitionForm.Store = s.storeService.ToForm(plantID, requisitionModel.Store)
	}
	requisitionForm.Items = s.ToItemFormSlice(plantID, requisitionModel.Items)
	return requisitionForm
}

func (s *DefaultRequisitionService) ToFormSlice(plantID uint, requisitionModels []*models.Requisition) []*requisition.Form {
	data := make([]*requisition.Form, 0)
	for _, requisitionModel := range requisitionModels {
		data = append(data, s.ToForm(plantID, requisitionModel))
	}
	return data
}

func (s *DefaultRequisitionService) ToModelSlice(plantID uint, requisitionForms []*requisition.Form) []*models.Requisition {
	data := make([]*models.Requisition, 0)
	for _, requisitionForm := range requisitionForms {
		data = append(data, s.ToModel(plantID, requisitionForm))
	}
	return data
}

func (s *DefaultRequisitionService) ToItemForm(plantID uint, itemModel *models.RequisitionItem) *requisition.ItemsForm {
	ingredientForm := &requisition.ItemsForm{
		ID:        itemModel.ID,
		ProductID: itemModel.ProductID,
		Product:   s.productService.ToForm(itemModel.Product),
		Quantity:  itemModel.Quantity,
	}
	return ingredientForm
}

func (s *DefaultRequisitionService) ToItemFormSlice(plantID uint, itemModels []*models.RequisitionItem) []*requisition.ItemsForm {
	data := make([]*requisition.ItemsForm, 0)
	for _, itemModel := range itemModels {
		data = append(data, s.ToItemForm(plantID, itemModel))
	}
	return data
}

func (s *DefaultRequisitionService) ToItemModel(plantID uint, itemForm *requisition.ItemsForm) *models.RequisitionItem {
	product := &models.Product{}
	product.ID = itemForm.ProductID

	requisitionItemModel := &models.RequisitionItem{
		ProductID: itemForm.ProductID,
		Product:   product,
		Quantity:  itemForm.Quantity,
	}
	requisitionItemModel.ID = itemForm.ID
	return requisitionItemModel
}

func (s *DefaultRequisitionService) ToItemModelSlice(plantID uint, itemForms []*requisition.ItemsForm) []*models.RequisitionItem {
	data := make([]*models.RequisitionItem, 0)
	for _, itemForm := range itemForms {
		data = append(data, s.ToItemModel(plantID, itemForm))
	}
	return data
}
