package stockapproval

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"star-wms/app/admin/dto/plant"
	"star-wms/app/base/dto/container"
	baseService "star-wms/app/base/service"
	"star-wms/core/auth"
	"star-wms/core/common/requests"
	"star-wms/core/common/responses"
	"star-wms/core/validation"
	"strconv"
)

type Handler struct {
	storeService     baseService.StoreService
	containerService baseService.ContainerService
}

func NewStockapprovalHandler(storeService baseService.StoreService, containerService baseService.ContainerService) *Handler {
	return &Handler{
		storeService:     storeService,
		containerService: containerService,
	}
}

// List stockapprovals with filter, pagination, and sorting
func (ph *Handler) List(c *gin.Context) {
	//authValue, _ := c.Get(auth.AuthUserKey)
	//authUser := authValue.(user.Form)

	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	//stores, err := ph.storeService.GetAllStoresByApprover(plantForm.ID, authUser.ID)
	//if err != nil {
	//	resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
	//	c.JSON(resp.Status, resp)
	//	return
	//}

	var filter container.Filter
	paginationValue, _ := c.Get("pagination")
	pagination, _ := paginationValue.(requests.Pagination)

	sortingValue, _ := c.Get("sorting")
	sorting, _ := sortingValue.(requests.Sorting)

	if err := c.ShouldBindQuery(&filter); err != nil {
		resp := responses.NewValidationErrorResponse(err, filter)
		c.JSON(resp.Status, resp)
		return
	}
	containers, totalRecords, err := ph.containerService.GetAllContainersRequiredApproval(plantForm.ID, filter, pagination, sorting)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewPageResponse(containers, totalRecords, pagination.Page, pagination.PageSize))
}

func (ph *Handler) Approve(c *gin.Context) {
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
	err = ph.containerService.ApproveContainer(plantForm.ID, id)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Container approved successfully"))
}

// ApproveBulk a container
func (ph *Handler) ApproveBulk(c *gin.Context) {
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

	err := ph.containerService.ApproveContainers(plantForm.ID, idsForm.IDs)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Containers approved successfully"))
}
