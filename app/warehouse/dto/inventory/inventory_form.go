package inventory

import (
	"star-wms/app/base/dto/container"
	"star-wms/app/base/dto/joborder"
	"star-wms/app/base/dto/product"
	"star-wms/app/base/dto/store"
	"star-wms/app/warehouse/dto/batchlabel"
	"star-wms/core/types"
)

type Form struct {
	ID       uint          `json:"id" binding:"-"`
	Store    *store.Form   `json:"store" validationTag:"store.id" validate:"required,validRelationID,structonly"`
	Product  *product.Form `json:"product" validationTag:"product.id" validate:"required,validRelationID,structonly"`
	Quantity float64       `json:"quantity" validate:"required,lte=10000"`
}

type InventoryTransactionForm struct {
	ID             uint                           `json:"id" binding:"-"`
	Store          *store.Form                    `json:"store" validationTag:"store.id" validate:"required,validRelationID,structonly"`
	Product        *product.Form                  `json:"product" validationTag:"product.id" validate:"required,validRelationID,structonly"`
	Quantity       string                         `json:"quantity" validate:"required,min=1,max=100"`
	Joborder       *joborder.Form                 `json:"joborder" validationTag:"joborder.id" validate:"omitempty,validRelationID,structonly"`
	Container      *container.Form                `json:"container" validationTag:"container.id" validate:"omitempty,validRelationID,structonly"`
	Batchlabel     *batchlabel.Form               `json:"batchlabel" validationTag:"batchlabel.id" validate:"omitempty,validRelationID,structonly"`
	BarcodeSticker *batchlabel.BarcodeStickerForm `json:"barcode_sticker" validationTag:"barcode_sticker.id" validate:"omitempty,validRelationID,structonly"`
	Status         types.InventoryStatus          `json:"status" validate:"required,gt=0"`
}
