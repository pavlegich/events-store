package event

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pavlegich/events-store/internal/entities"
	errs "github.com/pavlegich/events-store/internal/errors"
	"github.com/pavlegich/events-store/internal/mocks"
	repo "github.com/pavlegich/events-store/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestNewCommandService(t *testing.T) {
	ctx := context.Background()

	type args struct {
		repo repo.Repository
	}
	tests := []struct {
		name string
		args args
		want *EventService
	}{
		{
			name: "ok",
			args: args{
				repo: nil,
			},
			want: &EventService{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewEventService(ctx, tt.args.repo)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCommandService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventService_Create(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	mockRepo := mocks.NewMockRepository(mockCtrl)
	s := NewEventService(ctx, mockRepo)

	type expected struct {
		err error
	}
	type args struct {
		event *entities.Event
	}
	tests := []struct {
		name     string
		args     args
		expected expected
		wantErr  error
	}{
		{
			name: "success",
			args: args{
				event: &entities.Event{
					EventType: "login",
					UserID:    1,
					EventTime: entities.CustomTime{
						Time: time.Now(),
					},
					Payload: `{"some_field": "some_value"}`,
				},
			},
			expected: expected{
				err: nil,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().CreateEvent(gomock.Any(), gomock.Any()).
				Return(tt.expected.err).Times(1)

			err := s.Create(ctx, tt.args.event)

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestEventService_Unload(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	mockRepo := mocks.NewMockRepository(mockCtrl)
	s := NewEventService(ctx, mockRepo)

	currTime := time.Now()

	type expected struct {
		events []*entities.Event
		err    error
	}
	type args struct {
		eventType string
		startTime time.Time
		endTime   time.Time
	}
	tests := []struct {
		name     string
		args     args
		expected expected
		wantErr  error
		want     []*entities.Event
	}{
		{
			name: "success",
			args: args{
				eventType: "login",
				startTime: currTime.Add(-time.Duration(1) * time.Hour),
				endTime:   currTime,
			},
			expected: expected{
				events: []*entities.Event{
					{
						EventType: "login",
						UserID:    1,
						EventTime: entities.CustomTime{
							Time: currTime.Add(-time.Duration(40) * time.Minute),
						},
						Payload: `{"some_field": "some_value"}`,
					},
					{
						EventType: "login",
						UserID:    2,
						EventTime: entities.CustomTime{
							Time: currTime.Add(-time.Duration(30) * time.Minute),
						},
						Payload: `{"some_field": "some_value"}`,
					},
				},
				err: nil,
			},
			wantErr: nil,
			want: []*entities.Event{
				{
					EventType: "login",
					UserID:    1,
					EventTime: entities.CustomTime{
						Time: currTime.Add(-time.Duration(40) * time.Minute),
					},
					Payload: `{"some_field": "some_value"}`,
				},
				{
					EventType: "login",
					UserID:    2,
					EventTime: entities.CustomTime{
						Time: currTime.Add(-time.Duration(30) * time.Minute),
					},
					Payload: `{"some_field": "some_value"}`,
				},
			},
		},
		{
			name: "events_not_found",
			args: args{
				eventType: "login",
				startTime: currTime.Add(-time.Duration(1) * time.Hour),
				endTime:   currTime,
			},
			expected: expected{
				events: nil,
				err:    errs.ErrEventNotFound,
			},
			wantErr: errs.ErrEventNotFound,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().GetEventsByFilter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(tt.expected.events, tt.expected.err).Times(1)

			got, err := s.Unload(ctx, tt.args.eventType, tt.args.startTime, tt.args.endTime)

			require.ErrorIs(t, err, tt.wantErr)
			require.ElementsMatch(t, got, tt.want)
		})
	}
}
