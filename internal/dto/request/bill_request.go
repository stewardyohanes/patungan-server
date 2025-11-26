package request

import "github.com/google/uuid"

type CreateBillRequest struct {
	Title        string      `json:"title" validate:"required"`
	Description  string      `json:"description"`
	TotalAmount  float64     `json:"total_amount" validate:"required,min=0"`
	Currency     string      `json:"currency" validate:"required,len=3"`
	SplitMethod  string      `json:"split_method" validate:"required,oneof=equal items custom"`
	Participants []uuid.UUID `json:"participants"`
}

type UpdateBillRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	TotalAmount float64 `json:"total_amount" validate:"omitempty,gt=0"`
	Status      string  `json:"status" validate:"omitempty,oneof=pending settled cancelled"`
}