package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
	"gofiber-template/interfaces/api/routes"
	websocketHandler "gofiber-template/interfaces/api/websocket"
	"gofiber-template/pkg/di"
)

func main() {
	// Initialize DI container
	container := di.NewContainer()

	// Initialize all dependencies
	if err := container.Initialize(); err != nil {
		log.Fatal("Failed to initialize container:", err)
	}

	// Setup graceful shutdown
	setupGracefulShutdown(container)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler(),
		AppName:      container.GetConfig().App.Name,
		BodyLimit:    300 * 1024 * 1024, // 300 MB
	})

	// Setup middleware
	app.Use(middleware.LoggerMiddleware())
	app.Use(middleware.CorsMiddleware())

	// Create handlers from services
	services := container.GetHandlerServices()

	// Create ChatWebSocketHandler
	chatWSHandler := websocketHandler.NewChatWebSocketHandler(container.ChatHub)

	h := handlers.NewHandlers(services, container.GetConfig(), chatWSHandler, container.ChatHub, container.ConversationRepository, container.MediaUploadService)

	// Setup routes
	routes.SetupRoutes(app, h)

	// Start server
	port := container.GetConfig().App.Port
	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("🌍 Environment: %s", container.GetConfig().App.Env)
	log.Printf("📚 Health check: http://localhost:%s/health", port)
	log.Printf("📖 API docs: http://localhost:%s/api/v1", port)
	log.Printf("🔌 WebSocket: ws://localhost:%s/ws", port)
	log.Printf("💬 Chat WebSocket: ws://localhost:%s/chat/ws", port)

	log.Fatal(app.Listen(":" + port))
}

func setupGracefulShutdown(container *di.Container) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("\n🛑 Gracefully shutting down...")

		if err := container.Cleanup(); err != nil {
			log.Printf("❌ Error during cleanup: %v", err)
		}

		log.Println("👋 Shutdown complete")
		os.Exit(0)
	}()
}