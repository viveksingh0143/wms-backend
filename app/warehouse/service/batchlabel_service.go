package service

import (
	"fmt"
	baseModels "star-wms/app/base/models"
	baseService "star-wms/app/base/service"
	"star-wms/app/warehouse/dto/batchlabel"
	"star-wms/app/warehouse/models"
	"star-wms/app/warehouse/repository"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/common/responses"
	"star-wms/core/types"
	"time"
)

type BatchlabelService interface {
	GetAllBatchlabels(plantID uint, filter batchlabel.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*batchlabel.Form, int64, error)
	CreateBatchlabel(plantID uint, batchlabelForm *batchlabel.Form) error
	GetBatchlabelByID(plantID uint, id uint) (*batchlabel.Form, error)
	GetBatchlabelByBatchNo(plantID uint, batchNo string, needCustomer bool, needProduct bool, needMachine bool, needJoborder bool, needJoborderItem bool, needStickers bool) (*batchlabel.Form, error)
	UpdateBatchlabel(plantID uint, id uint, batchlabelForm *batchlabel.Form) error
	DeleteBatchlabel(plantID uint, id uint) error
	DeleteBatchlabels(plantID uint, ids []uint) error
	ExistsBatchlabelById(plantID uint, ID uint) bool
	ExistsBatchlabelByBatchNo(plantID uint, batchNo string, ID uint) bool
	ToBatchlabelModel(plantID uint, batchlabelForm *batchlabel.Form) *models.Batchlabel
	FormToBatchlabelModel(plantID uint, batchlabelForm *batchlabel.Form, batchlabelModel *models.Batchlabel)
	ToBatchlabelForm(plantID uint, batchlabelModel *models.Batchlabel) *batchlabel.Form
	ToBatchlabelFormSlice(plantID uint, batchlabelModels []*models.Batchlabel) []*batchlabel.Form
	ToBatchlabelModelSlice(plantID uint, batchlabelForms []*batchlabel.Form) []*models.Batchlabel

	GetStickersCountForBatchlabel(plantID uint, batchlabelID uint) (int64, error)
	GetAllStickers(plantID uint, batchlabelID uint, filter batchlabel.StickerFilter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*batchlabel.StickerForm, int64, error)
	CreateSticker(plantID uint, batchlabelID uint, multiStickerForm *batchlabel.MultiStickerForm) error
	GetStickerByID(plantID uint, batchlabelID uint, id uint) (*batchlabel.StickerForm, error)
	GetStickerByBarcodePlantwise(plantID uint, barcode string) (*batchlabel.StickerForm, error)
	ExistsStickerById(plantID uint, batchlabelID uint, ID uint) bool
	ExistsStickerByBarcode(plantID uint, batchlabelID uint, barcode string, ID uint) bool
	ExistsStickerByBarcodePlantwise(plantID uint, barcode string) bool
	ToStickerForm(plantID uint, batchlabelID uint, stickerModel *models.Sticker) *batchlabel.StickerForm
	ToStickerFormSlice(plantID uint, batchlabelID uint, stickerModels []*models.Sticker) []*batchlabel.StickerForm
	ToStickerModel(plantID uint, stickerForm *batchlabel.StickerForm) *models.Sticker
	ToStickerModelFormSlice(plantID uint, stickersForm []*batchlabel.StickerForm) []*models.Sticker
}

type DefaultBatchlabelService struct {
	batchlabelRepository repository.BatchlabelRepository
	stickerRepository    repository.StickerRepository
	customerService      baseService.CustomerService
	productService       baseService.ProductService
	machineService       baseService.MachineService
	joborderService      baseService.JoborderService
}

func NewBatchlabelService(batchlabelRepository repository.BatchlabelRepository, stickerRepository repository.StickerRepository, customerService baseService.CustomerService, productService baseService.ProductService, machineService baseService.MachineService, joborderService baseService.JoborderService) BatchlabelService {
	return &DefaultBatchlabelService{batchlabelRepository: batchlabelRepository, stickerRepository: stickerRepository, customerService: customerService, productService: productService, machineService: machineService, joborderService: joborderService}
}

func (s *DefaultBatchlabelService) GetAllBatchlabels(plantID uint, filter batchlabel.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*batchlabel.Form, int64, error) {
	data, count, err := s.batchlabelRepository.GetAll(plantID, filter, pagination, sorting)
	if err != nil {
		return nil, count, err
	}
	return s.ToBatchlabelFormSlice(plantID, data), count, err
}

func (s *DefaultBatchlabelService) CreateBatchlabel(plantID uint, batchlabelForm *batchlabel.Form) error {
	if s.ExistsBatchlabelByBatchNo(plantID, batchlabelForm.BatchNo, 0) {
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

	batchlabelModel := s.ToBatchlabelModel(plantID, batchlabelForm)
	return s.batchlabelRepository.Create(plantID, batchlabelModel)
}

func (s *DefaultBatchlabelService) GetBatchlabelByID(plantID uint, id uint) (*batchlabel.Form, error) {
	data, err := s.batchlabelRepository.GetByID(plantID, id)
	if err != nil {
		return nil, err
	}

	batchlabelForm := s.ToBatchlabelForm(plantID, data)
	stickerCount, err := s.GetStickersCountForBatchlabel(plantID, id)
	batchlabelForm.TotalPrinted = stickerCount
	batchlabelForm.LabelsToPrint = batchlabelForm.GetStickerCountToPrint()
	if err != nil {
		return nil, err
	}
	return batchlabelForm, nil
}

func (s *DefaultBatchlabelService) GetBatchlabelByBatchNo(plantID uint, batchNo string, needCustomer bool, needProduct bool, needMachine bool, needJoborder bool, needJoborderItem bool, needStickers bool) (*batchlabel.Form, error) {
	data, err := s.batchlabelRepository.GetByBatchNo(plantID, batchNo, needCustomer, needProduct, needMachine, needJoborder, needJoborderItem, needStickers)
	if err != nil {
		return nil, err
	}
	return s.ToBatchlabelForm(plantID, data), nil
}

func (s *DefaultBatchlabelService) UpdateBatchlabel(plantID uint, id uint, batchlabelForm *batchlabel.Form) error {
	if s.ExistsBatchlabelByBatchNo(plantID, batchlabelForm.BatchNo, id) {
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
	batchlabelModel, err := s.batchlabelRepository.GetByID(plantID, id)
	if err != nil {
		return err
	}
	s.FormToBatchlabelModel(plantID, batchlabelForm, batchlabelModel)
	return s.batchlabelRepository.Update(plantID, batchlabelModel)
}

func (s *DefaultBatchlabelService) DeleteBatchlabel(plantID uint, id uint) error {
	return s.batchlabelRepository.Delete(plantID, id)
}

func (s *DefaultBatchlabelService) DeleteBatchlabels(plantID uint, ids []uint) error {
	return s.batchlabelRepository.DeleteMulti(plantID, ids)
}

func (s *DefaultBatchlabelService) ExistsBatchlabelById(plantID uint, ID uint) bool {
	return s.batchlabelRepository.ExistsByID(plantID, ID)
}

func (s *DefaultBatchlabelService) ExistsBatchlabelByBatchNo(plantID uint, batchNo string, ID uint) bool {
	return s.batchlabelRepository.ExistsByBatchNo(plantID, batchNo, ID)
}

func (s *DefaultBatchlabelService) ToBatchlabelModel(plantID uint, batchlabelForm *batchlabel.Form) *models.Batchlabel {
	if batchlabelForm == nil {
		return nil
	}

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

func (s *DefaultBatchlabelService) FormToBatchlabelModel(plantID uint, batchlabelForm *batchlabel.Form, batchlabelModel *models.Batchlabel) {
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

func (s *DefaultBatchlabelService) ToBatchlabelForm(plantID uint, batchlabelModel *models.Batchlabel) *batchlabel.Form {
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
	if batchlabelModel.Stickers != nil {
		batchlabelForm.Stickers = s.ToStickerFormSlice(plantID, batchlabelModel.ID, batchlabelModel.Stickers)
	}
	return batchlabelForm
}

func (s *DefaultBatchlabelService) ToBatchlabelFormSlice(plantID uint, batchlabelModels []*models.Batchlabel) []*batchlabel.Form {
	data := make([]*batchlabel.Form, 0)
	for _, batchlabelModel := range batchlabelModels {
		data = append(data, s.ToBatchlabelForm(plantID, batchlabelModel))
	}
	return data
}

func (s *DefaultBatchlabelService) ToBatchlabelModelSlice(plantID uint, batchlabelForms []*batchlabel.Form) []*models.Batchlabel {
	data := make([]*models.Batchlabel, 0)
	for _, batchlabelForm := range batchlabelForms {
		data = append(data, s.ToBatchlabelModel(plantID, batchlabelForm))
	}
	return data
}

func (s *DefaultBatchlabelService) GetStickersCountForBatchlabel(plantID uint, batchlabelID uint) (int64, error) {
	count, err := s.stickerRepository.GetCountForBatchlabel(plantID, batchlabelID)
	return count, err
}

func (s *DefaultBatchlabelService) GetAllStickers(plantID uint, batchlabelID uint, filter batchlabel.StickerFilter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*batchlabel.StickerForm, int64, error) {
	data, count, err := s.stickerRepository.GetAll(plantID, batchlabelID, filter, pagination, sorting)
	if err != nil {
		return nil, count, err
	}
	return s.ToStickerFormSlice(plantID, batchlabelID, data), count, err
}

func (s *DefaultBatchlabelService) CreateSticker(plantID uint, batchlabelID uint, multiStickerForm *batchlabel.MultiStickerForm) error {
	batchlabelForm, err := s.GetBatchlabelByID(plantID, batchlabelID)
	if err != nil {
		return responses.NewInputError("batchlabel", "not exists", batchlabelID)
	}

	totalCount, err := s.stickerRepository.GetCount(plantID, batchlabelID)
	if err != nil {
		return responses.NewInputError("sticker_count", "failed to fetch stickers count", batchlabelID)
	}

	stickersAvailable := batchlabelForm.GetStickerCountToPrint() - totalCount
	if stickersAvailable < multiStickerForm.StickersToGenerate {
		errMsg := fmt.Sprintf("only %d stickers are left, please reduce the sticker count", stickersAvailable)
		return responses.NewInputError("stickers_to_generate", errMsg, multiStickerForm.StickersToGenerate)
	}

	shiftTotalCount, err := s.stickerRepository.GetCountForShift(plantID, batchlabelID, multiStickerForm.Shift, time.Now())
	if err != nil {
		return responses.NewInputError("shift", "failed to fetch stickers count", batchlabelID)
	}

	var stickers = make([]*models.Sticker, 0, multiStickerForm.StickersToGenerate)
	for i := int64(0); i < multiStickerForm.StickersToGenerate; i++ {
		var stickerItems []*models.StickerItem
		if multiStickerForm.Items != nil && len(multiStickerForm.Items) > 0 {
			stickerItems = make([]*models.StickerItem, 0, len(multiStickerForm.Items))
			for _, item := range multiStickerForm.Items {
				itemBatchlabel, _ := s.GetBatchlabelByBatchNo(plantID, item.BatchNo, false, false, false, false, false, false)
				stickerItem := &models.StickerItem{
					Product:    s.productService.ToModel(item.Product),
					Quantity:   item.Quantity,
					BatchNo:    item.BatchNo,
					Batchlabel: s.ToBatchlabelModel(plantID, itemBatchlabel),
				}
				stickerItems = append(stickerItems, stickerItem)
			}
		}

		nextSerialNumber := fmt.Sprintf("%04d", totalCount+i+1)
		packetNo := shiftTotalCount + i + 1
		packetNoString := fmt.Sprintf("%04d", packetNo)
		stickerModel := &models.Sticker{
			Barcode:        "",
			PacketNo:       packetNoString,
			PrintCount:     0,
			Shift:          multiStickerForm.Shift,
			ProductLine:    batchlabelForm.Product.Name,
			BatchNo:        batchlabelForm.BatchNo,
			UnitWeightLine: fmt.Sprintf("%f", batchlabelForm.UnitWeight),
			QuantityLine:   fmt.Sprintf("%f", batchlabelForm.PackageQuantity),
			MachineNo:      batchlabelForm.Machine.Code,
			IsUsed:         false,
			Quantity:       batchlabelForm.PackageQuantity,
			Supervisor:     multiStickerForm.Supervisor,
			ProductID:      batchlabelForm.Product.ID,
			BatchlabelID:   batchlabelID,
			PlantID:        plantID,
		}
		plantNo := fmt.Sprintf("%02d", plantID)
		stickerModel.Barcode = fmt.Sprintf("%s%s%s%s", plantNo, stickerModel.BatchNo, nextSerialNumber, stickerModel.MachineNo)
		if stickerItems != nil {
			stickerModel.StickerItems = stickerItems
		}
		stickers = append(stickers, stickerModel)
	}
	err = s.stickerRepository.CreateAll(plantID, batchlabelID, stickers)
	if err != nil {
		return err
	}
	totalCount, err = s.stickerRepository.GetCount(plantID, batchlabelID)
	if err != nil {
		return responses.NewInputError("sticker_count", "failed to fetch stickers count", batchlabelID)
	}
	stickersAvailable = batchlabelForm.GetStickerCountToPrint() - totalCount
	if stickersAvailable > 0 {
		return s.batchlabelRepository.MarkProcessStatusAs(plantID, batchlabelID, types.ProcessStarted)
	} else {
		return s.batchlabelRepository.MarkProcessStatusAs(plantID, batchlabelID, types.ProcessClosed)
	}
}

func (s *DefaultBatchlabelService) GetStickerByID(plantID uint, batchlabelID uint, id uint) (*batchlabel.StickerForm, error) {
	data, err := s.stickerRepository.GetByID(plantID, batchlabelID, id)
	if err != nil {
		return nil, err
	}
	return s.ToStickerForm(plantID, batchlabelID, data), nil
}

func (s *DefaultBatchlabelService) GetStickerByBarcodePlantwise(plantID uint, barcode string) (*batchlabel.StickerForm, error) {
	data, err := s.stickerRepository.GetByBarcodePlantwise(plantID, barcode)
	if err != nil {
		return nil, err
	}
	return s.ToStickerForm(plantID, data.BatchlabelID, data), nil
}

func (s *DefaultBatchlabelService) ExistsStickerById(plantID uint, batchlabelID uint, ID uint) bool {
	return s.stickerRepository.ExistsByID(plantID, batchlabelID, ID)
}

func (s *DefaultBatchlabelService) ExistsStickerByBarcode(plantID uint, batchlabelID uint, barcode string, ID uint) bool {
	return s.stickerRepository.ExistsByBarcode(plantID, batchlabelID, barcode, ID)
}

func (s *DefaultBatchlabelService) ExistsStickerByBarcodePlantwise(plantID uint, barcode string) bool {
	return s.stickerRepository.ExistsByBarcodePlantwise(plantID, barcode)
}

func (s *DefaultBatchlabelService) ToStickerForm(plantID uint, batchlabelID uint, stickerModel *models.Sticker) *batchlabel.StickerForm {
	stickerForm := &batchlabel.StickerForm{
		ID:           stickerModel.ID,
		Barcode:      stickerModel.Barcode,
		PacketNo:     stickerModel.PacketNo,
		PrintCount:   stickerModel.PrintCount,
		Shift:        stickerModel.Shift,
		ProductLine:  stickerModel.ProductLine,
		BatchNo:      stickerModel.BatchNo,
		MachineNo:    stickerModel.MachineNo,
		IsUsed:       stickerModel.IsUsed,
		UnitWeight:   stickerModel.UnitWeightLine,
		QuantityLine: stickerModel.QuantityLine,
		Quantity:     stickerModel.Quantity,
		Supervisor:   stickerModel.Supervisor,
	}
	stickerForm.ProductID = stickerModel.ProductID
	if stickerModel.Product != nil {
		stickerForm.Product = s.productService.ToForm(stickerModel.Product)
	}
	if stickerModel.Batchlabel != nil {
		stickerForm.Batchlabel = s.ToBatchlabelForm(plantID, stickerModel.Batchlabel)
	}
	if stickerModel.StickerItems != nil {
		stickerForm.Items = s.ToStickerItemFormSlice(plantID, batchlabelID, stickerModel.StickerItems)
	}
	return stickerForm
}

func (s *DefaultBatchlabelService) ToStickerFormSlice(plantID uint, batchlabelID uint, stickerModels []*models.Sticker) []*batchlabel.StickerForm {
	data := make([]*batchlabel.StickerForm, 0)
	for _, stickerModel := range stickerModels {
		data = append(data, s.ToStickerForm(plantID, batchlabelID, stickerModel))
	}
	return data
}

func (s *DefaultBatchlabelService) ToStickerItemForm(plantID uint, batchlabelID uint, stickerItemModel *models.StickerItem) *batchlabel.StickerItem {
	stickerItemForm := &batchlabel.StickerItem{
		ID:       stickerItemModel.ID,
		Quantity: stickerItemModel.Quantity,
		BatchNo:  stickerItemModel.BatchNo,
	}
	if stickerItemModel.Product != nil {
		stickerItemForm.Product = s.productService.ToForm(stickerItemModel.Product)
	}

	if stickerItemModel.Batchlabel != nil {
		stickerItemForm.Batchlabel = s.ToBatchlabelForm(plantID, stickerItemModel.Batchlabel)
	}
	return stickerItemForm
}

func (s *DefaultBatchlabelService) ToStickerItemFormSlice(plantID uint, batchlabelID uint, stickerItemModels []*models.StickerItem) []*batchlabel.StickerItem {
	data := make([]*batchlabel.StickerItem, 0)
	for _, stickerItemModel := range stickerItemModels {
		data = append(data, s.ToStickerItemForm(plantID, batchlabelID, stickerItemModel))
	}
	return data
}

func (s *DefaultBatchlabelService) ToStickerModel(plantID uint, stickerForm *batchlabel.StickerForm) *models.Sticker {
	stickerModel := &models.Sticker{
		Barcode:      stickerForm.Barcode,
		PacketNo:     stickerForm.PacketNo,
		PrintCount:   stickerForm.PrintCount,
		Shift:        stickerForm.Shift,
		ProductLine:  stickerForm.ProductLine,
		BatchNo:      stickerForm.BatchNo,
		MachineNo:    stickerForm.MachineNo,
		IsUsed:       stickerForm.IsUsed,
		Supervisor:   stickerForm.Supervisor,
		Quantity:     stickerForm.Quantity,
		QuantityLine: stickerForm.QuantityLine,
	}
	stickerModel.ID = stickerForm.ID

	stickerModel.ProductID = stickerForm.ProductID
	if stickerForm.Product != nil {
		stickerModel.Product = s.productService.ToModel(stickerForm.Product)
	}
	if stickerForm.Batchlabel != nil {
		stickerModel.Batchlabel = s.ToBatchlabelModel(plantID, stickerForm.Batchlabel)
	}
	return stickerModel
}

func (s *DefaultBatchlabelService) ToStickerModelFormSlice(plantID uint, stickersForm []*batchlabel.StickerForm) []*models.Sticker {
	data := make([]*models.Sticker, 0)
	for _, stickerForm := range stickersForm {
		data = append(data, s.ToStickerModel(plantID, stickerForm))
	}
	return data
}
