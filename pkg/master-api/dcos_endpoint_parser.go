package master

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"github.com/turbonomic/mesosturbo/pkg/util"
	"io/ioutil"
	"net/http"
)

// ============================================ DCOS Login Parser ======================================================

type DCOSLoginParser struct {
	Message string
}

const DCOSLoginParserClass = "[DCOSLoginParser]"

func (parser *DCOSLoginParser) parse(resp *http.Response) error {
	glog.V(2).Infof("[DCOSLoginParser] in parseAPIStateResponse")
	if resp == nil {
		return ErrorEmptyResponse(DCOSLoginParserClass) //errors.New("Response sent from mesos/DCOS master is nil")
	}

	// Get token if response if OK
	if resp.Status != "200" {
		return errors.New("[DCOSLoginParser] Invalid response status " + fmt.Sprintf("%s", resp))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("[DCOSLoginParser] Error in ioutil.ReadAll: " + err.Error())
	}
	byteContent := []byte(body)
	var tokenResp = new(util.TokenResponse)
	err = json.Unmarshal(byteContent, &tokenResp)
	if err != nil {
		return errors.New("[DCOSLoginParser] error in json unmarshalling login response :" + err.Error())
	}
	parser.Message = tokenResp.Token
	return nil
	//return errors.New("[DCOSLoginParser] DCOS authorization credentials are not correct : " + string(content)))

}

func (parser *DCOSLoginParser) GetMessage() interface{} {
	glog.V(2).Infof("[DCOSLoginParser] DCOS Login Token %s\n", parser.Message)
	return parser.Message
}
