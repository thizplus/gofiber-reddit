package handlers

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	chatWebsocket "gofiber-template/infrastructure/websocket"
	"gofiber-template/infrastructure/storage"
	"gofiber-template/pkg/utils"
)

type MessageHandler struct {
	messageService     services.MessageService
	mediaUploadService *storage.MediaUploadService
	chatHub            *chatWebsocket.ChatHub
}

func NewMessageHandler(messageService services.MessageService, mediaUploadService *storage.MediaUploadService, chatHub *chatWebsocket.ChatHub) *MessageHandler {
	return &MessageHandler{
		messageService:     messageService,
		mediaUploadService: mediaUploadService,
		chatHub:            chatHub,
	}
}

// SendMessage sends a new message in a conversation
// POST /conversations/:conversationId/messages
func (h *MessageHandler) SendMessage(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	// Get conversationId from URL path parameter
	conversationID, err := uuid.Parse(c.Params("conversationId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid conversation ID")
	}

	// Check Content-Type to determine handler
	contentType := string(c.Request().Header.ContentType())

	if strings.Contains(contentType, "multipart/form-data") {
		// Handle media upload
		return h.sendMediaMessage(c, userID, conversationID)
	} else {
		// Handle text message (existing logic)
		return h.sendTextMessage(c, userID, conversationID)
	}
}

// sendTextMessage handles JSON text messages
func (h *MessageHandler) sendTextMessage(c *fiber.Ctx, userID uuid.UUID, conversationID uuid.UUID) error {
	var req dto.SendMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	// Set conversationId from path parameter
	req.ConversationID = conversationID

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	message, err := h.messageService.SendMessage(c.Context(), userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to send message", err)
	}

	// Send WebSocket notification to receiver
	h.sendWebSocketNotification(message)

	return utils.SuccessResponse(c, "Message sent successfully", message)
}

// sendMediaMessage handles multipart/form-data media messages
func (h *MessageHandler) sendMediaMessage(c *fiber.Ctx, userID uuid.UUID, conversationID uuid.UUID) error {
	// 1. Parse form data
	messageType := c.FormValue("type") // "image", "video", "file"
	content := c.FormValue("content")  // Optional caption

	// 2. Get uploaded files
	form, err := c.MultipartForm()
	if err != nil {
		return utils.ValidationErrorResponse(c, "Failed to parse form data")
	}

	files := form.File["media[]"]

	// If no files, treat as text-only message
	if len(files) == 0 {
		// Must have content for text message
		if content == "" {
			return utils.ValidationErrorResponse(c, "Either files or text content is required")
		}

		// Send as text message
		req := &dto.SendMessageRequest{
			ConversationID: conversationID,
			Type:           "text",
			Content:        &content,
		}

		message, err := h.messageService.SendMessage(c.Context(), userID, req)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to send message", err)
		}

		return utils.SuccessResponse(c, "Message sent successfully", message)
	}

	// 3. Validate message type (for media messages)
	validTypes := map[string]bool{"image": true, "video": true, "file": true}
	if messageType == "" || !validTypes[messageType] {
		return utils.ValidationErrorResponse(c, "Invalid message type. Must be one of: image, video, file")
	}

	// 4. Validate file count
	maxFiles := map[string]int{"image": 10, "video": 1, "file": 5}
	if len(files) > maxFiles[messageType] {
		return utils.ValidationErrorResponse(c, fmt.Sprintf("Maximum %d files allowed for type %s", maxFiles[messageType], messageType))
	}

	// 5. Process each file
	mediaItems := make([]dto.MessageMedia, 0, len(files))

	for _, fileHeader := range files {
		// Open file
		file, err := fileHeader.Open()
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to open file", err)
		}
		defer file.Close()

		// Validate file
		if err := h.validateFile(file, fileHeader, messageType); err != nil {
			return utils.ValidationErrorResponse(c, err.Error())
		}

		// Reset file pointer after validation
		file.Seek(0, 0)

		// Upload to Bunny Storage
		mediaItem, err := h.uploadFile(c.Context(), file, fileHeader, messageType)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to upload file", err)
		}

		mediaItems = append(mediaItems, mediaItem)
	}

	// 6. Create message via service
	var contentPtr *string
	if content != "" {
		contentPtr = &content
	}

	req := &dto.SendMessageRequest{
		ConversationID: conversationID,
		Type:           messageType,
		Content:        contentPtr,
		Media:          mediaItems,
	}

	message, err := h.messageService.SendMessage(c.Context(), userID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to send message", err)
	}

	// Send WebSocket notification to receiver
	h.sendWebSocketNotification(message)

	return utils.SuccessResponse(c, "Message sent successfully", message)
}

// ListMessages retrieves messages in a conversation
// GET /conversations/:conversationId/messages?cursor=xxx&limit=50
func (h *MessageHandler) ListMessages(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	// Get conversationID from URL path parameter
	conversationID, err := uuid.Parse(c.Params("conversationId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid conversation ID")
	}

	// Get cursor from query params
	cursor := c.Query("cursor")
	var cursorPtr *string
	if cursor != "" {
		cursorPtr = &cursor
	}

	// Get limit from query params (default: 50, max: 100)
	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			if parsedLimit > 0 && parsedLimit <= 100 {
				limit = parsedLimit
			}
		}
	}

	messages, err := h.messageService.ListMessages(c.Context(), conversationID, userID, cursorPtr, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve messages", err)
	}

	return utils.SuccessResponse(c, "Messages retrieved successfully", messages)
}

// GetMessage retrieves a single message by ID
// GET /messages/:id
func (h *MessageHandler) GetMessage(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	messageID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid message ID")
	}

	message, err := h.messageService.GetMessage(c.Context(), messageID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve message", err)
	}

	return utils.SuccessResponse(c, "Message retrieved successfully", message)
}

// GetMessageContext retrieves a message with surrounding context
// GET /messages/:id/context
func (h *MessageHandler) GetMessageContext(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	messageID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid message ID")
	}

	context, err := h.messageService.GetMessageContext(c.Context(), messageID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve message context", err)
	}

	return utils.SuccessResponse(c, "Message context retrieved successfully", context)
}

// Note: MarkAsRead is handled by ConversationHandler
// POST /conversations/:conversationId/read

// GetConversationMedia retrieves media messages (images/videos) in a conversation
// GET /conversations/:conversationId/media?type=image&cursor=xxx&limit=50
func (h *MessageHandler) GetConversationMedia(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	// Get conversationID from URL path parameter
	conversationID, err := uuid.Parse(c.Params("conversationId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid conversation ID")
	}

	// Get optional mediaType filter (image or video)
	mediaType := c.Query("type")
	var mediaTypePtr *string
	if mediaType != "" && (mediaType == "image" || mediaType == "video") {
		mediaTypePtr = &mediaType
	}

	// Get cursor from query params
	cursor := c.Query("cursor")
	var cursorPtr *string
	if cursor != "" {
		cursorPtr = &cursor
	}

	// Get limit from query params (default: 50, max: 100)
	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			if parsedLimit > 0 && parsedLimit <= 100 {
				limit = parsedLimit
			}
		}
	}

	messages, err := h.messageService.ListMediaMessages(c.Context(), conversationID, userID, mediaTypePtr, cursorPtr, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve media messages", err)
	}

	return utils.SuccessResponse(c, "Media messages retrieved successfully", messages)
}

// GetConversationLinks retrieves messages containing links in a conversation
// GET /conversations/:conversationId/links?cursor=xxx&limit=50
func (h *MessageHandler) GetConversationLinks(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	// Get conversationID from URL path parameter
	conversationID, err := uuid.Parse(c.Params("conversationId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid conversation ID")
	}

	// Get cursor from query params
	cursor := c.Query("cursor")
	var cursorPtr *string
	if cursor != "" {
		cursorPtr = &cursor
	}

	// Get limit from query params (default: 50, max: 100)
	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			if parsedLimit > 0 && parsedLimit <= 100 {
				limit = parsedLimit
			}
		}
	}

	messages, err := h.messageService.ListMessagesWithLinks(c.Context(), conversationID, userID, cursorPtr, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve messages with links", err)
	}

	return utils.SuccessResponse(c, "Messages with links retrieved successfully", messages)
}

// GetConversationFiles retrieves file messages in a conversation
// GET /conversations/:conversationId/files?cursor=xxx&limit=50
func (h *MessageHandler) GetConversationFiles(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	// Get conversationID from URL path parameter
	conversationID, err := uuid.Parse(c.Params("conversationId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid conversation ID")
	}

	// Get cursor from query params
	cursor := c.Query("cursor")
	var cursorPtr *string
	if cursor != "" {
		cursorPtr = &cursor
	}

	// Get limit from query params (default: 50, max: 100)
	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			if parsedLimit > 0 && parsedLimit <= 100 {
				limit = parsedLimit
			}
		}
	}

	messages, err := h.messageService.ListFileMessages(c.Context(), conversationID, userID, cursorPtr, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve file messages", err)
	}

	return utils.SuccessResponse(c, "File messages retrieved successfully", messages)
}

// validateFile validates file size and MIME type
func (h *MessageHandler) validateFile(file multipart.File, fileHeader *multipart.FileHeader, messageType string) error {
	// Check file size
	maxSizes := map[string]int64{
		"image": 10 * 1024 * 1024,  // 10MB
		"video": 100 * 1024 * 1024, // 100MB
		"file":  50 * 1024 * 1024,  // 50MB
	}

	if fileHeader.Size > maxSizes[messageType] {
		return fmt.Errorf("file size exceeds maximum %d MB", maxSizes[messageType]/(1024*1024))
	}

	// Check MIME type
	allowedMimeTypes := map[string][]string{
		"image": {"image/jpeg", "image/png", "image/gif", "image/webp"},
		"video": {"video/mp4", "video/quicktime", "video/x-matroska"},
		"file": {
			"application/pdf",
			"application/msword",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			"application/vnd.ms-excel",
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			"application/zip",
			"text/plain",
		},
	}

	// Detect MIME type from file content
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	mimeType := http.DetectContentType(buffer[:n])

	// Check if MIME type is allowed
	allowed := false
	for _, allowedType := range allowedMimeTypes[messageType] {
		if strings.HasPrefix(mimeType, allowedType) {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("file type %s not allowed for %s messages", mimeType, messageType)
	}

	return nil
}

// uploadFile uploads file to Bunny Storage and returns MessageMedia
func (h *MessageHandler) uploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, messageType string) (dto.MessageMedia, error) {
	switch messageType {
	case "image":
		result, err := h.mediaUploadService.UploadImage(ctx, file, fileHeader.Filename)
		if err != nil {
			return dto.MessageMedia{}, err
		}

		return dto.MessageMedia{
			URL:       result.URL,
			Thumbnail: &result.Thumbnail,
			Type:      "image",
			Filename:  &fileHeader.Filename,
			MimeType:  &result.MimeType,
			Size:      &result.Size,
			Width:     &result.Width,
			Height:    &result.Height,
		}, nil

	case "video":
		result, err := h.mediaUploadService.UploadVideo(ctx, file, fileHeader.Filename)
		if err != nil {
			return dto.MessageMedia{}, err
		}

		item := dto.MessageMedia{
			URL:      result.URL,
			Type:     "video",
			Filename: &fileHeader.Filename,
			MimeType: &result.MimeType,
			Size:     &result.Size,
		}

		if result.Thumbnail != "" {
			item.Thumbnail = &result.Thumbnail
		}
		if result.Width > 0 {
			item.Width = &result.Width
			item.Height = &result.Height
		}
		if result.Duration > 0 {
			item.Duration = &result.Duration
		}

		// Add Bunny Stream HLS fields
		if result.VideoID != "" {
			item.VideoID = &result.VideoID
			item.HLSURL = &result.HLSURL
			item.EncodingStatus = &result.EncodingStatus
			item.EncodingProgress = &result.EncodingProgress
		}

		return item, nil

	case "file":
		// Detect MIME type
		buffer := make([]byte, 512)
		file.Read(buffer)
		file.Seek(0, 0) // Reset
		mimeType := http.DetectContentType(buffer)

		result, err := h.mediaUploadService.UploadFile(ctx, file, fileHeader.Filename, mimeType)
		if err != nil {
			return dto.MessageMedia{}, err
		}

		return dto.MessageMedia{
			URL:      result.URL,
			Type:     "file",
			Filename: &fileHeader.Filename,
			MimeType: &result.MimeType,
			Size:     &result.Size,
		}, nil

	default:
		return dto.MessageMedia{}, fmt.Errorf("unsupported message type: %s", messageType)
	}
}

// sendWebSocketNotification sends WebSocket notification to receiver
func (h *MessageHandler) sendWebSocketNotification(message *dto.MessageResponse) {
	if h.chatHub == nil {
		log.Printf("⚠️ ChatHub is nil, skipping WebSocket notification")
		return
	}

	receiverID := message.Receiver.ID

	// Send new message notification to receiver (message.new)
	h.chatHub.SendToUser(receiverID, &chatWebsocket.ChatMessage{
		Type: "message.new",
		Payload: map[string]interface{}{
			"message": message,
		},
	})

	log.Printf("📤 WebSocket notification sent to receiver: %s (message: %s)", receiverID, message.ID)
}
