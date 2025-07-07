package listdisplays

import (
	"context"

	"github.com/LexBokun/ControlAgent/internal/domain/display"
)

type Query struct {
}

type Result struct {
	Displays []display.Display
}

type QueryHandler struct {
	displayRepository DisplayRepository
}

func NewQueryHandler(r DisplayRepository) *QueryHandler {
	return &QueryHandler{displayRepository: r}
}

func (h *QueryHandler) Handle(ctx context.Context, _ Query) (*Result, error) {
	displays, err := h.displayRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	return &Result{Displays: displays}, nil
}
