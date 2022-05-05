package main

import (
	"net/rpc"

	"github.com/go-dmux/plugins"
	"github.com/go-dmux/sideline"
	gplugin "github.com/hashicorp/go-plugin"
)

type SidelineEm struct{}

func (SidelineEm) CheckMessageSideline(byte interface{}) bool {
	return true
}

func (SidelineEm) SidelineMessage(msg interface{}) {
	// do nothing
}

type SidelineEmPlugin struct{}

func (SidelineEmPlugin) Server(*gplugin.MuxBroker) (interface{}, error) {
	return &sideline.CheckMessageSidelineRPCServer{Impl: new(SidelineEm)}, nil
}

func (SidelineEmPlugin) Client(b *gplugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &sideline.CheckMessageSidelineRPC{Client: c}, nil
}

func main() {
	// We're a plugin! Serve the plugin. We set the handshake config
	// so that the host and our plugin can verify they can talk to each other.
	// Then we set the plugin map to say what plugins we're serving.
	gplugin.Serve(&gplugin.ServeConfig{
		HandshakeConfig: plugins.HandshakeConfig,
		Plugins:         pluginMap,
	})
}

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]gplugin.Plugin{
	"em": new(SidelineEmPlugin),
}
