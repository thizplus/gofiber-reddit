package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	videoEncodingQueueKey = "video:encoding:queue"
	videoEncodingPrefix   = "video:encoding:"
)

// VideoEncodingJob represents a video encoding job in the queue
type VideoEncodingJob struct {
	MediaID     uuid.UUID `json:"media_id"`
	VideoID     string    `json:"video_id"`      // Bunny Stream video ID
	Status      string    `json:"status"`        // pending, processing, completed, failed
	Progress    int       `json:"progress"`      // 0-100
	QueuedAt    time.Time `json:"queued_at"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Error       string    `json:"error,omitempty"`
}

// EnqueueVideoEncoding adds a video to the encoding queue
func (r *RedisService) EnqueueVideoEncoding(ctx context.Context, mediaID uuid.UUID, videoID string) error {
	job := VideoEncodingJob{
		MediaID:  mediaID,
		VideoID:  videoID,
		Status:   "pending",
		Progress: 0,
		QueuedAt: time.Now(),
	}

	jobData, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	// Store job data in Redis hash
	jobKey := fmt.Sprintf("%s%s", videoEncodingPrefix, mediaID.String())
	if err := r.client.Set(ctx, jobKey, jobData, 24*time.Hour).Err(); err != nil {
		return fmt.Errorf("failed to store job: %w", err)
	}

	// Add to processing queue (FIFO)
	if err := r.client.RPush(ctx, videoEncodingQueueKey, mediaID.String()).Err(); err != nil {
		return fmt.Errorf("failed to enqueue job: %w", err)
	}

	return nil
}

// DequeueVideoEncoding retrieves and removes the next video from the queue
func (r *RedisService) DequeueVideoEncoding(ctx context.Context) (*VideoEncodingJob, error) {
	// Pop from the left (FIFO)
	result, err := r.client.LPop(ctx, videoEncodingQueueKey).Result()
	if err != nil {
		return nil, err // redis.Nil means queue is empty
	}

	mediaID, err := uuid.Parse(result)
	if err != nil {
		return nil, fmt.Errorf("invalid media ID in queue: %w", err)
	}

	// Retrieve job data
	jobKey := fmt.Sprintf("%s%s", videoEncodingPrefix, mediaID.String())
	jobData, err := r.client.Get(ctx, jobKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get job data: %w", err)
	}

	var job VideoEncodingJob
	if err := json.Unmarshal([]byte(jobData), &job); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}

	// Mark as processing
	now := time.Now()
	job.Status = "processing"
	job.StartedAt = &now

	// Update job in Redis
	updatedData, _ := json.Marshal(job)
	r.client.Set(ctx, jobKey, updatedData, 24*time.Hour)

	return &job, nil
}

// UpdateVideoEncodingStatus updates the status and progress of a video encoding job
func (r *RedisService) UpdateVideoEncodingStatus(ctx context.Context, mediaID uuid.UUID, status string, progress int, errorMsg string) error {
	jobKey := fmt.Sprintf("%s%s", videoEncodingPrefix, mediaID.String())

	// Get existing job
	jobData, err := r.client.Get(ctx, jobKey).Result()
	if err != nil {
		return fmt.Errorf("job not found: %w", err)
	}

	var job VideoEncodingJob
	if err := json.Unmarshal([]byte(jobData), &job); err != nil {
		return fmt.Errorf("failed to unmarshal job: %w", err)
	}

	// Update job
	job.Status = status
	job.Progress = progress
	if errorMsg != "" {
		job.Error = errorMsg
	}

	if status == "completed" || status == "failed" {
		now := time.Now()
		job.CompletedAt = &now
	}

	// Save updated job
	updatedData, _ := json.Marshal(job)
	ttl := 24 * time.Hour
	if status == "completed" || status == "failed" {
		ttl = 1 * time.Hour // Keep completed/failed jobs for 1 hour
	}

	if err := r.client.Set(ctx, jobKey, updatedData, ttl).Err(); err != nil {
		return fmt.Errorf("failed to update job: %w", err)
	}

	return nil
}

// GetVideoEncodingJob retrieves a video encoding job by media ID
func (r *RedisService) GetVideoEncodingJob(ctx context.Context, mediaID uuid.UUID) (*VideoEncodingJob, error) {
	jobKey := fmt.Sprintf("%s%s", videoEncodingPrefix, mediaID.String())

	jobData, err := r.client.Get(ctx, jobKey).Result()
	if err != nil {
		return nil, err
	}

	var job VideoEncodingJob
	if err := json.Unmarshal([]byte(jobData), &job); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}

	return &job, nil
}

// GetPendingEncodingCount returns the number of videos waiting in the queue
func (r *RedisService) GetPendingEncodingCount(ctx context.Context) (int64, error) {
	return r.client.LLen(ctx, videoEncodingQueueKey).Result()
}

// ClearVideoEncodingJob removes a job from Redis (for cleanup)
func (r *RedisService) ClearVideoEncodingJob(ctx context.Context, mediaID uuid.UUID) error {
	jobKey := fmt.Sprintf("%s%s", videoEncodingPrefix, mediaID.String())
	return r.client.Del(ctx, jobKey).Err()
}
