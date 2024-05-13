package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/pavlegich/events-store/internal/entities"
	errs "github.com/pavlegich/events-store/internal/errors"
	"github.com/pavlegich/events-store/internal/infra/config"
	"github.com/pavlegich/events-store/internal/infra/logger"
	"github.com/pavlegich/events-store/internal/repository"
	"github.com/pavlegich/events-store/internal/service/event"
	"go.uber.org/zap"
)

// EventHandler contains objects for work with command handlers.
type EventHandler struct {
	Config  *config.Config
	Service event.Service
}

type eventFilter struct {
	eventType string
	startTime time.Time
	endTime   time.Time
}

// eventsActivate activates handler for command object.
func eventsActivate(ctx context.Context, r *http.ServeMux, repo repository.Repository, cfg *config.Config) {
	s := event.NewEventService(ctx, repo)
	newHandler(r, cfg, s)
}

// newHandler initializes handler for command object.
func newHandler(r *http.ServeMux, cfg *config.Config, s event.Service) {
	h := &EventHandler{
		Config:  cfg,
		Service: s,
	}

	r.HandleFunc("/api/event", h.HandleEvent)
}

// HandleEvent handles request to interact with event.
func (h *EventHandler) HandleEvent(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.HandleCreateEvent(w, r)
	case http.MethodGet:
		h.HandleGetEvent(w, r)
	default:
		logger.Log.Error("HandleEvent: incorrect method",
			zap.String("method", r.Method))

		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// HandleCreateEvent handles request to create new event.
func (h *EventHandler) HandleCreateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req entities.Event
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		logger.Log.Error("HandleCreateEvent: read request body failed",
			zap.Error(err))

		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		logger.Log.Error("HandleCreateEvent: request unmarshal failed",
			zap.String("body", buf.String()),
			zap.Error(err))

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.EventType == "" || req.UserID == 0 {
		logger.Log.Error("HandleCreateEvent: event type or user id are empty",
			zap.String("eventType", req.EventType),
			zap.Int64("userID", req.UserID),
			zap.Error(err))

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.Service.Create(ctx, &req)
	if err != nil {
		logger.Log.Error("HandleCreateEvent: create event failed",
			zap.Error(err))

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// HandleGetEvent handles request to get events by filter.
func (h *EventHandler) HandleGetEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req eventFilter
	want := map[string]struct{}{
		"eventType": {},
		"startTime": {},
		"endTime":   {},
	}

	queries := r.URL.Query()
	for val := range queries {
		_, ok := want[val]
		if !ok {
			logger.Log.Error("HandleGetEvent: incorrect query",
				zap.String("query", val))

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(queries[val]) != 1 {
			logger.Log.Error("HandleGetEvent: incorrect number of queries",
				zap.Int("queries_count", len(queries)))

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		switch val {
		case "eventType":
			req.eventType = queries[val][0]
		case "startTime":
			sTime, err := time.Parse(entities.Layout, queries[val][0])
			if err != nil {
				logger.Log.Error("HandleGetEvent: incorrect time format",
					zap.String("query", queries[val][0]),
					zap.String("layout", entities.Layout))

				w.WriteHeader(http.StatusBadRequest)
				return
			}
			req.startTime = sTime
		case "endTime":
			eTime, err := time.Parse(entities.Layout, queries[val][0])
			if err != nil {
				logger.Log.Error("HandleGetEvent: incorrect time format",
					zap.String("query", queries[val][0]),
					zap.String("layout", entities.Layout))

				w.WriteHeader(http.StatusBadRequest)
				return
			}
			req.endTime = eTime
		}
	}

	events, err := h.Service.Unload(ctx, req.eventType, req.startTime, req.endTime)
	if err != nil {
		logger.Log.Error("HandleGetEvent: unload events failed",
			zap.Error(err))

		if errors.Is(err, errs.ErrEventNotFound) {
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}
