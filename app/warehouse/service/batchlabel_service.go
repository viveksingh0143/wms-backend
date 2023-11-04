package service

import (
	baseModels "star-wms/app/base/models"
	baseService "star-wms/app/base/service"
	"star-wms/app/warehouse/dto/batchlabel"
	"star-wms/app/warehouse/models"
	"star-wms/app/warehouse/repository"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/common/responses"
)

type BatchlabelService interface {
	GetAllBatchlabels(plantID uint, filter batchlabel.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*batchlabel.Form, int64, error)
	CreateBatchlabel(plantID uint, batchlabelForm *batchlabel.Form) error
	GetBatchlabelByID(plantID uint, id uint) (*batchlabel.Form, error)
	UpdateBatchlabel(plantID uint, id uint, batchlabelForm *batchlabel.Form) error
	DeleteBatchlabel(plantID uint, id uint) error
	DeleteBatchlabels(plantID uint, ids []uint) error
	ExistsById(plantID uint, ID uint) bool
	ExistsByBatchNo(plantID uint, batchNo string, ID uint) bool
	ToModel(plantID uint, batchlabelForm *batchlabel.Form) *models.Batchlabel
	FormToModel(plantID uint, batchlabelForm *batchlabel.Form, batchlabelModel *models.Batchlabel)
	ToForm(plantID uint, batchlabelModel *models.Batchlabel) *batchlabel.Form
	ToFormSlice(plantID uint, batchlabelModels []*models.Batchlabel) []*batchlabel.Form
	ToModelSlice(plantID uint, batchlabelForms []*batchlabel.Form) []*models.Batchlabel
}

type DefaultBatchlabelService struct {
	repo            repository.BatchlabelRepository
	customerService baseService.CustomerService
	productService  baseService.ProductService
	machineService  baseService.MachineService
	joborderService baseService.JoborderService
}

func NewBatchlabelService(repo repository.BatchlabelRepository, customerService baseService.CustomerService, productService baseService.ProductService, machineService baseService.MachineService, joborderService baseService.JoborderService) BatchlabelService {
	return &DefaultBatchlabelService{repo: repo, customerService: customerService, productService: productService, machineService: machineService, joborderService: joborderService}
}

func (s *DefaultBatchlabelService) GetAllBatchlabels(plantID uint, filter batchlabel.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*batchlabel.Form, int64, error) {
	data, count, err := s.repo.GetAll(plantID, filter, pagination, sorting)
	if err != nil {
		return nil, count, err
	}
	return s.ToFormSlice(plantID, data), count, err
}

func (s *DefaultBatchlabelService) CreateBatchlabel(plantID uint, batchlabelForm *batchlabel.Form) error {
	if s.ExistsByBatchNo(plantID, batchlabelForm.BatchNo, 0) {
		return responses.NewInputError("batch_no", "already exists", batchlabelForm.BatchNo)
	}
	if !s.customerService.ExistsById(plantID, batchlabelForm.Customer.ID) {
		return responses.NewInputError("customer.id", "customer not exists", batchlabelForm.Customer.ID)
	}
	if !s.productService.ExistsById(batchlabelForm.Product.ID) {
		return responses.NewInputError("product.id", "product not exists", batchlabelForm.Product.ID)
	}
	if !s.machineService.ExistsById(plantID, batchlabelForm.Machine.ID) {
		return responses.NewInputError("machine.id", "machine not exists", batchlabelForm.Machine.ID)
	}

	if batchlabelForm.Joborder != nil && batchlabelForm.Joborder.ID > 0 {
		if !s.joborderService.ExistsById(plantID, batchlabelForm.Joborder.ID) {
			return responses.NewInputError("joborder.id", "joborder not exists", batchlabelForm.Joborder.ID)
		}
		if batchlabelForm.JoborderItem != nil && batchlabelForm.JoborderItem.ID > 0 && !s.joborderService.ExistsByItemId(batchlabelForm.Joborder.ID, batchlabelForm.JoborderItem.ID) {
			return responses.NewInputError("joborder_item.id", "joborder item not exists", batchlabelForm.JoborderItem.ID)
		}
	}

	batchlabelModel := s.ToModel(plantID, batchlabelForm)
	return s.repo.Create(plantID, batchlabelModel)
}

func (s *DefaultBatchlabelService) GetBatchlabelByID(plantID uint, id uint) (*batchlabel.Form, error) {
	data, err := s.repo.GetByID(plantID, id)
	if err != nil {
		return nil, err
	}
	return s.ToForm(plantID, data), nil
}

func (s *DefaultBatchlabelService) UpdateBatchlabel(plantID uint, id uint, batchlabelForm *batchlabel.Form) error {
	if s.ExistsByBatchNo(plantID, batchlabelForm.BatchNo, id) {
		return responses.NewInputError("batch_no", "already exists", batchlabelForm.BatchNo)
	}
	if !s.customerService.ExistsById(plantID, batchlabelForm.Customer.ID) {
		return responses.NewInputError("customer.id", "customer not exists", batchlabelForm.Customer.ID)
	}
	if !s.productService.ExistsById(batchlabelForm.Product.ID) {
		return responses.NewInputError("product.id", "product not exists", batchlabelForm.Product.ID)
	}
	if !s.machineService.ExistsById(plantID, batchlabelForm.Machine.ID) {
		return responses.NewInputError("machine.id", "machine not exists", batchlabelForm.Machine.ID)
	}
	if batchlabelForm.Joborder != nil && batchlabelForm.Joborder.ID > 0 {
		if !s.joborderService.ExistsById(plantID, batchlabelForm.Joborder.ID) {
			return responses.NewInputError("joborder.id", "joborder not exists", batchlabelForm.Joborder.ID)
		}
		if batchlabelForm.JoborderItem != nil && batchlabelForm.JoborderItem.ID > 0 && !s.joborderService.ExistsByItemId(batchlabelForm.Joborder.ID, batchlabelForm.JoborderItem.ID) {
			return responses.NewInputError("joborder_item.id", "joborder item not exists", batchlabelForm.JoborderItem.ID)
		}
	}
	batchlabelModel, err := s.repo.GetByID(plantID, id)
	if err != nil {
		return err
	}
	s.FormToModel(plantID, batchlabelForm, batchlabelModel)
	return s.repo.Update(plantID, batchlabelModel)
}

func (s *DefaultBatchlabelService) DeleteBatchlabel(plantID uint, id uint) error {
	return s.repo.Delete(plantID, id)
}

func (s *DefaultBatchlabelService) DeleteBatchlabels(plantID uint, ids []uint) error {
	return s.repo.DeleteMulti(plantID, ids)
}

func (s *DefaultBatchlabelService) ExistsById(plantID uint, ID uint) bool {
	return s.repo.ExistsByID(plantID, ID)
}

func (s *DefaultBatchlabelService) ExistsByBatchNo(plantID uint, batchNo string, ID uint) bool {
	return s.repo.ExistsByBatchNo(plantID, batchNo, ID)
}

func (s *DefaultBatchlabelService) ToModel(plantID uint, batchlabelForm *batchlabel.Form) *models.Batchlabel {
	batchlabelModel := &models.Batchlabel{
		BatchDate:       batchlabelForm.BatchDate,
		BatchNo:         batchlabelForm.BatchNo,
		SoNumber:        batchlabelForm.SoNumber,
		POCategory:      baseModels.POCategory(batchlabelForm.POCategory),
		UnitType:        baseModels.UnitType(batchlabelForm.UnitType),
		UnitWeight:      batchlabelForm.UnitWeight,
		UnitValue:       baseModels.UnitValue(batchlabelForm.UnitValue),
		TargetQuantity:  batchlabelForm.TargetQuantity,
		PackageQuantity: batchlabelForm.PackageQuantity,
		Status:          batchlabelForm.Status,
		ProcessStatus:   batchlabelForm.ProcessStatus,
	}
	batchlabelModel.ID = batchlabelForm.ID

	if batchlabelForm.Product != nil {
		batchlabelModel.Product = s.productService.ToModel(batchlabelForm.Product)
	}
	if batchlabelForm.Customer != nil {
		batchlabelModel.Customer = s.customerService.ToModel(plantID, batchlabelForm.Customer)
	}
	if batchlabelForm.Machine != nil {
		batchlabelModel.Machine = s.machineService.ToModel(plantID, batchlabelForm.Machine)
	}
	if batchlabelForm.Joborder != nil {
		batchlabelModel.Joborder = s.joborderService.ToModel(plantID, batchlabelForm.Joborder)
	}
	if batchlabelForm.JoborderItem != nil {
		batchlabelModel.JoborderItem = s.joborderService.ToItemModel(plantID, batchlabelForm.JoborderItem)
	}
	return batchlabelModel
}

func (s *DefaultBatchlabelService) FormToModel(plantID uint, batchlabelForm *batchlabel.Form, batchlabelModel *models.Batchlabel) {
	batchlabelModel.BatchDate = batchlabelForm.BatchDate
	batchlabelModel.BatchNo = batchlabelForm.BatchNo
	batchlabelModel.SoNumber = batchlabelForm.SoNumber
	batchlabelModel.POCategory = baseModels.POCategory(batchlabelForm.POCategory)
	batchlabelModel.UnitType = baseModels.UnitType(batchlabelForm.UnitType)
	batchlabelModel.UnitWeight = batchlabelForm.UnitWeight
	batchlabelModel.UnitValue = baseModels.UnitValue(batchlabelForm.UnitValue)
	batchlabelModel.TargetQuantity = batchlabelForm.TargetQuantity
	batchlabelModel.PackageQuantity = batchlabelForm.PackageQuantity
	batchlabelModel.Status = batchlabelForm.Status
	batchlabelModel.ProcessStatus = batchlabelForm.ProcessStatus

	if batchlabelForm.Product != nil {
		batchlabelModel.Product = s.productService.ToModel(batchlabelForm.Product)
	}
	if batchlabelForm.Customer != nil {
		batchlabelModel.Customer = s.customerService.ToModel(plantID, batchlabelForm.Customer)
	}
	if batchlabelForm.Machine != nil {
		batchlabelModel.Machine = s.machineService.ToModel(plantID, batchlabelForm.Machine)
	}
	if batchlabelForm.Joborder != nil {
		batchlabelModel.Joborder = s.joborderService.ToModel(plantID, batchlabelForm.Joborder)
	}
	if batchlabelForm.JoborderItem != nil {
		batchlabelModel.JoborderItem = s.joborderService.ToItemModel(plantID, batchlabelForm.JoborderItem)
	}
}

func (s *DefaultBatchlabelService) ToForm(plantID uint, batchlabelModel *models.Batchlabel) *batchlabel.Form {
	batchlabelForm := &batchlabel.Form{
		ID:              batchlabelModel.ID,
		BatchDate:       batchlabelModel.BatchDate,
		BatchNo:         batchlabelModel.BatchNo,
		SoNumber:        batchlabelModel.SoNumber,
		POCategory:      string(batchlabelModel.POCategory),
		UnitType:        string(batchlabelModel.UnitType),
		UnitWeight:      batchlabelModel.UnitWeight,
		UnitValue:       string(batchlabelModel.UnitValue),
		TargetQuantity:  batchlabelModel.TargetQuantity,
		PackageQuantity: batchlabelModel.PackageQuantity,
		Status:          batchlabelModel.Status,
		ProcessStatus:   batchlabelModel.ProcessStatus,
	}

	if batchlabelModel.Product != nil {
		batchlabelForm.Product = s.productService.ToForm(batchlabelModel.Product)
	}
	if batchlabelModel.Customer != nil {
		batchlabelForm.Customer = s.customerService.ToForm(plantID, batchlabelModel.Customer)
	}
	if batchlabelModel.Machine != nil {
		batchlabelForm.Machine = s.machineService.ToForm(plantID, batchlabelModel.Machine)
	}
	if batchlabelModel.Joborder != nil {
		batchlabelForm.Joborder = s.joborderService.ToForm(plantID, batchlabelModel.Joborder)
	}
	if batchlabelModel.JoborderItem != nil {
		batchlabelForm.JoborderItem = s.joborderService.ToItemForm(plantID, batchlabelModel.JoborderItem)
	}
	return batchlabelForm
}

func (s *DefaultBatchlabelService) ToFormSlice(plantID uint, batchlabelModels []*models.Batchlabel) []*batchlabel.Form {
	data := make([]*batchlabel.Form, 0)
	for _, batchlabelModel := range batchlabelModels {
		data = append(data, s.ToForm(plantID, batchlabelModel))
	}
	return data
}

func (s *DefaultBatchlabelService) ToModelSlice(plantID uint, batchlabelForms []*batchlabel.Form) []*models.Batchlabel {
	data := make([]*models.Batchlabel, 0)
	for _, batchlabelForm := range batchlabelForms {
		data = append(data, s.ToModel(plantID, batchlabelForm))
	}
	return data
}
