package master

import (
	"bytes"
	"encoding/json"
	"github.com/golang/glog"
	"github.com/turbonomic/mesosturbo/pkg/util"
	conf "github.com/turbonomic/mesosturbo/pkg/conf"
	"io/ioutil"
	"net/http"
	"fmt"
	"errors"
)


type DCOSMesosRestClient struct {
	MesosMasterIP string
	MesosMasterPort string
	Username string
	Password string

	LoginToken string
	EndpointStore *DCOSMesosEndpointStore
}

func NewDCOSMesosRestClient(masterIP, masterPort, username, password string) conf.MasterRestClient {
	return &DCOSMesosRestClient {
		MesosMasterIP: masterIP,
		MesosMasterPort: masterPort,
		Username: username,
		Password: password,
		EndpointStore: NewDCOSMesosEndpointStore(),
	}
}

func (dcosRestClient *DCOSMesosRestClient) Login() (string, error) {

	var jsonStr []byte
	url := "http://" + dcosRestClient.MesosMasterIP  + "/acs/api/v1/auth/login"

	jsonStr = []byte(`{"uid":"` + dcosRestClient.Username + `","password":"` + dcosRestClient.Password + `"}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	if err != nil {
		fmt.Println("[DCOSMesosRestClient] Error in Login request: %s \n", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("[DCOSMesosRestClient] Error in Login request: %s \n", err)
		return "", err
	} else {
		// Get token if response if OK
		defer resp.Body.Close()
		if resp.Status == "" {
			fmt.Println("[DCOSMesosRestClient] Empty response status \n")
			return "", errors.New("Empty response status \n")
		}

		glog.Infof(" Status is : %s \n", resp.Status)

		if resp.StatusCode == 200 {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error in ioutil.ReadAll: %s \n", err)
				return "", err
			}
			byteContent := []byte(body)
			var tokenResp = new(util.TokenResponse)
			err = json.Unmarshal(byteContent, &tokenResp)
			if err != nil {
				fmt.Println("[DCOSMesosRestClient] error in json unmarshal : %s . \r\nLogin failed , please try again with correct credentials.\n", err)
				return "", err
			}
			dcosRestClient.LoginToken = tokenResp.Token
			return tokenResp.Token, nil
		} else {
			fmt.Println("[DCOSMesosRestClient] Please check DCOS credentials and start mesosturbo again.\n")
			return "", errors.New("DCOS authorization credentials are not correct, check mesosturbo arguments --dcos-uid , --dcos-pwd, or --token! \n")
		}

	}
}

func (mesosRestClient *DCOSMesosRestClient) GetState() (*conf.MesosState, error) {
	fmt.Println("[DCOSMesosRestClient] Get State ...")
	endpoint, _ := mesosRestClient.EndpointStore.EndpointMap[State]
	request, err := mesosRestClient.createRequest(endpoint)
	if err != nil {
		fmt.Println("[DCOSMesosRestClient] Error in GetState request: %s\n", err)
		return nil, err
	}
	fmt.Println("[DCOSMesosRestClient] Send Request to " , request)
	client := &http.Client{}
	resp, err := client.Do(request)

	if err != nil {
		fmt.Println("[DCOSMesosRestClient] Error in GetState request to mesos master: %s\n", err)
		return nil, err
	}

	defer resp.Body.Close()

	parser := &MesosStateParser{}
	err = parser.parse(resp)

	if err != nil {
		fmt.Println("DCOSMesosRestClient] Error parsing Mesos master state response ", err)
		return nil, err
	}

	msg := parser.GetMessage()
	st, ok :=  msg.(*conf.MesosState)
	if ok {
		return st, nil
	}
	return nil, errors.New("[DCOSMesosRestClient] Error converting response to MesosState")
}


func (mesosRestClient *DCOSMesosRestClient) createRequest(endpoint DCOSEndpointPath) (*http.Request, error) {
	var fullUrl string
	if mesosRestClient.MesosMasterPort == "" {
		fullUrl = "http://" + mesosRestClient.MesosMasterIP + string(endpoint)
	} else {
		fullUrl = "http://" + mesosRestClient.MesosMasterIP + ":" + mesosRestClient.MesosMasterPort + string(endpoint)
	}
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Authorization", "token="+mesosRestClient.LoginToken)
	fmt.Println("[DCOSMesosRestClient] Created Request %+v", req)
	return req, nil
}

func (mesosRestClient *DCOSMesosRestClient) GetNodes() {


}

func (mesosRestClient *DCOSMesosRestClient) GetFrameworks() {
}

