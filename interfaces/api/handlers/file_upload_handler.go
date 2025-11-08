package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/services"
)

type FileUploadHandler struct {
	fileUploadService services.FileUploadService
}

func NewFileUploadHandler(fileUploadService services.FileUploadService) *FileUploadHandler {
	return &FileUploadHandler{
		fileUploadService: fileUploadService,
	}
}

// UploadFile handles file upload endpoint
// POST /upload/file
func (h *FileUploadHandler) UploadFile(c *fiber.Ctx) error {
	// Get authenticated user ID
	userID := c.Locals("userID").(uuid.UUID)

	// Get uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	// Open file
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open file",
		})
	}
	defer file.Close()

	// Upload file using service
	media, err := h.fileUploadService.UploadFile(c.Context(), userID, file, fileHeader)
	if err != nil {
		// Check for validation errors
		if err.Error() == "file size exceeds 50MB limit" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err.Error()[:19] == "file type not allowed" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload file",
			"detail": err.Error(),
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":        media.ID,
		"filename":  media.FileName,
		"extension": media.Extension,
		"size":      media.Size,
		"mime_type": media.MimeType,
		"url":       media.URL,
		"uploaded_at": media.CreatedAt,
	})
}
