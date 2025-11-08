package postgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gofiber-template/domain/models"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewDatabase(config DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// Core models (enhanced/new)
		&models.User{},
		&models.Post{},
		&models.Comment{},
		&models.Media{},

		// Voting & Social
		&models.Vote{},
		&models.Follow{},
		&models.SavedPost{},

		// Notifications
		&models.Notification{},
		&models.NotificationSettings{},
		&models.PushSubscription{},

		// Tags
		&models.Tag{},

		// Search
		&models.SearchHistory{},

		// Chat System (Order matters: Conversation first, then Message, then Block)
		&models.Conversation{},
		&models.Block{},
		&models.Message{},

		// Legacy models (keep for now, can remove later)
		&models.Task{},
		&models.File{},
		&models.Job{},
	)
}