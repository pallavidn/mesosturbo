package master

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
	"github.com/turbonomic/mesosturbo/pkg/util"
)

// Parser interface for different server messages
type EndpointParser interface {
	parse(resp *http.Response) error
	GetMessage() interface{}
}

// ========================================= State Request Parser ===================================================

type GenericMasterStateParser struct {
	//Message *conf.MesosState
	Message *util.MesosAPIResponse
}

const GenericMasterStateParserClass = "[GenericMasterStateParser]"

func (parser *GenericMasterStateParser) parse(resp *http.Response) error {
	glog.V(2).Infof("[MesosStateParser] in parseAPIStateResponse")
	if resp == nil {
		return ErrorEmptyResponse(GenericMasterStateParserClass) //errors.New("Response sent from mesos/DCOS master is nil")
	}

	// Get token if response if OK
	if resp.Status == "" {
		return errors.New("[MesosStateParser] Empty response status\n")
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("[MesosStateParser] Error in ioutil.ReadAll: " + err.Error())
	}

	//fmt.Println("response content is %s", string(content))
	byteContent := []byte(content)
	var jsonMesosMaster util.MesosAPIResponse //conf.MesosState
	err = json.Unmarshal(byteContent, &jsonMesosMaster)
	if err != nil {
		return errors.New("[MesosStateParser] Error in json unmarshal for state response : " + err.Error())
	}
	parser.Message = &jsonMesosMaster
	return nil
}

func (parser *GenericMasterStateParser) GetMessage() interface{} {
	glog.V(2).Infof("[MesosStateParser] Mesos State %s\n", parser.Message)
	return parser.Message
}
