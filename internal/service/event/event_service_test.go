package event

import (
	"context"
	"reflect"
	"testing"

	repo "github.com/pavlegich/events-store/internal/repository"
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
