package framework

import (
	"github.com/turbonomic/mesosturbo/pkg/conf"
	"net/http"
	"errors"
	"github.com/turbonomic/mesosturbo/pkg/util"
	"github.com/golang/glog"
)


type FrameworkEndpointName string
const (
	Apps FrameworkEndpointName = "state"
)

type FrameworkAPIClient struct {
	MesosConf	*conf.MesosTargetConf
	EndpointStore *FrameworkEndpointStore
}

type FrameworkEndpoint struct {
	EndpointName string
	EndpointPath string
	Parser 	EndpointParser
}

type FrameworkEndpointStore struct {
	EndpointMap map[FrameworkEndpointName]*FrameworkEndpoint
}


func NewFrameworkAPIClient (mesosConf *conf.MesosTargetConf, epStore *FrameworkEndpointStore) conf.FrameworkRestClient {
	return &FrameworkAPIClient{
		MesosConf: mesosConf,
		EndpointStore: epStore,
	}
}

func (frameworkClient *FrameworkAPIClient) GetFrameworkApps() (*util.FrameworkApps, error) {
	config := frameworkClient.MesosConf
	glog.V(2).Infof("[MarathonRestClient] Get GetFrameworkApps ...%s", config)
	endpoint, _ := frameworkClient.EndpointStore.EndpointMap[Apps]
	request, err := frameworkClient.createRequest(config.FrameworkIP, config.FrameworkPort, endpoint.EndpointPath, config.Token)
	if err != nil {
		return nil, errors.New("[FrameworkAPIClient] Error creating GetFrameworkApps request: %s\n" + err.Error())
	}
	glog.Infof("[FrameworkAPIClient] Send Request to " , request)

	client := &http.Client{}
	resp, err := client.Do(request)

	if err != nil {
		return nil, errors.New("FrameworkAPIClient] Error executing GetFrameworkApps request: %s\n" + err.Error())
	}

	defer resp.Body.Close()

	parser := endpoint.Parser
	err = parser.parse(resp)

	if err != nil {
		return nil, errors.New("FrameworkAPIClient] Error parsing GetFrameworkApps response: %s\n" + err.Error())
	}

	msg := parser.GetMessage()
	st, ok :=  msg.(*util.FrameworkApps)
	if ok {
		return st, nil
	}
	return nil, errors.New("[FrameworkAPIClient] Error converting response to MarathonApps")
}

func (frameworkClient *FrameworkAPIClient) createRequest(ip, port, endpoint, loginToken string) (*http.Request, error) {
	fullUrl := "http://" + ip  + ":" + port + endpoint
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-type", "application/json")
	if loginToken != "" {
		req.Header.Add("Authorization", "token=" + loginToken)
	}
	glog.V(2).Infof("[FrameworkAPIClient] Created Request %+v", req)
	return req, nil
}