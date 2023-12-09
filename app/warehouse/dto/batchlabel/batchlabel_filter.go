package batchlabel

import (
	"star-wms/core/types"
)

type Filter struct {
	Query         string              `form:"query" db:"batch_no,po_category,so_number" whereType:"like" binding:"omitempty,max=100"`
	ID            uint                `form:"id" db:"id" binding:"omitempty,gt=0"`
	POCategory    string              `form:"po_category" db:"po_category" binding:"omitempty,max=100"`
	JoborderID    uint                `form:"joborder_id" db:"joborder_id" binding:"omitempty,gt=0"`
	CustomerID    uint                `form:"customer_id" db:"customer_id" binding:"omitempty,gt=0"`
	Status        types.Status        `form:"status" db:"status" binding:"omitempty,gt=0"`
	ProcessStatus types.ProcessStatus `form:"process_status" db:"process_status" validate:"omitempty,gt=0"`
}

type StickerFilter struct {
	Query        string `form:"query" db:"barcode,packet_no,shift,product_line,batch_no,machine_no" whereType:"like" binding:"omitempty,max=100"`
	ID           uint   `form:"id" db:"id" binding:"omitempty,gt=0"`
	Barcode      string `form:"barcode" db:"barcode" binding:"omitempty,max=100"`
	PacketNo     string `form:"packet_no" db:"packet_no" binding:"omitempty,max=100"`
	Shift        string `form:"shift" db:"shift" binding:"omitempty,max=100"`
	ProductLine  string `form:"product_line" db:"product_line" binding:"omitempty,max=100"`
	BatchNo      string `form:"batch_no" db:"batch_no" binding:"omitempty,max=100"`
	MachineNo    string `form:"machine_no" db:"machine_no" binding:"omitempty,max=100"`
	IsUsed       bool   `form:"is_used" db:"is_used" binding:"omitempty,max=100"`
	BatchLabelID *uint  `form:"batchlabel_id" db:"batchlabel_id" binding:"omitempty,gt=0"`
}
