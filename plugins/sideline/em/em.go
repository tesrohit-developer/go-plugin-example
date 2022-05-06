package pluginst

import (
	"net/rpc"

	gplugin "github.com/hashicorp/go-plugin"
	"github.com/tesrohit-developer/go-dmux/plugins"
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
	return &plugins.CheckMessageSidelineRPCServer{Impl: new(SidelineEm)}, nil
}

func (SidelineEmPlugin) Client(b *gplugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &plugins.CheckMessageSidelineRPC{Client: c}, nil
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
