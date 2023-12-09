package product

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"star-wms/app/base/dto/product"
	"star-wms/app/base/service"
	"star-wms/core/common/requests"
	"star-wms/core/common/responses"
	"star-wms/core/utils"
	"star-wms/core/validation"
	"strconv"
)

type Handler struct {
	service service.ProductService
}

func NewProductHandler(s service.ProductService) *Handler {
	return &Handler{
		service: s,
	}
}

// List products with filter, pagination, and sorting
func (ph *Handler) List(c *gin.Context) {
	var filter product.Filter
	paginationValue, _ := c.Get("pagination")
	pagination, _ := paginationValue.(requests.Pagination)

	sortingValue, _ := c.Get("sorting")
	sorting, _ := sortingValue.(requests.Sorting)

	if err := c.ShouldBindQuery(&filter); err != nil {
		resp := responses.NewValidationErrorResponse(err, filter)
		c.JSON(resp.Status, resp)
		return
	}
	products, totalRecords, err := ph.service.GetAllProducts(filter, pagination, sorting)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewPageResponse(products, totalRecords, pagination.Page, pagination.PageSize))
}

// Create a new product
func (ph *Handler) Create(c *gin.Context) {
	var productForm product.Form
	if err := c.ShouldBindJSON(&productForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, productForm)
		c.JSON(resp.Status, resp)
		return
	}

	validate := validation.GetValidator()

	productForm.Slug = utils.GenerateSlug(productForm.Name)
	if productForm.Category != nil && productForm.Category.ID == 0 {
		productForm.Category = nil
	}

	if err := validate.Struct(productForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, productForm)
		c.JSON(resp.Status, resp)
		return
	}

	err := ph.service.CreateProduct(&productForm)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusCreated, responses.NewSuccessResponse(http.StatusCreated, "Product created successfully"))
}

// Get a product
func (ph *Handler) Get(c *gin.Context) {
	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
		c.JSON(resp.Status, resp)
		return
	}

	productDto, err := ph.service.GetProductByID(id)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessDataResponse(http.StatusOK, "Product fetched successfully", productDto))
}

// Update a product
func (ph *Handler) Update(c *gin.Context) {
	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
		c.JSON(resp.Status, resp)
		return
	}

	var productForm product.Form
	if err := c.BindJSON(&productForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, productForm)
		c.JSON(resp.Status, resp)
		return
	}

	productForm.Slug = utils.GenerateSlug(productForm.Name)
	if productForm.Category != nil && productForm.Category.ID == 0 {
		productForm.Category = nil
	}

	validate := validation.GetValidator()
	if err := validate.Struct(productForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, productForm)
		c.JSON(resp.Status, resp)
		return
	}

	err = ph.service.UpdateProduct(id, &productForm)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Product updated successfully"))
}

// Delete a product
func (ph *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id") // assuming the ID is passed as a URL parameter
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt < 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
		c.JSON(resp.Status, resp)
		return
	}

	id := uint(idInt)
	err = ph.service.DeleteProduct(id)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Product deleted successfully"))
}

// DeleteBulk a product
func (ph *Handler) DeleteBulk(c *gin.Context) {
	var idsForm requests.RequestIDs
	if err := c.ShouldBindJSON(&idsForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, idsForm)
		c.JSON(resp.Status, resp)
		return
	}

	validate := validation.GetValidator()
	if err := validate.Struct(idsForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, idsForm)
		c.JSON(resp.Status, resp)
		return
	}

	err := ph.service.DeleteProducts(idsForm.IDs)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Products deleted successfully"))
}
