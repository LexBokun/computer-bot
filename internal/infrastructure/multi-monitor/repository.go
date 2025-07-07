package multimonitor

import (
	"context"

	"github.com/LexBokun/ControlAgent/internal/domain/display"
)

type repository struct {
	monitorTool *Config
}

func (r *repository) SetEnabled(ctx context.Context, id string, enabled bool) error {
	return r.monitorTool.SetEnabled(ctx, id, enabled)
}

func (r *repository) List(ctx context.Context) ([]display.Display, error) {
	return r.monitorTool.List(ctx)
}

func NewRepository(monitorTool *Config) display.ControlDisplay {
	return &repository{monitorTool: monitorTool}
}
