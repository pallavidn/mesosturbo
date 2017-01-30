package master

import (
	"github.com/turbonomic/mesosturbo/pkg/conf"
	"fmt"
	"net/http"
	"errors"
)

type ApacheMesosRestClient struct {
	MesosMasterIP string
	MesosMasterPort string
	Username string
	Password string

	LoginToken string
	EndpointStore *ApacheMesosEndpointStore
}


func NewApacheMesosRestClient (masterIP, masterPort, username, password string) conf.MasterRestClient {
	return &ApacheMesosRestClient {
		MesosMasterIP: masterIP,
		MesosMasterPort: masterPort,
		Username: username,
		Password: password,
		EndpointStore: NewApacheMesosEndpointStore(),
	}
}

func (mesosRestClient *ApacheMesosRestClient) Login() (string, error) {
	return "", nil
}

func (mesosRestClient *ApacheMesosRestClient) GetState() (*conf.MesosState, error) {
	fmt.Println("[ApacheMesosRestClient] Get State ...")
	endpoint, _ := mesosRestClient.EndpointStore.EndpointMap[State]
	request, err := mesosRestClient.createRequest(endpoint)
	if err != nil {
		fmt.Println("[ApacheMesosRestClient] Error in GetState request: %s\n", err)
		return nil, err
	}
	fmt.Println("[ApacheMesosRestClient] Send Request to " , request)

	client := &http.Client{}
	resp, err := client.Do(request)

	if err != nil {
		fmt.Println("[ApacheMesosRestClient] Error in GetState request to mesos master: %s\n", err)
		return nil, err
	}

	defer resp.Body.Close()

	parser := &MesosStateParser{}
	err = parser.parse(resp)

	if err != nil {
		fmt.Println("ApacheMesosRestClient] Error parsing Mesos master state response ", err)
		return nil, err
	}

	msg := parser.GetMessage()
	st, ok :=  msg.(*conf.MesosState)
	if ok {
		return st, nil
	}
	return nil, errors.New("[ApacheMesosRestClient] Error converting response to MesosState")
}

func (mesosRestClient *ApacheMesosRestClient) createRequest(endpoint ApacheMesosEndpointPath) (*http.Request, error) {
	fullUrl := "http://" + mesosRestClient.MesosMasterIP + ":" + mesosRestClient.MesosMasterPort + string(endpoint)
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("[ApacheMesosRestClient] Created Request %+v", req)
	return req, nil
}

func (mesosRestClient *ApacheMesosRestClient) GetNodes() {

}

func (mesosRestClient *ApacheMesosRestClient) GetFrameworks() {

}