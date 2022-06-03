package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dkiser/go-plugin-example/plugin"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"time"

	gplugin "github.com/hashicorp/go-plugin"
	emclientmodels "github.fkinternal.com/Flipkart/entity-manager/modules/entity-manager-client-model/EntityManagerClientModel"
	emmodels "github.fkinternal.com/Flipkart/entity-manager/modules/entity-manager-model/EntityManagerModel"
	"google.golang.org/protobuf/proto"
)

type SidelineEm struct{}

func execute(method, url string, headers map[string]string,
	payload io.Reader) (bool, int, emclientmodels.ResponseCode, string) {
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
		return false, 1, -1, ""
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
		return false, 2, -1, ""
	}
	//TODO check if this can be avoided
	io.Copy(ioutil.Discard, response.Body)
	responseBytes, _ := ioutil.ReadAll(response.Body)
	var readResponse emclientmodels.ReadEntityResponse
	proto.Unmarshal(responseBytes, &readResponse)
	defer response.Body.Close()
	if emclientmodels.ResponseStatus_STATUS_SUCCESS.Number() == readResponse.ResponseMeta.ResponseStatus.Number() {
		return true, response.StatusCode, readResponse.ResponseMeta.ResponseCode, readResponse.String()
	}
	fmt.Println(readResponse.String())
	return false, response.StatusCode, readResponse.ResponseMeta.ResponseCode, readResponse.String()
}

func (SidelineEm) CheckMessageSideline(byte string) (bool, error) {
	fmt.Println("Checking message in EM")
	url := "http://10.24.19.136/entity-manager/v1/entity/read"
	headers := make(map[string]string)
	headers["Content-Type"] = "application/octet-stream"
	headers["X-IDEMPOTENCY-ID"] = time.Now().String()
	headers["X-CLIENT-ID"] = "go-dmux"
	//headers["X-PERF-TTL"] = "LONG_PERF"
	entityIdentifier := emmodels.EntityIdentifier{
		Namespace: "com.dmux",
		Name:      "SidelineMessage",
	}
	tenantIdentifier := emmodels.TenantIdentifier{
		Name: "OMSDMUX",
	}
	readEntityRequest := emclientmodels.ReadEntityRequest{
		EntityIdentifier: &entityIdentifier,
		TenantIdentifier: &tenantIdentifier,
		EntityId:         "OD39848785211959690",
		FieldsToRead:     nil,
	}
	b, e := proto.Marshal(&readEntityRequest)
	if e != nil {
		fmt.Println("error in ser ReadEntityRequest")
		return false, errors.New("error in ser ReadEntityRequest")
	}
	responseBoolean, responseCode, emResponseCode, readResponseString := execute("POST", url, headers, bytes.NewReader(b))
	if responseCode < 300 {
		fmt.Println("Success ")
		return true, nil
	}
	
	if !responseBoolean {
		if emclientmodels.ResponseCode_ENTITY_NOT_FOUND.Number() == emResponseCode.Number() {
			fmt.Println("Not sidelined message ")
			return true, nil
		}
		fmt.Println("error in reading Sideline Table")
		errStr := "error in reading Sideline Table, ResponseCode: " + strconv.Itoa(responseCode) +
			" EmResponseCode: " + emResponseCode.String() +
			" ResponseBoolean: " + strconv.FormatBool(responseBoolean) +
			" ReadResponseString: " + readResponseString
		return false, errors.New(errStr)
	}

	return false, errors.New("error in reading Sideline Table")
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
