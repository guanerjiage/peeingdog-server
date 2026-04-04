package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"peeingdog-server/sql/queries/generated"
)

type Message struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Text      string    `json:"text"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type CreateMessageRequest struct {
	Text      string  `json:"text"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type MessageService struct {
	queries *generated.Queries
}

func NewMessageService(q *generated.Queries) *MessageService {
	return &MessageService{queries: q}
}

// CreateMessage creates a new message
func (s *MessageService) CreateMessage(ctx context.Context, userID int, req CreateMessageRequest) (*Message, error) {
	if req.Text == "" {
		return nil, fmt.Errorf("message text is required")
	}

	if len(req.Text) > 180 {
		return nil, fmt.Errorf("message text exceeds maximum length of 180 characters")
	}

	if req.Latitude < -90 || req.Latitude > 90 {
		return nil, fmt.Errorf("invalid latitude")
	}

	if req.Longitude < -180 || req.Longitude > 180 {
		return nil, fmt.Errorf("invalid longitude")
	}

	msg, err := s.queries.CreateMessage(ctx, generated.CreateMessageParams{
		UserID:    int32(userID),
		Text:      req.Text,
		Latitude:  fmt.Sprintf("%f", req.Latitude),
		Longitude: fmt.Sprintf("%f", req.Longitude),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	lat, _ := ParseFloat(msg.Latitude)
	lon, _ := ParseFloat(msg.Longitude)

	return &Message{
		ID:        int(msg.ID),
		UserID:    int(msg.UserID),
		Text:      msg.Text,
		Latitude:  lat,
		Longitude: lon,
		CreatedAt: msg.CreatedAt,
		ExpiresAt: msg.ExpiresAt,
	}, nil
}

// GetNearbyMessages retrieves active messages within a radius
// radius is in kilometers
func (s *MessageService) GetNearbyMessages(ctx context.Context, latitude, longitude, radiusKm float64) ([]Message, error) {
	if radiusKm <= 0 {
		return nil, fmt.Errorf("radius must be positive")
	}

	// Simple distance calculation using degrees (rough approximation)
	// 1 degree ≈ 111 km
	latDelta := radiusKm / 111.0
	lonDelta := radiusKm / (111.0 * cosine(latitude))

	minLat := latitude - latDelta
	maxLat := latitude + latDelta
	minLon := longitude - lonDelta
	maxLon := longitude + lonDelta

	messages, err := s.queries.GetNearbyMessages(ctx, generated.GetNearbyMessagesParams{
		Latitude:  fmt.Sprintf("%f", minLat),
		Latitude_2: fmt.Sprintf("%f", maxLat),
		Longitude: fmt.Sprintf("%f", minLon),
		Longitude_2: fmt.Sprintf("%f", maxLon),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch nearby messages: %w", err)
	}

	result := make([]Message, len(messages))
	for i, msg := range messages {
		lat, _ := ParseFloat(msg.Latitude)
		lon, _ := ParseFloat(msg.Longitude)
		result[i] = Message{
			ID:        int(msg.ID),
			UserID:    int(msg.UserID),
			Text:      msg.Text,
			Latitude:  lat,
			Longitude: lon,
			CreatedAt: msg.CreatedAt.Time,
			ExpiresAt: msg.ExpiresAt.Time,
		}
	}
	return result, nil
}

// GetUserMessages retrieves all active messages from a specific user
func (s *MessageService) GetUserMessages(ctx context.Context, userID int) ([]Message, error) {
	messages, err := s.queries.GetUserMessages(ctx, int32(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user messages: %w", err)
	}

	result := make([]Message, len(messages))
	for i, msg := range messages {
		lat, _ := ParseFloat(msg.Latitude)
		lon, _ := ParseFloat(msg.Longitude)
		result[i] = Message{
			ID:        int(msg.ID),
			UserID:    int(msg.UserID),
			Text:      msg.Text,
			Latitude:  lat,
			Longitude: lon,
			CreatedAt: msg.CreatedAt.Time,
			ExpiresAt: msg.ExpiresAt.Time,
		}
	}
	return result, nil
}

// ArchiveExpiredMessages moves expired messages to archive table
func (s *MessageService) ArchiveExpiredMessages(ctx context.Context) error {
	err := s.queries.ArchiveExpiredMessages(ctx)
	if err != nil {
		return fmt.Errorf("failed to archive expired messages: %w", err)
	}
	return nil
}

// Helper functions
func ParseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

func cosine(latDegrees float64) float64 {
	const pi = 3.14159265359
	const deg2rad = pi / 180.0
	lat := latDegrees * deg2rad
	return math.Cos(lat)
}
