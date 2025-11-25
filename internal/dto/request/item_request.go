package request

import "github.com/google/uuid"

type CreateItemRequest struct {
	Name         string      `json:"name" validate:"required"`
	Price        float64     `json:"price" validate:"required,gt=0"`
	Quantity     int         `json:"quantity" validate:"required,min=1"`
	Participants []uuid.UUID `json:"participants" validate:"required,min=1"` // Users sharing this item
}

type AddItemsRequest struct {
	Items []CreateItemRequest `json:"items" validate:"required,min=1,dive"`
}

type UpdateItemRequest struct {
	Name string `json:"name,omitempty"`
	Price float64 `json:"price,omitempty" validate:"omitempty,gt=0"`
	Quantity int `json:"quantity,omitempty" validate:"omitempty,min=1"`
}