package requisition

import (
	"star-wms/app/base/dto/product"
	"star-wms/app/base/dto/store"
	"star-wms/core/types"
	"time"
)

type Form struct {
	ID         uint         `json:"id" binding:"-"`
	IssuedDate time.Time    `json:"issued_date" validate:"required"`
	OrderNo    string       `json:"order_no" validate:"required,min=4,max=100"`
	Department string       `json:"department" validate:"required"`
	Status     types.Status `json:"status" validate:"required,gt=0"`
	Approved   bool         `json:"approved" binding:"-"`
	Store      *store.Form  `json:"store" validationTag:"store.id" validate:"omitempty,validRelationID,structonly"`
	Items      []*ItemsForm `json:"items"`
}

type ItemsForm struct {
	ID        uint          `json:"id" binding:"-"`
	ProductID uint          `json:"product_id" validate:"required"`
	Product   *product.Form `json:"product"`
	Quantity  float64       `json:"quantity"`
}
