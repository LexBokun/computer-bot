package setenable

import (
	"context"
)

type Command struct {
	ID     string
	Enable bool
}

type Result struct {
}

type CommandHandler struct {
	displayRepository DisplayRepository
}

func NewCommandHandler(r DisplayRepository) *CommandHandler {
	return &CommandHandler{displayRepository: r}
}

func (h *CommandHandler) Handle(ctx context.Context, cmd Command) error {
	return h.displayRepository.SetEnabled(ctx, cmd.ID, cmd.Enable)
}
