package sticker

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"star-wms/app/admin/dto/plant"
	"star-wms/app/warehouse/dto/batchlabel"
	"star-wms/app/warehouse/service"
	"star-wms/core/auth"
	"star-wms/core/common/requests"
	"star-wms/core/common/responses"
	"star-wms/core/validation"
)

type Handler struct {
	service service.BatchlabelService
}

func NewStickerHandler(s service.BatchlabelService) *Handler {
	return &Handler{
		service: s,
	}
}

func (ph *Handler) GetBatchlabelFromRoute(c *gin.Context) (uint, *responses.APIResponse) {
	id, err := requests.StringToUInt(c.Param("batchlabelID"))
	if err != nil || id <= 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid BatchlabelID format", err)
		return 0, &resp
	}
	return id, nil
}

// List stickers with filter, pagination, and sorting
func (ph *Handler) List(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	batchlabelID, errResp := ph.GetBatchlabelFromRoute(c)
	if errResp != nil {
		c.JSON(errResp.Status, errResp)
		return
	}

	var filter batchlabel.StickerFilter
	paginationValue, _ := c.Get("pagination")
	pagination, _ := paginationValue.(requests.Pagination)

	sortingValue, _ := c.Get("sorting")
	sorting, _ := sortingValue.(requests.Sorting)

	if err := c.ShouldBindQuery(&filter); err != nil {
		resp := responses.NewValidationErrorResponse(err, filter)
		c.JSON(resp.Status, resp)
		return
	}
	stickers, totalRecords, err := ph.service.GetAllStickers(plantForm.ID, batchlabelID, filter, pagination, sorting)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewPageResponse(stickers, totalRecords, pagination.Page, pagination.PageSize))
}

// Create a new sticker
func (ph *Handler) Create(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	batchlabelID, errResp := ph.GetBatchlabelFromRoute(c)
	if errResp != nil {
		c.JSON(errResp.Status, errResp)
		return
	}

	var stickerForm batchlabel.MultiStickerForm
	if err := c.ShouldBindJSON(&stickerForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, stickerForm)
		c.JSON(resp.Status, resp)
		return
	}

	validate := validation.GetValidator()

	if err := validate.Struct(stickerForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, stickerForm)
		c.JSON(resp.Status, resp)
		return
	}

	err := ph.service.CreateSticker(plantForm.ID, batchlabelID, &stickerForm)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusCreated, responses.NewSuccessResponse(http.StatusCreated, "Sticker created successfully"))
}

// Get a sticker
func (ph *Handler) Get(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	batchlabelID, errResp := ph.GetBatchlabelFromRoute(c)
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

	stickerDto, err := ph.service.GetStickerByID(plantForm.ID, batchlabelID, id)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessDataResponse(http.StatusOK, "Sticker fetched successfully", stickerDto))
}
