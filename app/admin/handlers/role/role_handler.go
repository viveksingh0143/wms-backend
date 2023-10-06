package role

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"star-wms/app/admin/dto/role"
	"star-wms/app/admin/service"
	"star-wms/core/common/requests"
	"star-wms/core/common/responses"
	"star-wms/core/validation"
	"strconv"
)

type Handler struct {
	service service.RoleService
}

func NewRoleHandler(s service.RoleService) *Handler {
	return &Handler{
		service: s,
	}
}

// List roles with filter, pagination, and sorting
func (ph *Handler) List(c *gin.Context) {
	var filter role.Filter
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
	roles, totalRecords, err := ph.service.GetAllRoles(filter, pagination, sorting)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err))
		return
	}
	c.JSON(http.StatusOK, responses.NewPageResponse(roles, totalRecords, pagination.Page, pagination.PageSize))
}

// Create a new role
func (ph *Handler) Create(c *gin.Context) {
	var perm role.Form
	if err := c.ShouldBindJSON(&perm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, perm))
		return
	}

	validate := validation.GetValidator()
	if err := validate.Struct(perm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, perm))
		return
	}

	err := ph.service.CreateRole(&perm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err))
		return
	}
	c.JSON(http.StatusCreated, responses.NewSuccessResponse(http.StatusCreated, "Role created successfully"))
}

// Get a role
func (ph *Handler) Get(c *gin.Context) {
	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err))
		return
	}

	roleDto, err := ph.service.GetRoleByID(id)
	if err != nil {
		errResponse := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(errResponse.Status, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessDataResponse(http.StatusOK, "Role fetched successfully", roleDto))
}

// Update a role
func (ph *Handler) Update(c *gin.Context) {
	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err))
		return
	}

	var perm role.Form
	if err := c.BindJSON(&perm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, perm))
		return
	}

	validate := validation.GetValidator()
	if err := validate.Struct(perm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, perm))
		return
	}

	err = ph.service.UpdateRole(id, &perm)
	if err != nil {
		errResponse := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(errResponse.Status, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Role updated successfully"))
}

// Delete a role
func (ph *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id") // assuming the ID is passed as a URL parameter
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt < 0 {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err))
		return
	}

	id := uint(idInt)
	err = ph.service.DeleteRole(id)
	if err != nil {
		errResponse := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(errResponse.Status, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Role deleted successfully"))
}

// DeleteBulk a role
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

	err := ph.service.DeleteRoles(idsForm.IDs)
	if err != nil {
		errResponse := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(errResponse.Status, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Roles deleted successfully"))
}
