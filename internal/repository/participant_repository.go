package repository

import (
	"patungan-server/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ParticipantRepository interface {
	Create(participant *models.BillParticipant) error
	CreateBatch(participants []models.BillParticipant) error
	FindByBillID(billID uuid.UUID) ([]models.BillParticipant, error)
	FindByID(id uuid.UUID) (*models.BillParticipant, error)
	Update(participant *models.BillParticipant) error
	Delete(id uuid.UUID) error
	DeleteByBillAndUser(billID, userID uuid.UUID) error
	CheckParticipant(billID, userID uuid.UUID) (bool, error)
}

type participantRepository struct {
	db *gorm.DB
}

func NewParticipantRepository(db *gorm.DB) ParticipantRepository {
	return &participantRepository{db: db}
}

func (r *participantRepository) Create(participant *models.BillParticipant) error {
	return r.db.Create(participant).Error
}

func (r *participantRepository) CreateBatch(participants []models.BillParticipant) error {
	return r.db.Create(&participants).Error
}

func (r *participantRepository) FindByBillID(billID uuid.UUID) ([]models.BillParticipant, error) {
	var participants []models.BillParticipant
	err := r.db.Preload("User").Where("bill_id = ?", billID).Find(&participants).Error
	return participants, err
}

func (r *participantRepository) FindByID(id uuid.UUID) (*models.BillParticipant, error) {
	var participant models.BillParticipant
	err := r.db.Preload("User").First(&participant, "id = ?", id).Error
	return &participant, err
}

func (r *participantRepository) Update(participant *models.BillParticipant) error {
	return r.db.Save(participant).Error
}

func (r *participantRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.BillParticipant{}, "id = ?", id).Error
}

func (r *participantRepository) DeleteByBillAndUser(billID, userID uuid.UUID) error {
	return r.db.Delete(&models.BillParticipant{}, "bill_id = ? AND user_id = ?", billID, userID).Error
}

func (r *participantRepository) CheckParticipant(billID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.BillParticipant{}).
		Where("bill_id = ? AND user_id = ?", billID, userID).
		Count(&count).Error
	return count > 0, err
}