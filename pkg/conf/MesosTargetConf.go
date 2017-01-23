package conf

import (
	"github.com/golang/glog"
	"io/ioutil"
	"os"
	"encoding/json"
	"fmt"
)

// Configuration Parameters to connect to a Mesos Target
type MesosTargetConf struct {
	DCOS               bool
	DCOS_Username      string
	DCOS_Password      string
	MarathonIP         string
	MarathonPort       string
	MesosIP            string
	MesosPort          string
	ActionIP           string
	ActionPort         string
	ActionAPI          string
	SlavePort          string
	Token              string
}

// Create a new MesosClientConf from file. Other fields have default values and can be overrided.
func NewMesosTargetConf(targetConfigFilePath string) (*MesosTargetConf, error) {

	fmt.Println("[MesosClientConf] Target configuration from %s", targetConfigFilePath)
	metaConfig := readConfig(targetConfigFilePath)

	if metaConfig.DCOS_Username != "" && metaConfig.DCOS_Password != "" {
		metaConfig.DCOS = true
	}

	//if metaConfig.ActionIP != "" {
	//	meta.ActionIP = metaConfig.ActionIP
	//} else {
	//	glog.V(4).Infof("Error getting LayerX Master\n")
	//	return nil, errors.New("Error getting LayerX Master.")
	//}
	//
	//if metaConfig.ActionPort != "" {
	//	meta.ActionPort = metaConfig.ActionPort
	//} else {
	//	glog.V(4).Infof("Error getting LayerX Master.\n")
	//	return nil, errors.New("error getting LayerX Master\n")
	//}

	return metaConfig, nil
}

// Get the config from file.
func readConfig(path string) *MesosTargetConf {
	file, e := ioutil.ReadFile(path)
	if e != nil {
		glog.Errorf("File error: %v\n", e)
		os.Exit(1)
	}
	fmt.Println(string(file))

	var config MesosTargetConf
	err := json.Unmarshal(file, &config)

	if err != nil {
		fmt.Printf("[MesosTargetConf] Unmarshall error :%v\n", err)
	}
	fmt.Printf("[MesosTargetConf] Results: %+v\n", config)

	return &config
}
