package machine

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"star-wms/app/admin/dto/plant"
	"star-wms/app/base/dto/machine"
	"star-wms/app/base/service"
	"star-wms/core/auth"
	"star-wms/core/common/requests"
	"star-wms/core/common/responses"
	"star-wms/core/validation"
	"strconv"
)

type Handler struct {
	service service.MachineService
}

func NewMachineHandler(s service.MachineService) *Handler {
	return &Handler{
		service: s,
	}
}

// List machines with filter, pagination, and sorting
func (ph *Handler) List(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	var filter machine.Filter
	paginationValue, _ := c.Get("pagination")
	pagination, _ := paginationValue.(requests.Pagination)

	sortingValue, _ := c.Get("sorting")
	sorting, _ := sortingValue.(requests.Sorting)

	if err := c.ShouldBindQuery(&filter); err != nil {
		resp := responses.NewValidationErrorResponse(err, filter)
		c.JSON(resp.Status, resp)
		return
	}
	machines, totalRecords, err := ph.service.GetAllMachines(plantForm.ID, filter, pagination, sorting)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewPageResponse(machines, totalRecords, pagination.Page, pagination.PageSize))
}

// Create a new machine
func (ph *Handler) Create(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	var machineForm machine.Form
	if err := c.ShouldBindJSON(&machineForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, machineForm)
		c.JSON(resp.Status, resp)
		return
	}

	validate := validation.GetValidator()
	if err := validate.Struct(machineForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, machineForm)
		c.JSON(resp.Status, resp)
		return
	}

	err := ph.service.CreateMachine(plantForm.ID, &machineForm)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusCreated, responses.NewSuccessResponse(http.StatusCreated, "Machine created successfully"))
}

// Get a machine
func (ph *Handler) Get(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
		c.JSON(resp.Status, resp)
		return
	}

	machineDto, err := ph.service.GetMachineByID(plantForm.ID, id)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessDataResponse(http.StatusOK, "Machine fetched successfully", machineDto))
}

// Update a machine
func (ph *Handler) Update(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
		c.JSON(resp.Status, resp)
		return
	}

	var machineForm machine.Form
	if err := c.BindJSON(&machineForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, machineForm)
		c.JSON(resp.Status, resp)
		return
	}

	validate := validation.GetValidator()
	if err := validate.Struct(machineForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, machineForm)
		c.JSON(resp.Status, resp)
		return
	}

	err = ph.service.UpdateMachine(plantForm.ID, id, &machineForm)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Machine updated successfully"))
}

// Delete a machine
func (ph *Handler) Delete(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	idStr := c.Param("id") // assuming the ID is passed as a URL parameter
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt < 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
		c.JSON(resp.Status, resp)
		return
	}

	id := uint(idInt)
	err = ph.service.DeleteMachine(plantForm.ID, id)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Machine deleted successfully"))
}

// DeleteBulk a machine
func (ph *Handler) DeleteBulk(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

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

	err := ph.service.DeleteMachines(plantForm.ID, idsForm.IDs)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Machines deleted successfully"))
}
