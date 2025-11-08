package serviceimpl

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	"gofiber-template/infrastructure/redis"
	"gofiber-template/pkg/utils"
)

type MessageServiceImpl struct {
	messageRepo      repositories.MessageRepository
	conversationRepo repositories.ConversationRepository
	blockRepo        repositories.BlockRepository
	userRepo         repositories.UserRepository
	redisService     *redis.RedisService
}

func NewMessageService(
	messageRepo repositories.MessageRepository,
	conversationRepo repositories.ConversationRepository,
	blockRepo repositories.BlockRepository,
	userRepo repositories.UserRepository,
	redisService *redis.RedisService,
) services.MessageService {
	return &MessageServiceImpl{
		messageRepo:      messageRepo,
		conversationRepo: conversationRepo,
		blockRepo:        blockRepo,
		userRepo:         userRepo,
		redisService:     redisService,
	}
}

func (s *MessageServiceImpl) SendMessage(ctx context.Context, userID uuid.UUID, req *dto.SendMessageRequest) (*dto.MessageResponse, error) {
	// Validate: Content OR Media must be provided
	if (req.Content == nil || *req.Content == "") && len(req.Media) == 0 {
		return nil, errors.New("either content or media must be provided")
	}

	// Get conversation
	conversation, err := s.conversationRepo.GetByID(ctx, req.ConversationID)
	if err != nil {
		return nil, errors.New("conversation not found")
	}

	// Check if user is participant
	if conversation.User1ID != userID && conversation.User2ID != userID {
		return nil, errors.New("access denied: not a participant")
	}

	// Determine sender and receiver
	var receiverID uuid.UUID
	if conversation.User1ID == userID {
		receiverID = conversation.User2ID
	} else {
		receiverID = conversation.User1ID
	}

	// Check block status
	blocked, blockedBy, err := s.blockRepo.GetBlockStatus(ctx, userID, receiverID)
	if err != nil {
		return nil, err
	}

	if blocked || blockedBy {
		return nil, errors.New("cannot send message: user is blocked")
	}

	// Convert MessageType string to enum
	messageType := models.MessageType(req.Type)

	// Convert Media array to JSONB
	var mediaJSON datatypes.JSON
	if len(req.Media) > 0 {
		mediaBytes, err := json.Marshal(req.Media)
		if err != nil {
			return nil, errors.New("failed to process media")
		}
		mediaJSON = datatypes.JSON(mediaBytes)
	}

	// Create message
	now := time.Now()
	message := &models.Message{
		ConversationID: req.ConversationID,
		SenderID:       userID,
		ReceiverID:     receiverID,
		Type:           messageType,
		Content:        req.Content,
		Media:          mediaJSON,
		IsRead:         false,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.messageRepo.Create(ctx, message); err != nil {
		return nil, err
	}

	// Update conversation last message and increment unread count
	_ = s.conversationRepo.UpdateLastMessage(ctx, req.ConversationID, message.ID, now)
	_ = s.conversationRepo.IncrementUnreadCount(ctx, req.ConversationID, receiverID)

	// Increment unread counts in Redis
	_ = s.redisService.IncrementTotalUnread(ctx, receiverID)
	_ = s.redisService.IncrementConversationUnread(ctx, receiverID, req.ConversationID)

	// Cache last message in Redis
	_ = s.redisService.CacheLastMessage(ctx, req.ConversationID, message.ID, message.SenderID, message.Content, string(message.Type), message.CreatedAt)

	// Load message with full relations
	fullMessage, err := s.messageRepo.GetByID(ctx, message.ID)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	resp := dto.MessageToMessageResponse(fullMessage)
	if req.TempID != nil {
		resp.TempID = req.TempID
	}

	return resp, nil
}

func (s *MessageServiceImpl) GetMessage(ctx context.Context, messageID uuid.UUID, userID uuid.UUID) (*dto.MessageResponse, error) {
	message, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return nil, errors.New("message not found")
	}

	// Check if user is participant
	if message.SenderID != userID && message.ReceiverID != userID {
		return nil, errors.New("access denied")
	}

	return dto.MessageToMessageResponse(message), nil
}

func (s *MessageServiceImpl) ListMessages(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, cursorStr *string, limit int) (*dto.MessageListResponse, error) {
	// Verify user is participant
	conversation, err := s.conversationRepo.GetByID(ctx, conversationID)
	if err != nil {
		return nil, errors.New("conversation not found")
	}

	if conversation.User1ID != userID && conversation.User2ID != userID {
		return nil, errors.New("access denied: not a participant")
	}

	// Decode cursor if provided
	var cursor *time.Time
	if cursorStr != nil && *cursorStr != "" {
		decoded, err := utils.DecodeCursor(*cursorStr)
		if err != nil {
			return nil, errors.New("invalid cursor")
		}
		cursor = decoded
	}

	// Set default limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	// Fetch messages (limit + 1 to check for more)
	messages, err := s.messageRepo.ListByConversation(ctx, conversationID, cursor, limit+1)
	if err != nil {
		return nil, err
	}

	// Check if there are more results
	hasMore := len(messages) > limit
	if hasMore {
		messages = messages[:limit]
	}

	// Convert to DTOs
	messageResponses := make([]dto.MessageResponse, len(messages))
	for i, msg := range messages {
		messageResponses[i] = *dto.MessageToMessageResponse(msg)
	}

	// Generate next cursor
	var nextCursor *string
	if hasMore && len(messages) > 0 {
		lastMsg := messages[len(messages)-1]
		encoded, err := utils.EncodeCursor(lastMsg.CreatedAt)
		if err == nil {
			nextCursor = &encoded
		}
	}

	return &dto.MessageListResponse{
		Messages:   messageResponses,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (s *MessageServiceImpl) GetMessageContext(ctx context.Context, messageID uuid.UUID, userID uuid.UUID) (*dto.MessageContextResponse, error) {
	// Get target message
	targetMessage, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return nil, errors.New("message not found")
	}

	// Verify user is participant
	if targetMessage.SenderID != userID && targetMessage.ReceiverID != userID {
		return nil, errors.New("access denied")
	}

	// Get messages before (20 messages)
	messagesBefore, err := s.messageRepo.GetMessagesBeforeTimestamp(ctx, targetMessage.ConversationID, targetMessage.CreatedAt, 20)
	if err != nil {
		return nil, err
	}

	// Get messages after (20 messages)
	messagesAfter, err := s.messageRepo.GetMessagesAfterTimestamp(ctx, targetMessage.ConversationID, targetMessage.CreatedAt, 20)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	beforeDTOs := make([]dto.MessageResponse, len(messagesBefore))
	for i, msg := range messagesBefore {
		beforeDTOs[i] = *dto.MessageToMessageResponse(msg)
	}

	afterDTOs := make([]dto.MessageResponse, len(messagesAfter))
	for i, msg := range messagesAfter {
		afterDTOs[i] = *dto.MessageToMessageResponse(msg)
	}

	// Generate cursors
	var beforeCursor, afterCursor *string
	if len(messagesBefore) > 0 {
		encoded, err := utils.EncodeCursor(messagesBefore[0].CreatedAt)
		if err == nil {
			beforeCursor = &encoded
		}
	}
	if len(messagesAfter) > 0 {
		encoded, err := utils.EncodeCursor(messagesAfter[len(messagesAfter)-1].CreatedAt)
		if err == nil {
			afterCursor = &encoded
		}
	}

	return &dto.MessageContextResponse{
		TargetMessage:  *dto.MessageToMessageResponse(targetMessage),
		Before:         beforeDTOs,
		After:          afterDTOs,
		BeforeCursor:   beforeCursor,
		AfterCursor:    afterCursor,
		HasMoreBefore:  len(messagesBefore) == 20,
		HasMoreAfter:   len(messagesAfter) == 20,
	}, nil
}

func (s *MessageServiceImpl) MarkAsRead(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
	// Verify user is participant
	conversation, err := s.conversationRepo.GetByID(ctx, conversationID)
	if err != nil {
		return errors.New("conversation not found")
	}

	if conversation.User1ID != userID && conversation.User2ID != userID {
		return errors.New("access denied: not a participant")
	}

	// Mark all messages as read
	if err := s.messageRepo.MarkAllAsRead(ctx, conversationID, userID); err != nil {
		return err
	}

	// Reset unread count
	return s.conversationRepo.ResetUnreadCount(ctx, conversationID, userID)
}

// Phase 2: Media/Links/Files Queries

func (s *MessageServiceImpl) ListMediaMessages(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, mediaType *string, cursor *string, limit int) (*dto.MessageListResponse, error) {
	// Verify access
	conversation, err := s.conversationRepo.GetByID(ctx, conversationID)
	if err != nil {
		return nil, errors.New("conversation not found")
	}

	if conversation.User1ID != userID && conversation.User2ID != userID {
		return nil, errors.New("access denied: not a participant")
	}

	// Check if blocked
	isBlocked, err := s.blockRepo.IsBlocked(ctx, userID, conversation.User1ID)
	if err == nil && isBlocked {
		return nil, errors.New("access denied: blocked")
	}
	isBlocked, err = s.blockRepo.IsBlocked(ctx, userID, conversation.User2ID)
	if err == nil && isBlocked {
		return nil, errors.New("access denied: blocked")
	}

	// Parse cursor
	var beforeCursor *time.Time
	if cursor != nil && *cursor != "" {
		decoded, err := utils.DecodeCursor(*cursor)
		if err == nil {
			beforeCursor = decoded
		}
	}

	// Get media messages
	messages, err := s.messageRepo.ListMediaMessages(ctx, conversationID, mediaType, beforeCursor, limit)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	responseDTOs := make([]dto.MessageResponse, 0, len(messages))
	for _, message := range messages {
		responseDTOs = append(responseDTOs, *dto.MessageToMessageResponse(message))
	}

	// Generate next cursor
	var nextCursor *string
	if len(messages) == limit && len(messages) > 0 {
		encoded, err := utils.EncodeCursor(messages[len(messages)-1].CreatedAt)
		if err == nil {
			nextCursor = &encoded
		}
	}

	return &dto.MessageListResponse{
		Messages:   responseDTOs,
		NextCursor: nextCursor,
		HasMore:    len(messages) == limit,
	}, nil
}

func (s *MessageServiceImpl) ListMessagesWithLinks(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, cursor *string, limit int) (*dto.MessageListResponse, error) {
	// Verify access
	conversation, err := s.conversationRepo.GetByID(ctx, conversationID)
	if err != nil {
		return nil, errors.New("conversation not found")
	}

	if conversation.User1ID != userID && conversation.User2ID != userID {
		return nil, errors.New("access denied: not a participant")
	}

	// Check if blocked
	isBlocked, err := s.blockRepo.IsBlocked(ctx, userID, conversation.User1ID)
	if err == nil && isBlocked {
		return nil, errors.New("access denied: blocked")
	}
	isBlocked, err = s.blockRepo.IsBlocked(ctx, userID, conversation.User2ID)
	if err == nil && isBlocked {
		return nil, errors.New("access denied: blocked")
	}

	// Parse cursor
	var beforeCursor *time.Time
	if cursor != nil && *cursor != "" {
		decoded, err := utils.DecodeCursor(*cursor)
		if err == nil {
			beforeCursor = decoded
		}
	}

	// Get messages with links
	messages, err := s.messageRepo.ListMessagesWithLinks(ctx, conversationID, beforeCursor, limit)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	responseDTOs := make([]dto.MessageResponse, 0, len(messages))
	for _, message := range messages {
		responseDTOs = append(responseDTOs, *dto.MessageToMessageResponse(message))
	}

	// Generate next cursor
	var nextCursor *string
	if len(messages) == limit && len(messages) > 0 {
		encoded, err := utils.EncodeCursor(messages[len(messages)-1].CreatedAt)
		if err == nil {
			nextCursor = &encoded
		}
	}

	return &dto.MessageListResponse{
		Messages:   responseDTOs,
		NextCursor: nextCursor,
		HasMore:    len(messages) == limit,
	}, nil
}

func (s *MessageServiceImpl) ListFileMessages(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID, cursor *string, limit int) (*dto.MessageListResponse, error) {
	// Verify access
	conversation, err := s.conversationRepo.GetByID(ctx, conversationID)
	if err != nil {
		return nil, errors.New("conversation not found")
	}

	if conversation.User1ID != userID && conversation.User2ID != userID {
		return nil, errors.New("access denied: not a participant")
	}

	// Check if blocked
	isBlocked, err := s.blockRepo.IsBlocked(ctx, userID, conversation.User1ID)
	if err == nil && isBlocked {
		return nil, errors.New("access denied: blocked")
	}
	isBlocked, err = s.blockRepo.IsBlocked(ctx, userID, conversation.User2ID)
	if err == nil && isBlocked {
		return nil, errors.New("access denied: blocked")
	}

	// Parse cursor
	var beforeCursor *time.Time
	if cursor != nil && *cursor != "" {
		decoded, err := utils.DecodeCursor(*cursor)
		if err == nil {
			beforeCursor = decoded
		}
	}

	// Get file messages
	messages, err := s.messageRepo.ListFileMessages(ctx, conversationID, beforeCursor, limit)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	responseDTOs := make([]dto.MessageResponse, 0, len(messages))
	for _, message := range messages {
		responseDTOs = append(responseDTOs, *dto.MessageToMessageResponse(message))
	}

	// Generate next cursor
	var nextCursor *string
	if len(messages) == limit && len(messages) > 0 {
		encoded, err := utils.EncodeCursor(messages[len(messages)-1].CreatedAt)
		if err == nil {
			nextCursor = &encoded
		}
	}

	return &dto.MessageListResponse{
		Messages:   responseDTOs,
		NextCursor: nextCursor,
		HasMore:    len(messages) == limit,
	}, nil
}

// Ensure interface compliance
var _ services.MessageService = (*MessageServiceImpl)(nil)
