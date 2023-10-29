package service

import (
	"star-wms/app/base/dto/category"
	"star-wms/app/base/models"
	"star-wms/app/base/repository"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/common/responses"
)

type CategoryService interface {
	GetAllCategories(filter category.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*category.Form, int64, error)
	CreateCategory(categoryForm *category.Form) error
	GetCategoryByID(id uint) (*category.Form, error)
	UpdateCategory(id uint, categoryForm *category.Form) error
	DeleteCategory(id uint) error
	DeleteCategories(ids []uint) error
	ExistsById(ID uint) bool
	ExistsByName(name string, ID uint) bool
	ExistsBySlug(slug string, ID uint) bool
	ToModel(categoryForm *category.Form) *models.Category
	FormToModel(categoryForm *category.Form, categoryModel *models.Category)
	ToForm(categoryModel *models.Category) *category.Form
	ToFormSlice(categoryModels []*models.Category) []*category.Form
	ToModelSlice(categoryForms []*category.Form) []*models.Category
}

type DefaultCategoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &DefaultCategoryService{repo: repo}
}

func (s *DefaultCategoryService) GetAllCategories(filter category.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*category.Form, int64, error) {
	data, count, err := s.repo.GetAll(filter, pagination, sorting)
	if err != nil {
		return nil, count, err
	}
	return s.ToFormSlice(data), count, err
}

func (s *DefaultCategoryService) CreateCategory(categoryForm *category.Form) error {
	if s.ExistsByName(categoryForm.Name, 0) {
		return responses.NewInputError("name", "already exists", categoryForm.Name)
	}
	if s.ExistsBySlug(categoryForm.Slug, 0) {
		return responses.NewInputError("slug", "already exists", categoryForm.Slug)
	}
	if categoryForm.Parent != nil {
		if !s.ExistsById(categoryForm.Parent.ID) {
			return responses.NewInputError("parent.id", "parent not exists", categoryForm.Parent.ID)
		}
	}
	categoryModel := s.ToModel(categoryForm)
	return s.repo.Create(categoryModel)
}

func (s *DefaultCategoryService) GetCategoryByID(id uint) (*category.Form, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.ToForm(data), nil
}

func (s *DefaultCategoryService) UpdateCategory(id uint, categoryForm *category.Form) error {
	if s.ExistsByName(categoryForm.Name, id) {
		return responses.NewInputError("name", "already exists", categoryForm.Name)
	}
	if s.ExistsBySlug(categoryForm.Slug, id) {
		return responses.NewInputError("slug", "already exists", categoryForm.Slug)
	}
	if categoryForm.Parent != nil {
		if !s.ExistsById(categoryForm.Parent.ID) {
			return responses.NewInputError("parent.id", "parent not exists", categoryForm.Parent.ID)
		}
	}
	categoryModel, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	s.FormToModel(categoryForm, categoryModel)
	return s.repo.Update(categoryModel)
}

func (s *DefaultCategoryService) DeleteCategory(id uint) error {
	return s.repo.Delete(id)
}

func (s *DefaultCategoryService) DeleteCategories(ids []uint) error {
	return s.repo.DeleteMulti(ids)
}

func (s *DefaultCategoryService) ExistsById(ID uint) bool {
	return s.repo.ExistsByID(ID)
}

func (s *DefaultCategoryService) ExistsByName(name string, ID uint) bool {
	return s.repo.ExistsByName(name, ID)
}

func (s *DefaultCategoryService) ExistsBySlug(slug string, ID uint) bool {
	return s.repo.ExistsBySlug(slug, ID)
}

func (s *DefaultCategoryService) ToModel(categoryForm *category.Form) *models.Category {
	categoryModel := &models.Category{
		Name:   categoryForm.Name,
		Slug:   categoryForm.Slug,
		Status: categoryForm.Status,
	}
	categoryModel.ID = categoryForm.ID

	if categoryForm.Parent != nil {
		categoryModel.Parent = s.ToModel(categoryForm.Parent)
	}
	return categoryModel
}

func (s *DefaultCategoryService) FormToModel(categoryForm *category.Form, categoryModel *models.Category) {
	categoryModel.Name = categoryForm.Name
	categoryModel.Status = categoryForm.Status
	categoryModel.Slug = categoryForm.Slug

	if categoryForm.Parent != nil {
		categoryModel.Parent = s.ToModel(categoryForm.Parent)
	} else {
		categoryModel.Parent = nil
		categoryModel.ParentID = nil
	}
}

func (s *DefaultCategoryService) ToForm(categoryModel *models.Category) *category.Form {
	categoryForm := &category.Form{
		ID:       categoryModel.ID,
		Name:     categoryModel.Name,
		Slug:     categoryModel.Slug,
		FullPath: categoryModel.FullPath,
		Status:   categoryModel.Status,
	}
	if categoryModel.Children != nil {
		children := make([]*category.Form, 0)
		if len(categoryModel.Children) > 0 {
			children = s.ToFormSlice(categoryModel.Children)
		}
		categoryForm.Children = children
	}
	if categoryModel.Parent != nil {
		categoryForm.Parent = s.ToForm(categoryModel.Parent)
	}
	return categoryForm
}

func (s *DefaultCategoryService) ToFormSlice(categoryModels []*models.Category) []*category.Form {
	data := make([]*category.Form, 0)
	for _, categoryModel := range categoryModels {
		data = append(data, s.ToForm(categoryModel))
	}
	return data
}

func (s *DefaultCategoryService) ToModelSlice(categoryForms []*category.Form) []*models.Category {
	data := make([]*models.Category, 0)
	for _, categoryForm := range categoryForms {
		data = append(data, s.ToModel(categoryForm))
	}
	return data
}
