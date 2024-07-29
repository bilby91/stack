package grpc

import (
	"context"
	"plugin"

	"github.com/formancehq/stack/components/paymentsv3/internal/plugins/models"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type PSPGRPCPlugin struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl models.Plugin
}

func (p *PSPGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterPSPServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *PSPGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewPSPClient(c)}, nil
}

var PluginMap = map[string]plugin.Plugin{
	"psp": &PSPGRPCPlugin{},
}

var _ plugin.GRPCPlugin = &PSPGRPCPlugin{}
