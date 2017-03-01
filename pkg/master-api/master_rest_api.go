package master

import (
	"bytes"
	"github.com/golang/glog"
	"net/http"

	"github.com/turbonomic/mesosturbo/pkg/conf"
	"github.com/turbonomic/mesosturbo/pkg/util"
)

// Represents the generic client used to connect to a Mesos Master. Implements the MasterRestClient interface
type GenericMasterAPIClient struct {
	// Mesos target configuration
	MesosConf     *conf.MesosTargetConf
	// Endpoint store with the endpoint paths for different rest api calls
	EndpointStore *MasterEndpointStore
}

type MasterEndpointName string

const (
	Login      MasterEndpointName = "login"
	State      MasterEndpointName = "state"
	Frameworks MasterEndpointName = "frameworks"
	Tasks      MasterEndpointName = "tasks"
)

// The endpoints used for making RestAPI calls to the Mesos Master
type MasterEndpoint struct {
	EndpointName string
	EndpointPath string
	Parser       EndpointParser
}

// Store containing the Rest API endpoints for communicating with the Mesos Master
type MasterEndpointStore struct {
	EndpointMap map[MasterEndpointName]*MasterEndpoint
}

// Create a new instance of the GenericMasterAPIClient
// @param mesosConf the conf.MesosTargetConf that contains the configuration information for the Mesos Target
// @param epStore    the Endpoint store containing the Rest API endpoints for the Mesos Master
func NewGenericMasterAPIClient(mesosConf *conf.MesosTargetConf, epStore *MasterEndpointStore) conf.MasterRestClient {
	return &GenericMasterAPIClient{
		MesosConf:     mesosConf,
		EndpointStore: epStore,
	}
}

const MesosMasterAPIClientClass = "MesosMasterAPIClient"

// Handle login to the Mesos Client using the path specified for the MasterEndpointName.Login endpoint
func (mesosRestClient *GenericMasterAPIClient) Login() (string, error) {
	glog.V(2).Infof("[GenericMasterAPIClient] Login ...")
	endpoint, _ := mesosRestClient.EndpointStore.EndpointMap[Login]

	if endpoint == nil {
		return "", nil
	}
	request, err := mesosRestClient.createLoginRequest(endpoint.EndpointPath)

	if err != nil {
		return "", ErrorCreateRequest(MesosMasterAPIClientClass, err)
	}
	glog.Infof("[GenericMasterAPIClient] Send Request to ", request)

	client := &http.Client{}
	resp, err := client.Do(request)

	if err != nil {
		return "", ErrorExecuteRequest(MesosMasterAPIClientClass, err)
	}

	defer resp.Body.Close()

	parser := endpoint.Parser
	err = parser.parse(resp)

	if err != nil {
		return "", ErrorParseRequest(MesosMasterAPIClientClass, err)
	}

	msg := parser.GetMessage()
	st, ok := msg.(string)
	if ok {
		return st, nil
	}
	return "", ErrorConvertResponse(MesosMasterAPIClientClass, err)
}

// Make a RestAPI call to get the Mesos State using the path specified for the MasterEndpointName.Login endpoint
func (mesosRestClient *GenericMasterAPIClient) GetState() (*util.MesosAPIResponse, error) {
	glog.V(2).Infof("[GenericMasterAPIClient] Get State ...")
	endpoint, _ := mesosRestClient.EndpointStore.EndpointMap[State]
	request, err := mesosRestClient.createRequest(endpoint.EndpointPath)
	if err != nil {
		return nil, ErrorCreateRequest(MesosMasterAPIClientClass, err)
	}
	glog.Infof("[GenericMasterAPIClient] Send Request to ", request)

	client := &http.Client{}
	resp, err := client.Do(request)

	if err != nil {
		return nil, ErrorExecuteRequest(MesosMasterAPIClientClass, err)
	}

	defer resp.Body.Close()

	parser := endpoint.Parser
	err = parser.parse(resp)

	if err != nil {
		return nil, ErrorParseRequest(MesosMasterAPIClientClass, err)
	}

	msg := parser.GetMessage()
	st, ok := msg.(*util.MesosAPIResponse)
	if ok {
		return st, nil
	}
	return nil, ErrorConvertResponse(MesosMasterAPIClientClass, err)
}

func (mesosRestClient *GenericMasterAPIClient) createRequest(endpoint string) (*http.Request, error) {
	fullUrl := "http://" + mesosRestClient.MesosConf.MasterIP + ":" + mesosRestClient.MesosConf.MasterPort + endpoint
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-type", "application/json")
	if mesosRestClient.MesosConf.Token != "" {
		req.Header.Add("Authorization", "token="+mesosRestClient.MesosConf.Token)
	}
	glog.V(2).Infof("[GenericMasterAPIClient] Created Request %+v", req)
	return req, nil
}

func (mesosRestClient *GenericMasterAPIClient) createLoginRequest(endpoint string) (*http.Request, error) {
	var jsonStr []byte
	url := "http://" + mesosRestClient.MesosConf.MasterIP + endpoint //"/acs/api/v1/auth/login"

	// Send user and password
	jsonStr = []byte(`{"uid":"` + mesosRestClient.MesosConf.MasterUsername + `","password":"` + mesosRestClient.MesosConf.MasterPassword + `"}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	glog.V(2).Infof("[GenericMasterAPIClient] Created Request %+v", req)
	return req, nil
}

func (mesosRestClient *GenericMasterAPIClient) GetNodes() {

}

func (mesosRestClient *GenericMasterAPIClient) GetFrameworks() {

}
