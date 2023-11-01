package service

import (
	"star-wms/app/base/dto/customer"
	"star-wms/app/base/models"
	"star-wms/app/base/repository"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/common/responses"
)

type CustomerService interface {
	GetAllCustomers(plantID uint, filter customer.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*customer.Form, int64, error)
	CreateCustomer(plantID uint, customerForm *customer.Form) error
	GetCustomerByID(plantID uint, id uint) (*customer.Form, error)
	UpdateCustomer(plantID uint, id uint, customerForm *customer.Form) error
	DeleteCustomer(plantID uint, id uint) error
	DeleteCustomers(plantID uint, ids []uint) error
	ExistsById(plantID uint, ID uint) bool
	ExistsByName(plantID uint, name string, ID uint) bool
	ExistsByCode(plantID uint, code string, ID uint) bool
	ToModel(plantID uint, customerForm *customer.Form) *models.Customer
	FormToModel(plantID uint, customerForm *customer.Form, customerModel *models.Customer)
	ToForm(plantID uint, customerModel *models.Customer) *customer.Form
	ToFormSlice(plantID uint, customerModels []*models.Customer) []*customer.Form
	ToModelSlice(plantID uint, customerForms []*customer.Form) []*models.Customer
}

type DefaultCustomerService struct {
	repo repository.CustomerRepository
}

func NewCustomerService(repo repository.CustomerRepository) CustomerService {
	return &DefaultCustomerService{repo: repo}
}

func (s *DefaultCustomerService) GetAllCustomers(plantID uint, filter customer.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*customer.Form, int64, error) {
	data, count, err := s.repo.GetAll(plantID, filter, pagination, sorting)
	if err != nil {
		return nil, count, err
	}
	return s.ToFormSlice(plantID, data), count, err
}

func (s *DefaultCustomerService) CreateCustomer(plantID uint, customerForm *customer.Form) error {
	if s.ExistsByName(plantID, customerForm.Name, 0) {
		return responses.NewInputError("name", "already exists", customerForm.Name)
	}
	if s.ExistsByCode(plantID, customerForm.Code, 0) {
		return responses.NewInputError("code", "already exists", customerForm.Code)
	}
	customerModel := s.ToModel(plantID, customerForm)
	return s.repo.Create(plantID, customerModel)
}

func (s *DefaultCustomerService) GetCustomerByID(plantID uint, id uint) (*customer.Form, error) {
	data, err := s.repo.GetByID(plantID, id)
	if err != nil {
		return nil, err
	}
	return s.ToForm(plantID, data), nil
}

func (s *DefaultCustomerService) UpdateCustomer(plantID uint, id uint, customerForm *customer.Form) error {
	if s.ExistsByName(plantID, customerForm.Name, id) {
		return responses.NewInputError("name", "already exists", customerForm.Name)
	}
	if s.ExistsByCode(plantID, customerForm.Code, id) {
		return responses.NewInputError("code", "already exists", customerForm.Code)
	}
	customerModel, err := s.repo.GetByID(plantID, id)
	if err != nil {
		return err
	}
	s.FormToModel(plantID, customerForm, customerModel)
	return s.repo.Update(plantID, customerModel)
}

func (s *DefaultCustomerService) DeleteCustomer(plantID uint, id uint) error {
	return s.repo.Delete(plantID, id)
}

func (s *DefaultCustomerService) DeleteCustomers(plantID uint, ids []uint) error {
	return s.repo.DeleteMulti(plantID, ids)
}

func (s *DefaultCustomerService) ExistsById(plantID uint, ID uint) bool {
	return s.repo.ExistsByID(plantID, ID)
}

func (s *DefaultCustomerService) ExistsByName(plantID uint, name string, ID uint) bool {
	return s.repo.ExistsByName(plantID, name, ID)
}

func (s *DefaultCustomerService) ExistsByCode(plantID uint, code string, ID uint) bool {
	return s.repo.ExistsByCode(plantID, code, ID)
}

func (s *DefaultCustomerService) ToModel(plantID uint, customerForm *customer.Form) *models.Customer {
	customerModel := &models.Customer{
		Name:             customerForm.Name,
		Code:             customerForm.Code,
		ContactPerson:    customerForm.ContactPerson,
		BillingAddress1:  customerForm.BillingAddress1,
		BillingAddress2:  customerForm.BillingAddress2,
		BillingState:     customerForm.BillingState,
		BillingCountry:   customerForm.BillingCountry,
		BillingPincode:   customerForm.BillingPincode,
		ShippingAddress1: customerForm.ShippingAddress1,
		ShippingAddress2: customerForm.ShippingAddress2,
		ShippingState:    customerForm.ShippingState,
		ShippingCountry:  customerForm.ShippingCountry,
		ShippingPincode:  customerForm.ShippingPincode,
		Status:           customerForm.Status,
	}
	customerModel.ID = customerForm.ID
	customerModel.PlantID = plantID
	return customerModel
}

func (s *DefaultCustomerService) FormToModel(plantID uint, customerForm *customer.Form, customerModel *models.Customer) {
	customerModel.Name = customerForm.Name
	customerModel.Code = customerForm.Code
	customerModel.ContactPerson = customerForm.ContactPerson
	customerModel.BillingAddress1 = customerForm.BillingAddress1
	customerModel.BillingAddress2 = customerForm.BillingAddress2
	customerModel.BillingState = customerForm.BillingState
	customerModel.BillingCountry = customerForm.BillingCountry
	customerModel.BillingPincode = customerForm.BillingPincode
	customerModel.ShippingAddress1 = customerForm.ShippingAddress1
	customerModel.ShippingAddress2 = customerForm.ShippingAddress2
	customerModel.ShippingState = customerForm.ShippingState
	customerModel.ShippingCountry = customerForm.ShippingCountry
	customerModel.ShippingPincode = customerForm.ShippingPincode
	customerModel.Status = customerForm.Status
}

func (s *DefaultCustomerService) ToForm(plantID uint, customerModel *models.Customer) *customer.Form {
	customerForm := &customer.Form{
		ID:               customerModel.ID,
		Name:             customerModel.Name,
		Code:             customerModel.Code,
		ContactPerson:    customerModel.ContactPerson,
		BillingAddress1:  customerModel.BillingAddress1,
		BillingAddress2:  customerModel.BillingAddress2,
		BillingState:     customerModel.BillingState,
		BillingCountry:   customerModel.BillingCountry,
		BillingPincode:   customerModel.BillingPincode,
		ShippingAddress1: customerModel.ShippingAddress1,
		ShippingAddress2: customerModel.ShippingAddress2,
		ShippingState:    customerModel.ShippingState,
		ShippingCountry:  customerModel.ShippingCountry,
		ShippingPincode:  customerModel.ShippingPincode,
		Status:           customerModel.Status,
	}
	customerForm.PlantID = customerModel.PlantID
	return customerForm
}

func (s *DefaultCustomerService) ToFormSlice(plantID uint, customerModels []*models.Customer) []*customer.Form {
	data := make([]*customer.Form, 0)
	for _, customerModel := range customerModels {
		data = append(data, s.ToForm(plantID, customerModel))
	}
	return data
}

func (s *DefaultCustomerService) ToModelSlice(plantID uint, customerForms []*customer.Form) []*models.Customer {
	data := make([]*models.Customer, 0)
	for _, customerForm := range customerForms {
		data = append(data, s.ToModel(plantID, customerForm))
	}
	return data
}
