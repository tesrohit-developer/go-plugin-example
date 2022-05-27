package main

import (
	"fmt"
	"github.com/dkiser/go-plugin-example/plugin"
	"net/rpc"

	gplugin "github.com/hashicorp/go-plugin"
)

type SidelineEm struct{}

func (SidelineEm) CheckMessageSideline(byte interface{}) (bool, error) {
	fmt.Println("Checking message")
	return true, nil
}

func (SidelineEm) SidelineMessage(msg interface{}) error {
	// do nothing
	fmt.Println("Sidelining message")
	return nil
}

type SidelineEmPlugin struct{}

func (SidelineEmPlugin) Server(*gplugin.MuxBroker) (interface{}, error) {
	return &plugin.CheckMessageSidelineRPCServer{Impl: new(SidelineEm)}, nil
}

func (SidelineEmPlugin) Client(b *gplugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &plugin.CheckMessageSidelineRPC{Client: c}, nil
}

func main() {
	// We're a plugin! Serve the plugin. We set the handshake config
	// so that the host and our plugin can verify they can talk to each other.
	// Then we set the plugin map to say what plugins we're serving.
	gplugin.Serve(&gplugin.ServeConfig{
		HandshakeConfig: gplugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "BASIC_PLUGIN",
			MagicCookieValue: "hello",
		},
		Plugins: pluginMap,
	})
}

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]gplugin.Plugin{
	"em": new(SidelineEmPlugin),
}
