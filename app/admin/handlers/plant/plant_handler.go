package plant

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"star-wms/app/admin/dto/plant"
	"star-wms/app/admin/service"
	"star-wms/core/common/requests"
	"star-wms/core/common/responses"
	"star-wms/core/validation"
	"strconv"
)

type Handler struct {
	service service.PlantService
}

func NewPlantHandler(s service.PlantService) *Handler {
	return &Handler{
		service: s,
	}
}

// List plants with filter, pagination, and sorting
func (ph *Handler) List(c *gin.Context) {
	var filter plant.Filter
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
	plants, totalRecords, err := ph.service.GetAllPlants(filter, pagination, sorting)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err))
		return
	}
	c.JSON(http.StatusOK, responses.NewPageResponse(plants, totalRecords, pagination.Page, pagination.PageSize))
}

// Create a new plant
func (ph *Handler) Create(c *gin.Context) {
	var perm plant.Form
	if err := c.ShouldBindJSON(&perm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, perm))
		return
	}

	validate := validation.GetValidator()
	if err := validate.Struct(perm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, perm))
		return
	}

	err := ph.service.CreatePlant(&perm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err))
		return
	}
	c.JSON(http.StatusCreated, responses.NewSuccessResponse(http.StatusCreated, "Plant created successfully"))
}

// Get a plant
func (ph *Handler) Get(c *gin.Context) {
	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err))
		return
	}

	plantDto, err := ph.service.GetPlantByID(id)
	if err != nil {
		errResponse := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(errResponse.Status, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessDataResponse(http.StatusOK, "Plant fetched successfully", plantDto))
}

// Update a plant
func (ph *Handler) Update(c *gin.Context) {
	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err))
		return
	}

	var perm plant.Form
	if err := c.BindJSON(&perm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, perm))
		return
	}

	validate := validation.GetValidator()
	if err := validate.Struct(perm); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewValidationErrorResponse(err, perm))
		return
	}

	err = ph.service.UpdatePlant(id, &perm)
	if err != nil {
		errResponse := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(errResponse.Status, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Plant updated successfully"))
}

// Delete a plant
func (ph *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id") // assuming the ID is passed as a URL parameter
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt < 0 {
		c.JSON(http.StatusBadRequest, responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err))
		return
	}

	id := uint(idInt)
	err = ph.service.DeletePlant(id)
	if err != nil {
		errResponse := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(errResponse.Status, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Plant deleted successfully"))
}

// DeleteBulk a plant
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

	err := ph.service.DeletePlants(idsForm.IDs)
	if err != nil {
		errResponse := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(errResponse.Status, errResponse)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Plants deleted successfully"))
}
