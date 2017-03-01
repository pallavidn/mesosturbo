package conf

import (
	"io/ioutil"
	"encoding/json"
	"github.com/golang/glog"
	"errors"
	"fmt"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

// Configuration Parameters to connect to a Mesos Target
type MesosTargetConf struct {
	// Master related - Apache or DCOS Mesos
	Master         MesosMasterType  `json:"master"`
	MasterIP       string		`json:"master-ip"`
	MasterPort     string		`json:"master-port"`
	MasterUsername string		`json:"master-user"`
	MasterPassword string		`json:"master-pwd"`

	// Scheduler or Framework related
	Framework      	   MesosFrameworkType   `json:"framework"`
	FrameworkIP        string		`json:"framework-ip"`
	FrameworkPort      string		`json:"framework-port"`
	FrameworkUser      string		`json:"framework-user"`
	FrameworkPassword  string	        `json:"framework-pwd"`

	// Action Executor related is using Layer-X
	ActionIP       string
	ActionPort     string
	ActionAPI      string

	// Others ?
	SlavePort      string
	// Login Token obtained from the Mesos Master
	Token          string
}

// Create a new MesosClientConf from a json file.
// Return null config if the there are errors loading or parsing the file
func NewMesosTargetConf(targetConfigFilePath string) (*MesosTargetConf, error) {

	glog.Infof("[MesosClientConf] Target configuration from %s", targetConfigFilePath)
	config, err := readConfig(targetConfigFilePath)
	if err != nil {
		return nil, errors.New("[[MesosTargetConf] Error reading config "+ targetConfigFilePath + " : " + err.Error())
	}
	if config == nil {
		return nil, errors.New("[[MesosTargetConf] Null config "+ targetConfigFilePath)
	}

	if (config.FrameworkIP == "") {
		config.FrameworkIP = config.MasterIP
		config.FrameworkPort = config.MasterPort
		config.Framework = DCOS_Marathon
	}
	ok, err := config.validate()
	if !ok {
		return nil, errors.New("[MesosTargetConf] Invalid config : " + err.Error() )
	}
	glog.Infof("[MesosTargetConf] Mesos Target Config: %+v\n", config)
	return config, nil
}
//
//func CreateMesosTargetConf(targetType string, accountValues []*proto.AccountValue) probe.TurboTargetConf {
//	var mesosMasterType MesosMasterType
//	if targetType == string(Apache) {
//		mesosMasterType = Apache
//	} else if targetType == string(DCOS) {
//		mesosMasterType = DCOS
//	} else {
//		glog.Errorf("Unknown Mesos Master Type " , targetType)
//		return nil
//	}
//	config := &MesosTargetConf{
//		Master: mesosMasterType,
//	}
//	for _, accVal := range accountValues {
//		if *accVal.Key ==  string(MasterIP) {
//			config.MasterIP = *accVal.StringValue
//		}
//		if *accVal.Key ==  string(MasterPort) {
//			config.MasterPort = *accVal.StringValue
//		}
//		if *accVal.Key ==  string(MasterUsername) {
//			config.MasterUsername = *accVal.StringValue
//		}
//		if *accVal.Key ==  string(MasterPassword) {
//			config.MasterPassword = *accVal.StringValue
//		}
//		if *accVal.Key ==  string(FrameworkIP) {
//			config.FrameworkIP = *accVal.StringValue
//		}
//		if *accVal.Key ==  string(FrameworkPort) {
//			config.FrameworkPort = *accVal.StringValue
//		}
//		if *accVal.Key ==  string(FrameworkUsername) {
//			config.FrameworkUser = *accVal.StringValue
//		}
//		if *accVal.Key ==  string(FrameworkPassword) {
//			config.FrameworkPassword = *accVal.StringValue
//		}
//	}
//
//	ok, err := config.validate()
//	if !ok {
//		glog.Errorf("Invalid config : " + err.Error())
//		return nil //, errors.New("Invalid config : " + err.Error() )
//	}
//	return config
//}


// Get the Account Values to create VMTTarget in the turbo server corresponding to this client
func (mesosConf *MesosTargetConf) GetAccountValues() []*proto.AccountValue {
	var accountValues []*proto.AccountValue
	// Convert all parameters in clientConf to AccountValue list
	prop1 := string(MasterIP)
	accVal := &proto.AccountValue{
		Key: &prop1,
		StringValue: &mesosConf.MasterIP,
	}
	accountValues = append(accountValues, accVal)

	prop2 := string(MasterPort)
	accVal = &proto.AccountValue{
		Key: &prop2,
		StringValue: &mesosConf.MasterPort,
	}
	accountValues = append(accountValues, accVal)

	prop3 := string(MasterUsername)
	accVal = &proto.AccountValue{
		Key: &prop3,
		StringValue: &mesosConf.MasterUsername,
	}
	accountValues = append(accountValues, accVal)

	prop4 := string(MasterPassword)
	accVal = &proto.AccountValue{
		Key: &prop4,
		StringValue: &mesosConf.MasterPassword,
	}
	accountValues = append(accountValues, accVal)

	if mesosConf.Master == Apache {
		prop5 := string(FrameworkIP)
		accVal = &proto.AccountValue{
			Key: &prop5,
			StringValue: &mesosConf.FrameworkIP,
		}
		accountValues = append(accountValues, accVal)

		prop6 := string(FrameworkPort)
		accVal = &proto.AccountValue{
			Key: &prop6,
			StringValue: &mesosConf.FrameworkPort,
		}
		accountValues = append(accountValues, accVal)

		prop7 := string(FrameworkUsername)
		accVal = &proto.AccountValue{
			Key: &prop7,
			StringValue: &mesosConf.FrameworkUser,
		}
		accountValues = append(accountValues, accVal)

		prop8 := string(FrameworkPassword)
		accVal = &proto.AccountValue{
			Key: &prop8,
			StringValue: &mesosConf.FrameworkPassword,
		}
		accountValues = append(accountValues, accVal)
	}

	glog.Infof("[MesosDiscoveryClient] account values %s\n",  accountValues)

	return accountValues
}


func (conf *MesosTargetConf) validate() (bool, error) {
	if (conf.MasterIP == "") {
		return false, errors.New("Mesos Master IP is required "+ fmt.Sprint(conf))
	}
	return true, nil
}


// Get the config from file.
func readConfig(path string) (*MesosTargetConf, error) {
	file, e := ioutil.ReadFile(path)
	if e != nil {
		return nil, errors.New("File error: " + e.Error())
	}
	var config MesosTargetConf
	err := json.Unmarshal(file, &config)

	if err != nil {
		return nil, errors.New(string(file) + " \nUnmarshall error : " + err.Error())
	}

	return &config, nil
}
