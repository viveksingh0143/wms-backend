package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"star-wms/app/admin/dto/user"
	"star-wms/app/admin/service"
	"star-wms/core/common/requests"
	"star-wms/core/common/responses"
	"star-wms/core/validation"
	"strconv"
)

type Handler struct {
	service service.UserService
}

func NewUserHandler(s service.UserService) *Handler {
	return &Handler{
		service: s,
	}
}

// List users with filter, pagination, and sorting
func (ph *Handler) List(c *gin.Context) {
	var filter user.Filter
	paginationValue, _ := c.Get("pagination")
	pagination, _ := paginationValue.(requests.Pagination)

	sortingValue, _ := c.Get("sorting")
	sorting, _ := sortingValue.(requests.Sorting)

	if err := c.ShouldBindQuery(&filter); err != nil {
		resp := responses.NewValidationErrorResponse(err, filter)
		c.JSON(resp.Status, resp)
		return
	}
	users, totalRecords, err := ph.service.GetAllUsers(filter, pagination, sorting)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewPageResponse(users, totalRecords, pagination.Page, pagination.PageSize))
}

// Create a new user
func (ph *Handler) Create(c *gin.Context) {
	var userForm user.Form
	if err := c.ShouldBindJSON(&userForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, userForm)
		c.JSON(resp.Status, resp)
		return
	}

	validate := validation.GetValidator()
	if userForm.Plant != nil && userForm.Plant.ID == 0 {
		userForm.Plant = nil
	}
	if err := validate.Struct(userForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, userForm)
		c.JSON(resp.Status, resp)
		return
	}

	err := ph.service.CreateUser(&userForm)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusCreated, responses.NewSuccessResponse(http.StatusCreated, "User created successfully"))
}

// Get a user
func (ph *Handler) Get(c *gin.Context) {
	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
		c.JSON(resp.Status, resp)
		return
	}

	userDto, err := ph.service.GetUserByID(id)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessDataResponse(http.StatusOK, "User fetched successfully", userDto))
}

// Update a user
func (ph *Handler) Update(c *gin.Context) {
	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
		c.JSON(resp.Status, resp)
		return
	}

	var userForm user.Form
	if err := c.BindJSON(&userForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, userForm)
		c.JSON(resp.Status, resp)
		return
	}

	if userForm.Plant != nil && userForm.Plant.ID == 0 {
		userForm.Plant = nil
	}

	validate := validation.GetValidator()
	if err := validate.StructExcept(userForm, "Password"); err != nil {
		resp := responses.NewValidationErrorResponse(err, userForm)
		c.JSON(resp.Status, resp)
		return
	}

	err = ph.service.UpdateUser(id, &userForm)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "User updated successfully"))
}

// Delete a user
func (ph *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id") // assuming the ID is passed as a URL parameter
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt < 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
		c.JSON(resp.Status, resp)
		return
	}

	id := uint(idInt)
	err = ph.service.DeleteUser(id)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "User deleted successfully"))
}

// DeleteBulk a user
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

	err := ph.service.DeleteUsers(idsForm.IDs)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Users deleted successfully"))
}
