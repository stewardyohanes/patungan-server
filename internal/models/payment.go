package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BillID uuid.UUID `gorm:"type:uuid;not null" json:"bill_id"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Amount float64 `gorm:"type:decimal(12,2);not null" json:"amount"`
	PaymentMethod string `gorm:"type:varchar(50)" json:"payment_method"`
	ProofURL string `gorm:"type:text" json:"proof_url"`
	Status string `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, verified, rejected
	CreatedAt time.Time `json:"created_at"`

	// Relations
	Bill Bill `gorm:"foreignKey:BillID" json:"bill,omitempty"`
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}