package container

import (
	"star-wms/app/base/dto/product"
)

type ContentForm struct {
	PlantID   uint          `json:"plant_id" binding:"-"`
	ID        uint          `json:"id" binding:"-"`
	Product   *product.Form `json:"product" binding:"-"`
	Quantity  float64       `json:"quantity"`
	Container *Form         `json:"container"`
	Barcode   string        `json:"barcode"`
}
