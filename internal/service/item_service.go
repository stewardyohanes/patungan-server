package service

import (
	"errors"
	"patungan-server/internal/dto/request"
	"patungan-server/internal/dto/response"
	"patungan-server/internal/models"
	"patungan-server/internal/repository"

	"github.com/google/uuid"
)

type ItemService interface {
	AddItems(billID uuid.UUID, req *request.AddItemsRequest, creatorID uuid.UUID) error
	GetItems(billID uuid.UUID) ([]response.ItemResponse, error)
	UpdateItem(itemID uuid.UUID, req *request.UpdateItemRequest, userID uuid.UUID) (*response.ItemResponse, error)
	DeleteItem(itemID uuid.UUID, userID uuid.UUID) error
}

type itemService struct {
	itemRepo repository.ItemRepository
	billRepo repository.BillRepository
	itemParticipantRepo repository.ItemParticipantRepository
	participantService ParticipantService
}

func NewItemService(
	itemRepo repository.ItemRepository,
	billRepo repository.BillRepository,
	itemParticipantRepo repository.ItemParticipantRepository,
	participantService ParticipantService,
) ItemService {
	return &itemService{
		itemRepo: itemRepo,
		billRepo: billRepo,
		itemParticipantRepo: itemParticipantRepo,
		participantService: participantService,
	}
}

func (s *itemService) AddItems(billID uuid.UUID, req *request.AddItemsRequest, creatorID uuid.UUID) error {
	// Verifiy bill exists
	bill, err := s.billRepo.FindByID(billID)
	if err != nil {
		return errors.New("bill not found")
	}

	// Only allow for item-based split
	if bill.SplitMethod != "items" {
		return errors.New("bill is not configured for item-based split")
	}

	// Create items
	for _, itemReq := range req.Items {
		item := models.BillItem{
			BillID: billID,
			Name: itemReq.Name,
			Price: itemReq.Price,
			Quantity: itemReq.Quantity,
		}

		if err := s.itemRepo.Create(&item); err != nil {
			return err
		}

		// Add item participants
		var itemParticipants []models.ItemParticipant
		for _, participantID := range itemReq.Participants {
			itemParticipants = append(itemParticipants, models.ItemParticipant{
				ItemID: item.ID,
				UserID: participantID,
				SplitRatio: 1.0,
			})
		}

		if len(itemParticipants) > 0 {
			if err := s.itemParticipantRepo.CreateBatch(itemParticipants); err != nil {
				return err
			}
		}
	}

	// Recalculate bill splits
	return s.participantService.RecalculateSplits(billID)
}

func (s *itemService) GetItems(billID uuid.UUID) ([]response.ItemResponse, error) {
	items, err := s.itemRepo.FindByBillID(billID)
	if err != nil {
		return nil, err
	}

	var responses []response.ItemResponse
	for _, item := range items {
		itemResp := response.ItemResponse{
			ID: item.ID,
			BillID: item.BillID,
			Name: item.Name,
			Price: item.Price,
			Quantity: item.Quantity,
			CreatedAt: item.CreatedAt,
		}

		// Add participants
		for _, ip := range item.Participants {
			totalItemCost := item.Price * float64(item.Quantity)
			participantCount := len(item.Participants)
			amount := totalItemCost / float64(participantCount)

			itemResp.Participants = append(itemResp.Participants, response.ItemParticipantResponse{
				ID: ip.ID,
				UserID: ip.UserID,
				SplitRatio: ip.SplitRatio,
				Amount: amount,
				User: response.UserResponse{
					ID: ip.User.ID,
					Email: ip.User.Email,
					FullName: ip.User.FullName,
				},
			})
		}

		responses = append(responses, itemResp)
	}

	return responses, nil
}

func (s *itemService) UpdateItem(itemID uuid.UUID, req *request.UpdateItemRequest, userID uuid.UUID) (*response.ItemResponse, error) {
	item, err := s.itemRepo.FindByID(itemID)
	if err != nil {
		return nil, errors.New("item not found")
	}

	// Update fields
	if req.Name != "" {
		item.Name = req.Name
	}
	if req.Price > 0 {
		item.Price = req.Price
	}
	if req.Quantity > 0 {
		item.Quantity = req.Quantity
	}

	if err := s.itemRepo.Update(item); err != nil {
		return nil, err
	}

	// Recalculate splits
	if err := s.participantService.RecalculateSplits(item.BillID); err != nil {
		return nil, err
	}

	return &response.ItemResponse{
		ID: item.ID,
		BillID: item.BillID,
		Name: item.Name,
		Price: item.Price,
		Quantity: item.Quantity,
	}, nil
}

func (s *itemService) DeleteItem(itemID, userID uuid.UUID) error {
	item, err := s.itemRepo.FindByID(itemID)
	if err != nil {
		return errors.New("item not found")
	}

	billID := item.BillID

	// Delete item participants first
	if err := s.itemParticipantRepo.DeleteByItemID(itemID); err != nil {
		return err
	}

	// Delete item
	if err := s.itemRepo.Delete(itemID); err != nil {
		return err
	}

	// Recalculate splits
	return s.participantService.RecalculateSplits(billID)
}