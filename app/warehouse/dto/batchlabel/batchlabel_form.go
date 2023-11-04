package batchlabel

import (
	"star-wms/app/base/dto/customer"
	"star-wms/app/base/dto/joborder"
	"star-wms/app/base/dto/machine"
	"star-wms/app/base/dto/product"
	"star-wms/core/types"
	"time"
)

type Form struct {
	ID              uint                  `json:"id" binding:"-"`
	BatchDate       time.Time             `json:"batch_date" validate:"required"`
	BatchNo         string                `json:"batch_no" validate:"required,min=4,max=100"`
	SoNumber        string                `json:"so_number" validate:"omitempty,min=4,max=100"`
	POCategory      string                `json:"po_category" validate:"required,oneof='PRODUCTION' 'TRAILS' 'NPD' 'SAMPLES'"`
	UnitType        string                `json:"unit_type" validate:"required,oneof='WEIGHT' 'PIECE' 'LIQUID'"`
	UnitWeight      float64               `json:"unit_weight" validate:"omitempty,lte=10000"`
	UnitValue       string                `json:"unit_weight_type" validate:"required,oneof='Kilogram' 'Gram' 'Liter' 'Milliliter' 'Piece'"`
	TargetQuantity  float64               `json:"target_quantity" validate:"required,min=0"`
	PackageQuantity float64               `json:"package_quantity" validate:"required,min=0"`
	Status          types.Status          `json:"status" validate:"required,gt=0"`
	ProcessStatus   types.ProcessStatus   `json:"process_status" form:"default=1" validate:"omitempty,gt=0"`
	Joborder        *joborder.Form        `json:"joborder" validationTag:"joborder.id" validate:"omitempty,validRelationID,structonly"`
	JoborderItem    *joborder.ItemsForm   `json:"joborder_item" validationTag:"joborder_item.id" validate:"omitempty,validRelationID,structonly"`
	Customer        *customer.Form        `json:"customer" validationTag:"customer.id" validate:"required,validRelationID,structonly"`
	Product         *product.Form         `json:"product" validationTag:"product.id" validate:"required,validRelationID,structonly"`
	Machine         *machine.Form         `json:"machine" validationTag:"machine.id" validate:"required,validRelationID,structonly"`
	Stickers        []*BarcodeStickerForm `json:"stickers"`
}

type BarcodeStickerForm struct {
	ID          uint   `json:"id" binding:"-"`
	Barcode     string `json:"barcode" validate:"required,min=4,max=100"`
	PacketNo    string `json:"packet_no" validate:"required,min=1,max=100"`
	PrintCount  int32  `json:"print_count" binding:"-"`
	Shift       string `json:"shift" validate:"required,min=1,max=100"`
	ProductLine string `json:"product_line" validate:"required,min=1,max=100"`
	BatchNo     string `json:"batch_no" validate:"required,min=1,max=100"`
	MachineNo   string `json:"machine_no" validate:"required,min=1,max=100"`
	Quantity    string `json:"quantity" validate:"required,min=1,max=100"`
	IsUsed      bool   `json:"is_used"`
	BatchLabel  Form
}
