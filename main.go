package main

import (
	"github.com/dkiser/go-plugin-example/plgin"
	"github.com/dkiser/go-plugin-example/plugin"
	"log"
)

func playWithGreeters() {
	// Grab a new manager for dealing with "greeter" type plugins
	greeters := plugin.NewManager("greeter", "greeter-*", "./plugins/built", &plugin.GreeterPlugin{})
	defer greeters.Dispose()

	// Initialize greeters manager
	err := greeters.Init()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Launch all greeters binaries
	greeters.Launch()

	// Lets see what all the greeters say when you say "sup" to them
	for _, pluginName := range []string{"foo", "hello"} {
		// grab a plugin by its string id
		p, err := greeters.GetInterface(pluginName)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Printf("\n\n%s: %s plugin gives me: %s\n\n", greeters.Type, pluginName, p.(plugin.Greeter).Greet())
	}

}

func playWithClubbers() {
	// Grab a new manager for dealing with "clubber" type plugins
	clubbers := plugin.NewManager("clubber", "clubber-*", "./plugins/built", &plugin.ClubberPlugin{})
	defer clubbers.Dispose()

	// Initialize clubbers manager
	err := clubbers.Init()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Launch all clubbbers binaries
	clubbers.Launch()

	// Lets see what all the clubbers do when they party hardy!
	for _, pluginName := range []string{"raver", "milkboy"} {
		// grab a plugin by its string id
		p, err := clubbers.GetInterface(pluginName)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Printf("\n%s: %s plugin gives me: %s\n", clubbers.Type, pluginName, p.(plugin.Clubber).FistPump())
	}
}

func playWithDubbers() {
	// Grab a new manager for dealing with "clubber" type plugins
	dubbers := plugin.NewManager("dubber", "dubber-*", "./plugins/built", &plgin.DubberPlugin{})
	defer dubbers.Dispose()

	// Initialize clubbers manager
	err := dubbers.Init()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Launch all clubbbers binaries
	dubbers.Launch()

	// Lets see what all the clubbers do when they party hardy!
	for _, pluginName := range []string{"milkboy"} {
		// grab a plugin by its string id
		p, err := dubbers.GetInterface(pluginName)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Printf("\n%s: %s plugin gives me: %s\n", dubbers.Type, pluginName, p.(plgin.Dubber).FistPump("asd"))
	}
}

func playWithSidelinePlugin() {
	s := plugin.NewManager("sideline_plugin", "sideline-*", "./plugins/built", &plugin.CheckMessageSidelineImplPlugin{})
	defer s.Dispose()

	err := s.Init()

	if err != nil {
		log.Fatal(err.Error())
	}

	s.Launch()

	p, err := s.GetInterface("em")
	if err != nil {
		log.Fatal(err.Error())
	}

	ch := make([]chan interface{}, 10)
	ch[0] = make(chan interface{}, 10)
	//p.(plugin.CheckMessageSidelineImpl).SidelineMessage(new(interface{}))
	log.Printf("\n%s: %s plugin gives me: %s\n", s.Type, "sideline-em", p.(plugin.CheckMessageSidelineImpl).CheckMessageSideline(new(interface{})))
	for msg := range ch[0] {
		log.Printf("\n%s: %s plugin gives me: %s\n", s.Type, "sideline-em", p.(plugin.CheckMessageSidelineImpl).CheckMessageSideline(msg))
	}
}

/*func playWithDmuxSidelinePlugin() {
	s := plugins.NewManager("sideline_plugin", "sideline-*", "./plugins/built", &plugins.CheckMessageSidelineImplPlugin{})
	defer s.Dispose()

	err := s.Init()

	if err != nil {
		log.Fatal(err.Error())
	}

	s.Launch()

	p, err := s.GetInterface("em")
	if err != nil {
		log.Fatal(err.Error())
	}

	ch := make([]chan interface{}, 10)
	ch[0] = make(chan interface{}, 10)
	//p.(plugin.CheckMessageSidelineImpl).SidelineMessage(new(interface{}))
	log.Printf("\n%s: %s plugin gives me: %s\n", s.Type, "sideline-em", p.(plugins.CheckMessageSidelineImpl).CheckMessageSideline(new(interface{})))
	/*for msg := range ch[0] {
		log.Printf("\n%s: %s plugin gives me: %s\n", s.Type, "sideline-em", p.(sideline.CheckMessageSidelineImpl).CheckMessageSideline(msg))
	}
}*/

func main() {

	// excercise some greeter plugins
	//playWithGreeters()

	// excercise some clubber plugins
	//playWithClubbers()

	playWithDubbers()
	//playWithSidelinePlugin()

	//playWithDmuxSidelinePlugin()

}
