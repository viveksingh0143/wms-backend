package models

import (
	adminModels "star-wms/app/admin/models"
	"star-wms/core/common/models"
	"star-wms/core/types"
)

type Customer struct {
	models.MyModel
	Name             string            `gorm:"type:varchar(255);uniqueIndex;not null;column:name"`
	Code             string            `gorm:"type:varchar(255);uniqueIndex;not null;column:code"`
	ContactPerson    string            `gorm:"type:varchar(255)"`
	BillingAddress1  string            `gorm:"type:varchar(255)"`
	BillingAddress2  string            `gorm:"type:varchar(255)"`
	BillingState     string            `gorm:"type:varchar(255)"`
	BillingCountry   string            `gorm:"type:varchar(255)"`
	BillingPincode   string            `gorm:"type:varchar(255)"`
	ShippingAddress1 string            `gorm:"type:varchar(255)"`
	ShippingAddress2 string            `gorm:"type:varchar(255)"`
	ShippingState    string            `gorm:"type:varchar(255)"`
	ShippingCountry  string            `gorm:"type:varchar(255)"`
	ShippingPincode  string            `gorm:"type:varchar(255)"`
	Status           types.Status      `gorm:"type:int;default:1"`
	PlantID          uint              `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Plant            adminModels.Plant `gorm:"foreignKey:PlantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
