package rmbatch

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"star-wms/app/admin/dto/plant"
	"star-wms/app/warehouse/dto/rmbatch"
	"star-wms/app/warehouse/service"
	"star-wms/core/auth"
	"star-wms/core/common/requests"
	"star-wms/core/common/responses"
)

type Handler struct {
	service service.RMBatchService
}

func NewRMBatchHandler(s service.RMBatchService) *Handler {
	return &Handler{
		service: s,
	}
}

// List inventories with filter, pagination, and sorting
func (ph *Handler) List(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	var filter rmbatch.Filter
	paginationValue, _ := c.Get("pagination")
	pagination, _ := paginationValue.(requests.Pagination)

	sortingValue, _ := c.Get("sorting")
	sorting, _ := sortingValue.(requests.Sorting)

	if err := c.ShouldBindQuery(&filter); err != nil {
		resp := responses.NewValidationErrorResponse(err, filter)
		c.JSON(resp.Status, resp)
		return
	}
	inventories, totalRecords, err := ph.service.GetAllRMBatchs(plantForm.ID, filter, pagination, sorting)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewPageResponse(inventories, totalRecords, pagination.Page, pagination.PageSize))
}

// Get a rmbatch
func (ph *Handler) Get(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
		c.JSON(resp.Status, resp)
		return
	}

	rmbatchDto, err := ph.service.GetRMBatchByID(plantForm.ID, id)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessDataResponse(http.StatusOK, "RMBatch fetched successfully", rmbatchDto))
}
