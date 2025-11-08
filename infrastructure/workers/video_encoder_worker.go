package workers

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	"gofiber-template/infrastructure/redis"
	"gofiber-template/infrastructure/storage"
	"gofiber-template/infrastructure/websocket"
)

const (
	workerInterval = 10 * time.Second // Poll every 10 seconds
)

type VideoEncoderWorker struct {
	redisService    *redis.RedisService
	bunnyStream     *storage.BunnyStreamService
	mediaRepo       repositories.MediaRepository
	notifService    services.NotificationService
	running         bool
	stopChan        chan struct{}
}

func NewVideoEncoderWorker(
	redisService *redis.RedisService,
	bunnyStream *storage.BunnyStreamService,
	mediaRepo repositories.MediaRepository,
	notifService services.NotificationService,
) *VideoEncoderWorker {
	return &VideoEncoderWorker{
		redisService: redisService,
		bunnyStream:  bunnyStream,
		mediaRepo:    mediaRepo,
		notifService: notifService,
		stopChan:     make(chan struct{}),
	}
}

// Start begins the worker's processing loop
func (w *VideoEncoderWorker) Start() {
	if w.running {
		log.Println("⚠️  VideoEncoderWorker is already running")
		return
	}

	w.running = true
	log.Println("🎬 VideoEncoderWorker started")

	go w.processLoop()
}

// Stop gracefully stops the worker
func (w *VideoEncoderWorker) Stop() {
	if !w.running {
		return
	}

	log.Println("🛑 Stopping VideoEncoderWorker...")
	w.running = false
	close(w.stopChan)
}

// processLoop is the main worker loop
func (w *VideoEncoderWorker) processLoop() {
	ticker := time.NewTicker(workerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.processNextJob()
		case <-w.stopChan:
			log.Println("✓ VideoEncoderWorker stopped")
			return
		}
	}
}

// processNextJob processes the next video in the queue
func (w *VideoEncoderWorker) processNextJob() {
	ctx := context.Background()

	// Check queue size
	queueSize, err := w.redisService.GetPendingEncodingCount(ctx)
	if err != nil {
		log.Printf("❌ Failed to get queue size: %v", err)
		return
	}

	if queueSize == 0 {
		// Queue is empty, skip this iteration
		return
	}

	// Dequeue the next job
	job, err := w.redisService.DequeueVideoEncoding(ctx)
	if err != nil {
		if err.Error() != "redis: nil" { // Ignore empty queue error
			log.Printf("❌ Failed to dequeue job: %v", err)
		}
		return
	}

	log.Printf("🎬 Processing video encoding job: mediaID=%s, videoID=%s", job.MediaID, job.VideoID)

	// Get video status from Bunny Stream
	videoInfo, err := w.bunnyStream.GetVideo(job.VideoID)
	if err != nil {
		log.Printf("❌ Failed to get video info from Bunny Stream: %v", err)
		// Update status to failed
		w.redisService.UpdateVideoEncodingStatus(ctx, job.MediaID, "failed", 0, err.Error())
		w.updateMediaStatus(ctx, job.MediaID, "failed", 0, "")
		return
	}

	// Map Bunny Stream status to our status
	// Status: 0: Queued, 1: Processing, 2: Encoding, 3: Finished, 4: Resolution Finished, 5: Failed
	var status string
	var progress int
	var hlsURL string

	switch videoInfo.Status {
	case 0, 1, 2: // Queued, Processing, Encoding
		status = "processing"
		progress = videoInfo.EncodeProgress
		// Re-queue for next check
		if err := w.redisService.EnqueueVideoEncoding(ctx, job.MediaID, job.VideoID); err != nil {
			log.Printf("❌ Failed to re-queue job: %v", err)
		}
	case 3, 4: // Finished, Resolution Finished
		status = "completed"
		progress = 100
		hlsURL = w.bunnyStream.GetHLSURL(job.VideoID)
	case 5: // Failed
		status = "failed"
		progress = 0
	default:
		log.Printf("⚠️  Unknown video status: %d", videoInfo.Status)
		status = "processing"
		progress = videoInfo.EncodeProgress
	}

	// Update Redis
	w.redisService.UpdateVideoEncodingStatus(ctx, job.MediaID, status, progress, "")

	// Update database
	w.updateMediaStatus(ctx, job.MediaID, status, progress, hlsURL)

	// Send notification if completed or failed
	if status == "completed" || status == "failed" {
		w.sendEncodingNotification(ctx, job, status)
	}

	log.Printf("✓ Job processed: mediaID=%s, status=%s, progress=%d%%", job.MediaID, status, progress)
}

// updateMediaStatus updates the media record in the database
func (w *VideoEncoderWorker) updateMediaStatus(ctx context.Context, mediaID uuid.UUID, status string, progress int, hlsURL string) {
	media, err := w.mediaRepo.GetByID(ctx, mediaID)
	if err != nil {
		log.Printf("❌ Failed to get media: %v", err)
		return
	}

	media.EncodingStatus = status
	media.EncodingProgress = progress
	if hlsURL != "" {
		media.HLSURL = hlsURL
	}

	if err := w.mediaRepo.Update(ctx, media); err != nil {
		log.Printf("❌ Failed to update media: %v", err)
	}
}

// sendEncodingNotification sends a notification to the user when encoding completes
func (w *VideoEncoderWorker) sendEncodingNotification(ctx context.Context, job *redis.VideoEncodingJob, status string) {
	// Get media to find user ID
	media, err := w.mediaRepo.GetByID(ctx, job.MediaID)
	if err != nil {
		log.Printf("❌ Failed to get media for notification: %v", err)
		return
	}

	// Send WebSocket notification
	var message string

	if status == "completed" {
		message = "Your video has been processed and is ready to view!"
	} else {
		message = "Video processing failed. Please try uploading again."
	}

	websocket.Manager.BroadcastToUser(media.UserID, "video_encoding", map[string]interface{}{
		"media_id": job.MediaID,
		"video_id": job.VideoID,
		"status":   status,
		"message":  message,
	})

	log.Printf("📢 Notification sent: user=%s, status=%s", media.UserID, status)
}
