package conf

import (
	"io/ioutil"
	"os"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

// Configuration Parameters to connect to a Mesos Target
type MesosTargetConf struct {
	// Master IP - Apache or DCOS Mesos
	Master         MesosMasterType  `json:"master"`
	MasterIP       string		`json:"master-ip"`
	MasterPort     string		`json:"master-port"`
	MasterUsername string		`json:"master-user"`
	MasterPassword string		`json:"master-pwd"`

	//DCOS               bool
	//DCOS_Username      string
	//DCOS_Password      string
	//MesosIP            string
	//MesosPort          string

	// Scheduler IP
	Framework      	   MesosFrameworkType   `json:"framework"`
	FrameworkIP        string		`json:"framework-ip"`
	FrameworkPort      string		`json:"framework-port"`
	FrameworkUser      string		`json:"framework-user"`
	FrameworkPassword  string	        `json:"framework-pwd"`
	//MarathonIP         string
	//MarathonPort       string

	// Action Executor IP
	ActionIP       string
	ActionPort     string
	ActionAPI      string

	// Others ?
	SlavePort      string
	Token          string
}

// Create a new MesosClientConf from file. Other fields have default values and can be overrided.
func NewMesosTargetConf(targetConfigFilePath string) (*MesosTargetConf, error) {

	fmt.Println("[MesosClientConf] Target configuration from %s", targetConfigFilePath)
	metaConfig := readConfig(targetConfigFilePath)

	// TODO: validate conf parameters

	return metaConfig, nil
}

func NewMesosTargetConfFromAccountValues([]*proto.AccountValue) (*MesosTargetConf, error) {
	// TODO:
	return nil, nil
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
