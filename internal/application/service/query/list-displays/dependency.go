package listdisplays

import (
	"context"

	"github.com/LexBokun/ControlAgent/internal/domain/display"
)

type DisplayRepository interface {
	List(ctx context.Context) ([]display.Display, error)
}
