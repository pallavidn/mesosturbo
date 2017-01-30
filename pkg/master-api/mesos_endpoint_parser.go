package master

import (
	"net/http"
	"errors"
	"io/ioutil"
	"encoding/json"
	"fmt"

	"github.com/turbonomic/mesosturbo/pkg/conf"
)

type MesosEndpointName string
const (
	State MesosEndpointName = "state"
	Frameworks MesosEndpointName = "frameworks"
	Tasks MesosEndpointName = "tasks"
)

// Parser interface for different server messages
type EndpointParser interface {
	parse(resp *http.Response) error
	GetMessage() interface{}
}

// ============================================================================================================

type MesosStateParser struct{
	Message *conf.MesosState
}

func (parser *MesosStateParser) parse(resp *http.Response) error {
	fmt.Println("[MesosStateParser] in parseAPIStateResponse")
	if resp == nil {
		return errors.New("Response sent from mesos/DCOS master is nil")
	}

	// Get token if response if OK
	if resp.Status == "" {
		fmt.Println("Empty response status\n")
		return errors.New("Empty response status\n")
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error after ioutil.ReadAll: %s", err)
		return err
	}

	fmt.Println("response content is %s", string(content))
	byteContent := []byte(content)
	var jsonMesosMaster conf.MesosState	//util.MesosAPIResponse)
	err = json.Unmarshal(byteContent, &jsonMesosMaster)
	if err != nil {
		fmt.Println("error in json unmarshal : %s", err)
		return errors.New("Error in json unmarshal")
	}
	parser.Message = &jsonMesosMaster
	return  nil
}

func (parser *MesosStateParser) GetMessage() interface{}  {
	fmt.Printf("[MesosStateParser] Mesos State %s\n", parser.Message)
	return parser.Message
}
