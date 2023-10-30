package service

import (
	"star-wms/app/base/dto/product"
	"star-wms/app/base/models"
	"star-wms/app/base/repository"
	commonModels "star-wms/core/common/requests"
	"star-wms/core/common/responses"
)

type ProductService interface {
	GetAllProducts(filter product.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*product.Form, int64, error)
	CreateProduct(productForm *product.Form) error
	GetProductByID(id uint) (*product.Form, error)
	UpdateProduct(id uint, productForm *product.Form) error
	DeleteProduct(id uint) error
	DeleteProducts(ids []uint) error
	ExistsById(ID uint) bool
	ExistsByName(name string, ID uint) bool
	ExistsBySlug(slug string, ID uint) bool
	ExistsByCode(code string, ID uint) bool
	ExistsByCmsCode(cmsCode string, ID uint) bool
	ToModel(productForm *product.Form) *models.Product
	FormToModel(productForm *product.Form, productModel *models.Product)
	ToForm(productModel *models.Product) *product.Form
	ToFormSlice(productModels []*models.Product) []*product.Form
	ToModelSlice(productForms []*product.Form) []*models.Product
}

type DefaultProductService struct {
	repo            repository.ProductRepository
	categoryService CategoryService
}

func NewProductService(repo repository.ProductRepository, categoryService CategoryService) ProductService {
	return &DefaultProductService{repo: repo, categoryService: categoryService}
}

func (s *DefaultProductService) GetAllProducts(filter product.Filter, pagination commonModels.Pagination, sorting commonModels.Sorting) ([]*product.Form, int64, error) {
	data, count, err := s.repo.GetAll(filter, pagination, sorting)
	if err != nil {
		return nil, count, err
	}
	return s.ToFormSlice(data), count, err
}

func (s *DefaultProductService) CreateProduct(productForm *product.Form) error {
	if s.ExistsByName(productForm.Name, 0) {
		return responses.NewInputError("name", "already exists", productForm.Name)
	}
	if s.ExistsBySlug(productForm.Slug, 0) {
		return responses.NewInputError("slug", "already exists", productForm.Slug)
	}
	if s.ExistsByCode(productForm.Code, 0) {
		return responses.NewInputError("code", "already exists", productForm.Code)
	}
	if productForm.CmsCode != "" && s.ExistsByCmsCode(productForm.CmsCode, 0) {
		return responses.NewInputError("cms_code", "already exists", productForm.CmsCode)
	}
	if productForm.Category != nil {
		if !s.categoryService.ExistsById(productForm.Category.ID) {
			return responses.NewInputError("category.id", "category not exists", productForm.Category.ID)
		}
	}
	productModel := s.ToModel(productForm)
	return s.repo.Create(productModel)
}

func (s *DefaultProductService) GetProductByID(id uint) (*product.Form, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.ToForm(data), nil
}

func (s *DefaultProductService) UpdateProduct(id uint, productForm *product.Form) error {
	if s.ExistsByName(productForm.Name, id) {
		return responses.NewInputError("name", "already exists", productForm.Name)
	}
	if s.ExistsBySlug(productForm.Slug, id) {
		return responses.NewInputError("slug", "already exists", productForm.Slug)
	}
	if s.ExistsByCode(productForm.Code, id) {
		return responses.NewInputError("code", "already exists", productForm.Code)
	}
	if productForm.CmsCode != "" && s.ExistsByCmsCode(productForm.CmsCode, id) {
		return responses.NewInputError("cms_code", "already exists", productForm.CmsCode)
	}
	if productForm.Category != nil {
		if !s.categoryService.ExistsById(productForm.Category.ID) {
			return responses.NewInputError("category.id", "category not exists", productForm.Category.ID)
		}
	}
	productModel, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	s.FormToModel(productForm, productModel)
	return s.repo.Update(productModel)
}

func (s *DefaultProductService) DeleteProduct(id uint) error {
	return s.repo.Delete(id)
}

func (s *DefaultProductService) DeleteProducts(ids []uint) error {
	return s.repo.DeleteMulti(ids)
}

func (s *DefaultProductService) ExistsById(ID uint) bool {
	return s.repo.ExistsByID(ID)
}

func (s *DefaultProductService) ExistsByName(name string, ID uint) bool {
	return s.repo.ExistsByName(name, ID)
}

func (s *DefaultProductService) ExistsBySlug(slug string, ID uint) bool {
	return s.repo.ExistsBySlug(slug, ID)
}

func (s *DefaultProductService) ExistsByCode(code string, ID uint) bool {
	return s.repo.ExistsByCode(code, ID)
}

func (s *DefaultProductService) ExistsByCmsCode(cmsCode string, ID uint) bool {
	return s.repo.ExistsByCmsCode(cmsCode, ID)
}

func (s *DefaultProductService) ToModel(productForm *product.Form) *models.Product {
	productModel := &models.Product{
		ProductType: models.ProductType(productForm.ProductType),
		Name:        productForm.Name,
		Slug:        productForm.Slug,
		Code:        productForm.Code,
		CmsCode:     productForm.CmsCode,
		Description: productForm.Description,
		UnitType:    models.UnitType(productForm.UnitType),
		UnitWeight:  productForm.UnitWeight,
		UnitValue:   models.UnitValue(productForm.UnitValue),
		Status:      productForm.Status,
	}
	productModel.ID = productForm.ID

	if productForm.Category != nil {
		productModel.Category = s.categoryService.ToModel(productForm.Category)
	}
	return productModel
}

func (s *DefaultProductService) FormToModel(productForm *product.Form, productModel *models.Product) {
	productModel.ProductType = models.ProductType(productForm.ProductType)
	productModel.Name = productForm.Name
	productModel.Slug = productForm.Slug
	productModel.Code = productForm.Code
	productModel.CmsCode = productForm.CmsCode
	productModel.Description = productForm.Description
	productModel.UnitType = models.UnitType(productForm.UnitType)
	productModel.UnitWeight = productForm.UnitWeight
	productModel.UnitValue = models.UnitValue(productForm.UnitValue)
	productModel.Status = productForm.Status

	if productForm.Category != nil {
		productModel.Category = s.categoryService.ToModel(productForm.Category)
	} else {
		productModel.Category = nil
		productModel.CategoryID = nil
	}
}

func (s *DefaultProductService) ToForm(productModel *models.Product) *product.Form {
	productForm := &product.Form{
		ID:          productModel.ID,
		ProductType: string(productModel.ProductType),
		Name:        productModel.Name,
		Slug:        productModel.Slug,
		Code:        productModel.Code,
		CmsCode:     productModel.CmsCode,
		Description: productModel.Description,
		UnitType:    string(productModel.UnitType),
		UnitWeight:  productModel.UnitWeight,
		UnitValue:   string(productModel.UnitValue),
		Status:      productModel.Status,
	}
	if productModel.Category != nil {
		productForm.Category = s.categoryService.ToForm(productModel.Category)
	}
	return productForm
}

func (s *DefaultProductService) ToFormSlice(productModels []*models.Product) []*product.Form {
	data := make([]*product.Form, 0)
	for _, productModel := range productModels {
		data = append(data, s.ToForm(productModel))
	}
	return data
}

func (s *DefaultProductService) ToModelSlice(productForms []*product.Form) []*models.Product {
	data := make([]*models.Product, 0)
	for _, productForm := range productForms {
		data = append(data, s.ToModel(productForm))
	}
	return data
}
