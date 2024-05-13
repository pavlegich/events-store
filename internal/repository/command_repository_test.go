package repository

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
)

func TestNewEventRepository(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx context.Context
		db  *sql.DB
	}
	tests := []struct {
		name string
		args args
		want *EventRepository
	}{
		{
			name: "ok",
			args: args{
				ctx: ctx,
				db:  nil,
			},
			want: &EventRepository{
				db: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewEventRepository(tt.args.ctx, tt.args.db)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCommandRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}
