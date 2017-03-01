package framework

import (
	"net/http"
	"errors"
	"io/ioutil"
	"encoding/json"
	"github.com/turbonomic/mesosturbo/pkg/util"
	"github.com/golang/glog"
)


// Parser interface for different server messages
type EndpointParser interface {
	parse(resp *http.Response) error
	GetMessage() interface{}
}


// ============================================================================================================

type FrameworkAppsParser struct{
	Message *util.FrameworkApps	//MarathonApps
}

func (parser *FrameworkAppsParser) parse(resp *http.Response) error {
	glog.V(2).Infof("[FrameworkAppsParser] in parse Frameworks Apps Response")
	if resp == nil {
		return errors.New("Response sent from framework master is nil")
	}

	// Get token if response if OK
	if resp.Status == "" {
		//fmt.Println("Empty response status\n")
		return errors.New("Empty response status\n")
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("[FrameworkAppsParser] Error in ioutil.ReadAll: " + err.Error())
	}

	//fmt.Println("[FrameworkAppsParser] response content is %s", string(content))
	byteContent := []byte(content)
	var jsonMarathonMaster = new(util.FrameworkApps)
	err = json.Unmarshal(byteContent, &jsonMarathonMaster)
	if err != nil {
		return errors.New("[FrameworkAppsParser] Error in json unmarshal for state response : " + err.Error())
	}
	for i, app := range jsonMarathonMaster.Apps {
		//fmt.Println("[FrameworkAppsParser] Id : ", app.Id)
		newN := app.Id[1:len(app.Id)]
		jsonMarathonMaster.Apps[i].Id = newN
	}

	parser.Message = jsonMarathonMaster
	return nil
}

func (parser *FrameworkAppsParser) GetMessage() interface{}  {
	glog.V(2).Infof("[FrameworkAppsParser] Framework Apps %s\n", parser.Message)
	return parser.Message
}