package plgin

import "net/rpc"
import gplugin "github.com/hashicorp/go-plugin"

// Dubber is the interface that we're exposing as a plugin.
type Dubber interface {
	FistPump(s string) (string, error)
}

// Here is an implementation that talks over RPC
type DubberRPC struct {
	Client *rpc.Client
}

func (g *DubberRPC) FistPump(s string) string {
	var resp string
	err := g.Client.Call("Plugin.FistPump", s, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// Here is the RPC server that DubberRPC talks to, conforming to
// the requirements of net/rpc
type DubberRPCServer struct {
	// This is the real implementation
	Impl Dubber
}

func (s *DubberRPCServer) FistPump(st string, resp *string) string {
	*resp, _ = s.Impl.FistPump(st)
	return *resp
}

// Dummy implementation of a plugin.Plugin interface for use in PluginMap.
// At runtime, a real implementation from a plugin implementation overwrides
// this.
type DubberPlugin struct{}

func (DubberPlugin) Server(*gplugin.MuxBroker) (interface{}, error) {
	return &DubberRPCServer{}, nil
}

func (DubberPlugin) Client(b *gplugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &DubberRPC{Client: c}, nil
}
