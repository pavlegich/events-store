package handlers

import (
	"context"
	"reflect"
	"testing"

	"github.com/pavlegich/events-store/internal/infra/config"
)

func TestNewController(t *testing.T) {
	ctx := context.Background()
	cfg := config.NewConfig(ctx)
	want := &Controller{cfg: cfg}

	t.Run("success", func(t *testing.T) {
		got := NewController(ctx, cfg)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NewController() = %v, want %v", got, want)
		}
	})
}
