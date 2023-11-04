package dto

import (
	"time"
)

type BaseDto struct {
	ID        uint      `json:"id" binding:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
}
