package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

// BunnyStreamService handles video upload and encoding with Bunny Stream
type BunnyStreamService struct {
	apiKey    string
	libraryID string
	cdnURL    string
	client    *http.Client
}

// CreateVideoResponse represents the response from Bunny Stream CreateVideo API
type CreateVideoResponse struct {
	VideoID string `json:"guid"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	// Full response fields
	VideoLibraryID      int64  `json:"videoLibraryId"`
	Title               string `json:"title"`
	DateUploaded        string `json:"dateUploaded"`
	Views               int    `json:"views"`
	IsPublic            bool   `json:"isPublic"`
	Length              int    `json:"length"`
	Status              int    `json:"status"`
	Framerate           int    `json:"framerate"`
	Width               int    `json:"width"`
	Height              int    `json:"height"`
	AvailableResolutions string `json:"availableResolutions"`
	ThumbnailCount      int    `json:"thumbnailCount"`
	EncodeProgress      int    `json:"encodeProgress"`
	StorageSize         int64  `json:"storageSize"`
	HasMP4Fallback      bool   `json:"hasMP4Fallback"`
}

// GetVideoResponse represents the response from Bunny Stream GetVideo API
type GetVideoResponse struct {
	VideoID              string `json:"guid"`
	VideoLibraryID       int64  `json:"videoLibraryId"`
	Title                string `json:"title"`
	DateUploaded         string `json:"dateUploaded"`
	Views                int    `json:"views"`
	IsPublic             bool   `json:"isPublic"`
	Length               int    `json:"length"`
	Status               int    `json:"status"` // 0: Queued, 1: Processing, 2: Encoding, 3: Finished, 4: Resolution Finished, 5: Failed
	Framerate            int    `json:"framerate"`
	Width                int    `json:"width"`
	Height               int    `json:"height"`
	AvailableResolutions string `json:"availableResolutions"`
	ThumbnailCount       int    `json:"thumbnailCount"`
	EncodeProgress       int    `json:"encodeProgress"`
	StorageSize          int64  `json:"storageSize"`
	HasMP4Fallback       bool   `json:"hasMP4Fallback"`
}

func NewBunnyStreamService(apiKey, libraryID, cdnURL string) *BunnyStreamService {
	return &BunnyStreamService{
		apiKey:    apiKey,
		libraryID: libraryID,
		cdnURL:    cdnURL,
		client: &http.Client{
			Timeout: 300 * time.Second, // 5 minutes for large uploads
		},
	}
}

// CreateVideo creates a new video in Bunny Stream and uploads the file
func (s *BunnyStreamService) CreateVideo(file multipart.File, filename string) (*CreateVideoResponse, error) {
	// Step 1: Create video metadata
	createURL := fmt.Sprintf("https://video.bunnycdn.com/library/%s/videos", s.libraryID)

	createBody := map[string]interface{}{
		"title": filename,
	}
	createData, _ := json.Marshal(createBody)

	req, err := http.NewRequest("POST", createURL, bytes.NewReader(createData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("AccessKey", s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create video metadata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bunny API error (status %d): %s", resp.StatusCode, string(body))
	}

	var createResp CreateVideoResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Step 2: Upload video file
	uploadURL := fmt.Sprintf("https://video.bunnycdn.com/library/%s/videos/%s", s.libraryID, createResp.VideoID)

	uploadReq, err := http.NewRequest("PUT", uploadURL, file)
	if err != nil {
		return nil, fmt.Errorf("failed to create upload request: %w", err)
	}

	uploadReq.Header.Set("AccessKey", s.apiKey)
	uploadReq.Header.Set("Content-Type", "application/octet-stream")

	uploadResp, err := s.client.Do(uploadReq)
	if err != nil {
		return nil, fmt.Errorf("failed to upload video: %w", err)
	}
	defer uploadResp.Body.Close()

	if uploadResp.StatusCode != http.StatusOK && uploadResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(uploadResp.Body)
		return nil, fmt.Errorf("upload failed (status %d): %s", uploadResp.StatusCode, string(body))
	}

	return &createResp, nil
}

// GetVideo retrieves video information from Bunny Stream
func (s *BunnyStreamService) GetVideo(videoID string) (*GetVideoResponse, error) {
	url := fmt.Sprintf("https://video.bunnycdn.com/library/%s/videos/%s", s.libraryID, videoID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("AccessKey", s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get video: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bunny API error (status %d): %s", resp.StatusCode, string(body))
	}

	var videoResp GetVideoResponse
	if err := json.NewDecoder(resp.Body).Decode(&videoResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &videoResp, nil
}

// DeleteVideo deletes a video from Bunny Stream
func (s *BunnyStreamService) DeleteVideo(videoID string) error {
	url := fmt.Sprintf("https://video.bunnycdn.com/library/%s/videos/%s", s.libraryID, videoID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("AccessKey", s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete video: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete failed (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetHLSURL generates the HLS playlist URL for a video
func (s *BunnyStreamService) GetHLSURL(videoID string) string {
	return fmt.Sprintf("%s/%s/playlist.m3u8", s.cdnURL, videoID)
}

// GetThumbnailURL generates thumbnail URL for a video
func (s *BunnyStreamService) GetThumbnailURL(videoID string) string {
	return fmt.Sprintf("%s/%s/thumbnail.jpg", s.cdnURL, videoID)
}
