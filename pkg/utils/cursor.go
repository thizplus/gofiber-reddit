package utils

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

// Cursor represents a pagination cursor
type Cursor struct {
	Timestamp time.Time `json:"timestamp"`
	ID        string    `json:"id,omitempty"`
}

// EncodeCursor encodes a timestamp into a base64 cursor string
func EncodeCursor(timestamp time.Time) (string, error) {
	cursor := Cursor{
		Timestamp: timestamp,
	}

	jsonBytes, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(jsonBytes), nil
}

// DecodeCursor decodes a base64 cursor string into a timestamp
func DecodeCursor(cursorStr string) (*time.Time, error) {
	if cursorStr == "" {
		return nil, nil
	}

	jsonBytes, err := base64.StdEncoding.DecodeString(cursorStr)
	if err != nil {
		return nil, err
	}

	var cursor Cursor
	if err := json.Unmarshal(jsonBytes, &cursor); err != nil {
		return nil, err
	}

	return &cursor.Timestamp, nil
}
