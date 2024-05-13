// Package event contains event service object and methods for interaction between
// handlers and repositories.
package event

import (
	"context"
	"fmt"
	"time"

	"github.com/pavlegich/events-store/internal/entities"
	repo "github.com/pavlegich/events-store/internal/repository"
)

// Service describes methods for communication between
// handlers and repositories.
//
//go:generate mockgen -destination=../../mocks/mock_Service.go -package=mocks github.com/pavlegich/events-store/internal/service/event Service
type Service interface {
	Create(ctx context.Context, event *entities.Event) error
	Unload(ctx context.Context, eventType string, startTime time.Time, endTime time.Time) ([]*entities.Event, error)
}

// EventService contains objects for event service.
type EventService struct {
	repo repo.Repository
}

// NewEventService returns new event service.
func NewEventService(ctx context.Context, repo repo.Repository) *EventService {
	return &EventService{
		repo: repo,
	}
}

// Create creates new requested event and requests repository to put it into the storage.
func (s *EventService) Create(ctx context.Context, c *entities.Event) error {
	err := s.repo.CreateEvent(ctx, c)
	if err != nil {
		return fmt.Errorf("Create: create event failed %w", err)
	}

	return nil
}

// Unload gets and returns events by event's filter.
func (s *EventService) Unload(ctx context.Context, eventType string, startTime time.Time, endTime time.Time) ([]*entities.Event, error) {
	events, err := s.repo.GetEventsByFilter(ctx, eventType, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("Unload: unload events failed %w", err)
	}

	return events, nil
}
