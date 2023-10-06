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
	var pagination requests.Pagination
	var sorting requests.Sorting

	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, filter))
		return
	}
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, pagination))
		return
	}
	if err := c.ShouldBindQuery(&sorting); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, sorting))
		return
	}
	permissions, totalRecords, err := ph.service.GetAllPermissions(filter, pagination, sorting)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err))
		return
	}
	c.JSON(http.StatusOK, responses.NewPageResponse(permissions, totalRecords, pagination.Page, pagination.PageSize))
}

// Create a new permission
func (ph *Handler) Create(c *gin.Context) {
	var perm permission.Form
	if err := c.ShouldBindJSON(&perm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, perm))
		return
	}

	validate := validation.GetValidator()
	if err := validate.Struct(perm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, perm))
		return
	}

	err := ph.service.CreatePermission(&perm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err))
		return
	}
	c.JSON(http.StatusCreated, responses.NewSuccessResponse(http.StatusCreated, "Permission created successfully"))
}

// Get a permission
func (ph *Handler) Get(c *gin.Context) {
	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err))
		return
	}

	permissionDto, err := ph.service.GetPermissionByID(id)
	if err != nil {
		errResponse := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(errResponse.Status, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessDataResponse(http.StatusOK, "Permission fetched successfully", permissionDto))
}

// Update a permission
func (ph *Handler) Update(c *gin.Context) {
	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err))
		return
	}

	var perm permission.Form
	if err := c.BindJSON(&perm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, perm))
		return
	}

	validate := validation.GetValidator()
	if err := validate.Struct(perm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, perm))
		return
	}

	err = ph.service.UpdatePermission(id, &perm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err))
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Permission updated successfully"))
}

// Delete a permission
func (ph *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id") // assuming the ID is passed as a URL parameter
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt < 0 {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err))
		return
	}

	id := uint(idInt)
	err = ph.service.DeletePermission(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err))
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Permission deleted successfully"))
}

// DeleteBulk a permission
func (ph *Handler) DeleteBulk(c *gin.Context) {
	var idsForm requests.RequestIDs
	if err := c.ShouldBindJSON(&idsForm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, idsForm))
		return
	}

	validate := validation.GetValidator()
	if err := validate.Struct(idsForm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, idsForm))
		return
	}

	err := ph.service.DeletePermissions(idsForm.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err))
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Permissions deleted successfully"))
}
