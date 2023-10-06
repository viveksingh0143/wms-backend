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
	users, totalRecords, err := ph.service.GetAllUsers(filter, pagination, sorting)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err))
		return
	}
	c.JSON(http.StatusOK, responses.NewPageResponse(users, totalRecords, pagination.Page, pagination.PageSize))
}

// Create a new user
func (ph *Handler) Create(c *gin.Context) {
	var userForm user.Form
	if err := c.ShouldBindJSON(&userForm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, userForm))
		return
	}

	validate := validation.GetValidator()
	if userForm.Plant != nil {
		if err := validate.Var(userForm.Plant.ID, "required,gt=0"); err != nil {
			iErr := responses.NewInputError("plant.id", "invalid ID for Plant", userForm.Plant.ID)
			c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, "Invalid Plant ID format", iErr))
		}
	}
	if err := validate.Struct(userForm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, userForm))
		return
	}

	err := ph.service.CreateUser(&userForm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err))
		return
	}
	c.JSON(http.StatusCreated, responses.NewSuccessResponse(http.StatusCreated, "User created successfully"))
}

// Get a user
func (ph *Handler) Get(c *gin.Context) {
	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err))
		return
	}

	userDto, err := ph.service.GetUserByID(id)
	if err != nil {
		errResponse := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(errResponse.Status, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessDataResponse(http.StatusOK, "User fetched successfully", userDto))
}

// Update a user
func (ph *Handler) Update(c *gin.Context) {
	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err))
		return
	}

	var userForm user.Form
	if err := c.BindJSON(&userForm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, userForm))
		return
	}

	validate := validation.GetValidator()
	if err := validate.Struct(userForm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, userForm))
		return
	}

	err = ph.service.UpdateUser(id, &userForm)
	if err != nil {
		errResponse := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(errResponse.Status, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "User updated successfully"))
}

// Delete a user
func (ph *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id") // assuming the ID is passed as a URL parameter
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt < 0 {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err))
		return
	}

	id := uint(idInt)
	err = ph.service.DeleteUser(id)
	if err != nil {
		errResponse := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(errResponse.Status, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "User deleted successfully"))
}

// DeleteBulk a user
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

	err := ph.service.DeleteUsers(idsForm.IDs)
	if err != nil {
		errResponse := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(errResponse.Status, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Users deleted successfully"))
}
