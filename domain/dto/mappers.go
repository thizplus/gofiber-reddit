package dto

import (
	"encoding/json"

	"github.com/google/uuid"
	"gofiber-template/domain/models"
)

func UserToUserResponse(user *models.User) *UserResponse {
	if user == nil {
		return nil
	}
	return &UserResponse{
		ID:             user.ID,
		Email:          user.Email,
		Username:       user.Username,
		DisplayName:    user.DisplayName,
		Avatar:         user.Avatar,
		Bio:            user.Bio,
		Location:       user.Location,
		Website:        user.Website,
		Karma:          user.Karma,
		FollowersCount: user.FollowersCount,
		FollowingCount: user.FollowingCount,
		Role:           user.Role,
		IsActive:       user.IsActive,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}

func CreateUserRequestToUser(req *CreateUserRequest) *models.User {
	return &models.User{
		Email:       req.Email,
		Username:    req.Username,
		Password:    req.Password,
		DisplayName: req.DisplayName,
	}
}

func UpdateUserRequestToUser(req *UpdateUserRequest) *models.User {
	return &models.User{
		DisplayName: req.DisplayName,
		Bio:         req.Bio,
		Location:    req.Location,
		Website:     req.Website,
		Avatar:      req.Avatar,
	}
}

func TaskToTaskResponse(task *models.Task, user *models.User) *TaskResponse {
	if task == nil {
		return nil
	}
	taskResp := &TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		DueDate:     task.DueDate,
		UserID:      task.UserID,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
	if user != nil {
		taskResp.User = *UserToUserResponse(user)
	}
	return taskResp
}

func CreateTaskRequestToTask(req *CreateTaskRequest) *models.Task {
	return &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
	}
}

func UpdateTaskRequestToTask(req *UpdateTaskRequest) *models.Task {
	return &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
	}
}

func JobToJobResponse(job *models.Job) *JobResponse {
	if job == nil {
		return nil
	}
	return &JobResponse{
		ID:        job.ID,
		Name:      job.Name,
		CronExpr:  job.CronExpr,
		Payload:   job.Payload,
		Status:    job.Status,
		LastRun:   job.LastRun,
		NextRun:   job.NextRun,
		IsActive:  job.IsActive,
		CreatedAt: job.CreatedAt,
		UpdatedAt: job.UpdatedAt,
	}
}

func CreateJobRequestToJob(req *CreateJobRequest) *models.Job {
	return &models.Job{
		Name:     req.Name,
		CronExpr: req.CronExpr,
		Payload:  req.Payload,
	}
}

func UpdateJobRequestToJob(req *UpdateJobRequest) *models.Job {
	return &models.Job{
		Name:     req.Name,
		CronExpr: req.CronExpr,
		Payload:  req.Payload,
		IsActive: req.IsActive,
	}
}

func FileToFileResponse(file *models.File) *FileResponse {
	if file == nil {
		return nil
	}
	return &FileResponse{
		ID:        file.ID,
		FileName:  file.FileName,
		FileSize:  file.FileSize,
		MimeType:  file.MimeType,
		URL:       file.URL,
		CDNPath:   file.CDNPath,
		UserID:    file.UserID,
		CreatedAt: file.CreatedAt,
		UpdatedAt: file.UpdatedAt,
	}
}

// Post mappers
func PostToPostResponse(post *models.Post) *PostResponse {
	if post == nil {
		return nil
	}

	resp := &PostResponse{
		ID:           post.ID,
		Title:        post.Title,
		Content:      post.Content,
		Author:       *UserToUserResponse(&post.Author),
		Votes:        post.Votes,
		CommentCount: post.CommentCount,
		CreatedAt:    post.CreatedAt,
		UpdatedAt:    post.UpdatedAt,
	}

	// Map media
	if len(post.Media) > 0 {
		resp.Media = make([]MediaResponse, len(post.Media))
		for i, media := range post.Media {
			resp.Media[i] = *MediaToMediaResponse(&media)
		}
	}

	// Map tags
	if len(post.Tags) > 0 {
		resp.Tags = make([]TagResponse, len(post.Tags))
		for i, tag := range post.Tags {
			resp.Tags[i] = *TagToTagResponse(&tag)
		}
	}

	// Map source post (for crossposts)
	if post.SourcePost != nil {
		resp.SourcePost = PostToPostResponse(post.SourcePost)
	}

	return resp
}

// Comment mappers
func CommentToCommentResponse(comment *models.Comment) *CommentResponse {
	if comment == nil {
		return nil
	}

	resp := &CommentResponse{
		ID:        comment.ID,
		PostID:    comment.PostID,
		ParentID:  comment.ParentID,
		Author:    *UserToUserResponse(&comment.Author),
		Content:   comment.Content,
		Votes:     comment.Votes,
		Depth:     comment.Depth,
		IsDeleted: comment.IsDeleted,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}

	// Map post summary if available
	if comment.Post.ID != (uuid.UUID{}) {
		resp.Post = PostToPostSummaryResponse(&comment.Post)
	}

	return resp
}

// PostToPostSummaryResponse converts a Post model to a lightweight PostSummaryResponse
func PostToPostSummaryResponse(post *models.Post) *PostSummaryResponse {
	if post == nil {
		return nil
	}

	return &PostSummaryResponse{
		ID:        post.ID,
		Title:     post.Title,
		Author:    *UserToUserResponse(&post.Author),
		CreatedAt: post.CreatedAt,
	}
}

func CommentToCommentWithReplies(comment *models.Comment, replies []models.Comment) *CommentWithRepliesResponse {
	if comment == nil {
		return nil
	}

	resp := &CommentWithRepliesResponse{
		CommentResponse: *CommentToCommentResponse(comment),
	}

	if len(replies) > 0 {
		resp.Replies = make([]CommentWithRepliesResponse, len(replies))
		for i, reply := range replies {
			replyResp := CommentToCommentWithReplies(&reply, nil)
			if replyResp != nil {
				resp.Replies[i] = *replyResp
			}
		}
	}

	return resp
}

// Vote mappers
func VoteToVoteResponse(vote *models.Vote) *VoteResponse {
	if vote == nil {
		return nil
	}

	return &VoteResponse{
		TargetID:   vote.TargetID,
		TargetType: vote.TargetType,
		VoteType:   vote.VoteType,
		CreatedAt:  vote.CreatedAt,
	}
}

// Tag mappers
func TagToTagResponse(tag *models.Tag) *TagResponse {
	if tag == nil {
		return nil
	}

	return &TagResponse{
		ID:        tag.ID,
		Name:      tag.Name,
		PostCount: tag.PostCount,
		CreatedAt: tag.CreatedAt,
	}
}

// Media mappers
func MediaToMediaResponse(media *models.Media) *MediaResponse {
	if media == nil {
		return nil
	}

	return &MediaResponse{
		ID:        media.ID,
		Type:      media.Type,
		FileName:  media.FileName,
		MimeType:  media.MimeType,
		Size:      media.Size,
		URL:       media.URL,
		Thumbnail: media.Thumbnail,
		Width:     media.Width,
		Height:    media.Height,
		Duration:  media.Duration,
		CreatedAt: media.CreatedAt,
	}
}

func MediaToMediaUploadResponse(media *models.Media) *MediaUploadResponse {
	if media == nil {
		return nil
	}

	return &MediaUploadResponse{
		ID:        media.ID,
		Type:      media.Type,
		FileName:  media.FileName,
		MimeType:  media.MimeType,
		Size:      media.Size,
		URL:       media.URL,
		Thumbnail: media.Thumbnail,
		Width:     media.Width,
		Height:    media.Height,
		Duration:  media.Duration,
		CreatedAt: media.CreatedAt,
	}
}

// Notification mappers
func NotificationToNotificationResponse(notification *models.Notification) *NotificationResponse {
	if notification == nil {
		return nil
	}

	return &NotificationResponse{
		ID:        notification.ID,
		User:      *UserToUserResponse(&notification.User),
		Sender:    *UserToUserResponse(&notification.Sender),
		Type:      notification.Type,
		Message:   notification.Message,
		PostID:    notification.PostID,
		CommentID: notification.CommentID,
		IsRead:    notification.IsRead,
		CreatedAt: notification.CreatedAt,
	}
}

func NotificationSettingsToResponse(settings *models.NotificationSettings) *NotificationSettingsResponse {
	if settings == nil {
		return nil
	}

	return &NotificationSettingsResponse{
		UserID:             settings.UserID,
		Replies:            settings.Replies,
		Mentions:           settings.Mentions,
		Votes:              settings.Votes,
		Follows:            settings.Follows,
		EmailNotifications: settings.EmailNotifications,
		UpdatedAt:          settings.UpdatedAt,
	}
}

// SearchHistory mappers
func SearchHistoryToResponse(history *models.SearchHistory) *SearchHistoryResponse {
	if history == nil {
		return nil
	}

	return &SearchHistoryResponse{
		ID:         history.ID,
		Query:      history.Query,
		Type:       history.Type,
		SearchedAt: history.SearchedAt,
	}
}

// ============================================================================
// Chat mappers
// ============================================================================

// MessageToMessageResponse converts Message model to MessageResponse DTO
func MessageToMessageResponse(message *models.Message) *MessageResponse {
	if message == nil {
		return nil
	}

	resp := &MessageResponse{
		ID:             message.ID,
		ConversationID: message.ConversationID,
		Sender:         *UserToUserResponse(&message.Sender),
		Receiver:       *UserToUserResponse(&message.Receiver),
		Type:           string(message.Type),
		Content:        message.Content,
		IsRead:         message.IsRead,
		ReadAt:         message.ReadAt,
		CreatedAt:      message.CreatedAt,
		UpdatedAt:      message.UpdatedAt,

		// Helper fields
		SenderId:   message.SenderID,
	}

	// Unmarshal Media JSONB to []MessageMedia
	if message.Media != nil && len(message.Media) > 0 {
		var mediaList []MessageMedia
		if err := json.Unmarshal(message.Media, &mediaList); err == nil {
			resp.Media = mediaList
		}
	}

	return resp
}

// ConversationToConversationResponse converts Conversation model to ConversationResponse DTO
// currentUserID is needed to determine who the "other user" is and which unread count to show
func ConversationToConversationResponse(conversation *models.Conversation, currentUserID uuid.UUID) *ConversationResponse {
	if conversation == nil {
		return nil
	}

	// Determine who is the "other user" and get their unread count
	var otherUser models.User
	var unreadCount int

	if conversation.User1ID == currentUserID {
		otherUser = conversation.User2
		unreadCount = conversation.User1UnreadCount
	} else {
		otherUser = conversation.User1
		unreadCount = conversation.User2UnreadCount
	}

	resp := &ConversationResponse{
		ID:            conversation.ID,
		OtherUser:     *UserToUserResponse(&otherUser),
		LastMessageAt: conversation.LastMessageAt,
		UnreadCount:   unreadCount,
		CreatedAt:     conversation.CreatedAt,
		UpdatedAt:     conversation.UpdatedAt,
	}

	return resp
}

// BlockToBlockedUserResponse converts Block model to BlockedUserResponse DTO
func BlockToBlockedUserResponse(block *models.Block) *BlockedUserResponse {
	if block == nil {
		return nil
	}

	return &BlockedUserResponse{
		User:      *UserToUserResponse(&block.Blocked),
		BlockedAt: block.CreatedAt,
	}
}