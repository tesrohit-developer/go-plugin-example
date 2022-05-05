package plugin

import (
	gplugin "github.com/hashicorp/go-plugin"
	"net/rpc"
)

// CheckMessageSidelineImpl is the interface that we're exposing as a plugin.
type CheckMessageSidelineImpl interface {
	CheckMessageSideline(key interface{}) bool
	SidelineMessage(msg interface{})
}

// Here is an implementation that talks over RPC
type CheckMessageSidelineRPC struct {
	Client *rpc.Client
}

func (g *CheckMessageSidelineRPC) CheckMessageSideline(key interface{}) bool {
	var resp bool
	err := g.Client.Call("Plugin.CheckMessageSideline", key, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return false
}

func (g *CheckMessageSidelineRPC) SidelineMessage(msg interface{}) {
	var resp bool
	err := g.Client.Call("Plugin.SidelineMessage", msg, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}
}

// Here is the RPC server that CheckMessageSidelineRPC talks to, conforming to
// the requirements of net/rpc
type CheckMessageSidelineRPCServer struct {
	// This is the real implementation
	Impl CheckMessageSidelineImpl
}

func (s *CheckMessageSidelineRPCServer) CheckMessageSideline(args interface{}, resp *bool) error {
	b := []byte("asd")
	*resp = s.Impl.CheckMessageSideline(b)
	return nil
}

func (s *CheckMessageSidelineRPCServer) SidelineMessage(args interface{}, resp *bool) error {
	s.Impl.SidelineMessage(args)
	return nil
}

// Dummy implementation of a plugin.Plugin interface for use in PluginMap.
// At runtime, a real implementation from a plugin implementation overwrides
// this.
type CheckMessageSidelineImplPlugin struct{}

func (CheckMessageSidelineImplPlugin) Server(*gplugin.MuxBroker) (interface{}, error) {
	return &CheckMessageSidelineRPCServer{}, nil
	//return interface{}, nil
}

func (CheckMessageSidelineImplPlugin) Client(b *gplugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &CheckMessageSidelineRPC{Client: c}, nil
	//return interface{}, nil
}
