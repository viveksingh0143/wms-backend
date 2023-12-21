package inventory

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"star-wms/app/admin/dto/plant"
	"star-wms/app/warehouse/dto/inventory"
	"star-wms/app/warehouse/service"
	"star-wms/core/auth"
	"star-wms/core/common/requests"
	"star-wms/core/common/responses"
	"star-wms/core/validation"
)

type Handler struct {
	service service.InventoryService
}

func NewInventoryHandler(s service.InventoryService) *Handler {
	return &Handler{
		service: s,
	}
}

// List inventories with filter, pagination, and sorting
func (ph *Handler) List(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	var filter inventory.Filter
	paginationValue, _ := c.Get("pagination")
	pagination, _ := paginationValue.(requests.Pagination)

	sortingValue, _ := c.Get("sorting")
	sorting, _ := sortingValue.(requests.Sorting)

	if err := c.ShouldBindQuery(&filter); err != nil {
		resp := responses.NewValidationErrorResponse(err, filter)
		c.JSON(resp.Status, resp)
		return
	}
	inventories, totalRecords, err := ph.service.GetAllInventorys(plantForm.ID, filter, pagination, sorting)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewPageResponse(inventories, totalRecords, pagination.Page, pagination.PageSize))
}

// Get a inventory
func (ph *Handler) Get(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	id, err := requests.StringToUInt(c.Param("id"))
	if err != nil || id <= 0 {
		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
		c.JSON(resp.Status, resp)
		return
	}

	inventoryDto, err := ph.service.GetInventoryByID(plantForm.ID, id)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusOK, responses.NewSuccessDataResponse(http.StatusOK, "Inventory fetched successfully", inventoryDto))
}

// RawMaterialStockin Create a new inventory for raw material
func (ph *Handler) RawMaterialStockin(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	var rawMaterialStockInForm inventory.RawMaterialStockInForm
	if err := c.ShouldBindJSON(&rawMaterialStockInForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, rawMaterialStockInForm)
		c.JSON(resp.Status, resp)
		return
	}

	validate := validation.GetValidator()

	if err := validate.Struct(rawMaterialStockInForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, rawMaterialStockInForm)
		c.JSON(resp.Status, resp)
		return
	}

	err := ph.service.CreateRawMaterialStockIn(plantForm.ID, &rawMaterialStockInForm)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusCreated, responses.NewSuccessResponse(http.StatusCreated, "Raw material stockin successfully"))
}

func (ph *Handler) FinishedGoodsStockin(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	var finishedGoodsStockInForm inventory.FinishedGoodsStockInForm
	if err := c.ShouldBindJSON(&finishedGoodsStockInForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, finishedGoodsStockInForm)
		c.JSON(resp.Status, resp)
		return
	}

	validate := validation.GetValidator()

	if err := validate.Struct(finishedGoodsStockInForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, finishedGoodsStockInForm)
		c.JSON(resp.Status, resp)
		return
	}

	err := ph.service.CreateFinishedGoodsStockIn(plantForm.ID, &finishedGoodsStockInForm)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusCreated, responses.NewSuccessResponse(http.StatusCreated, "Goods stockin successfully"))
}

func (ph *Handler) FinishedGoodStockin(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	var finishedGoodStockInForm inventory.FinishedGoodStockInForm
	if err := c.ShouldBindJSON(&finishedGoodStockInForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, finishedGoodStockInForm)
		c.JSON(resp.Status, resp)
		return
	}

	validate := validation.GetValidator()

	if err := validate.Struct(finishedGoodStockInForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, finishedGoodStockInForm)
		c.JSON(resp.Status, resp)
		return
	}

	stickerForm, err := ph.service.CreateFinishedGoodStockIn(plantForm.ID, &finishedGoodStockInForm)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusCreated, responses.NewSuccessDataResponse(http.StatusCreated, "Goods stockin successfully", stickerForm))
}

func (ph *Handler) AttachContainer(c *gin.Context) {
	plantValue, _ := c.Get(auth.AuthPlantKey)
	plantForm, _ := plantValue.(plant.Form)

	var attachContainerForm inventory.AttachContainerForm
	if err := c.ShouldBindJSON(&attachContainerForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, attachContainerForm)
		c.JSON(resp.Status, resp)
		return
	}

	validate := validation.GetValidator()

	if err := validate.Struct(attachContainerForm); err != nil {
		resp := responses.NewValidationErrorResponse(err, attachContainerForm)
		c.JSON(resp.Status, resp)
		return
	}

	err := ph.service.AttachContainerToLocation(plantForm.ID, &attachContainerForm)
	if err != nil {
		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
		c.JSON(resp.Status, resp)
		return
	}
	c.JSON(http.StatusCreated, responses.NewSuccessResponse(http.StatusCreated, "Container attached successfully"))
}

//// Create a new inventory
//func (ph *Handler) Create(c *gin.Context) {
//	plantValue, _ := c.Get(auth.AuthPlantKey)
//	plantForm, _ := plantValue.(plant.Form)
//
//	var inventoryForm inventory.Form
//	if err := c.ShouldBindJSON(&inventoryForm); err != nil {
//		resp := responses.NewValidationErrorResponse(err, inventoryForm)
//		c.JSON(resp.Status, resp)
//		return
//	}
//
//	validate := validation.GetValidator()
//
//	if err := validate.Struct(inventoryForm); err != nil {
//		resp := responses.NewValidationErrorResponse(err, inventoryForm)
//		c.JSON(resp.Status, resp)
//		return
//	}
//
//	err := ph.service.CreateInventory(plantForm.ID, &inventoryForm)
//	if err != nil {
//		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
//		c.JSON(resp.Status, resp)
//		return
//	}
//	c.JSON(http.StatusCreated, responses.NewSuccessResponse(http.StatusCreated, "Inventory created successfully"))
//}
//
//// Update a inventory
//func (ph *Handler) Update(c *gin.Context) {
//	plantValue, _ := c.Get(auth.AuthPlantKey)
//	plantForm, _ := plantValue.(plant.Form)
//
//	id, err := requests.StringToUInt(c.Param("id"))
//	if err != nil || id <= 0 {
//		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
//		c.JSON(resp.Status, resp)
//		return
//	}
//
//	var inventoryForm inventory.Form
//	if err := c.BindJSON(&inventoryForm); err != nil {
//		resp := responses.NewValidationErrorResponse(err, inventoryForm)
//		c.JSON(resp.Status, resp)
//		return
//	}
//
//	validate := validation.GetValidator()
//	if err := validate.Struct(inventoryForm); err != nil {
//		resp := responses.NewValidationErrorResponse(err, inventoryForm)
//		c.JSON(resp.Status, resp)
//		return
//	}
//
//	err = ph.service.UpdateInventory(plantForm.ID, id, &inventoryForm)
//	if err != nil {
//		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
//		c.JSON(resp.Status, resp)
//		return
//	}
//	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Inventory updated successfully"))
//}
//
//// Delete a inventory
//func (ph *Handler) Delete(c *gin.Context) {
//	plantValue, _ := c.Get(auth.AuthPlantKey)
//	plantForm, _ := plantValue.(plant.Form)
//
//	idStr := c.Param("id") // assuming the ID is passed as a URL parameter
//	idInt, err := strconv.Atoi(idStr)
//	if err != nil || idInt < 0 {
//		resp := responses.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err)
//		c.JSON(resp.Status, resp)
//		return
//	}
//
//	id := uint(idInt)
//	err = ph.service.DeleteInventory(plantForm.ID, id)
//	if err != nil {
//		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
//		c.JSON(resp.Status, resp)
//		return
//	}
//	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Inventory deleted successfully"))
//}
//
//// DeleteBulk a inventory
//func (ph *Handler) DeleteBulk(c *gin.Context) {
//	plantValue, _ := c.Get(auth.AuthPlantKey)
//	plantForm, _ := plantValue.(plant.Form)
//
//	var idsForm requests.RequestIDs
//	if err := c.ShouldBindJSON(&idsForm); err != nil {
//		resp := responses.NewValidationErrorResponse(err, idsForm)
//		c.JSON(resp.Status, resp)
//		return
//	}
//
//	validate := validation.GetValidator()
//	if err := validate.Struct(idsForm); err != nil {
//		resp := responses.NewValidationErrorResponse(err, idsForm)
//		c.JSON(resp.Status, resp)
//		return
//	}
//
//	err := ph.service.DeleteInventorys(plantForm.ID, idsForm.IDs)
//	if err != nil {
//		resp := responses.NewErrorResponse(http.StatusInternalServerError, "Something went wrong at server", err)
//		c.JSON(resp.Status, resp)
//		return
//	}
//	c.JSON(http.StatusOK, responses.NewSuccessResponse(http.StatusOK, "Inventorys deleted successfully"))
//}
//
