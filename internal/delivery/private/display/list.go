package display

import (
	"context"

	listDisplays "github.com/LexBokun/ControlAgent/internal/application/service/query/list-displays"
	displayv1 "github.com/LexBokun/ControlAgent/pkg/api/agent/private/display/v1"
)

func (i *Implementation) List(ctx context.Context, req *displayv1.ListRequest) (*displayv1.ListResponse, error) {
	result, err := i.list.Handle(ctx, listDisplays.Query{})
	if err != nil {
		return nil, err
	}

	var displays []*displayv1.Display
	for _, d := range result.Displays {
		displays = append(displays, &displayv1.Display{
			Id:     d.ID,
			Name:   d.Name,
			Enable: d.Enable,
		})
	}

	return &displayv1.ListResponse{Displays: displays}, nil
}
