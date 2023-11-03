package joborder

import (
	"star-wms/app/base/dto/customer"
	"star-wms/app/base/dto/product"
	"star-wms/core/types"
	"time"
)

type Form struct {
	ID            uint                `json:"id" binding:"-"`
	IssuedDate    time.Time           `json:"issued_date" validate:"required"`
	OrderNo       string              `json:"order_no" validate:"required,min=4,max=100"`
	POCategory    string              `json:"po_category" validate:"required,oneof='PRODUCTION' 'TRAILS' 'NPD' 'SAMPLES'"`
	Status        types.Status        `json:"status" validate:"required,gt=0"`
	ProcessStatus types.ProcessStatus `json:"process_status" form:"default=1" validate:"omitempty,gt=0"`
	Customer      *customer.Form      `json:"customer" validationTag:"customer.id" validate:"omitempty,validRelationID,structonly"`
	Items         []*ItemsForm        `json:"items"`
}

type ItemsForm struct {
	ID        uint          `json:"id" binding:"-"`
	ProductID uint          `json:"product_id" validate:"required"`
	Product   *product.Form `json:"product"`
	Quantity  float64       `json:"quantity"`
}
