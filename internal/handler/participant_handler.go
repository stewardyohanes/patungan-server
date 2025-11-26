package handler

import (
	"patungan-server/internal/dto/request"
	"patungan-server/internal/service"
	"patungan-server/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ParticipantHandler struct {
	participantService service.ParticipantService
	validator *validator.Validate
}

func NewParticipantHandler(participantService service.ParticipantService) *ParticipantHandler {
	return &ParticipantHandler{
		participantService: participantService,
		validator: validator.New(),
	}
}

func (h *ParticipantHandler) AddParticipants(c fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	billID, err := uuid.Parse(c.Params("billId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid bill id", nil)
	}

	var req request.AddParticipantsRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request", err.Error())
	}

	if err := h.validator.Struct(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation failed", err.Error())
	}

	if err := h.participantService.AddParticipants(billID, &req, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "participants added successfully", nil)
}

func (h *ParticipantHandler) RemoveParticipant(c fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	billID, err := uuid.Parse(c.Params("billId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid bill id", nil)
	}

	participantID, err := uuid.Parse(c.Params("participantId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid participant id", nil)
	}

	if err := h.participantService.RemoveParticipant(billID, participantID, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "participant removed successfully", nil)
}

func (h *ParticipantHandler) UpdateParticipant(c fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	participantID, err := uuid.Parse(c.Params("participantId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid participant id", nil)
	}

	var req request.UpdateParticipantRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request", err.Error())
	}

	result, err := h.participantService.UpdateParticipant(participantID, &req, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "participant updated successfully", result)
}

func (h *ParticipantHandler) GetParticipants(c fiber.Ctx) error {
	billID, err := uuid.Parse(c.Params("billId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid bill id", nil)
	}

	result, err := h.participantService.GetParticipants(billID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "participants retrieved successfully", result)
}

func (h *ParticipantHandler) RecalculateSplits(c fiber.Ctx) error {
	billID, err := uuid.Parse(c.Params("billId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid bill id", nil)
	}

	if err := h.participantService.RecalculateSplits(billID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "splits recalculated successfully", nil)
}