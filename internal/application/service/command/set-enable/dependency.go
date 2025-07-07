package setenable

import (
	"context"
)

type DisplayRepository interface {
	SetEnabled(ctx context.Context, id string, enable bool) error
}
