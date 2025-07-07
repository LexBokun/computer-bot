package display

import (
	setenable "github.com/LexBokun/ControlAgent/internal/application/service/command/set-enable"
	listdisplays "github.com/LexBokun/ControlAgent/internal/application/service/query/list-displays"
	displayv1 "github.com/LexBokun/ControlAgent/pkg/api/agent/private/display/v1"

	"github.com/utrack/clay/v3/transport"
)

type Implementation struct {
	list      *listdisplays.QueryHandler
	setEnable *setenable.CommandHandler
}

func NewImplementation(
	list *listdisplays.QueryHandler,
	setEnable *setenable.CommandHandler,
) *Implementation {
	return &Implementation{
		list:      list,
		setEnable: setEnable,
	}
}

// GetDescription is a simple alias to the ServiceDesc constructor.
// It makes it possible to register the service implementation @ the server.
func (i *Implementation) GetDescription() transport.ServiceDesc {
	return displayv1.NewPrivateDisplayServiceServiceDesc(i)
}
