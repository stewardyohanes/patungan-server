package request

import "github.com/google/uuid"

type AddParticipantRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	AmountOwed float64 `json:"amount_owed,omitempty"` // For custom split
}

type AddParticipantsRequest struct {
	Participants []AddParticipantRequest `json:"participants" validate:"required,min=1,dive"`
}

type UpdateParticipantRequest struct {
	AmountOwed float64 `json:"amount_owed" validate:"required,gt=0"`
	AmountPaid float64 `json:"amount_paid,omitempty"`
	IsSettled *bool `json:"is_settled,omitempty"`
}