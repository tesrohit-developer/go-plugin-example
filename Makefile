GREETER_PLUGIN_DIRS=$(wildcard ./plugins/greeter/*)
CLUBBER_PLUGIN_DIRS=$(wildcard ./plugins/clubber/*)
DUBBER_PLUGIN_DIRS=$(wildcard ./plugins/dubber/*)
#SIDELINE_PLUGIN_DIRS=$(wildcard ./plugins/sideline/*)

all: build-plugins

clean: clean-plugins

build-plugins: $(GREETER_PLUGIN_DIRS) $(CLUBBER_PLUGIN_DIRS) $(DUBBER_PLUGIN_DIRS)

clean-plugins: 
	rm -f ./plugins/built/*

$(GREETER_PLUGIN_DIRS): 
	$(info Greeter plugins at: $(GREETER_PLUGIN_DIRS))
	$(MAKE) -C $@

$(CLUBBER_PLUGIN_DIRS): 
	$(info Clubber plugins at: $(CLUBBER_PLUGIN_DIRS))
	$(MAKE) -C $@

$(DUBBER_PLUGIN_DIRS):
	$(info Dubber plugins at: $(DUBBER_PLUGIN_DIRS))
	$(MAKE) -C $@

#$(SIDELINE_PLUGIN_DIRS):
#	$(info Sideline plugins at: $(SIDELINE_PLUGIN_DIRS))
#	$(MAKE) -C $@

#.PHONY: all $(GREETER_PLUGIN_DIRS) $(CLUBBER_PLUGIN_DIRS) $(SIDELINE_PLUGIN_DIRS)
.PHONY: all $(GREETER_PLUGIN_DIRS) $(CLUBBER_PLUGIN_DIRS) $(DUBBER_PLUGIN_DIRS)