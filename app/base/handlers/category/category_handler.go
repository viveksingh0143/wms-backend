package category

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"star-wms/app/base/dto/category"
	"star-wms/app/base/service"
	"star-wms/core/common/requests"
	"star-wms/core/common/responses"
	"star-wms/core/utils"
	"star-wms/core/validation"
	"strconv"
)

type Handler struct {
	service service.CategoryService
}

func NewCategoryHandler(s service.CategoryService) *Handler {
	return &Handler{
		service: s,
	}
}

// List categories with filter, pagination, and sorting
func (ph *Handler) List(c *gin.Context) {
	var filter category.Filter
	paginationValue, _ := c.Get("pagination")
	pagination, _ := paginationValue.(requests.Pagination)

	sortingValue, _ := c.Get("sorting")
	sorting, _ := sortingValue.(requests.Sorting)

	if err := c.ShouldBindQuery(&filter); err != nil {
		resp := responses.NewValidationErrorResponse(err, filter)
		c.JSON(resp.Status, resp)
		return
	}
	categories, totalRecords, err := ph.service.GetAllCategories(filter, pagination, sorting)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewPageResponse(categories, totalRecords, pagination.Page, pagination.PageSize))
}

// Create a new category
func (ph *Handler) Create(c *gin.Context) {
	var categoryForm category.Form
	if err := c.ShouldBindJSON(&categoryForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, categoryForm)
		c.JSON(resp.Status, resp)
		return
	}

	validate := validation.GetValidator()

	categoryForm.Slug = utils.GenerateSlug(categoryForm.Name)
	if categoryForm.Parent != nil && categoryForm.Parent.ID == 0 {
		categoryForm.Parent = nil
	}

	if err := validate.Struct(categoryForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, categoryForm)
		c.JSON(resp.Status, resp)
		return
	}

	err := ph.service.CreateCategory(&categoryForm)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusCreated, responses.NewSuccessResponse(http.StatusCreated, "Category created successfully"))
}

// Get a category
func (ph *Handler) Get(c *gin.Context) {
	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
		c.JSON(resp.Status, resp)
		return
	}

	categoryDto, err := ph.service.GetCategoryByID(id)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessDataResponse(http.StatusOK, "Category fetched successfully", categoryDto))
}

// Update a category
func (ph *Handler) Update(c *gin.Context) {
	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
		c.JSON(resp.Status, resp)
		return
	}

	var categoryForm category.Form
	if err := c.BindJSON(&categoryForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, categoryForm)
		c.JSON(resp.Status, resp)
		return
	}

	categoryForm.Slug = utils.GenerateSlug(categoryForm.Name)
	if categoryForm.Parent != nil && categoryForm.Parent.ID == 0 {
		categoryForm.Parent = nil
	}

	validate := validation.GetValidator()
	if err := validate.Struct(categoryForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, categoryForm)
		c.JSON(resp.Status, resp)
		return
	}

	err = ph.service.UpdateCategory(id, &categoryForm)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Category updated successfully"))
}

// Delete a category
func (ph *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id") // assuming the ID is passed as a URL parameter
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt < 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
		c.JSON(resp.Status, resp)
		return
	}

	id := uint(idInt)
	err = ph.service.DeleteCategory(id)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Category deleted successfully"))
}

// DeleteBulk a category
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

	err := ph.service.DeleteCategories(idsForm.IDs)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Categories deleted successfully"))
}
