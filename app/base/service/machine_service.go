package service

import (
	"star-wms/app/base/dto/machine"
	"star-wms/app/base/models"
	"star-wms/app/base/repository"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/common/responses"
)

type MachineService interface {
	GetAllMachines(plantID uint, filter machine.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*machine.Form, int64, error)
	CreateMachine(plantID uint, machineForm *machine.Form) error
	GetMachineByID(plantID uint, id uint) (*machine.Form, error)
	GetMachineByCode(plantID uint, code string) (*machine.Form, error)
	UpdateMachine(plantID uint, id uint, machineForm *machine.Form) error
	DeleteMachine(plantID uint, id uint) error
	DeleteMachines(plantID uint, ids []uint) error
	ExistsById(plantID uint, ID uint) bool
	ExistsByName(plantID uint, name string, ID uint) bool
	ExistsByCode(plantID uint, code string, ID uint) bool
	ToModel(plantID uint, machineForm *machine.Form) *models.Machine
	FormToModel(plantID uint, machineForm *machine.Form, machineModel *models.Machine)
	ToForm(plantID uint, machineModel *models.Machine) *machine.Form
	ToFormSlice(plantID uint, machineModels []*models.Machine) []*machine.Form
	ToModelSlice(plantID uint, machineForms []*machine.Form) []*models.Machine
}

type DefaultMachineService struct {
	repo repository.MachineRepository
}

func NewMachineService(repo repository.MachineRepository) MachineService {
	return &DefaultMachineService{repo: repo}
}

func (s *DefaultMachineService) GetAllMachines(plantID uint, filter machine.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*machine.Form, int64, error) {
	data, count, err := s.repo.GetAll(plantID, filter, pagination, sorting)
	if err != nil {
		return nil, count, err
	}
	return s.ToFormSlice(plantID, data), count, err
}

func (s *DefaultMachineService) CreateMachine(plantID uint, machineForm *machine.Form) error {
	if s.ExistsByCode(plantID, machineForm.Code, 0) {
		return responses.NewInputError("code", "already exists", machineForm.Code)
	}
	machineModel := s.ToModel(plantID, machineForm)
	return s.repo.Create(plantID, machineModel)
}

func (s *DefaultMachineService) GetMachineByID(plantID uint, id uint) (*machine.Form, error) {
	data, err := s.repo.GetByID(plantID, id)
	if err != nil {
		return nil, err
	}
	return s.ToForm(plantID, data), nil
}

func (s *DefaultMachineService) GetMachineByCode(plantID uint, code string) (*machine.Form, error) {
	data, err := s.repo.GetByCode(plantID, code)
	if err != nil {
		return nil, err
	}
	return s.ToForm(plantID, data), nil
}

func (s *DefaultMachineService) UpdateMachine(plantID uint, id uint, machineForm *machine.Form) error {
	//if s.ExistsByName(plantID, machineForm.Name, id) {
	//	return responses.NewInputError("name", "already exists", machineForm.Name)
	//}
	if s.ExistsByCode(plantID, machineForm.Code, id) {
		return responses.NewInputError("code", "already exists", machineForm.Code)
	}
	machineModel, err := s.repo.GetByID(plantID, id)
	if err != nil {
		return err
	}
	s.FormToModel(plantID, machineForm, machineModel)
	return s.repo.Update(plantID, machineModel)
}

func (s *DefaultMachineService) DeleteMachine(plantID uint, id uint) error {
	return s.repo.Delete(plantID, id)
}

func (s *DefaultMachineService) DeleteMachines(plantID uint, ids []uint) error {
	return s.repo.DeleteMulti(plantID, ids)
}

func (s *DefaultMachineService) ExistsById(plantID uint, ID uint) bool {
	return s.repo.ExistsByID(plantID, ID)
}

func (s *DefaultMachineService) ExistsByName(plantID uint, name string, ID uint) bool {
	return s.repo.ExistsByName(plantID, name, ID)
}

func (s *DefaultMachineService) ExistsByCode(plantID uint, code string, ID uint) bool {
	return s.repo.ExistsByCode(plantID, code, ID)
}

func (s *DefaultMachineService) ToModel(plantID uint, machineForm *machine.Form) *models.Machine {
	machineModel := &models.Machine{
		Name:   machineForm.Name,
		Code:   machineForm.Code,
		Status: machineForm.Status,
	}
	machineModel.ID = machineForm.ID
	machineModel.PlantID = plantID
	return machineModel
}

func (s *DefaultMachineService) FormToModel(plantID uint, machineForm *machine.Form, machineModel *models.Machine) {
	machineModel.Name = machineForm.Name
	machineModel.Code = machineForm.Code
	machineModel.Status = machineForm.Status
}

func (s *DefaultMachineService) ToForm(plantID uint, machineModel *models.Machine) *machine.Form {
	machineForm := &machine.Form{
		ID:     machineModel.ID,
		Name:   machineModel.Name,
		Code:   machineModel.Code,
		Status: machineModel.Status,
	}
	machineForm.PlantID = machineModel.PlantID
	return machineForm
}

func (s *DefaultMachineService) ToFormSlice(plantID uint, machineModels []*models.Machine) []*machine.Form {
	data := make([]*machine.Form, 0)
	for _, machineModel := range machineModels {
		data = append(data, s.ToForm(plantID, machineModel))
	}
	return data
}

func (s *DefaultMachineService) ToModelSlice(plantID uint, machineForms []*machine.Form) []*models.Machine {
	data := make([]*models.Machine, 0)
	for _, machineForm := range machineForms {
		data = append(data, s.ToModel(plantID, machineForm))
	}
	return data
}
