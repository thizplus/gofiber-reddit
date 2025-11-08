package di

import (
	"context"
	"log"
	"gofiber-template/application/serviceimpl"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	"gofiber-template/infrastructure/postgres"
	"gofiber-template/infrastructure/redis"
	"gofiber-template/infrastructure/storage"
	"gofiber-template/infrastructure/websocket"
	"gofiber-template/infrastructure/workers"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/pkg/config"
	"gofiber-template/pkg/scheduler"
	"gorm.io/gorm"
)

type Container struct {
	// Configuration
	Config *config.Config

	// Infrastructure
	DB                 *gorm.DB
	RedisClient        *redis.RedisClient
	RedisService       *redis.RedisService
	BunnyStorage       storage.BunnyStorage
	BunnyStreamService *storage.BunnyStreamService
	MediaUploadService *storage.MediaUploadService
	EventScheduler     scheduler.EventScheduler
	ChatHub            *websocket.ChatHub
	VideoEncoderWorker *workers.VideoEncoderWorker

	// Repositories - Legacy
	UserRepository repositories.UserRepository
	TaskRepository repositories.TaskRepository
	FileRepository repositories.FileRepository
	JobRepository  repositories.JobRepository

	// Repositories - Social Media
	PostRepository                 repositories.PostRepository
	CommentRepository              repositories.CommentRepository
	VoteRepository                 repositories.VoteRepository
	FollowRepository               repositories.FollowRepository
	SavedPostRepository            repositories.SavedPostRepository
	NotificationRepository         repositories.NotificationRepository
	NotificationSettingsRepository repositories.NotificationSettingsRepository
	PushSubscriptionRepository     repositories.PushSubscriptionRepository
	TagRepository                  repositories.TagRepository
	SearchHistoryRepository        repositories.SearchHistoryRepository
	MediaRepository                repositories.MediaRepository

	// Repositories - Chat System
	ConversationRepository repositories.ConversationRepository
	MessageRepository      repositories.MessageRepository
	BlockRepository        repositories.BlockRepository

	// Services - Legacy
	UserService services.UserService
	TaskService services.TaskService
	FileService services.FileService
	JobService  services.JobService

	// Services - Social Media
	PostService         services.PostService
	CommentService      services.CommentService
	VoteService         services.VoteService
	FollowService       services.FollowService
	SavedPostService    services.SavedPostService
	NotificationService services.NotificationService
	PushService         services.PushService
	TagService          services.TagService
	SearchService       services.SearchService
	MediaService        services.MediaService
	OAuthService        services.OAuthService

	// Services - Chat System
	ConversationService services.ConversationService
	MessageService      services.MessageService
	BlockService        services.BlockService

	// Services - Upload
	FileUploadService services.FileUploadService
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) Initialize() error {
	if err := c.initConfig(); err != nil {
		return err
	}

	if err := c.initInfrastructure(); err != nil {
		return err
	}

	if err := c.initRepositories(); err != nil {
		return err
	}

	if err := c.initServices(); err != nil {
		return err
	}

	if err := c.initScheduler(); err != nil {
		return err
	}

	if err := c.initChatHub(); err != nil {
		return err
	}

	if err := c.initVideoEncoderWorker(); err != nil {
		return err
	}

	return nil
}

func (c *Container) initConfig() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	c.Config = cfg
	log.Println("✓ Configuration loaded")
	return nil
}

func (c *Container) initInfrastructure() error {
	// Initialize Database
	dbConfig := postgres.DatabaseConfig{
		Host:     c.Config.Database.Host,
		Port:     c.Config.Database.Port,
		User:     c.Config.Database.User,
		Password: c.Config.Database.Password,
		DBName:   c.Config.Database.DBName,
		SSLMode:  c.Config.Database.SSLMode,
	}

	db, err := postgres.NewDatabase(dbConfig)
	if err != nil {
		return err
	}
	c.DB = db
	log.Println("✓ Database connected")

	// Run migrations
	if err := postgres.Migrate(db); err != nil {
		return err
	}
	log.Println("✓ Database migrated")

	// Initialize Redis
	redisConfig := redis.RedisConfig{
		Host:     c.Config.Redis.Host,
		Port:     c.Config.Redis.Port,
		Password: c.Config.Redis.Password,
		DB:       c.Config.Redis.DB,
	}
	c.RedisClient = redis.NewRedisClient(redisConfig)

	// Test Redis connection
	if err := c.RedisClient.Ping(context.Background()); err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
	} else {
		log.Println("✓ Redis connected")
	}

	// Initialize RedisService
	c.RedisService = redis.NewRedisService(c.RedisClient)
	log.Println("✓ RedisService initialized")

	// Initialize Bunny Storage
	bunnyConfig := storage.BunnyConfig{
		StorageZone: c.Config.Bunny.StorageZone,
		AccessKey:   c.Config.Bunny.AccessKey,
		BaseURL:     c.Config.Bunny.BaseURL,
		CDNUrl:      c.Config.Bunny.CDNUrl,
	}
	c.BunnyStorage = storage.NewBunnyStorage(bunnyConfig)
	log.Println("✓ Bunny Storage initialized")

	// Initialize Bunny Stream
	c.BunnyStreamService = storage.NewBunnyStreamService(
		c.Config.Bunny.StreamAPIKey,
		c.Config.Bunny.StreamLibraryID,
		c.Config.Bunny.StreamCDNURL,
	)
	log.Println("✓ Bunny Stream initialized")

	// Initialize MediaUploadService
	c.MediaUploadService = storage.NewMediaUploadService(c.BunnyStorage, c.BunnyStreamService)
	log.Println("✓ MediaUploadService initialized")

	return nil
}

func (c *Container) initRepositories() error {
	// Legacy repositories
	c.UserRepository = postgres.NewUserRepository(c.DB)
	c.TaskRepository = postgres.NewTaskRepository(c.DB)
	c.FileRepository = postgres.NewFileRepository(c.DB)
	c.JobRepository = postgres.NewJobRepository(c.DB)

	// Social media repositories
	c.PostRepository = postgres.NewPostRepository(c.DB)
	c.CommentRepository = postgres.NewCommentRepository(c.DB)
	c.VoteRepository = postgres.NewVoteRepository(c.DB)
	c.FollowRepository = postgres.NewFollowRepository(c.DB)
	c.SavedPostRepository = postgres.NewSavedPostRepository(c.DB)
	c.NotificationRepository = postgres.NewNotificationRepository(c.DB)
	c.NotificationSettingsRepository = postgres.NewNotificationSettingsRepository(c.DB)
	c.PushSubscriptionRepository = postgres.NewPushSubscriptionRepository(c.DB)
	c.TagRepository = postgres.NewTagRepository(c.DB)
	c.SearchHistoryRepository = postgres.NewSearchHistoryRepository(c.DB)
	c.MediaRepository = postgres.NewMediaRepository(c.DB)

	// Chat system repositories
	c.ConversationRepository = postgres.NewConversationRepository(c.DB)
	c.MessageRepository = postgres.NewMessageRepository(c.DB)
	c.BlockRepository = postgres.NewBlockRepository(c.DB)

	log.Println("✓ Repositories initialized (18 repositories)")
	return nil
}

func (c *Container) initServices() error {
	// Legacy services
	c.UserService = serviceimpl.NewUserService(c.UserRepository, c.FollowRepository, c.Config.JWT.Secret)
	c.TaskService = serviceimpl.NewTaskService(c.TaskRepository, c.UserRepository)
	c.FileService = serviceimpl.NewFileService(c.FileRepository, c.UserRepository, c.BunnyStorage)

	// OAuth service
	c.OAuthService = serviceimpl.NewOAuthService(c.UserRepository, c.Config)

	// Social media services (order matters due to dependencies)
	// 1. No service dependencies
	c.TagService = serviceimpl.NewTagService(c.TagRepository)
	c.NotificationService = serviceimpl.NewNotificationService(
		c.NotificationRepository,
		c.NotificationSettingsRepository,
		c.UserRepository,
	)
	c.PushService = serviceimpl.NewPushService(
		c.PushSubscriptionRepository,
		c.Config,
	)

	// 2. Depends on TagService
	c.PostService = serviceimpl.NewPostService(
		c.PostRepository,
		c.UserRepository,
		c.VoteRepository,
		c.SavedPostRepository,
		c.TagService,
		c.MediaRepository,
	)

	// 3. Depends on NotificationService
	c.CommentService = serviceimpl.NewCommentService(
		c.CommentRepository,
		c.PostRepository,
		c.VoteRepository,
		c.NotificationService,
	)
	c.VoteService = serviceimpl.NewVoteService(
		c.VoteRepository,
		c.PostRepository,
		c.CommentRepository,
		c.UserRepository,
		c.NotificationService,
	)
	c.FollowService = serviceimpl.NewFollowService(
		c.FollowRepository,
		c.UserRepository,
		c.NotificationService,
	)

	// 4. Independent services
	c.SavedPostService = serviceimpl.NewSavedPostService(
		c.SavedPostRepository,
		c.PostRepository,
		c.VoteRepository,
	)
	c.SearchService = serviceimpl.NewSearchService(
		c.PostRepository,
		c.UserRepository,
		c.TagRepository,
		c.SearchHistoryRepository,
		c.VoteRepository,
		c.SavedPostRepository,
	)
	c.MediaService = serviceimpl.NewMediaService(
		c.MediaRepository,
		c.BunnyStorage,
		c.BunnyStreamService,
		c.RedisService,
	)

	// 5. Chat system services
	c.ConversationService = serviceimpl.NewConversationService(
		c.ConversationRepository,
		c.MessageRepository,
		c.BlockRepository,
		c.UserRepository,
		c.FollowRepository,
		c.RedisService,
	)
	c.MessageService = serviceimpl.NewMessageService(
		c.MessageRepository,
		c.ConversationRepository,
		c.BlockRepository,
		c.UserRepository,
		c.RedisService,
	)
	c.BlockService = serviceimpl.NewBlockService(
		c.BlockRepository,
		c.UserRepository,
	)

	// 6. Upload services
	c.FileUploadService = serviceimpl.NewFileUploadService(
		c.MediaRepository,
		c.MediaUploadService,
	)

	// Set push service for notification service (to avoid circular dependency)
	if notifService, ok := c.NotificationService.(*serviceimpl.NotificationServiceImpl); ok {
		notifService.SetPushService(c.PushService)
	}

	log.Println("✓ Services initialized (19 services)")
	return nil
}

func (c *Container) initScheduler() error {
	c.EventScheduler = scheduler.NewEventScheduler()
	c.JobService = serviceimpl.NewJobService(c.JobRepository, c.EventScheduler)

	// Start the scheduler
	c.EventScheduler.Start()
	log.Println("✓ Event scheduler started")

	// Load and schedule existing active jobs
	ctx := context.Background()
	jobs, _, err := c.JobService.ListJobs(ctx, 0, 1000)
	if err != nil {
		log.Printf("Warning: Failed to load existing jobs: %v", err)
		return nil
	}

	activeJobCount := 0
	for _, job := range jobs {
		if job.IsActive {
			err := c.EventScheduler.AddJob(job.ID.String(), job.CronExpr, func() {
				c.JobService.ExecuteJob(ctx, job)
			})
			if err != nil {
				log.Printf("Warning: Failed to schedule job %s: %v", job.Name, err)
			} else {
				activeJobCount++
			}
		}
	}

	if activeJobCount > 0 {
		log.Printf("✓ Scheduled %d active jobs", activeJobCount)
	}

	return nil
}

func (c *Container) initChatHub() error {
	c.ChatHub = websocket.NewChatHub(
		c.MessageService,
		c.ConversationService,
		c.BlockService,
		c.RedisService,
		c.ConversationRepository,
		c.FollowRepository,
		c.PushService,
	)

	// Start ChatHub in background
	go c.ChatHub.Run()
	log.Println("✓ ChatHub started")

	return nil
}

func (c *Container) initVideoEncoderWorker() error {
	c.VideoEncoderWorker = workers.NewVideoEncoderWorker(
		c.RedisService,
		c.BunnyStreamService,
		c.MediaRepository,
		c.NotificationService,
	)

	// Start worker in background
	c.VideoEncoderWorker.Start()

	return nil
}

func (c *Container) Cleanup() error {
	log.Println("Starting cleanup...")

	// Stop VideoEncoderWorker
	if c.VideoEncoderWorker != nil {
		c.VideoEncoderWorker.Stop()
		log.Println("✓ VideoEncoderWorker stopped")
	}

	// Stop ChatHub
	if c.ChatHub != nil {
		c.ChatHub.Stop()
		log.Println("✓ ChatHub stopped")
	}

	// Stop scheduler
	if c.EventScheduler != nil {
		if c.EventScheduler.IsRunning() {
			c.EventScheduler.Stop()
			log.Println("✓ Event scheduler stopped")
		} else {
			log.Println("✓ Event scheduler was already stopped")
		}
	}

	// Close Redis connection
	if c.RedisClient != nil {
		if err := c.RedisClient.Close(); err != nil {
			log.Printf("Warning: Failed to close Redis connection: %v", err)
		} else {
			log.Println("✓ Redis connection closed")
		}
	}

	// Close database connection
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Printf("Warning: Failed to close database connection: %v", err)
			} else {
				log.Println("✓ Database connection closed")
			}
		}
	}

	log.Println("✓ Cleanup completed")
	return nil
}

func (c *Container) GetServices() (services.UserService, services.TaskService, services.FileService, services.JobService) {
	return c.UserService, c.TaskService, c.FileService, c.JobService
}

func (c *Container) GetConfig() *config.Config {
	return c.Config
}

func (c *Container) GetHandlerServices() *handlers.Services {
	return &handlers.Services{
		// Legacy services
		UserService: c.UserService,
		TaskService: c.TaskService,
		FileService: c.FileService,
		JobService:  c.JobService,

		// Social media services
		PostService:         c.PostService,
		CommentService:      c.CommentService,
		VoteService:         c.VoteService,
		FollowService:       c.FollowService,
		SavedPostService:    c.SavedPostService,
		NotificationService: c.NotificationService,
		PushService:         c.PushService,
		TagService:          c.TagService,
		SearchService:       c.SearchService,
		MediaService:        c.MediaService,
		OAuthService:        c.OAuthService,

		// Chat system services
		ConversationService: c.ConversationService,
		MessageService:      c.MessageService,
		BlockService:        c.BlockService,

		// Upload services
		FileUploadService: c.FileUploadService,
	}
}