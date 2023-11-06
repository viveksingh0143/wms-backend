package service

type StockapprovalService interface {
	//CreateApproval(plantID uint, inventoryForm *inventory.Form) error
	//GetAllApprovals(plantID uint, filter inventory.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*inventory.Form, int64, error)
	//GetApprovalByID(plantID uint, id uint) (*inventory.Form, error)
	//UpdateApproval(plantID uint, id uint, inventoryForm *inventory.Form) error
	//DeleteApproval(plantID uint, id uint) error
	//DeleteApprovals(plantID uint, ids []uint) error
	//ExistsById(plantID uint, ID uint) bool
	//ToModel(plantID uint, inventoryForm *inventory.Form) *models.Approval
	//FormToModel(plantID uint, inventoryForm *inventory.Form, inventoryModel *models.Approval)
	//ToForm(plantID uint, inventoryModel *models.Approval) *inventory.Form
	//ToFormSlice(plantID uint, inventoryModels []*models.Approval) []*inventory.Form
	//ToModelSlice(plantID uint, inventoryForms []*inventory.Form) []*models.Approval
	//
	//CreateRawMaterialStockIn(plantID uint, inventoryForm *inventory.RawMaterialStockInForm) error
}

type DefaultApprovalService struct {
	//repo             repository.ApprovalRepository
	//productService   baseService.ProductService
	//storeService     baseService.StoreService
	//containerService baseService.ContainerService
}

//func NewApprovalService(repo repository.ApprovalRepository, productService baseService.ProductService, storeService baseService.StoreService, containerService baseService.ContainerService) StockapprovalService {
//	return &DefaultApprovalService{repo: repo, productService: productService, storeService: storeService, containerService: containerService}
//}
//
//func (s *DefaultApprovalService) GetAllApprovals(plantID uint, filter inventory.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*inventory.Form, int64, error) {
//	data, count, err := s.repo.GetAll(plantID, filter, pagination, sorting)
//	if err != nil {
//		return nil, count, err
//	}
//	return s.ToFormSlice(plantID, data), count, err
//}
//
//func (s *DefaultApprovalService) CreateApproval(plantID uint, inventoryForm *inventory.Form) error {
//	if !s.storeService.ExistsById(plantID, inventoryForm.Store.ID) {
//		return responses.NewInputError("store.id", "store not exists", inventoryForm.Store.ID)
//	}
//	if !s.productService.ExistsById(inventoryForm.Product.ID) {
//		return responses.NewInputError("product.id", "product not exists", inventoryForm.Product.ID)
//	}
//	inventoryModel := s.ToModel(plantID, inventoryForm)
//	return s.repo.Create(plantID, inventoryModel)
//}
//
//func (s *DefaultApprovalService) CreateRawMaterialStockIn(plantID uint, inventoryForm *inventory.RawMaterialStockInForm) error {
//	if !s.storeService.ExistsById(plantID, inventoryForm.Store.ID) {
//		return responses.NewInputError("store.id", "store not exists", inventoryForm.Store.ID)
//	}
//	if !s.productService.ExistsById(inventoryForm.Product.ID) {
//		return responses.NewInputError("product.id", "product not exists", inventoryForm.Product.ID)
//	}
//	if inventoryForm.Container.Code == "" {
//		return responses.NewInputError("container.code", "code is required", inventoryForm.Container.Code)
//	}
//	var containerForm *container.Form
//	if !s.containerService.ExistsByCode(plantID, inventoryForm.Container.Code, 0) {
//		containerForm = &container.Form{
//			PlantID:       plantID,
//			ContainerType: string(baseModels.Pallet),
//			Name:          inventoryForm.Container.Code,
//			Code:          inventoryForm.Container.Code,
//		}
//		err := s.containerService.CreateContainer(plantID, containerForm)
//		if err != nil {
//			return err
//		}
//	}
//	containerForm, err := s.containerService.GetContainerByCode(plantID, inventoryForm.Container.Code, false, false, false, false)
//	if err != nil {
//		return err
//	}
//	if containerForm.StockLevel != baseModels.Empty {
//		return responses.NewInputError("container.code", "is not empty", inventoryForm.Container.Code)
//	}
//	storeModel := s.storeService.ToModel(plantID, inventoryForm.Store)
//	containerModel := s.containerService.ToModel(plantID, containerForm)
//	contentModel := &baseModels.ContainerContent{
//		ProductID: inventoryForm.Product.ID,
//		Product:   s.productService.ToModel(inventoryForm.Product),
//		Quantity:  inventoryForm.Quantity,
//	}
//	return s.repo.CreateRawMaterialStockIn(plantID, storeModel, containerModel, contentModel)
//}
//
//func (s *DefaultApprovalService) GetApprovalByID(plantID uint, id uint) (*inventory.Form, error) {
//	data, err := s.repo.GetByID(plantID, id)
//	if err != nil {
//		return nil, err
//	}
//	return s.ToForm(plantID, data), nil
//}
//
//func (s *DefaultApprovalService) UpdateApproval(plantID uint, id uint, inventoryForm *inventory.Form) error {
//	if !s.storeService.ExistsById(plantID, inventoryForm.Store.ID) {
//		return responses.NewInputError("store.id", "store not exists", inventoryForm.Store.ID)
//	}
//	if !s.productService.ExistsById(inventoryForm.Product.ID) {
//		return responses.NewInputError("product.id", "product not exists", inventoryForm.Product.ID)
//	}
//	inventoryModel, err := s.repo.GetByID(plantID, id)
//	if err != nil {
//		return err
//	}
//	s.FormToModel(plantID, inventoryForm, inventoryModel)
//	return s.repo.Update(plantID, inventoryModel)
//}
//
//func (s *DefaultApprovalService) DeleteApproval(plantID uint, id uint) error {
//	return s.repo.Delete(plantID, id)
//}
//
//func (s *DefaultApprovalService) DeleteApprovals(plantID uint, ids []uint) error {
//	return s.repo.DeleteMulti(plantID, ids)
//}
//
//func (s *DefaultApprovalService) ExistsById(plantID uint, ID uint) bool {
//	return s.repo.ExistsByID(plantID, ID)
//}
//
//func (s *DefaultApprovalService) ToModel(plantID uint, inventoryForm *inventory.Form) *models.Approval {
//	inventoryModel := &models.Approval{
//		Quantity: inventoryForm.Quantity,
//	}
//	inventoryModel.ID = inventoryForm.ID
//
//	if inventoryForm.Product != nil {
//		inventoryModel.Product = s.productService.ToModel(inventoryForm.Product)
//	}
//	if inventoryForm.Store != nil {
//		inventoryModel.Store = s.storeService.ToModel(plantID, inventoryForm.Store)
//	}
//	return inventoryModel
//}
//
//func (s *DefaultApprovalService) FormToModel(plantID uint, inventoryForm *inventory.Form, inventoryModel *models.Approval) {
//	inventoryModel.Quantity = inventoryForm.Quantity
//
//	if inventoryForm.Product != nil {
//		inventoryModel.Product = s.productService.ToModel(inventoryForm.Product)
//	}
//	if inventoryForm.Store != nil {
//		inventoryModel.Store = s.storeService.ToModel(plantID, inventoryForm.Store)
//	}
//}
//
//func (s *DefaultApprovalService) ToForm(plantID uint, inventoryModel *models.Approval) *inventory.Form {
//	inventoryForm := &inventory.Form{
//		ID:       inventoryModel.ID,
//		Quantity: inventoryModel.Quantity,
//	}
//
//	if inventoryModel.Product != nil {
//		inventoryForm.Product = s.productService.ToForm(inventoryModel.Product)
//	}
//	if inventoryModel.Store != nil {
//		inventoryForm.Store = s.storeService.ToForm(plantID, inventoryModel.Store)
//	}
//	return inventoryForm
//}
//
//func (s *DefaultApprovalService) ToFormSlice(plantID uint, inventoryModels []*models.Approval) []*inventory.Form {
//	data := make([]*inventory.Form, 0)
//	for _, inventoryModel := range inventoryModels {
//		data = append(data, s.ToForm(plantID, inventoryModel))
//	}
//	return data
//}
//
//func (s *DefaultApprovalService) ToModelSlice(plantID uint, inventoryForms []*inventory.Form) []*models.Approval {
//	data := make([]*models.Approval, 0)
//	for _, inventoryForm := range inventoryForms {
//		data = append(data, s.ToModel(plantID, inventoryForm))
//	}
//	return data
//}
