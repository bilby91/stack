package main

import (
	"github.com/formancehq/paymentsv3/internal/grpc"
	"github.com/formancehq/paymentsv3/internal/plugins/public/moneycorp"
	"github.com/hashicorp/go-plugin"
)

func main() {
	// TODO(logger, metrics etc...)
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: grpc.Handshake,
		Plugins: map[string]plugin.Plugin{
			"psp": &grpc.PSPGRPCPlugin{Impl: &moneycorp.Plugin{}},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
