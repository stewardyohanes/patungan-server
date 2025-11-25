package repository

import (
	"patungan-server/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ItemParticipantRepository interface {
	Create(itemParticipant *models.ItemParticipant) error
	CreateBatch(itemParticipants []models.ItemParticipant) error
	FindByItemID(itemID uuid.UUID) ([]models.ItemParticipant, error)
	DeleteByItemID(itemID uuid.UUID) error
}

type itemParticipantRepository struct {
	db *gorm.DB
}

func NewItemParticipantRepository(db *gorm.DB) ItemParticipantRepository {
	return &itemParticipantRepository{db: db}
}

func (r *itemParticipantRepository) Create(itemParticipant *models.ItemParticipant) error {
	return r.db.Create(itemParticipant).Error
}

func (r *itemParticipantRepository) CreateBatch(itemParticipants []models.ItemParticipant) error {
	return r.db.Create(&itemParticipants).Error
}

func (r *itemParticipantRepository) FindByItemID(itemID uuid.UUID) ([]models.ItemParticipant, error) {
	var itemParticipants []models.ItemParticipant
	err := r.db.Preload("User").Where("item_id = ?", itemID).Find(&itemParticipants).Error
	return itemParticipants, err
}

func (r *itemParticipantRepository) DeleteByItemID(itemID uuid.UUID) error {
	return r.db.Delete(&models.ItemParticipant{}, "item_id = ?", itemID).Error
}