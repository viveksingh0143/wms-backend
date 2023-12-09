package permission

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"star-wms/app/admin/dto/permission"
	"star-wms/app/admin/service"
	"star-wms/core/common/requests"
	"star-wms/core/common/responses"
	"star-wms/core/validation"
	"strconv"
)

type Handler struct {
	service service.PermissionService
}

func NewPermissionHandler(s service.PermissionService) *Handler {
	return &Handler{
		service: s,
	}
}

// List permissions with filter, pagination, and sorting
func (ph *Handler) List(c *gin.Context) {
	var filter permission.Filter
	paginationValue, _ := c.Get("pagination")
	pagination, _ := paginationValue.(requests.Pagination)

	sortingValue, _ := c.Get("sorting")
	sorting, _ := sortingValue.(requests.Sorting)

	if err := c.ShouldBindQuery(&filter); err != nil {
		resp := responses.NewValidationErrorResponse(err, filter)
		c.JSON(resp.Status, resp)
		return
	}

	permissions, totalRecords, err := ph.service.GetAllPermissions(filter, pagination, sorting)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewPageResponse(permissions, totalRecords, pagination.Page, pagination.PageSize))
}

// Create a new permission
func (ph *Handler) Create(c *gin.Context) {
	var perm permission.Form
	if err := c.ShouldBindJSON(&perm); err != nil {
		resp := responses.NewValidationErrorResponse(err, perm)
		c.JSON(resp.Status, resp)
		return
	}

	validate := validation.GetValidator()
	if err := validate.Struct(perm); err != nil {
		resp := responses.NewValidationErrorResponse(err, perm)
		c.JSON(resp.Status, resp)
		return
	}

	err := ph.service.CreatePermission(&perm)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusCreated, responses.NewSuccessResponse(http.StatusCreated, "Permission created successfully"))
}

// Get a permission
func (ph *Handler) Get(c *gin.Context) {
	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
		c.JSON(resp.Status, resp)
		return
	}

	permissionDto, err := ph.service.GetPermissionByID(id)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessDataResponse(http.StatusOK, "Permission fetched successfully", permissionDto))
}

// Update a permission
func (ph *Handler) Update(c *gin.Context) {
	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
		c.JSON(resp.Status, resp)
		return
	}

	var perm permission.Form
	if err := c.BindJSON(&perm); err != nil {
		resp := responses.NewValidationErrorResponse(err, perm)
		c.JSON(resp.Status, resp)
		return
	}

	validate := validation.GetValidator()
	if err := validate.Struct(perm); err != nil {
		resp := responses.NewValidationErrorResponse(err, perm)
		c.JSON(resp.Status, resp)
		return
	}

	err = ph.service.UpdatePermission(id, &perm)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Permission updated successfully"))
}

// Delete a permission
func (ph *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id") // assuming the ID is passed as a URL parameter
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt < 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
		c.JSON(resp.Status, resp)
		return
	}

	id := uint(idInt)
	err = ph.service.DeletePermission(id)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Permission deleted successfully"))
}

// DeleteBulk a permission
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

	err := ph.service.DeletePermissions(idsForm.IDs)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Permissions deleted successfully"))
}
