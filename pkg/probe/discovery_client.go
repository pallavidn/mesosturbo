package probe

import (
	"time"
	"net/http"
	"errors"
	"io/ioutil"
	"encoding/json"
	"strings"
	"strconv"

	"github.com/golang/glog"

	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/turbonomic/turbo-go-sdk/pkg/probe"

	"github.com/turbonomic/mesosturbo/pkg/conf"
	"github.com/turbonomic/mesosturbo/pkg/util"
	"github.com/turbonomic/mesosturbo/pkg/factory"
)

var (
	DEFAULT_NAMESPACE string = "DEFAULT"
)

// Discovery Client for the Mesos Probe
// Implements the TurboDiscoveryClient interface
type MesosDiscoveryClient struct {
	mesosMasterType 	 conf.MesosMasterType
	clientConf        *conf.MesosTargetConf
	masterRestClient	conf.MasterRestClient
	frameworkRestClient	conf.FrameworkRestClient

	targetIdentifier  string
	username          string
	pwd               string

	builderMap 	map[*proto.EntityDTO_EntityType]EntityBuilder

	lastDiscoveryTime *time.Time
	slaveUseMap       map[string]*util.CalculatedUse
	taskUseMap        map[string]*util.CalculatedUse
}

func NewDiscoveryClient(mesosMasterType conf.MesosMasterType, clientConf *conf.MesosTargetConf) (probe.TurboDiscoveryClient, error) {
	if clientConf == nil {
		return nil, errors.New("[MesosDiscoveryClient] Null config")
	}

	glog.V(2).Infof("[MesosDiscoveryClient] Target Conf ", clientConf)

	client := &MesosDiscoveryClient{
		mesosMasterType: mesosMasterType,
		targetIdentifier: clientConf.MasterIP,
		clientConf: clientConf,
		builderMap: make(map[*proto.EntityDTO_EntityType]EntityBuilder),
	}
	err := client.initDiscoveryClient()
	if err != nil {
		return nil, errors.New("MesosDiscoveryClient] " + err.Error())
	}

	return client, nil
}

func (discoveryClient *MesosDiscoveryClient) initDiscoveryClient() error {
	clientConf := discoveryClient.clientConf
	// Based on the Mesos vendor, instantiate the MesosRestClient
	masterRestClient := factory.GetMasterRestClient(clientConf.Master, clientConf)
	if masterRestClient == nil {
		return errors.New("Cannot find RestClient for Mesos : " + string(clientConf.Master))
	}

	// Login to the Mesos Master and save the login token
	token, err := masterRestClient.Login()
	if err != nil {
		return errors.New("Error logging to Mesos Master at " +
					clientConf.MasterIP + "::" + clientConf.MasterPort + " : " + err.Error())
	}

	clientConf.Token = token
	// Based on the Framework vendor, instantiate the Frameworks RestClient
	frameworkClient := factory.GetFrameworkRestClient(clientConf.Framework, clientConf)
	if frameworkClient == nil {
		glog.Errorf("[MesosDiscoveryClient] Cannot find framework Client for Mesos : ", clientConf.Framework)
	}

	discoveryClient.masterRestClient = masterRestClient
	discoveryClient.frameworkRestClient  = frameworkClient

	return nil
}


// Get the Account Values to create VMTTarget in the turbo server corresponding to this client
func (handler *MesosDiscoveryClient) GetAccountValues() *probe.TurboTargetInfo { //[]*proto.AccountValue {		//*probe.TurboTargetInfo {
	var accountValues []*proto.AccountValue
	clientConf := handler.clientConf
	// Convert all parameters in clientConf to AccountValue list
	prop1 := string(conf.MasterIP)
	accVal := &proto.AccountValue{
		Key: &prop1,
		StringValue: &clientConf.MasterIP,
	}
	accountValues = append(accountValues, accVal)

	prop2 := string(conf.MasterPort)
	accVal = &proto.AccountValue{
		Key: &prop2,
		StringValue: &clientConf.MasterPort,
	}
	accountValues = append(accountValues, accVal)

	prop3 := string(conf.MasterUsername)
	accVal = &proto.AccountValue{
		Key: &prop3,
		StringValue: &clientConf.MasterUsername,
	}
	accountValues = append(accountValues, accVal)

	prop4 := string(conf.MasterPassword)
	accVal = &proto.AccountValue{
		Key: &prop4,
		StringValue: &clientConf.MasterPassword,
	}
	accountValues = append(accountValues, accVal)

	if handler.mesosMasterType == conf.Apache {
		prop5 := string(conf.FrameworkIP)
		accVal = &proto.AccountValue{
			Key: &prop5,
			StringValue: &clientConf.FrameworkIP,
		}
		accountValues = append(accountValues, accVal)

		prop6 := string(conf.FrameworkPort)
		accVal = &proto.AccountValue{
			Key: &prop6,
			StringValue: &clientConf.FrameworkPort,
		}
		accountValues = append(accountValues, accVal)

		prop7 := string(conf.FrameworkUsername)
		accVal = &proto.AccountValue{
			Key: &prop7,
			StringValue: &clientConf.FrameworkUser,
		}
		accountValues = append(accountValues, accVal)

		prop8 := string(conf.FrameworkPassword)
		accVal = &proto.AccountValue{
			Key: &prop8,
			StringValue: &clientConf.FrameworkPassword,
		}
		accountValues = append(accountValues, accVal)
	}

	//prop9 := "ActionIP"
	//accVal = &proto.AccountValue{
	//	Key: &prop9,
	//	StringValue: &clientConf.ActionIP,
	//}
	//accountValues = append(accountValues, accVal)
	//
	//prop10 := "ActionPort"
	//accVal = &proto.AccountValue{
	//	Key: &prop10,
	//	StringValue: &clientConf.ActionPort,
	//}
	//accountValues = append(accountValues, accVal)

	glog.Infof("[MesosDiscoveryClient] account values %s\n",  accountValues)
	targetInfo := probe.NewTurboTargetInfoBuilder("CloudNative", "MesosProbe", string(conf.MasterIP), accountValues).Create()

	return targetInfo //accountValues
}

// Validate the Target
func (handler *MesosDiscoveryClient) Validate(accountValues[] *proto.AccountValue) *proto.ValidationResponse {
	glog.Infof("BEGIN Validation for MesosDiscoveryClient  %s\n", accountValues)
	// Login to the Mesos Master and save the login token
	token, err := handler.masterRestClient.Login()
	if err != nil {
		//TODO: throw exception to the calling layer
		glog.Errorf("[MesosDiscoveryClient] Error logging to Mesos Master at ",
					accountValues)
	}
	handler.clientConf.Token = token
	// TODO: login here and save the login token
	validationResponse := &proto.ValidationResponse{}

	glog.Infof("validation response %s\n", validationResponse)
	return validationResponse
}


// Discover the Target Topology
func (handler *MesosDiscoveryClient) Discover(accountValues[] *proto.AccountValue) (*proto.DiscoveryResponse) {
	glog.Infof("BEGIN Discovery for MesosDiscoveryClient %s\n", accountValues)
	//Discover the Mesos topology
	// Mesos Master state
	// TODO: update leader and reissue request
	mesosState, err := handler.masterRestClient.GetState()
	if err != nil {
		glog.Errorf("Error getting state from master : %s \n", err)
		return nil
	}
	glog.V(3).Infof("Mesos Get Succeeded: %v\n", mesosState)

	// Handler updated leader
	//if err != nil && err.Error() == "update leader" {
	//	mesosState, err = handler.masterRestClient.GetState()
	//	if err != nil {
	//		glog.Errorf("Error, need to update leader")
	//		return nil
	//	}
	//}

	// Framework response
	frameworkResp, err := handler.frameworkRestClient.GetFrameworkApps()

	if err != nil {
		glog.Errorf("Error parsing response from the Framework: %s \n", err)
		return nil
	}

	glog.V(3).Infof("Marathon Get Succeeded: %v\n", frameworkResp)
	handler.parseMesosState(mesosState, frameworkResp)	// to create convenience maps for slaves, tasks, convert units

	// 2. Build Entities
	nodeBuilder := &VMEntityBuilder{
		MasterState:mesosState,
	}
	nodeEntityDtos, err := nodeBuilder.BuildEntities()
	if err != nil {
		glog.Errorf("Error parsing nodes: %s. Will return.", err)
		return nil
	}

	containerBuilder := &ContainerEntityBuilder{
		MasterState:mesosState,
	}
	containerEntityDtos, err := containerBuilder.BuildEntities()	//handler.ParseTask(mesosState, handler.taskUseMap)
	if err != nil {
		// TODO, should here still send out msg to server? Or set errorDTO?
		glog.Errorf("Error parsing pods: %s. Will return.", err)
		return nil
	}

	appBuilder := &AppEntityBuilder{
		MasterState:mesosState,
	}
	appEntityDtos, err := appBuilder.BuildEntities()	//handler.ParseTask(mesosState, handler.taskUseMap)
	if err != nil {
		// TODO, should here still send out msg to server? Or set errorDTO?
		glog.Errorf("Error parsing pods: %s. Will return.", err)
		return nil
	}

	entityDtos := nodeEntityDtos
	entityDtos = append(entityDtos, containerEntityDtos...)
	entityDtos = append(entityDtos, appEntityDtos...)

	// 3. Discovery Response
	discoveryResponse := &proto.DiscoveryResponse{
		EntityDTO: entityDtos,
	}

	currtime := time.Now()
	handler.lastDiscoveryTime = &currtime
	glog.Infof("END Discovery for MesosDiscoveryClient %s", accountValues)
	return discoveryResponse
}

// ===================== Detailed Discovery ===========================================
func (handler *MesosDiscoveryClient) parseMesosState (stateResp *util.MesosAPIResponse, frameworkResp *util.FrameworkApps) (*util.MesosAPIResponse, error) {
	//glog.V(3).Infof("Get Succeed: %v\n", stateResp)
	// UPDATE RESOURCE UNITS AFTER HTTP REQUEST for each slave
	for idx := range stateResp.Slaves {
		glog.V(3).Infof("Number of slaves %d \n", len(stateResp.Slaves))
		s := &stateResp.Slaves[idx]
		s.Resources.Mem = s.Resources.Mem * float64(1024)		// TODO: convert units for mem only?
		s.UsedResources.Mem = s.UsedResources.Mem * float64(1024)
		s.OfferedResources.Mem = s.OfferedResources.Mem * float64(1024)
		glog.Infof("=======> SLAVE idk: %d name: %s, mem: %.2f, cpu: %.2f, disk: %.2f \n", idx, s.Name, s.Resources.Mem, s.Resources.CPUs, s.Resources.Disk)
	}

	if stateResp.Frameworks == nil {
		glog.Errorf("Error getting Frameworks response")
		return nil, errors.New("Error getting Frameworks response: %s")
	}

	// Map of Slave ID  and IP
	stateResp.SlaveIdIpMap = make(map[string]string)
	for _, slave := range stateResp.Slaves {
		slaveIP := util.GetSlaveIP(slave)
		stateResp.SlaveIdIpMap[slave.Id] = slaveIP
	}

	// Parse the Frameworks to get the list of all Tasks across all frameworks
	//We pass the entire http response as the respContent object
	taskContent, err := handler.parseAPITasksResponse(stateResp)
	if err != nil {
		glog.Errorf("Error getting response: %s", err)
		return nil, err
	}
	glog.V(3).Infof("Number of tasks \n", len(taskContent.Tasks))

	stateResp.TaskMasterAPI = *taskContent

	// Cluster
	stateResp.Cluster.MasterIP = handler.clientConf.MasterIP
	stateResp.Cluster.ClusterName = stateResp.ClusterName

	// stateResp.MApps = frameworkResp	//TODO:

	// TODO: related to metrics
	//// STATS
	//var mapTaskRes map[string]util.Statistics
	//mapTaskRes = make(map[string]util.Statistics)
	//var mapSlaveUse map[string]*util.CalculatedUse
	//mapSlaveUse = make(map[string]*util.CalculatedUse)
	//var mapTaskUse map[string]*util.CalculatedUse
	//mapTaskUse = make(map[string]*util.CalculatedUse)
	//var ports_slaves = []string{}
	//var allports []string
	//allports = make([]string, 1)
	//
	//// Get Metrics for each slave or agent
	//for i := range stateResp.Slaves {
	//	s := stateResp.Slaves[i]
	//	someports, err := handler.monitorSlaveStatistics(s, handler.taskUseMap, mapTaskRes, mapSlaveUse, mapTaskUse, ports_slaves)
	//	if err != nil {
	//		glog.Errorf("Error getting use data for slave %s \n", s.Name)
	//		continue
	//	}
	//	for _, someport := range someports {
	//		allports = append(allports, someport)
	//	}
	//} // slave loop
	//
	//glog.V(3).Info("--------------=======> ALL PORTS: %+v\n\n", allports)
	//stateResp.AllPorts = allports
	//// map task to resources
	//handler.taskUseMap = mapTaskUse
	//handler.slaveUseMap = mapSlaveUse
	//stateResp.MapTaskStatistics = mapTaskRes
	//stateResp.SlaveUseMap = mapSlaveUse

	for _, framework := range stateResp.Frameworks {
		glog.Infof("Framework : ", framework.Name + "::" + framework.Hostname)
		for _, task := range framework.Tasks {
			glog.Infof("	Task : %s", task.Name)
		}
	}

	for _, slave := range stateResp.Slaves {
		glog.Infof("Slave : %s", slave.Name + "::" + slave.Pid)
	}

	return stateResp, nil
}

func (handler *MesosDiscoveryClient) monitorSlaveStatistics(s util.Slave,
								previousUseMap map[string]*util.CalculatedUse,
								mapTaskRes map[string]util.Statistics,
								mapSlaveUse map[string]*util.CalculatedUse,
								mapTaskUse map[string]*util.CalculatedUse,
								ports_slaves []string) ([]string, error) {
	fullUrl := "http://" + util.GetSlaveIP(s) + ":" + handler.clientConf.SlavePort + "/monitor/statistics.json"

	req, err := http.NewRequest("GET", fullUrl, nil)

	if err != nil {
		return nil, err
	}

	//if handler.clientConf.DCOS { // for DCOS slave, send the session token
	//	req.Header.Add("content-type", "application/json")
	//	req.Header.Add("authorization", "token="+handler.clientConf.Token)
	//}

	req.Close = true
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		glog.Errorf("Error getting response: %s\n", err)
		return nil, err
	}
	defer resp.Body.Close()
	stringResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("error %s\n", err)
		return nil, err
	}
	byteContent := []byte(stringResp)
	var usedRes = new([]util.Executor)
	err = json.Unmarshal(byteContent, &usedRes)
	if err != nil {
		glog.Errorf("JSON error %s", err)
		return nil, err
	}

	var arrOfExec []util.Executor
	arrOfExec = *usedRes

	var allports []string
	allports = make([]string, 1)
	// Create port array and used port struct
	// "[31100-31100, 31250-31250, 31674-31674, 31766-31766, 31944-31944, 31978-31978]"
	var portsAtSlave map[string]util.PortUtil
	portsAtSlave = make(map[string]util.PortUtil)
	if s.UsedResources.Ports != "" {
		original := s.UsedResources.Ports
		glog.V(3).Infof("=========-------> used ports is %+v\n", original)
		portsStr := original[1 : len(original)-1]
		glog.V(3).Infof("=========-------> used ports is %+v\n", portsStr)
		portRanges := strings.Split(portsStr, ",")
		for _, prange := range portRanges {
			glog.V(3).Infof("=========-------> prange is %+v\n", prange)
			ports := strings.Split(prange, "-")
			glog.V(3).Infof("=========-------> port is %+v\n", ports[0])
			portStart, err := strconv.Atoi(strings.Trim(ports[0], " "))
			if err != nil {
				glog.V(3).Infof(" Error: %+v", err)
				return nil, err
			}
			if strings.Trim(ports[0], " ") == strings.Trim(ports[1], " ") {
				// all slaves
				allports = append(allports, strings.Trim(ports[0], " "))
				// single slave
				portsAtSlave[strings.Trim(ports[0], " ")] = util.PortUtil{
					Number:   float64(portStart),
					Capacity: float64(1.0),
					Used:     float64(1.0),
				}
			} else {
				//range from port start to end
				for _, p := range ports {
					allports = append(allports, strings.Trim(p, " "))
					port, err := strconv.Atoi(strings.Trim(p, " "))
					if err != nil {
						glog.V(3).Infof("Error getting port %+v", err)
						return nil, err
					}
					// single slave
					portsAtSlave[strings.Trim(p, " ")] = util.PortUtil{
						Number:   float64(port),
						Capacity: float64(1.0),
						Used:     float64(1.0),
					}
				}
			}
		}
	}
	mapSlaveUse[s.Id] = &util.CalculatedUse{
		CPUs:      float64(0.0),
		Mem:       float64(0.0),
		UsedPorts: portsAtSlave,
	}

	for j := range arrOfExec {
		executor := arrOfExec[j]
		// TODO check if this is taskId
		taskId := executor.Source
		mapTaskRes[taskId] = executor.Statistics

		// TASK MONITOR
		if _, ok := mapTaskUse[taskId]; !ok {
			var prevSecs float64

			// CPU use CALCULATION STARTS

			curSecs := executor.Statistics.CPUsystemTimeSecs + executor.Statistics.CPUuserTimeSecs
			_, ok := previousUseMap[taskId]
			if previousUseMap == nil || !ok {
				glog.V(4).Infof(" map was nil !!")
				prevSecs = curSecs

			} else {
				prevSecs = previousUseMap[taskId].CPUsumSystemUserSecs
				glog.V(4).Infof("previous system + user : %f ", prevSecs)
			}
			diffSecs := curSecs - prevSecs
			if diffSecs < 0 {
				diffSecs = float64(0.0)
			}
			glog.V(4).Infof(" t1 - t0 : %f \n", diffSecs)
			var lastTime time.Time
			if handler.lastDiscoveryTime == nil {
				lastTime = time.Now()
			} else {
				lastTime = *handler.lastDiscoveryTime
			}
			diffTime := time.Since(lastTime)
			diffT := diffTime.Seconds()
			usedCPUfraction := diffSecs / diffT
			// ratio * cores * 1000kHz
			glog.V(4).Infof("-------------> Fraction of CPU utilization: %f \n", usedCPUfraction)

			// s.Resources is # of cores
			// usedCPU is in MHz
			usedCPU := usedCPUfraction * s.Resources.CPUs * float64(1000)
			mapTaskUse[taskId] = &util.CalculatedUse{
				CPUs:                 usedCPU,
				CPUsumSystemUserSecs: curSecs,
			}
			glog.V(4).Infof("------------> Capacity in CPUs, directly from Mesos %f \n", s.Resources.CPUs)
			glog.V(4).Infof("------------->Used CPU in MHz : %f \n", usedCPU)

			// Sum the used CPU in MHz for each slave
			mapSlaveUse[s.Id].CPUs = usedCPU + mapSlaveUse[s.Id].CPUs
			// Mem is returned in B convert to KB
			// usedRes is reply from statistics.json
			usedMem_B := executor.Statistics.MemRSSBytes
			usedMem_KB := usedMem_B / float64(1024.0)
			mapSlaveUse[s.Id].Mem = mapSlaveUse[s.Id].Mem + usedMem_KB
		}
	} // task loop
	glog.V(3).Infof(" ------------------>>>>> ALLPORTS reutrning: %+v", allports)
	return allports, nil
}

func (handler *MesosDiscoveryClient) parseAPITasksResponse(resp *util.MesosAPIResponse) (*util.MasterTasks, error) {
	glog.V(4).Infof("----> in parseAPITasksResponse")
	if resp == nil {
		return nil, errors.New("Task information response received is nil")
	}
	glog.V(3).Infof("[MesosDiscoveryClient] Number of frameworks is %d\n", len(resp.Frameworks))

	allTasks := make([]util.Task, 0)
	for _, framework := range resp.Frameworks {
		if framework.Tasks != nil {
			ftasks := framework.Tasks
			for _, task := range ftasks {
				allTasks = append(allTasks, task)
			}
			glog.Infof("[MesosDiscoveryClient] Number of tasks in framework %s is %d", framework.Name, len(framework.Tasks))
		}
	}
	tasksObj := &util.MasterTasks{
		Tasks: allTasks,
	}

	// For each task, convert Mem units
	for j := range tasksObj.Tasks {
		t := tasksObj.Tasks[j]
		// MEM UNITS KB
		t.Resources.Mem = t.Resources.Mem * float64(1024)

	}
	return tasksObj, nil
}

