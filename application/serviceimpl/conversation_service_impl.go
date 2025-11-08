package serviceimpl

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	redisInfra "gofiber-template/infrastructure/redis"
	"gofiber-template/pkg/utils"
)

type ConversationServiceImpl struct {
	conversationRepo repositories.ConversationRepository
	messageRepo      repositories.MessageRepository
	blockRepo        repositories.BlockRepository
	userRepo         repositories.UserRepository
	followRepo       repositories.FollowRepository
	redisService     *redisInfra.RedisService
}

func NewConversationService(
	conversationRepo repositories.ConversationRepository,
	messageRepo repositories.MessageRepository,
	blockRepo repositories.BlockRepository,
	userRepo repositories.UserRepository,
	followRepo repositories.FollowRepository,
	redisService *redisInfra.RedisService,
) services.ConversationService {
	return &ConversationServiceImpl{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
		blockRepo:        blockRepo,
		userRepo:         userRepo,
		followRepo:       followRepo,
		redisService:     redisService,
	}
}

func (s *ConversationServiceImpl) GetOrCreateConversation(ctx context.Context, userID uuid.UUID, otherUsername string) (*dto.ConversationResponse, error) {
	// Get other user by username
	otherUser, err := s.userRepo.GetByUsername(ctx, otherUsername)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Can't chat with yourself
	if userID == otherUser.ID {
		return nil, errors.New("cannot chat with yourself")
	}

	// Check block status
	blocked, blockedBy, err := s.blockRepo.GetBlockStatus(ctx, userID, otherUser.ID)
	if err != nil {
		return nil, err
	}

	if blocked || blockedBy {
		return nil, errors.New("cannot start conversation: user is blocked")
	}

	// Get or create conversation
	conversation, _, err := s.conversationRepo.GetOrCreateByUsers(ctx, userID, otherUser.ID)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	resp := dto.ConversationToConversationResponse(conversation, userID)

	// Load last message if exists
	if conversation.LastMessageID != nil {
		lastMsg, err := s.messageRepo.GetByID(ctx, *conversation.LastMessageID)
		if err == nil {
			resp.LastMessage = dto.MessageToMessageResponse(lastMsg)
		}
	}

	return resp, nil
}

func (s *ConversationServiceImpl) GetConversation(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) (*dto.ConversationResponse, error) {
	conversation, err := s.conversationRepo.GetByID(ctx, conversationID)
	if err != nil {
		return nil, errors.New("conversation not found")
	}

	// Check if user is participant
	if conversation.User1ID != userID && conversation.User2ID != userID {
		return nil, errors.New("access denied: not a participant")
	}

	// Convert to DTO
	resp := dto.ConversationToConversationResponse(conversation, userID)

	// Load last message if exists
	if conversation.LastMessageID != nil {
		lastMsg, err := s.messageRepo.GetByID(ctx, *conversation.LastMessageID)
		if err == nil {
			resp.LastMessage = dto.MessageToMessageResponse(lastMsg)
		}
	}

	return resp, nil
}

func (s *ConversationServiceImpl) ListConversations(ctx context.Context, userID uuid.UUID, cursorStr *string, limit int) (*dto.ConversationListResponse, error) {
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
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	// Fetch conversations (limit + 1 to check for more)
	conversations, err := s.conversationRepo.ListByUser(ctx, userID, cursor, limit+1)
	if err != nil {
		return nil, err
	}

	// Check if there are more results
	hasMore := len(conversations) > limit
	if hasMore {
		conversations = conversations[:limit]
	}

	// Convert to DTOs
	conversationResponses := make([]dto.ConversationResponse, len(conversations))
	for i, conv := range conversations {
		resp := dto.ConversationToConversationResponse(conv, userID)

		// Fetch last message from database (always use DB for complete data)
		if conv.LastMessageID != nil {
			lastMsg, err := s.messageRepo.GetByID(ctx, *conv.LastMessageID)
			if err == nil {
				resp.LastMessage = dto.MessageToMessageResponse(lastMsg)
			}
		}

		// Get unread count from Redis
		unreadCount, _ := s.redisService.GetConversationUnreadCount(ctx, userID, conv.ID)
		resp.UnreadCount = unreadCount

		conversationResponses[i] = *resp
	}

	// Generate next cursor
	var nextCursor *string
	if hasMore && len(conversations) > 0 {
		lastConv := conversations[len(conversations)-1]
		encoded, err := utils.EncodeCursor(lastConv.LastMessageAt)
		if err == nil {
			nextCursor = &encoded
		}
	}

	return &dto.ConversationListResponse{
		Conversations: conversationResponses,
		NextCursor:    nextCursor,
		HasMore:       hasMore,
	}, nil
}

func (s *ConversationServiceImpl) GetUnreadCount(ctx context.Context, userID uuid.UUID) (*dto.UnreadCountResponse, error) {
	// Try Redis first
	count, err := s.redisService.GetTotalUnreadCount(ctx, userID)
	if err == nil && count >= 0 {
		return &dto.UnreadCountResponse{
			TotalUnread: count,
		}, nil
	}

	// Fallback to database if Redis empty/error
	dbCount, err := s.conversationRepo.GetTotalUnreadCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Note: If Redis is empty, we rely on database counts
	// Redis will be populated on next message send or mark as read

	return &dto.UnreadCountResponse{
		TotalUnread: dbCount,
	}, nil
}

func (s *ConversationServiceImpl) MarkAsRead(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
	// Verify user is participant
	conversation, err := s.conversationRepo.GetByID(ctx, conversationID)
	if err != nil {
		return errors.New("conversation not found")
	}

	if conversation.User1ID != userID && conversation.User2ID != userID {
		return errors.New("access denied: not a participant")
	}

	// Get current unread count from Redis before marking
	unreadCount, _ := s.redisService.GetConversationUnreadCount(ctx, userID, conversationID)

	// Mark all messages as read
	if err := s.messageRepo.MarkAllAsRead(ctx, conversationID, userID); err != nil {
		return err
	}

	// Reset unread count in conversation
	if err := s.conversationRepo.ResetUnreadCount(ctx, conversationID, userID); err != nil {
		return err
	}

	// Update Redis
	if unreadCount > 0 {
		// Decrement total unread
		_ = s.redisService.DecrementTotalUnread(ctx, userID, unreadCount)
		// Reset conversation unread (returns count but we ignore it)
		_, _ = s.redisService.ResetConversationUnread(ctx, userID, conversationID)
	}

	return nil
}

func (s *ConversationServiceImpl) SearchUsersForChat(ctx context.Context, userID uuid.UUID, query string, limit int) (*dto.ChatUserSearchResponse, error) {
	// Set default limit
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	// If query is empty, return suggested users only
	var searchResults []*models.User
	var err error

	if query != "" {
		// Search users by query
		searchResults, err = s.userRepo.SearchForChat(ctx, userID, query, limit)
		if err != nil {
			return nil, err
		}
	}

	// Get suggested users (followers/following)
	suggestedUsers, err := s.userRepo.GetSuggestedForChat(ctx, userID, 10)
	if err != nil {
		// Non-critical, continue without suggestions
		suggestedUsers = []*models.User{}
	}

	// Combine all user IDs for batch follow status check
	allUserIDs := make([]uuid.UUID, 0, len(searchResults)+len(suggestedUsers))
	for _, user := range searchResults {
		allUserIDs = append(allUserIDs, user.ID)
	}
	for _, user := range suggestedUsers {
		allUserIDs = append(allUserIDs, user.ID)
	}

	// Batch check follow status
	followStatusMap := make(map[uuid.UUID]bool)
	if len(allUserIDs) > 0 {
		followStatusMap, err = s.followRepo.GetFollowStatus(ctx, userID, allUserIDs)
		if err != nil {
			// Non-critical, continue without follow status
			followStatusMap = make(map[uuid.UUID]bool)
		}
	}

	// Convert to DTOs
	users := make([]dto.ChatUserSearchResult, len(searchResults))
	for i, user := range searchResults {
		isOnline, lastActiveTime, _ := s.redisService.IsUserOnline(ctx, user.ID)
		var lastActive *time.Time
		if !lastActiveTime.IsZero() {
			lastActive = &lastActiveTime
		}

		users[i] = dto.ChatUserSearchResult{
			ID:          user.ID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Avatar:      user.Avatar,
			Bio:         user.Bio,
			IsFollowing: followStatusMap[user.ID],
			IsOnline:    isOnline,
			LastActive:  lastActive,
		}
	}

	suggested := make([]dto.ChatUserSearchResult, len(suggestedUsers))
	for i, user := range suggestedUsers {
		isOnline, lastActiveTime, _ := s.redisService.IsUserOnline(ctx, user.ID)
		var lastActive *time.Time
		if !lastActiveTime.IsZero() {
			lastActive = &lastActiveTime
		}

		suggested[i] = dto.ChatUserSearchResult{
			ID:          user.ID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Avatar:      user.Avatar,
			Bio:         user.Bio,
			IsFollowing: followStatusMap[user.ID],
			IsOnline:    isOnline,
			LastActive:  lastActive,
		}
	}

	return &dto.ChatUserSearchResponse{
		Users:     users,
		Suggested: suggested,
	}, nil
}

// Ensure interface compliance
var _ services.ConversationService = (*ConversationServiceImpl)(nil)
