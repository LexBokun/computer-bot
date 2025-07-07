package display

import "context"

type ControlDisplay interface {
	List(ctx context.Context) ([]Display, error)
	SetEnabled(ctx context.Context, id string, enable bool) error
}

type Display struct {
	ID     string
	Name   string
	Enable bool
}
