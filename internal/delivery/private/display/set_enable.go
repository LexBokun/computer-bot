package display

import (
	"context"

	setenable "github.com/LexBokun/ControlAgent/internal/application/service/command/set-enable"
	displayv1 "github.com/LexBokun/ControlAgent/pkg/api/agent/private/display/v1"
)

func (i *Implementation) SetEnabled(ctx context.Context, request *displayv1.SetEnabledRequest) (*displayv1.SetEnabledResponse, error) {
	err := i.setEnable.Handle(ctx, setenable.Command{
		ID:     request.GetId(),
		Enable: request.GetEnable(),
	})
	if err != nil {
		return nil, err
	}

	return &displayv1.SetEnabledResponse{}, nil
}
