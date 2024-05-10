// Package repository contains repository object
// and methods for interaction with storage.
package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/pavlegich/events-store/internal/entities"
)

// Repository describes methods related with events
// for interaction with database.
//
//go:generate mockgen -destination=../mocks/mock_Repository.go -package=mocks github.com/pavlegich/events-store/internal/repository Repository
type Repository interface {
	CreateEvent(ctx context.Context, event *entities.Event) (*entities.Event, error)
	GetEventsByFilter(ctx context.Context, eventType string, start time.Time, end time.Time) ([]*entities.Event, error)
}

// EventRepository contains storage objects for storing the events.
type EventRepository struct {
	db *sql.DB
}

// NewEventRepository returns new events repository object.
func NewEventRepository(ctx context.Context, db *sql.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

// CreateEvent stores new event into the storage.
func (r *EventRepository) CreateEvent(ctx context.Context, e *entities.Event) (*entities.Event, error) {
	return &entities.Event{}, nil
}

// GetEventsByFilter gets and returns requested events by filter.
func (r *EventRepository) GetEventsByFilter(ctx context.Context, eventType string, start time.Time, end time.Time) ([]*entities.Event, error) {
	return []*entities.Event{}, nil
}
