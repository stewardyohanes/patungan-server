package response

import (
	"time"

	"github.com/google/uuid"
)

type BillResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	TotalAmount float64   `json:"total_amount"`
	Currency    string    `json:"currency"`
	SplitMethod string    `json:"split_method"`
	Status      string    `json:"status"`
	CreatedBy   uuid.UUID `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Creator     UserResponse `json:"creator"`
	Participants []ParticipantResponse `json:"participants,omitempty"`
}

type ParticipantResponse struct {
	ID uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	User UserResponse `json:"user"`
	AmountOwed float64 `json:"amount_owed"`
	AmountPaid float64 `json:"amount_paid"`
	IsSettled bool `json:"is_settled"`
}