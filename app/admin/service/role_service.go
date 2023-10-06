package service

import (
	"star-wms/app/admin/dto/ability"
	"star-wms/app/admin/dto/role"
	"star-wms/app/admin/models"
	"star-wms/app/admin/repository"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/common/responses"
)

type RoleService interface {
	GetAllRoles(filter role.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*role.Form, int64, error)
	CreateRole(roleForm *role.Form) error
	GetRoleByID(id uint) (*role.Form, error)
	UpdateRole(id uint, roleForm *role.Form) error
	DeleteRole(id uint) error
	DeleteRoles(ids []uint) error
	ExistsByName(name string, ID uint) bool
	FormToModel(roleForm *role.Form, roleModel *models.Role)
	ToModel(roleForm *role.Form) *models.Role
	ToForm(roleModel *models.Role) *role.Form
	ToFormSlice(roleModels []*models.Role) []*role.Form
	ToModelSlice(roleForms []*role.Form) []*models.Role
}

type DefaultRoleService struct {
	repo repository.RoleRepository
}

func NewRoleService(repo repository.RoleRepository) RoleService {
	return &DefaultRoleService{repo: repo}
}

func (s *DefaultRoleService) GetAllRoles(filter role.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*role.Form, int64, error) {
	data, count, err := s.repo.GetAll(filter, pagination, sorting)
	if err != nil {
		return nil, count, err
	}
	return s.ToFormSlice(data), count, err
}

func (s *DefaultRoleService) CreateRole(roleForm *role.Form) error {
	if s.ExistsByName(roleForm.Name, 0) {
		return responses.NewInputError("name", "already exists", roleForm.Name)
	}
	resultModel := s.ToModel(roleForm)
	return s.repo.Create(resultModel)
}

func (s *DefaultRoleService) GetRoleByID(id uint) (*role.Form, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.ToForm(data), nil
}

func (s *DefaultRoleService) UpdateRole(id uint, roleForm *role.Form) error {
	if s.ExistsByName(roleForm.Name, id) {
		return responses.NewInputError("name", "already exists", roleForm.Name)
	}
	roleModel, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	s.FormToModel(roleForm, roleModel)
	return s.repo.Update(roleModel)
}

func (s *DefaultRoleService) DeleteRole(id uint) error {
	return s.repo.Delete(id)
}

func (s *DefaultRoleService) DeleteRoles(ids []uint) error {
	return s.repo.DeleteMulti(ids)
}

func (s *DefaultRoleService) ExistsByName(moduleName string, ID uint) bool {
	return s.repo.ExistsByName(moduleName, ID)
}

func (s *DefaultRoleService) ToModel(roleForm *role.Form) *models.Role {
	roleModel := &models.Role{
		Name:   roleForm.Name,
		Status: roleForm.Status,
	}
	roleModel.ID = roleForm.ID
	if roleForm.Abilities != nil {
		abilities := make([]*models.Ability, 0)
		if len(roleForm.Abilities) > 0 {
			for _, abilityForm := range roleForm.Abilities {
				abilities = append(abilities, &models.Ability{
					Module: abilityForm.Module,
					Name:   abilityForm.Name,
				})
			}
		}
		roleModel.Abilities = abilities
	}
	return roleModel
}

func (s *DefaultRoleService) FormToModel(roleForm *role.Form, roleModel *models.Role) {
	roleModel.Name = roleForm.Name
	roleModel.Status = roleForm.Status
	abilities := make([]*models.Ability, 0)
	if len(roleForm.Abilities) > 0 {
		for _, abilityForm := range roleForm.Abilities {
			abilities = append(abilities, &models.Ability{
				Module: abilityForm.Module,
				Name:   abilityForm.Name,
			})
		}
	}
	roleModel.Abilities = abilities
}

func (s *DefaultRoleService) ToForm(roleModel *models.Role) *role.Form {
	roleForm := &role.Form{
		ID:     roleModel.ID,
		Name:   roleModel.Name,
		Status: roleModel.Status,
	}
	abilities := make([]*ability.Form, 0)
	if len(roleModel.Abilities) > 0 {
		for _, abilityModel := range roleModel.Abilities {
			abilities = append(abilities, &ability.Form{
				ID:     abilityModel.ID,
				Module: abilityModel.Module,
				Name:   abilityModel.Name,
			})
		}
	}
	roleForm.Abilities = abilities
	return roleForm
}

func (s *DefaultRoleService) ToFormSlice(roleModels []*models.Role) []*role.Form {
	data := make([]*role.Form, 0)
	for _, roleModel := range roleModels {
		data = append(data, s.ToForm(roleModel))
	}
	return data
}

func (s *DefaultRoleService) ToModelSlice(roleForms []*role.Form) []*models.Role {
	data := make([]*models.Role, 0)
	for _, roleForm := range roleForms {
		data = append(data, s.ToModel(roleForm))
	}
	return data
}
