package service

import (
	"star-wms/app/admin/dto/permission"
	"star-wms/app/admin/models"
	"star-wms/app/admin/repository"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/common/responses"
)

type PermissionService interface {
	GetAllPermissions(filter permission.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*permission.Form, int64, error)
	CreatePermission(permissionForm *permission.Form) error
	GetPermissionByID(id uint) (*permission.Form, error)
	UpdatePermission(id uint, permissionForm *permission.Form) error
	DeletePermission(id uint) error
	DeletePermissions(ids []uint) error
	ExistsByModuleName(moduleName string, ID uint) bool
	ToModel(permissionForm *permission.Form) *models.Permission
	ToForm(permModel *models.Permission) *permission.Form
	FormToModel(permissionForm *permission.Form, permModel *models.Permission)
	ToFormSlice(permModels []*models.Permission) []*permission.Form
	ToModelSlice(permModels []*permission.Form) []*models.Permission
}

type DefaultPermissionService struct {
	repo repository.PermissionRepository
}

func NewPermissionService(repo repository.PermissionRepository) PermissionService {
	return &DefaultPermissionService{repo: repo}
}

func (s *DefaultPermissionService) GetAllPermissions(filter permission.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*permission.Form, int64, error) {
	data, count, err := s.repo.GetAll(filter, pagination, sorting)
	if err != nil {
		return nil, count, err
	}
	return s.ToFormSlice(data), count, err
}

func (s *DefaultPermissionService) CreatePermission(permissionForm *permission.Form) error {
	if s.ExistsByModuleName(permissionForm.ModuleName, 0) {
		return responses.NewInputError("module_name", "already exists", permissionForm.ModuleName)
	}
	resultModel := s.ToModel(permissionForm)
	return s.repo.Create(resultModel)
}

func (s *DefaultPermissionService) GetPermissionByID(id uint) (*permission.Form, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.ToForm(data), nil
}

func (s *DefaultPermissionService) UpdatePermission(id uint, permissionForm *permission.Form) error {
	if s.ExistsByModuleName(permissionForm.ModuleName, id) {
		return responses.NewInputError("module_name", "already exists", permissionForm.ModuleName)
	}
	permModel, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	s.FormToModel(permissionForm, permModel)
	return s.repo.Update(permModel)
}

func (s *DefaultPermissionService) DeletePermission(id uint) error {
	return s.repo.Delete(id)
}

func (s *DefaultPermissionService) DeletePermissions(ids []uint) error {
	return s.repo.DeleteMulti(ids)
}

func (s *DefaultPermissionService) ExistsByModuleName(moduleName string, ID uint) bool {
	return s.repo.ExistsByModuleName(moduleName, ID)
}

func (s *DefaultPermissionService) ToModel(permissionForm *permission.Form) *models.Permission {
	result := &models.Permission{
		GroupName:  permissionForm.GroupName,
		ModuleName: permissionForm.ModuleName,
		ReadPerm:   permissionForm.ReadPerm,
		CreatePerm: permissionForm.CreatePerm,
		UpdatePerm: permissionForm.UpdatePerm,
		DeletePerm: permissionForm.DeletePerm,
		ImportPerm: permissionForm.ImportPerm,
		ExportPerm: permissionForm.ExportPerm,
	}
	result.ID = permissionForm.ID
	return result
}

func (s *DefaultPermissionService) FormToModel(permissionForm *permission.Form, permModel *models.Permission) {
	permModel.GroupName = permissionForm.GroupName
	permModel.ModuleName = permissionForm.ModuleName
	permModel.ReadPerm = permissionForm.ReadPerm
	permModel.CreatePerm = permissionForm.CreatePerm
	permModel.UpdatePerm = permissionForm.UpdatePerm
	permModel.DeletePerm = permissionForm.DeletePerm
	permModel.ImportPerm = permissionForm.ImportPerm
	permModel.ExportPerm = permissionForm.ExportPerm
}

func (s *DefaultPermissionService) ToForm(permModel *models.Permission) *permission.Form {
	result := &permission.Form{
		ID:         permModel.ID,
		GroupName:  permModel.GroupName,
		ModuleName: permModel.ModuleName,
		ReadPerm:   permModel.ReadPerm,
		CreatePerm: permModel.CreatePerm,
		UpdatePerm: permModel.UpdatePerm,
		DeletePerm: permModel.DeletePerm,
		ImportPerm: permModel.ImportPerm,
		ExportPerm: permModel.ExportPerm,
	}
	return result
}

func (s *DefaultPermissionService) ToFormSlice(permModels []*models.Permission) []*permission.Form {
	data := make([]*permission.Form, 0)
	for _, permModel := range permModels {
		data = append(data, s.ToForm(permModel))
	}
	return data
}

func (s *DefaultPermissionService) ToModelSlice(permForms []*permission.Form) []*models.Permission {
	data := make([]*models.Permission, 0)
	for _, permForm := range permForms {
		data = append(data, s.ToModel(permForm))
	}
	return data
}
