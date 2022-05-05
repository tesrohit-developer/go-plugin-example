package plugin

import (
	SidelinePlugin "github.com/go-dmux/sideline"
	"testing"
)

func TestManagerInit(t *testing.T) {
	manager := NewManager("greeter", "greeter-*", "../plugins/built", GreeterPlugin{})
	defer manager.Dispose()
	err := manager.Init()
	if err != nil {
		t.Fatal(err.Error())
	}
	manager.Launch()
	for id, info := range manager.Plugins {
		t.Log(id)
		t.Log(info.Path)
		t.Log(info.Client)
	}

	foo, err := manager.GetInterface("foo")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(foo.(Greeter).Greet())

	hello, err := manager.GetInterface("hello")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(hello.(Greeter).Greet())
}

func TestSidelinePluginInit(t *testing.T) {
	manager := NewManager("sideline_plugin", "sideline-*", "../plugins/built", &SidelinePlugin.CheckMessageSidelineImplPlugin{})
	defer manager.Dispose()
	err := manager.Init()
	if err != nil {
		t.Fatal(err.Error())
	}
	manager.Launch()
	for id, info := range manager.Plugins {
		t.Log(id)
		t.Log(info.Path)
		t.Log(info.Client)
	}

	foo, err := manager.GetInterface("em")
	if err != nil {
		t.Fatal(err.Error())
	}
	b := []byte("asd")
	t.Log(foo.(SidelinePlugin.CheckMessageSidelineImpl).CheckMessageSideline(b))
}
