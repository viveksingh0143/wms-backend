package inventory

import (
	"star-wms/core/types"
)

type Filter struct {
	Query     string `form:"query" db:"batch_no,po_category,so_number" whereType:"like" binding:"omitempty,max=100"`
	ID        uint   `form:"id" db:"id" binding:"omitempty,gt=0"`
	StoreID   uint   `form:"store_id" db:"store_id" binding:"omitempty,gt=0"`
	ProductID uint   `form:"product_id" db:"product_id" binding:"omitempty,gt=0"`
}

type InventoryTransactionFilter struct {
	Query            string                `form:"query" db:"barcode,packet_no,shift,product_line,batch_no,machine_no" whereType:"like" binding:"omitempty,max=100"`
	ID               uint                  `form:"id" db:"id" binding:"omitempty,gt=0"`
	Barcode          string                `form:"barcode" binding:"omitempty,max=100"`
	StoreID          uint                  `form:"Store_id" db:"Store_id" binding:"omitempty,gt=0"`
	ProductID        uint                  `form:"Product_id" db:"Product_id" binding:"omitempty,gt=0"`
	JoborderID       uint                  `form:"Joborder_id" db:"Joborder_id" binding:"omitempty,gt=0"`
	ContainerID      uint                  `form:"Container_id" db:"Container_id" binding:"omitempty,gt=0"`
	BatchlabelID     uint                  `form:"Batchlabel_id" db:"Batchlabel_id" binding:"omitempty,gt=0"`
	BarcodeStickerID uint                  `form:"BarcodeSticker_id" db:"BarcodeSticker_id" binding:"omitempty,gt=0"`
	Status           types.InventoryStatus `form:"status" db:"status" binding:"omitempty,gt=0"`
}
