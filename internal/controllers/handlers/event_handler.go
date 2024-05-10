package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

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

	eventID, err := h.Service.Create(ctx, &req)
	if err != nil {
		logger.Log.Error("HandleCreateCommand: create event failed",
			zap.Error(err))

		if errors.Is(err, errs.ErrEventAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"command_id": eventID})
}
