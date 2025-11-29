package response

import (
	"time"

	"github.com/google/uuid"
)

type ItemResponse struct {
	ID           uuid.UUID                 `json:"id"`
	BillID       uuid.UUID                 `json:"bill_id"`
	Name         string                    `json:"name"`
	Price        float64                   `json:"price"`
	Quantity     int                       `json:"quantity"`
	CreatedAt    time.Time                 `json:"created_at"`
	Participants []ItemParticipantResponse `json:"participants,omitempty"`
}

type ItemParticipantResponse struct {
	ID uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	User UserResponse `json:"user"`
	SplitRatio float64 `json:"split_ratio"`
	Amount float64 `json:"amount"` // Calculated amount this participant owes
}