// Package repository contains repository object
// and methods for interaction with storage.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/pavlegich/events-store/internal/entities"
	errs "github.com/pavlegich/events-store/internal/errors"
)

// Repository describes methods related with events
// for interaction with database.
//
//go:generate mockgen -destination=../mocks/mock_Repository.go -package=mocks github.com/pavlegich/events-store/internal/repository Repository
type Repository interface {
	CreateEvent(ctx context.Context, event *entities.Event) error
	GetEventsByFilter(ctx context.Context, eventType string, startTime time.Time, endTime time.Time) ([]*entities.Event, error)
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
func (r *EventRepository) CreateEvent(ctx context.Context, e *entities.Event) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO events (eventType, userID, eventTime, payload) 
	VALUES (?, ?, ?, ?)`, e.EventType, e.UserID, e.EventTime.UTC(), e.Payload)

	if err != nil {
		return fmt.Errorf("CreateEvent: insert event failed %w", err)
	}

	return nil
}

// GetEventsByFilter gets and returns requested events by filter.
func (r *EventRepository) GetEventsByFilter(ctx context.Context, eventType string, startTime time.Time, endTime time.Time) ([]*entities.Event, error) {
	query := `SELECT eventID, eventType, userID, eventTime, payload FROM events`
	args := []string{}
	if eventType != "" {
		args = append(args, fmt.Sprintf("eventType = '%s'", eventType))
	}
	if !startTime.IsZero() {
		args = append(args, fmt.Sprintf("eventTime >= '%s'", startTime.UTC().Format(entities.Layout)))
	}
	if !endTime.IsZero() {
		args = append(args, fmt.Sprintf("eventTime <= '%s'", endTime.UTC().Format(entities.Layout)))
	}
	if len(args) != 0 {
		query = query + " WHERE " + strings.Join(args, " AND ")
	}

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("GetEventsByFilter: read rows from table failed %w", err)
	}
	defer rows.Close()

	eventsList := make([]*entities.Event, 0)
	for rows.Next() {
		var e entities.Event
		err = rows.Scan(&e.EventID, &e.EventType, &e.UserID, &e.EventTime, &e.Payload)
		if err != nil {
			return nil, fmt.Errorf("GetEventsByFilter: scan row failed %w", err)
		}
		eventsList = append(eventsList, &e)
	}

	if len(eventsList) == 0 {
		return nil, fmt.Errorf("GetEventsByFilter: nothing to return %w", errs.ErrEventNotFound)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("GetEventsByFilter: rows.Err %w", err)
	}

	return eventsList, nil
}
