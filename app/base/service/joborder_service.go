package service

import (
	"star-wms/app/base/dto/joborder"
	"star-wms/app/base/models"
	"star-wms/app/base/repository"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/common/responses"
)

type JobOrderService interface {
	GetAllJobOrders(plantID uint, filter joborder.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*joborder.Form, int64, error)
	CreateJobOrder(plantID uint, joborderForm *joborder.Form) error
	GetJobOrderByID(plantID uint, id uint) (*joborder.Form, error)
	UpdateJobOrder(plantID uint, id uint, joborderForm *joborder.Form) error
	DeleteJobOrder(plantID uint, id uint) error
	DeleteJobOrders(plantID uint, ids []uint) error
	ExistsById(plantID uint, ID uint) bool
	ExistsByOrderNo(plantID uint, orderNo string, ID uint) bool
	ToModel(plantID uint, joborderForm *joborder.Form) *models.JobOrder
	FormToModel(plantID uint, joborderForm *joborder.Form, joborderModel *models.JobOrder)
	ToForm(plantID uint, joborderModel *models.JobOrder) *joborder.Form
	ToFormSlice(plantID uint, joborderModels []*models.JobOrder) []*joborder.Form
	ToModelSlice(plantID uint, joborderForms []*joborder.Form) []*models.JobOrder
}

type DefaultJobOrderService struct {
	repo            repository.JobOrderRepository
	customerService CustomerService
	productService  ProductService
}

func NewJobOrderService(repo repository.JobOrderRepository, customerService CustomerService, productService ProductService) JobOrderService {
	return &DefaultJobOrderService{repo: repo, customerService: customerService, productService: productService}
}

func (s *DefaultJobOrderService) GetAllJobOrders(plantID uint, filter joborder.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*joborder.Form, int64, error) {
	data, count, err := s.repo.GetAll(plantID, filter, pagination, sorting)
	if err != nil {
		return nil, count, err
	}
	return s.ToFormSlice(plantID, data), count, err
}

func (s *DefaultJobOrderService) CreateJobOrder(plantID uint, joborderForm *joborder.Form) error {
	if s.ExistsByOrderNo(plantID, joborderForm.OrderNo, 0) {
		return responses.NewInputError("order_no", "already exists", joborderForm.OrderNo)
	}
	if joborderForm.Customer != nil {
		if !s.customerService.ExistsById(plantID, joborderForm.Customer.ID) {
			return responses.NewInputError("customer.id", "customer not exists", joborderForm.Customer.ID)
		}
	}
	joborderModel := s.ToModel(plantID, joborderForm)
	return s.repo.Create(plantID, joborderModel)
}

func (s *DefaultJobOrderService) GetJobOrderByID(plantID uint, id uint) (*joborder.Form, error) {
	data, err := s.repo.GetByID(plantID, id)
	if err != nil {
		return nil, err
	}
	return s.ToForm(plantID, data), nil
}

func (s *DefaultJobOrderService) UpdateJobOrder(plantID uint, id uint, joborderForm *joborder.Form) error {
	if s.ExistsByOrderNo(plantID, joborderForm.OrderNo, id) {
		return responses.NewInputError("order_no", "already exists", joborderForm.OrderNo)
	}
	if joborderForm.Customer != nil {
		if !s.customerService.ExistsById(plantID, joborderForm.Customer.ID) {
			return responses.NewInputError("customer.id", "customer not exists", joborderForm.Customer.ID)
		}
	}
	joborderModel, err := s.repo.GetByID(plantID, id)
	if err != nil {
		return err
	}
	s.FormToModel(plantID, joborderForm, joborderModel)
	return s.repo.Update(plantID, joborderModel)
}

func (s *DefaultJobOrderService) DeleteJobOrder(plantID uint, id uint) error {
	return s.repo.Delete(plantID, id)
}

func (s *DefaultJobOrderService) DeleteJobOrders(plantID uint, ids []uint) error {
	return s.repo.DeleteMulti(plantID, ids)
}

func (s *DefaultJobOrderService) ExistsById(plantID uint, ID uint) bool {
	return s.repo.ExistsByID(plantID, ID)
}

func (s *DefaultJobOrderService) ExistsByOrderNo(plantID uint, orderNo string, ID uint) bool {
	return s.repo.ExistsByOrderNo(plantID, orderNo, ID)
}

func (s *DefaultJobOrderService) ToModel(plantID uint, joborderForm *joborder.Form) *models.JobOrder {
	joborderModel := &models.JobOrder{
		IssuedDate: joborderForm.IssuedDate,
		OrderNo:    joborderForm.OrderNo,
		POCategory: models.POCategory(joborderForm.POCategory),
		Status:     joborderForm.Status,
	}
	joborderModel.ID = joborderForm.ID

	if joborderForm.Customer != nil {
		joborderModel.Customer = s.customerService.ToModel(plantID, joborderForm.Customer)
	}

	if joborderForm.Items != nil {
		joborderModel.Items = s.ToItemModelSlice(plantID, joborderForm.Items)
	}
	return joborderModel
}

func (s *DefaultJobOrderService) FormToModel(plantID uint, joborderForm *joborder.Form, joborderModel *models.JobOrder) {
	joborderModel.IssuedDate = joborderForm.IssuedDate
	joborderModel.OrderNo = joborderForm.OrderNo
	joborderModel.POCategory = models.POCategory(joborderForm.POCategory)
	joborderModel.Status = joborderForm.Status

	if joborderForm.Customer != nil {
		joborderModel.Customer = s.customerService.ToModel(plantID, joborderForm.Customer)
	} else {
		joborderModel.Customer = nil
		joborderModel.CustomerID = nil
	}
	if joborderForm.Items != nil {
		joborderModel.Items = s.ToItemModelSlice(plantID, joborderForm.Items)
	} else {
		joborderModel.Items = make([]*models.JobOrderItem, 0)
	}
}

func (s *DefaultJobOrderService) ToForm(plantID uint, joborderModel *models.JobOrder) *joborder.Form {
	joborderForm := &joborder.Form{
		ID:         joborderModel.ID,
		IssuedDate: joborderModel.IssuedDate,
		OrderNo:    joborderModel.OrderNo,
		POCategory: string(joborderModel.POCategory),
		Status:     joborderModel.Status,
	}
	if joborderModel.Customer != nil {
		joborderForm.Customer = s.customerService.ToForm(plantID, joborderModel.Customer)
	}
	joborderForm.Items = s.ToItemFormSlice(plantID, joborderModel.Items)
	return joborderForm
}

func (s *DefaultJobOrderService) ToFormSlice(plantID uint, joborderModels []*models.JobOrder) []*joborder.Form {
	data := make([]*joborder.Form, 0)
	for _, joborderModel := range joborderModels {
		data = append(data, s.ToForm(plantID, joborderModel))
	}
	return data
}

func (s *DefaultJobOrderService) ToModelSlice(plantID uint, joborderForms []*joborder.Form) []*models.JobOrder {
	data := make([]*models.JobOrder, 0)
	for _, joborderForm := range joborderForms {
		data = append(data, s.ToModel(plantID, joborderForm))
	}
	return data
}

func (s *DefaultJobOrderService) ToItemForm(plantID uint, itemModel *models.JobOrderItem) *joborder.ItemsForm {
	ingredientForm := &joborder.ItemsForm{
		ProductID: itemModel.ProductID,
		Product:   s.productService.ToForm(itemModel.Product),
		Quantity:  itemModel.Quantity,
	}
	return ingredientForm
}

func (s *DefaultJobOrderService) ToItemFormSlice(plantID uint, itemModels []*models.JobOrderItem) []*joborder.ItemsForm {
	data := make([]*joborder.ItemsForm, 0)
	for _, itemModel := range itemModels {
		data = append(data, s.ToItemForm(plantID, itemModel))
	}
	return data
}

func (s *DefaultJobOrderService) ToItemModel(plantID uint, itemForm *joborder.ItemsForm) *models.JobOrderItem {
	product := &models.Product{}
	product.ID = itemForm.ProductID

	jobOrderItemModel := &models.JobOrderItem{
		ProductID: itemForm.ProductID,
		Product:   product,
		Quantity:  itemForm.Quantity,
	}
	return jobOrderItemModel
}

func (s *DefaultJobOrderService) ToItemModelSlice(plantID uint, itemForms []*joborder.ItemsForm) []*models.JobOrderItem {
	data := make([]*models.JobOrderItem, 0)
	for _, itemForm := range itemForms {
		data = append(data, s.ToItemModel(plantID, itemForm))
	}
	return data
}
