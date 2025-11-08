package config

import (
	"os"
	"strconv"
	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Bunny    BunnyConfig
	OAuth    OAuthConfig
	VAPID    VAPIDConfig
}

type AppConfig struct {
	Name        string
	Port        string
	Env         string
	FrontendURL string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret string
}

type BunnyConfig struct {
	// Bunny Storage (for images and files)
	StorageZone string
	AccessKey   string
	BaseURL     string
	CDNUrl      string

	// Bunny Stream (for videos)
	StreamAPIKey    string
	StreamLibraryID string
	StreamCDNURL    string
}

type OAuthConfig struct {
	Google GoogleOAuthConfig
}

type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type VAPIDConfig struct {
	PublicKey  string
	PrivateKey string
	Subject    string
}

func LoadConfig() (*Config, error) {
	// Load .env file if it exists (for local development)
	// In production/Docker, environment variables are set by the container
	_ = godotenv.Load()

	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	config := &Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "GoFiber Template"),
			Port:        getEnv("APP_PORT", "3000"),
			Env:         getEnv("APP_ENV", "development"),
			FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "gofiber_template"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-secret-key"),
		},
		Bunny: BunnyConfig{
			// Storage config
			StorageZone: getEnv("BUNNY_STORAGE_ZONE", ""),
			AccessKey:   getEnv("BUNNY_ACCESS_KEY", ""),
			BaseURL:     getEnv("BUNNY_BASE_URL", "https://storage.bunnycdn.com"),
			CDNUrl:      getEnv("BUNNY_CDN_URL", ""),

			// Stream config
			StreamAPIKey:    getEnv("BUNNY_STREAM_API_KEY", ""),
			StreamLibraryID: getEnv("BUNNY_STREAM_LIBRARY_ID", "533535"),
			StreamCDNURL:    getEnv("BUNNY_STREAM_CDN_URL", "https://vz-b1631ae0-4c8.b-cdn.net"),
		},
		OAuth: OAuthConfig{
			Google: GoogleOAuthConfig{
				ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
				ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/v1/auth/google/callback"),
			},
		},
		VAPID: VAPIDConfig{
			PublicKey:  getEnv("VAPID_PUBLIC_KEY", ""),
			PrivateKey: getEnv("VAPID_PRIVATE_KEY", ""),
			Subject:    getEnv("VAPID_SUBJECT", "mailto:admin@voobize.com"),
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}