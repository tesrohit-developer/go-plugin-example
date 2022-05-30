package main

import (
	"fmt"
	"github.com/dkiser/go-plugin-example/plugin"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	gplugin "github.com/hashicorp/go-plugin"
)

type SidelineEm struct{}

func execute(method, url string, headers map[string]string,
	payload io.Reader) (bool, int) {
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          2, // TODO
			MaxIdleConnsPerHost:   2, // TODO
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: 10 * time.Second, // TODO
	}

	//Never fail always recover
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recovered in execute %s %v", url, r)
		}
	}()

	//build request
	request, err := http.NewRequest(method, url, payload)
	if err != nil {
		log.Printf("failed in request build %s %s \n", url, err.Error())
		return false, 0
	}

	//set headers
	for key, val := range headers {
		request.Header.Set(key, val)
	}

	// if method != "GET" && !h.conf.CustomURL {
	// 	request.Header.Set("Content-Type", "application/octet-stream")
	// }

	//make request
	response, err := client.Do(request)
	if err != nil {
		log.Printf("failed in http call invoke %s %s \n", url, err.Error())
		return false, 0
	}
	//TODO check if this can be avoided
	io.Copy(ioutil.Discard, response.Body)
	defer response.Body.Close()

	return true, response.StatusCode
}

func (SidelineEm) CheckMessageSideline(byte interface{}) (bool, error) {
	fmt.Println("Checking message in EM")
	url := "http://10.47.101.183/entity-manager/v1/entity/read"
	headers := make(map[string]string)
	headers["Content-Type"] = "application/octet-stream"
	headers["X-IDEMPOTENCY-ID"] = time.Now().String()
	headers["X-CLIENT-ID"] = "go-dmux"
	//headers["X-PERF-TTL"] = "LONG_PERF"
	responseBoolean, responseCode := execute("POST", url, headers, nil)
	fmt.Println(responseCode)
	return responseBoolean, nil
}

func (SidelineEm) SidelineMessage(msg interface{}) error {
	// do nothing
	fmt.Println("Sidelining message in EM")
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
