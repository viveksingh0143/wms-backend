package inventory

import (
	"star-wms/app/base/dto/container"
	"star-wms/app/base/dto/joborder"
	"star-wms/app/base/dto/product"
	"star-wms/app/base/dto/store"
	"star-wms/app/warehouse/dto/batchlabel"
	"star-wms/core/types"
)

type RawMaterialStockInForm struct {
	Store      *store.Form    `json:"store" validationTag:"store.id" validate:"required,validRelationID,structonly"`
	Product    *product.Form  `json:"product" validationTag:"product.id" validate:"required,validRelationID,structonly"`
	Quantity   float64        `json:"quantity" validate:"required,lte=10000"`
	Container  container.Form `json:"container" validate:"required,structonly"`
	Shift      string         `json:"shift" validate:"required,min=1,max=100"`
	Supervisor string         `json:"supervisor" validate:"required,min=1,max=100"`
	BatchNo    string         `json:"batch_no" validate:"required,min=1,max=100"`
}

type FinishedGoodsStockInForm struct {
	Store         *store.Form `json:"store" validationTag:"store.id" validate:"required,validRelationID,structonly"`
	ContainerCode string      `json:"container_code" validate:"required,min=1,max=100"`
	Barcodes      []string    `json:"barcodes" validate:"required"`
}

type FinishedGoodStockInForm struct {
	ContainerCode string `json:"container_code" validate:"required,min=1,max=100"`
	Barcode       string `json:"barcode" validate:"required"`
}

type AttachContainerForm struct {
	ContainerCode string `json:"container_code" validate:"required"`
	LocationCode  string `json:"location_code" validate:"required"`
}

type Form struct {
	ID       uint          `json:"id" binding:"-"`
	Store    *store.Form   `json:"store" validationTag:"store.id" validate:"required,validRelationID,structonly"`
	Product  *product.Form `json:"product" validationTag:"product.id" validate:"required,validRelationID,structonly"`
	Quantity float64       `json:"quantity" validate:"required,lte=10000"`
}

type InventoryTransactionForm struct {
	ID         uint                    `json:"id" binding:"-"`
	Store      *store.Form             `json:"store" validationTag:"store.id" validate:"required,validRelationID,structonly"`
	Product    *product.Form           `json:"product" validationTag:"product.id" validate:"required,validRelationID,structonly"`
	Quantity   string                  `json:"quantity" validate:"required,min=1,max=100"`
	Joborder   *joborder.Form          `json:"joborder" validationTag:"joborder.id" validate:"omitempty,validRelationID,structonly"`
	Container  *container.Form         `json:"container" validationTag:"container.id" validate:"omitempty,validRelationID,structonly"`
	Batchlabel *batchlabel.Form        `json:"batchlabel" validationTag:"batchlabel.id" validate:"omitempty,validRelationID,structonly"`
	Sticker    *batchlabel.StickerForm `json:"barcode_sticker" validationTag:"barcode_sticker.id" validate:"omitempty,validRelationID,structonly"`
	Status     types.InventoryStatus   `json:"status" validate:"required,gt=0"`
}
