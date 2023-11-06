package outwardrequest

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"star-wms/app/admin/dto/plant"
	"star-wms/app/base/dto/outwardrequest"
	"star-wms/app/base/service"
	"star-wms/core/auth"
	"star-wms/core/common/requests"
	"star-wms/core/common/responses"
	"star-wms/core/validation"
	"strconv"
)

type Handler struct {
	service service.OutwardrequestService
}

func NewOutwardrequestHandler(s service.OutwardrequestService) *Handler {
	return &Handler{
		service: s,
	}
}

// List outwardrequests with filter, pagination, and sorting
func (ph *Handler) List(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	var filter outwardrequest.Filter
	paginationValue, _ := c.Get("pagination")
	pagination, _ := paginationValue.(requests.Pagination)

	sortingValue, _ := c.Get("sorting")
	sorting, _ := sortingValue.(requests.Sorting)

	if err := c.ShouldBindQuery(&filter); err != nil {
		resp := responses.NewValidationErrorResponse(err, filter)
		c.JSON(resp.Status, resp)
		return
	}
	outwardrequests, totalRecords, err := ph.service.GetAllOutwardrequests(plantForm.ID, filter, pagination, sorting)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewPageResponse(outwardrequests, totalRecords, pagination.Page, pagination.PageSize))
}

// Create a new outwardrequest
func (ph *Handler) Create(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	var outwardrequestForm outwardrequest.Form
	if err := c.ShouldBindJSON(&outwardrequestForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, outwardrequestForm)
		c.JSON(resp.Status, resp)
		return
	}

	validate := validation.GetValidator()

	if outwardrequestForm.Customer != nil && outwardrequestForm.Customer.ID == 0 {
		outwardrequestForm.Customer = nil
	}

	if err := validate.Struct(outwardrequestForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, outwardrequestForm)
		c.JSON(resp.Status, resp)
		return
	}

	err := ph.service.CreateOutwardrequest(plantForm.ID, &outwardrequestForm)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusCreated, responses.NewSuccessResponse(http.StatusCreated, "Outwardrequest created successfully"))
}

// Get a outwardrequest
func (ph *Handler) Get(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
		c.JSON(resp.Status, resp)
		return
	}

	outwardrequestDto, err := ph.service.GetOutwardrequestByID(plantForm.ID, id)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessDataResponse(http.StatusOK, "Outwardrequest fetched successfully", outwardrequestDto))
}

// Update a outwardrequest
func (ph *Handler) Update(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
		c.JSON(resp.Status, resp)
		return
	}

	var outwardrequestForm outwardrequest.Form
	if err := c.BindJSON(&outwardrequestForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, outwardrequestForm)
		c.JSON(resp.Status, resp)
		return
	}

	if outwardrequestForm.Customer != nil && outwardrequestForm.Customer.ID == 0 {
		outwardrequestForm.Customer = nil
	}

	validate := validation.GetValidator()
	if err := validate.Struct(outwardrequestForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, outwardrequestForm)
		c.JSON(resp.Status, resp)
		return
	}

	err = ph.service.UpdateOutwardrequest(plantForm.ID, id, &outwardrequestForm)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Outwardrequest updated successfully"))
}

// Delete a outwardrequest
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
	err = ph.service.DeleteOutwardrequest(plantForm.ID, id)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Outwardrequest deleted successfully"))
}

// DeleteBulk a outwardrequest
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

	err := ph.service.DeleteOutwardrequests(plantForm.ID, idsForm.IDs)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Outwardrequests deleted successfully"))
}
