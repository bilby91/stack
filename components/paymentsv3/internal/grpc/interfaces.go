package grpc

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/grpc/proto/services"
	"github.com/formancehq/paymentsv3/internal/plugins/models"
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
	services.RegisterPluginServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *PSPGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: services.NewPluginClient(c)}, nil
}

var PluginMap = map[string]plugin.Plugin{
	"psp": &PSPGRPCPlugin{},
}

var _ plugin.GRPCPlugin = &PSPGRPCPlugin{}

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "GRPC_PLUGIN",
	MagicCookieValue: "test",
}
