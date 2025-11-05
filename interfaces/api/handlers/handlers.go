package handlers

import (
	"gofiber-template/domain/services"
	"gofiber-template/pkg/config"
)

// Services contains all the services needed for handlers
type Services struct {
	UserService         services.UserService
	TaskService         services.TaskService
	FileService         services.FileService
	JobService          services.JobService
	PostService         services.PostService
	CommentService      services.CommentService
	VoteService         services.VoteService
	FollowService       services.FollowService
	SavedPostService    services.SavedPostService
	NotificationService services.NotificationService
	TagService          services.TagService
	SearchService       services.SearchService
	MediaService        services.MediaService
	OAuthService        services.OAuthService
}

// Handlers contains all HTTP handlers
type Handlers struct {
	UserHandler         *UserHandler
	TaskHandler         *TaskHandler
	FileHandler         *FileHandler
	JobHandler          *JobHandler
	PostHandler         *PostHandler
	CommentHandler      *CommentHandler
	VoteHandler         *VoteHandler
	FollowHandler       *FollowHandler
	SavedPostHandler    *SavedPostHandler
	NotificationHandler *NotificationHandler
	TagHandler          *TagHandler
	SearchHandler       *SearchHandler
	MediaHandler        *MediaHandler
	OAuthHandler        *OAuthHandler
}

// NewHandlers creates a new instance of Handlers with all dependencies
func NewHandlers(services *Services, cfg *config.Config) *Handlers {
	return &Handlers{
		UserHandler:         NewUserHandler(services.UserService),
		TaskHandler:         NewTaskHandler(services.TaskService),
		FileHandler:         NewFileHandler(services.FileService),
		JobHandler:          NewJobHandler(services.JobService),
		PostHandler:         NewPostHandler(services.PostService),
		CommentHandler:      NewCommentHandler(services.CommentService),
		VoteHandler:         NewVoteHandler(services.VoteService),
		FollowHandler:       NewFollowHandler(services.FollowService),
		SavedPostHandler:    NewSavedPostHandler(services.SavedPostService),
		NotificationHandler: NewNotificationHandler(services.NotificationService),
		TagHandler:          NewTagHandler(services.TagService),
		SearchHandler:       NewSearchHandler(services.SearchService),
		MediaHandler:        NewMediaHandler(services.MediaService),
		OAuthHandler:        NewOAuthHandler(services.OAuthService, cfg),
	}
}