package service

import (
	"errors"
	"patungan-server/internal/dto/request"
	"patungan-server/internal/dto/response"
	"patungan-server/internal/models"
	"patungan-server/internal/repository"
	"patungan-server/internal/utils"

	"github.com/google/uuid"
)

type ParticipantService interface {
	AddParticipants(billID uuid.UUID, req *request.AddParticipantsRequest, creatorID uuid.UUID) error
	RemoveParticipant(billID, participantID, creatorID uuid.UUID) error
	UpdateParticipant(participantID uuid.UUID, req *request.UpdateParticipantRequest, userID uuid.UUID) (*response.ParticipantResponse, error)
	RecalculateSplits(billID uuid.UUID) error
	GetParticipants(billID uuid.UUID) ([]response.ParticipantResponse, error)
}

type participantService struct {
	participantRepo repository.ParticipantRepository
	billRepo        repository.BillRepository
	userRepo        repository.UserRepository
	itemRepo        repository.ItemRepository
	calculator      *utils.SplitCalculator
}

func NewParticipantService(
	participantRepo repository.ParticipantRepository,
	billRepo repository.BillRepository,
	userRepo repository.UserRepository,
	itemRepo repository.ItemRepository,
) ParticipantService {
	return &participantService{
		participantRepo: participantRepo,
		billRepo: billRepo,
		userRepo: userRepo,
		itemRepo: itemRepo,
		calculator: utils.NewSplitCalculator(),
	}
}

func (s *participantService) AddParticipants(billID uuid.UUID, req *request.AddParticipantsRequest, creatorID uuid.UUID) error {
	// Verify bill exists and user is creator
	bill, err := s.billRepo.FindByID(billID)
	if err != nil {
		return errors.New("bill not found")
	}

	if bill.CreatedBy != creatorID {
		return errors.New("only bill creator can add participants")
	}

	// Validate all users exist
	var participants []models.BillParticipant
	for _, p := range req.Participants {
		// Check if user exists
		user, err := s.userRepo.FindByID(p.UserID)
		if err != nil {
			return errors.New("user not found: " + p.UserID.String())
		}

		// Check if already participant
		exists, _ := s.participantRepo.CheckParticipant(billID, user.ID)
		if exists {
			continue // Skip if already participant
		}

		participant := models.BillParticipant{
			BillID: billID,
			UserID: user.ID,
			AmountOwed: p.AmountOwed, // Will be calculated later if 0
			AmountPaid: 0,
			IsSettled: false,
		}
		participants = append(participants, participant)
	}

	if len(participants) == 0 {
		return errors.New("no new participants to add")
	}

	// Create participants
	if err := s.participantRepo.CreateBatch(participants); err != nil {
		return err
	}

	// Recalculate splits
	return s.RecalculateSplits(billID)
}

func (s *participantService) RemoveParticipant(billID, participantID, creatorID uuid.UUID) error {
	// Verify bill exists and user is creator
	bill, err := s.billRepo.FindByID(billID)
	if err != nil {
		return errors.New("bill not found")
	}

	if bill.CreatedBy != creatorID {
		return errors.New("only bill creator can remove participants")
	}

	// Verify participant exists
	participant, err := s.participantRepo.FindByID(participantID)
	if err != nil {
		return errors.New("participant not found")
	}

	if participant.BillID != billID {
		return errors.New("participant does not belong to this bill")
	}

	// Don't allow removing participant with payments
	if participant.AmountPaid > 0 {
		return errors.New("cannot remove participant with payments")
	}

	// Delete participant
	if err := s.participantRepo.Delete(participantID); err != nil {
		return err
	}

	// Recalculate splits
	return s.RecalculateSplits(billID)
}

func (s *participantService) UpdateParticipant(participantID uuid.UUID, req *request.UpdateParticipantRequest, userID uuid.UUID) (*response.ParticipantResponse, error) {
	participant, err := s.participantRepo.FindByID(participantID)
	if err != nil {
		return nil, errors.New("participant not found")
	}

	// Update fields
	if req.AmountOwed > 0 {
		participant.AmountOwed = req.AmountOwed
	}
	if req.AmountPaid >= 0 {
		participant.AmountPaid = req.AmountPaid
	}
	if req.IsSettled != nil {
		participant.IsSettled = *req.IsSettled
	}

	if err := s.participantRepo.Update(participant); err != nil {
		return nil, err
	}

	return &response.ParticipantResponse{
		ID: participant.ID,
		UserID: participant.UserID,
		AmountOwed: participant.AmountOwed,
		AmountPaid: participant.AmountPaid,
		IsSettled: participant.IsSettled,
		User: response.UserResponse{
			ID: participant.User.ID,
			Email: participant.User.Email,
			FullName: participant.User.FullName,
		},
	}, nil
}

func (s *participantService) RecalculateSplits(billID uuid.UUID) error {
	bill, err := s.billRepo.FindByID(billID)
	if err != nil {
		return err
	}

	participants, err := s.participantRepo.FindByBillID(billID)
	if err != nil {
		return err
	}

	if len(participants) == 0 {
		return nil // No participants, nothing to recalculate
	}

	var amounts map[uuid.UUID]float64

	switch bill.SplitMethod {
	case "equal":
		amount, err := s.calculator.CalculateEqualSplit(bill.TotalAmount, len(participants))
		if err != nil {
			return err
		}
		amounts = make(map[uuid.UUID]float64)
		for _, p := range participants {
			amounts[p.UserID] = amount
		}

	case "items":
		items, err := s.itemRepo.FindByBillID(billID)
		if err != nil {
			return err
		}
		amounts, err = s.calculator.CalculateItemBasedSplit(items)
		if err != nil {
			return err
		}

	case "custom" :
		return nil

	default:
		return errors.New("invalid split method")
	}

	// Distribute any remainder due to rounding
	amounts = s.calculator.DistributeRemainder(bill.TotalAmount, amounts)

	// Update Participants
	for i := range participants {
		if amount, ok := amounts[participants[i].UserID]; ok {
			participants[i].AmountOwed = amount
			if err := s.participantRepo.Update(&participants[i]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *participantService) GetParticipants(billID uuid.UUID) ([]response.ParticipantResponse, error) {
	participants, err := s.participantRepo.FindByBillID(billID)
	if err != nil {
		return nil, err
	}

	var responses []response.ParticipantResponse
	for _, p := range participants {
		responses = append(responses, response.ParticipantResponse{
			ID: p.ID,
			UserID: p.UserID,
			AmountOwed: p.AmountOwed,
			AmountPaid: p.AmountPaid,
			IsSettled: p.IsSettled,
			User: response.UserResponse{
				ID:       p.User.ID,
				Email:    p.User.Email,
				FullName: p.User.FullName,
				Phone:    p.User.Phone,
			},
		})
	}

	return responses, nil
}