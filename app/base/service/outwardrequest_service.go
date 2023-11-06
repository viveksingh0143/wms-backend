package service

import (
	"star-wms/app/base/dto/outwardrequest"
	"star-wms/app/base/models"
	"star-wms/app/base/repository"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/common/responses"
)

type OutwardrequestService interface {
	GetAllOutwardrequests(plantID uint, filter outwardrequest.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*outwardrequest.Form, int64, error)
	CreateOutwardrequest(plantID uint, outwardrequestForm *outwardrequest.Form) error
	GetOutwardrequestByID(plantID uint, id uint) (*outwardrequest.Form, error)
	UpdateOutwardrequest(plantID uint, id uint, outwardrequestForm *outwardrequest.Form) error
	DeleteOutwardrequest(plantID uint, id uint) error
	DeleteOutwardrequests(plantID uint, ids []uint) error
	ExistsByItemId(outwardrequestID uint, ID uint) bool
	ExistsById(plantID uint, ID uint) bool
	ExistsByOrderNo(plantID uint, orderNo string, ID uint) bool
	ToModel(plantID uint, outwardrequestForm *outwardrequest.Form) *models.Outwardrequest
	FormToModel(plantID uint, outwardrequestForm *outwardrequest.Form, outwardrequestModel *models.Outwardrequest)
	ToForm(plantID uint, outwardrequestModel *models.Outwardrequest) *outwardrequest.Form
	ToFormSlice(plantID uint, outwardrequestModels []*models.Outwardrequest) []*outwardrequest.Form
	ToModelSlice(plantID uint, outwardrequestForms []*outwardrequest.Form) []*models.Outwardrequest
	ToItemFormSlice(plantID uint, itemModels []*models.OutwardrequestItem) []*outwardrequest.ItemsForm
	ToItemForm(plantID uint, itemModel *models.OutwardrequestItem) *outwardrequest.ItemsForm
	ToItemModel(plantID uint, itemForm *outwardrequest.ItemsForm) *models.OutwardrequestItem
	ToItemModelSlice(plantID uint, itemForms []*outwardrequest.ItemsForm) []*models.OutwardrequestItem
}

type DefaultOutwardrequestService struct {
	repo            repository.OutwardrequestRepository
	customerService CustomerService
	productService  ProductService
}

func NewOutwardrequestService(repo repository.OutwardrequestRepository, customerService CustomerService, productService ProductService) OutwardrequestService {
	return &DefaultOutwardrequestService{repo: repo, customerService: customerService, productService: productService}
}

func (s *DefaultOutwardrequestService) GetAllOutwardrequests(plantID uint, filter outwardrequest.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*outwardrequest.Form, int64, error) {
	data, count, err := s.repo.GetAll(plantID, filter, pagination, sorting)
	if err != nil {
		return nil, count, err
	}
	return s.ToFormSlice(plantID, data), count, err
}

func (s *DefaultOutwardrequestService) CreateOutwardrequest(plantID uint, outwardrequestForm *outwardrequest.Form) error {
	if s.ExistsByOrderNo(plantID, outwardrequestForm.OrderNo, 0) {
		return responses.NewInputError("order_no", "already exists", outwardrequestForm.OrderNo)
	}
	if outwardrequestForm.Customer != nil {
		if !s.customerService.ExistsById(plantID, outwardrequestForm.Customer.ID) {
			return responses.NewInputError("customer.id", "customer not exists", outwardrequestForm.Customer.ID)
		}
	}
	outwardrequestModel := s.ToModel(plantID, outwardrequestForm)
	return s.repo.Create(plantID, outwardrequestModel)
}

func (s *DefaultOutwardrequestService) GetOutwardrequestByID(plantID uint, id uint) (*outwardrequest.Form, error) {
	data, err := s.repo.GetByID(plantID, id)
	if err != nil {
		return nil, err
	}
	return s.ToForm(plantID, data), nil
}

func (s *DefaultOutwardrequestService) UpdateOutwardrequest(plantID uint, id uint, outwardrequestForm *outwardrequest.Form) error {
	if s.ExistsByOrderNo(plantID, outwardrequestForm.OrderNo, id) {
		return responses.NewInputError("order_no", "already exists", outwardrequestForm.OrderNo)
	}
	if outwardrequestForm.Customer != nil {
		if !s.customerService.ExistsById(plantID, outwardrequestForm.Customer.ID) {
			return responses.NewInputError("customer.id", "customer not exists", outwardrequestForm.Customer.ID)
		}
	}
	outwardrequestModel, err := s.repo.GetByID(plantID, id)
	if err != nil {
		return err
	}
	s.FormToModel(plantID, outwardrequestForm, outwardrequestModel)
	return s.repo.Update(plantID, outwardrequestModel)
}

func (s *DefaultOutwardrequestService) DeleteOutwardrequest(plantID uint, id uint) error {
	return s.repo.Delete(plantID, id)
}

func (s *DefaultOutwardrequestService) DeleteOutwardrequests(plantID uint, ids []uint) error {
	return s.repo.DeleteMulti(plantID, ids)
}

func (s *DefaultOutwardrequestService) ExistsByItemId(outwardrequestID uint, ID uint) bool {
	return s.repo.ExistsByItemId(outwardrequestID, ID)
}

func (s *DefaultOutwardrequestService) ExistsById(plantID uint, ID uint) bool {
	return s.repo.ExistsByID(plantID, ID)
}

func (s *DefaultOutwardrequestService) ExistsByOrderNo(plantID uint, orderNo string, ID uint) bool {
	return s.repo.ExistsByOrderNo(plantID, orderNo, ID)
}

func (s *DefaultOutwardrequestService) ToModel(plantID uint, outwardrequestForm *outwardrequest.Form) *models.Outwardrequest {
	outwardrequestModel := &models.Outwardrequest{
		IssuedDate: outwardrequestForm.IssuedDate,
		OrderNo:    outwardrequestForm.OrderNo,
		Status:     outwardrequestForm.Status,
	}
	outwardrequestModel.ID = outwardrequestForm.ID

	if outwardrequestForm.Customer != nil {
		outwardrequestModel.Customer = s.customerService.ToModel(plantID, outwardrequestForm.Customer)
	}

	if outwardrequestForm.Items != nil {
		outwardrequestModel.Items = s.ToItemModelSlice(plantID, outwardrequestForm.Items)
	}
	return outwardrequestModel
}

func (s *DefaultOutwardrequestService) FormToModel(plantID uint, outwardrequestForm *outwardrequest.Form, outwardrequestModel *models.Outwardrequest) {
	outwardrequestModel.IssuedDate = outwardrequestForm.IssuedDate
	outwardrequestModel.OrderNo = outwardrequestForm.OrderNo
	outwardrequestModel.Status = outwardrequestForm.Status

	if outwardrequestForm.Customer != nil {
		outwardrequestModel.Customer = s.customerService.ToModel(plantID, outwardrequestForm.Customer)
	}
	if outwardrequestForm.Items != nil {
		outwardrequestModel.Items = s.ToItemModelSlice(plantID, outwardrequestForm.Items)
	} else {
		outwardrequestModel.Items = make([]*models.OutwardrequestItem, 0)
	}
}

func (s *DefaultOutwardrequestService) ToForm(plantID uint, outwardrequestModel *models.Outwardrequest) *outwardrequest.Form {
	outwardrequestForm := &outwardrequest.Form{
		ID:         outwardrequestModel.ID,
		IssuedDate: outwardrequestModel.IssuedDate,
		OrderNo:    outwardrequestModel.OrderNo,
		Status:     outwardrequestModel.Status,
	}
	if outwardrequestModel.Customer != nil {
		outwardrequestForm.Customer = s.customerService.ToForm(plantID, outwardrequestModel.Customer)
	}
	outwardrequestForm.Items = s.ToItemFormSlice(plantID, outwardrequestModel.Items)
	return outwardrequestForm
}

func (s *DefaultOutwardrequestService) ToFormSlice(plantID uint, outwardrequestModels []*models.Outwardrequest) []*outwardrequest.Form {
	data := make([]*outwardrequest.Form, 0)
	for _, outwardrequestModel := range outwardrequestModels {
		data = append(data, s.ToForm(plantID, outwardrequestModel))
	}
	return data
}

func (s *DefaultOutwardrequestService) ToModelSlice(plantID uint, outwardrequestForms []*outwardrequest.Form) []*models.Outwardrequest {
	data := make([]*models.Outwardrequest, 0)
	for _, outwardrequestForm := range outwardrequestForms {
		data = append(data, s.ToModel(plantID, outwardrequestForm))
	}
	return data
}

func (s *DefaultOutwardrequestService) ToItemForm(plantID uint, itemModel *models.OutwardrequestItem) *outwardrequest.ItemsForm {
	ingredientForm := &outwardrequest.ItemsForm{
		ID:        itemModel.ID,
		ProductID: itemModel.ProductID,
		Product:   s.productService.ToForm(itemModel.Product),
		Quantity:  itemModel.Quantity,
	}
	return ingredientForm
}

func (s *DefaultOutwardrequestService) ToItemFormSlice(plantID uint, itemModels []*models.OutwardrequestItem) []*outwardrequest.ItemsForm {
	data := make([]*outwardrequest.ItemsForm, 0)
	for _, itemModel := range itemModels {
		data = append(data, s.ToItemForm(plantID, itemModel))
	}
	return data
}

func (s *DefaultOutwardrequestService) ToItemModel(plantID uint, itemForm *outwardrequest.ItemsForm) *models.OutwardrequestItem {
	product := &models.Product{}
	product.ID = itemForm.ProductID

	outwardrequestItemModel := &models.OutwardrequestItem{
		ProductID: itemForm.ProductID,
		Product:   product,
		Quantity:  itemForm.Quantity,
	}
	outwardrequestItemModel.ID = itemForm.ID
	return outwardrequestItemModel
}

func (s *DefaultOutwardrequestService) ToItemModelSlice(plantID uint, itemForms []*outwardrequest.ItemsForm) []*models.OutwardrequestItem {
	data := make([]*models.OutwardrequestItem, 0)
	for _, itemForm := range itemForms {
		data = append(data, s.ToItemModel(plantID, itemForm))
	}
	return data
}
