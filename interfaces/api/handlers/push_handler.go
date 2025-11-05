package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type PushHandler struct {
	pushService services.PushService
}

func NewPushHandler(pushService services.PushService) *PushHandler {
	return &PushHandler{
		pushService: pushService,
	}
}

// Subscribe handles push notification subscription
func (h *PushHandler) Subscribe(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	var req dto.PushSubscriptionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	subscription, err := h.pushService.Subscribe(c.Context(), userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to save subscription", err)
	}

	return utils.SuccessResponse(c, "Subscription saved successfully", subscription)
}

// Unsubscribe handles push notification unsubscription
func (h *PushHandler) Unsubscribe(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	var req dto.PushSubscriptionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	err := h.pushService.Unsubscribe(c.Context(), userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to remove subscription", err)
	}

	return utils.SuccessResponse(c, "Subscription removed successfully", nil)
}

// GetPublicKey returns the VAPID public key for frontend
func (h *PushHandler) GetPublicKey(c *fiber.Ctx) error {
	publicKey := h.pushService.GetPublicKey()

	return c.JSON(fiber.Map{
		"publicKey": publicKey,
	})
}
