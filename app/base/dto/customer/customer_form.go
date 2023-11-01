package customer

import (
	"star-wms/core/types"
)

type Form struct {
	PlantID          uint         `json:"plant_id" binding:"-"`
	ID               uint         `json:"id" binding:"-"`
	Name             string       `json:"name" validate:"required,min=4,max=100"`
	Code             string       `json:"code" validate:"required,min=4,max=100"`
	ContactPerson    string       `json:"contact_person" validate:"required,min=4,max=100"`
	BillingAddress1  string       `json:"billing_address_1" validate:"required,min=4,max=100"`
	BillingAddress2  string       `json:"billing_address_2" validate:"required,min=4,max=100"`
	BillingState     string       `json:"billing_state" validate:"required,min=4,max=100"`
	BillingCountry   string       `json:"billing_country" validate:"required,min=4,max=100"`
	BillingPincode   string       `json:"billing_pincode" validate:"required,min=4,max=100"`
	ShippingAddress1 string       `json:"shipping_address_1" validate:"required,min=4,max=100"`
	ShippingAddress2 string       `json:"shipping_address_2" validate:"required,min=4,max=100"`
	ShippingState    string       `json:"shipping_state" validate:"required,min=4,max=100"`
	ShippingCountry  string       `json:"shipping_country" validate:"required,min=4,max=100"`
	ShippingPincode  string       `json:"shipping_pincode" validate:"required,min=4,max=100"`
	Status           types.Status `json:"status" validate:"required,gt=0"`
}
