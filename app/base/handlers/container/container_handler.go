package container

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"star-wms/app/admin/dto/plant"
	"star-wms/app/base/dto/container"
	"star-wms/app/base/models"
	"star-wms/app/base/service"
	"star-wms/core/auth"
	"star-wms/core/common/requests"
	"star-wms/core/common/responses"
	"star-wms/core/validation"
	"strconv"
)

type Handler struct {
	service service.ContainerService
}

func NewContainerHandler(s service.ContainerService) *Handler {
	return &Handler{
		service: s,
	}
}

// List containers with filter, pagination, and sorting
func (ph *Handler) List(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

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
	containers, totalRecords, err := ph.service.GetAllContainers(plantForm.ID, filter, pagination, sorting)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewPageResponse(containers, totalRecords, pagination.Page, pagination.PageSize))
}

// Create a new container
func (ph *Handler) Create(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	var containerForm container.Form
	if err := c.ShouldBindJSON(&containerForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, containerForm)
		c.JSON(resp.Status, resp)
		return
	}

	validate := validation.GetValidator()
	if containerForm.Product != nil && containerForm.Product.ID == 0 {
		containerForm.Product = nil
	}
	if containerForm.Store != nil && containerForm.Store.ID == 0 {
		containerForm.Store = nil
	}

	if err := validate.Struct(containerForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, containerForm)
		c.JSON(resp.Status, resp)
		return
	}

	err := ph.service.CreateContainer(plantForm.ID, &containerForm)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusCreated, responses.NewSuccessResponse(http.StatusCreated, "Container created successfully"))
}

// Get a container
func (ph *Handler) Get(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
		c.JSON(resp.Status, resp)
		return
	}

	containerDto, err := ph.service.GetContainerByID(plantForm.ID, id)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessDataResponse(http.StatusOK, "Container fetched successfully", containerDto))
}

// Update a container
func (ph *Handler) Update(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
		c.JSON(resp.Status, resp)
		return
	}

	var containerForm container.Form
	if err := c.BindJSON(&containerForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, containerForm)
		c.JSON(resp.Status, resp)
		return
	}

	if containerForm.Store != nil && containerForm.Store.ID == 0 {
		containerForm.Store = nil
	}

	if containerForm.Product != nil && containerForm.Product.ID == 0 {
		containerForm.Product = nil
	}

	validate := validation.GetValidator()
	if err := validate.Struct(containerForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, containerForm)
		c.JSON(resp.Status, resp)
		return
	}

	err = ph.service.UpdateContainer(plantForm.ID, id, &containerForm)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Container updated successfully"))
}

// Delete a container
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
	err = ph.service.DeleteContainer(plantForm.ID, id)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Container deleted successfully"))
}

// DeleteBulk a container
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

	err := ph.service.DeleteContainers(plantForm.ID, idsForm.IDs)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Containers deleted successfully"))
}

func (ph *Handler) MarkedFull(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	containerCode := c.Query("code")

	if containerCode == "" {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "No container code provided in query", nil)
		c.JSON(resp.Status, resp)
		return
	}

	containerDto, err := ph.service.GetContainerByCode(plantForm.ID, containerCode, false, false, false, false)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	} else if containerDto.StockLevel == models.Empty {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Container is empty, can not marked full", nil)
		c.JSON(resp.Status, resp)
		return
	}
	err = ph.service.MarkedContainerFull(plantForm.ID, containerDto)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Container marked full successfully"))
}

func (ph *Handler) GetContentsByCode(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	containerCode := c.Query("code")

	if containerCode == "" {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "No container code provided in query", nil)
		c.JSON(resp.Status, resp)
		return
	}

	containerDto, err := ph.service.GetContainerByCode(plantForm.ID, containerCode, true, true, false, false)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	if containerDto.Contents == nil {
		blankArray := make([]*container.ContentForm, 0)
		c.JSON(http.StatusOK, responses.NewSuccessDataResponse(http.StatusOK, "Container fetched successfully", blankArray))
	} else {
		c.JSON(http.StatusOK, responses.NewSuccessDataResponse(http.StatusOK, "Container fetched successfully", containerDto.Contents))
	}
}
