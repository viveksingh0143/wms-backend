package service

import (
	baseService "star-wms/app/base/service"
	"star-wms/app/warehouse/dto/rmbatch"
	"star-wms/app/warehouse/models"
	"star-wms/app/warehouse/repository"
	commonModels "star-wms/core/common/requests"
)

type RMBatchService interface {
	GetAllRMBatchs(plantID uint, filter rmbatch.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*rmbatch.Form, int64, error)
	GetRMBatchByID(plantID uint, id uint) (*rmbatch.Form, error)
	ToForm(plantID uint, rmbatchModel *models.RMBatch) *rmbatch.Form
	ToFormSlice(plantID uint, rmbatchModels []*models.RMBatch) []*rmbatch.Form
}

type DefaultRMBatchService struct {
	repo             repository.RMBatchRepository
	productService   baseService.ProductService
	storeService     baseService.StoreService
	containerService baseService.ContainerService
}

func NewRMBatchService(repo repository.RMBatchRepository, productService baseService.ProductService, storeService baseService.StoreService, containerService baseService.ContainerService) RMBatchService {
	return &DefaultRMBatchService{repo: repo, productService: productService, storeService: storeService, containerService: containerService}
}

func (s *DefaultRMBatchService) GetAllRMBatchs(plantID uint, filter rmbatch.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*rmbatch.Form, int64, error) {
	data, count, err := s.repo.GetAll(plantID, filter, pagination, sorting)
	if err != nil {
		return nil, count, err
	}
	return s.ToFormSlice(plantID, data), count, err
}

func (s *DefaultRMBatchService) GetRMBatchByID(plantID uint, id uint) (*rmbatch.Form, error) {
	data, err := s.repo.GetByID(plantID, id)
	if err != nil {
		return nil, err
	}
	return s.ToForm(plantID, data), nil
}

func (s *DefaultRMBatchService) ToForm(plantID uint, rmbatchModel *models.RMBatch) *rmbatch.Form {
	rmbatchForm := &rmbatch.Form{
		ID:          rmbatchModel.ID,
		BatchNumber: rmbatchModel.BatchNumber,
		Quantity:    rmbatchModel.Quantity,
		Unit:        rmbatchModel.Unit,
		Status:      rmbatchModel.Status,
	}
	if rmbatchModel.Product != nil {
		rmbatchForm.Product = s.productService.ToForm(rmbatchModel.Product)
	}
	if rmbatchModel.Container != nil {
		rmbatchForm.Container = s.containerService.ToForm(plantID, rmbatchModel.Container)
	}
	if rmbatchModel.Store != nil {
		rmbatchForm.Store = s.storeService.ToForm(plantID, rmbatchModel.Store)
	}
	if rmbatchModel.Transactions != nil {
		rmbatchForm.Transactions = s.ToTransactionFormSlice(plantID, rmbatchModel.Transactions)
	}

	return rmbatchForm
}

func (s *DefaultRMBatchService) ToFormSlice(plantID uint, rmbatchModels []*models.RMBatch) []*rmbatch.Form {
	data := make([]*rmbatch.Form, 0)
	for _, rmbatchModel := range rmbatchModels {
		data = append(data, s.ToForm(plantID, rmbatchModel))
	}
	return data
}

func (s *DefaultRMBatchService) ToTransactionForm(plantID uint, transaction *models.RMBatchTransaction) *rmbatch.Transaction {
	rmbatchForm := &rmbatch.Transaction{
		ID:              transaction.ID,
		TransactionType: transaction.TransactionType,
		Quantity:        transaction.Quantity,
		Notes:           transaction.Notes,
		RMBatchID:       transaction.RMBatchID,
		ProductID:       transaction.ProductID,
		CreatedAt:       transaction.CreatedAt,
	}
	return rmbatchForm
}

func (s *DefaultRMBatchService) ToTransactionFormSlice(plantID uint, transactions []*models.RMBatchTransaction) []*rmbatch.Transaction {
	data := make([]*rmbatch.Transaction, 0)
	for _, transaction := range transactions {
		data = append(data, s.ToTransactionForm(plantID, transaction))
	}
	return data
}
