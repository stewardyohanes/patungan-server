package handler

import (
	"patungan-server/internal/dto/request"
	"patungan-server/internal/service"
	"patungan-server/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	authService service.AuthService
	validator *validator.Validate
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator: validator.New(),
	}
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req request.RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body", err. Error())
	}

	if err := h.validator.Struct(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation error", err.Error())
	}
	result, err := h.authService.Register(&req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "registration successful", result)
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req request.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if err := h.validator.Struct(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation error", err.Error())
	}

	result, err := h.authService.Login(&req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "login successful", result)
}