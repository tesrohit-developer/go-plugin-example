package main

import (
	"github.com/dkiser/go-plugin-example/plgin"
	"net/rpc"

	plugin "github.com/dkiser/go-plugin-example/plugin"
	gplugin "github.com/hashicorp/go-plugin"
)

var fist = `
     .~~~~'\~~\
     ;       ~~ \
     |           ;
 ,--------,______|---.
/          \-----'    \
'.__________'-_______-'
`

// Here is a real implementation of Dubber
type DubberMilkboy struct{}

func (DubberMilkboy) FistPump() string { return fist }

// This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a DubberRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return DubberRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type DubberPlugin struct{}

func (DubberPlugin) Server(*gplugin.MuxBroker) (interface{}, error) {
	return &plgin.DubberRPCServer{Impl: new(DubberMilkboy)}, nil
}

func (DubberPlugin) Client(b *gplugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &plgin.DubberRPC{Client: c}, nil
}

func main() {
	// We're a plugin! Serve the plugin. We set the handshake config
	// so that the host and our plugin can verify they can talk to each other.
	// Then we set the plugin map to say what plugins we're serving.
	gplugin.Serve(&gplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig,
		Plugins:         pluginMap,
	})
}

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]gplugin.Plugin{
	"milkboy": new(DubberPlugin),
}
