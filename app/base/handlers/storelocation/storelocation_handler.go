package storelocation

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"star-wms/app/admin/dto/plant"
	"star-wms/app/base/dto/storelocation"
	"star-wms/app/base/service"
	"star-wms/core/auth"
	"star-wms/core/common/requests"
	"star-wms/core/common/responses"
	"star-wms/core/validation"
	"strconv"
)

type Handler struct {
	service service.StorelocationService
}

func NewStorelocationHandler(s service.StorelocationService) *Handler {
	return &Handler{
		service: s,
	}
}

func (ph *Handler) GetStoreFromRoute(c *gin.Context) (uint, *responses.APIResponse) {
	id, err := requests.StringToUInt(c.Param("storeID"))
	if err != nil || id <= 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid StoreID format", err)
		return 0, &resp
	}
	return id, nil
}

// List storelocations with filter, pagination, and sorting
func (ph *Handler) List(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	storeID, errResp := ph.GetStoreFromRoute(c)
	if errResp != nil {
		c.JSON(errResp.Status, errResp)
		return
	}

	var filter storelocation.Filter
	paginationValue, _ := c.Get("pagination")
	pagination, _ := paginationValue.(requests.Pagination)

	sortingValue, _ := c.Get("sorting")
	sorting, _ := sortingValue.(requests.Sorting)

	if err := c.ShouldBindQuery(&filter); err != nil {
		resp := responses.NewValidationErrorResponse(err, filter)
		c.JSON(resp.Status, resp)
		return
	}
	storelocations, totalRecords, err := ph.service.GetAllStorelocations(plantForm.ID, storeID, filter, pagination, sorting)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewPageResponse(storelocations, totalRecords, pagination.Page, pagination.PageSize))
}

// Create a new storelocation
func (ph *Handler) Create(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	storeID, errResp := ph.GetStoreFromRoute(c)
	if errResp != nil {
		c.JSON(errResp.Status, errResp)
		return
	}

	var storelocationForm storelocation.Form
	if err := c.ShouldBindJSON(&storelocationForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, storelocationForm)
		c.JSON(resp.Status, resp)
		return
	}

	validate := validation.GetValidator()
	if err := validate.Struct(storelocationForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, storelocationForm)
		c.JSON(resp.Status, resp)
		return
	}

	err := ph.service.CreateStorelocation(plantForm.ID, storeID, &storelocationForm)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusCreated, responses.NewSuccessResponse(http.StatusCreated, "Store location created successfully"))
}

// Get a storelocation
func (ph *Handler) Get(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	storeID, errResp := ph.GetStoreFromRoute(c)
	if errResp != nil {
		c.JSON(errResp.Status, errResp)
		return
	}

	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
		c.JSON(resp.Status, resp)
		return
	}

	storelocationDto, err := ph.service.GetStorelocationByID(plantForm.ID, storeID, id)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessDataResponse(http.StatusOK, "Store location fetched successfully", storelocationDto))
}

// Update a storelocation
func (ph *Handler) Update(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	storeID, errResp := ph.GetStoreFromRoute(c)
	if errResp != nil {
		c.JSON(errResp.Status, errResp)
		return
	}

	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
		c.JSON(resp.Status, resp)
		return
	}

	var storelocationForm storelocation.Form
	if err := c.BindJSON(&storelocationForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, storelocationForm)
		c.JSON(resp.Status, resp)
		return
	}

	validate := validation.GetValidator()
	if err := validate.Struct(storelocationForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, storelocationForm)
		c.JSON(resp.Status, resp)
		return
	}

	err = ph.service.UpdateStorelocation(plantForm.ID, storeID, id, &storelocationForm)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Store location updated successfully"))
}

// Delete a storelocation
func (ph *Handler) Delete(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	storeID, errResp := ph.GetStoreFromRoute(c)
	if errResp != nil {
		c.JSON(errResp.Status, errResp)
		return
	}

	idStr := c.Param("id") // assuming the ID is passed as a URL parameter
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt < 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
		c.JSON(resp.Status, resp)
		return
	}

	id := uint(idInt)
	err = ph.service.DeleteStorelocation(plantForm.ID, storeID, id)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Store location deleted successfully"))
}

// DeleteBulk a storelocation
func (ph *Handler) DeleteBulk(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	storeID, errResp := ph.GetStoreFromRoute(c)
	if errResp != nil {
		c.JSON(errResp.Status, errResp)
		return
	}

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

	err := ph.service.DeleteStorelocations(plantForm.ID, storeID, idsForm.IDs)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Store locations deleted successfully"))
}
