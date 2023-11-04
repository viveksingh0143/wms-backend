package service

import (
	"star-wms/app/base/dto/joborder"
	"star-wms/app/base/models"
	"star-wms/app/base/repository"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/common/responses"
)

type JoborderService interface {
	GetAllJoborders(plantID uint, filter joborder.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*joborder.Form, int64, error)
	CreateJoborder(plantID uint, joborderForm *joborder.Form) error
	GetJoborderByID(plantID uint, id uint) (*joborder.Form, error)
	UpdateJoborder(plantID uint, id uint, joborderForm *joborder.Form) error
	DeleteJoborder(plantID uint, id uint) error
	DeleteJoborders(plantID uint, ids []uint) error
	ExistsByItemId(joborderID uint, ID uint) bool
	ExistsById(plantID uint, ID uint) bool
	ExistsByOrderNo(plantID uint, orderNo string, ID uint) bool
	ToModel(plantID uint, joborderForm *joborder.Form) *models.Joborder
	FormToModel(plantID uint, joborderForm *joborder.Form, joborderModel *models.Joborder)
	ToForm(plantID uint, joborderModel *models.Joborder) *joborder.Form
	ToFormSlice(plantID uint, joborderModels []*models.Joborder) []*joborder.Form
	ToModelSlice(plantID uint, joborderForms []*joborder.Form) []*models.Joborder
	ToItemFormSlice(plantID uint, itemModels []*models.JoborderItem) []*joborder.ItemsForm
	ToItemForm(plantID uint, itemModel *models.JoborderItem) *joborder.ItemsForm
	ToItemModel(plantID uint, itemForm *joborder.ItemsForm) *models.JoborderItem
	ToItemModelSlice(plantID uint, itemForms []*joborder.ItemsForm) []*models.JoborderItem
}

type DefaultJoborderService struct {
	repo            repository.JoborderRepository
	customerService CustomerService
	productService  ProductService
}

func NewJoborderService(repo repository.JoborderRepository, customerService CustomerService, productService ProductService) JoborderService {
	return &DefaultJoborderService{repo: repo, customerService: customerService, productService: productService}
}

func (s *DefaultJoborderService) GetAllJoborders(plantID uint, filter joborder.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*joborder.Form, int64, error) {
	data, count, err := s.repo.GetAll(plantID, filter, pagination, sorting)
	if err != nil {
		return nil, count, err
	}
	return s.ToFormSlice(plantID, data), count, err
}

func (s *DefaultJoborderService) CreateJoborder(plantID uint, joborderForm *joborder.Form) error {
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

func (s *DefaultJoborderService) GetJoborderByID(plantID uint, id uint) (*joborder.Form, error) {
	data, err := s.repo.GetByID(plantID, id)
	if err != nil {
		return nil, err
	}
	return s.ToForm(plantID, data), nil
}

func (s *DefaultJoborderService) UpdateJoborder(plantID uint, id uint, joborderForm *joborder.Form) error {
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

func (s *DefaultJoborderService) DeleteJoborder(plantID uint, id uint) error {
	return s.repo.Delete(plantID, id)
}

func (s *DefaultJoborderService) DeleteJoborders(plantID uint, ids []uint) error {
	return s.repo.DeleteMulti(plantID, ids)
}

func (s *DefaultJoborderService) ExistsByItemId(joborderID uint, ID uint) bool {
	return s.repo.ExistsByItemId(joborderID, ID)
}

func (s *DefaultJoborderService) ExistsById(plantID uint, ID uint) bool {
	return s.repo.ExistsByID(plantID, ID)
}

func (s *DefaultJoborderService) ExistsByOrderNo(plantID uint, orderNo string, ID uint) bool {
	return s.repo.ExistsByOrderNo(plantID, orderNo, ID)
}

func (s *DefaultJoborderService) ToModel(plantID uint, joborderForm *joborder.Form) *models.Joborder {
	joborderModel := &models.Joborder{
		IssuedDate:    joborderForm.IssuedDate,
		OrderNo:       joborderForm.OrderNo,
		POCategory:    models.POCategory(joborderForm.POCategory),
		Status:        joborderForm.Status,
		ProcessStatus: joborderForm.ProcessStatus,
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

func (s *DefaultJoborderService) FormToModel(plantID uint, joborderForm *joborder.Form, joborderModel *models.Joborder) {
	joborderModel.IssuedDate = joborderForm.IssuedDate
	joborderModel.OrderNo = joborderForm.OrderNo
	joborderModel.POCategory = models.POCategory(joborderForm.POCategory)
	joborderModel.Status = joborderForm.Status
	joborderModel.ProcessStatus = joborderForm.ProcessStatus

	if joborderForm.Customer != nil {
		joborderModel.Customer = s.customerService.ToModel(plantID, joborderForm.Customer)
	} else {
		joborderModel.Customer = nil
		joborderModel.CustomerID = nil
	}
	if joborderForm.Items != nil {
		joborderModel.Items = s.ToItemModelSlice(plantID, joborderForm.Items)
	} else {
		joborderModel.Items = make([]*models.JoborderItem, 0)
	}
}

func (s *DefaultJoborderService) ToForm(plantID uint, joborderModel *models.Joborder) *joborder.Form {
	joborderForm := &joborder.Form{
		ID:            joborderModel.ID,
		IssuedDate:    joborderModel.IssuedDate,
		OrderNo:       joborderModel.OrderNo,
		POCategory:    string(joborderModel.POCategory),
		Status:        joborderModel.Status,
		ProcessStatus: joborderModel.ProcessStatus,
	}
	if joborderModel.Customer != nil {
		joborderForm.Customer = s.customerService.ToForm(plantID, joborderModel.Customer)
	}
	joborderForm.Items = s.ToItemFormSlice(plantID, joborderModel.Items)
	return joborderForm
}

func (s *DefaultJoborderService) ToFormSlice(plantID uint, joborderModels []*models.Joborder) []*joborder.Form {
	data := make([]*joborder.Form, 0)
	for _, joborderModel := range joborderModels {
		data = append(data, s.ToForm(plantID, joborderModel))
	}
	return data
}

func (s *DefaultJoborderService) ToModelSlice(plantID uint, joborderForms []*joborder.Form) []*models.Joborder {
	data := make([]*models.Joborder, 0)
	for _, joborderForm := range joborderForms {
		data = append(data, s.ToModel(plantID, joborderForm))
	}
	return data
}

func (s *DefaultJoborderService) ToItemForm(plantID uint, itemModel *models.JoborderItem) *joborder.ItemsForm {
	ingredientForm := &joborder.ItemsForm{
		ID:        itemModel.ID,
		ProductID: itemModel.ProductID,
		Product:   s.productService.ToForm(itemModel.Product),
		Quantity:  itemModel.Quantity,
	}
	return ingredientForm
}

func (s *DefaultJoborderService) ToItemFormSlice(plantID uint, itemModels []*models.JoborderItem) []*joborder.ItemsForm {
	data := make([]*joborder.ItemsForm, 0)
	for _, itemModel := range itemModels {
		data = append(data, s.ToItemForm(plantID, itemModel))
	}
	return data
}

func (s *DefaultJoborderService) ToItemModel(plantID uint, itemForm *joborder.ItemsForm) *models.JoborderItem {
	product := &models.Product{}
	product.ID = itemForm.ProductID

	joborderItemModel := &models.JoborderItem{
		ProductID: itemForm.ProductID,
		Product:   product,
		Quantity:  itemForm.Quantity,
	}
	joborderItemModel.ID = itemForm.ID
	return joborderItemModel
}

func (s *DefaultJoborderService) ToItemModelSlice(plantID uint, itemForms []*joborder.ItemsForm) []*models.JoborderItem {
	data := make([]*models.JoborderItem, 0)
	for _, itemForm := range itemForms {
		data = append(data, s.ToItemModel(plantID, itemForm))
	}
	return data
}
