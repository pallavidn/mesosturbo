package probe

import (
	"time"
	"net/http"
	"errors"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strings"
	"strconv"

	"github.com/golang/glog"

	proto "github.com/turbonomic/turbo-go-sdk/pkg/proto"
	probe "github.com/turbonomic/turbo-go-sdk/pkg/probe"
	builder "github.com/turbonomic/turbo-go-sdk/pkg/builder"

	conf "github.com/turbonomic/mesosturbo/pkg/conf"
	mesoshttp "github.com/turbonomic/mesosturbo/pkg/mesoshttp"
	"github.com/turbonomic/mesosturbo/pkg/util"
)

var (
	DEFAULT_NAMESPACE string = "DEFAULT"
)

// Discovery Client for the Mesos Probe
// Implements the TurboDiscoveryClient interface
type MesosDiscoveryClient struct {
	lastDiscoveryTime *time.Time
	slaveUseMap       map[string]*util.CalculatedUse
	taskUseMap        map[string]*util.CalculatedUse
	clientConf        *conf.MesosTargetConf
	targetIdentifier  string
	username          string
	pwd               string
}

func NewDiscoveryClient(targetIdentifier string, confFile string) *MesosDiscoveryClient {
	// Parse conf file to create clientConf
	clientConf, _ := conf.NewMesosTargetConf(confFile)
	fmt.Println("[MesosDiscoveryClient] Target Conf ", clientConf)
	// TODO: handle error
	client := &MesosDiscoveryClient{
		targetIdentifier: targetIdentifier,
		clientConf: clientConf,

	}
	return client
}

// Get the Account Values to create VMTTarget in the turbo server corresponding to this client
func (handler *MesosDiscoveryClient) GetAccountValues() *probe.TurboTarget {
	var accountValues []*proto.AccountValue
	// Convert all parameters in clientConf to AccountValue list
	prop := "MarathonIP"
	accVal := &proto.AccountValue{
		Key: &prop,
		StringValue: &handler.clientConf.MarathonIP,
	}
	accountValues = append(accountValues, accVal)

	prop = "MarathonPort"
	accVal = &proto.AccountValue{
		Key: &prop,
		StringValue: &handler.clientConf.MarathonPort,
	}
	accountValues = append(accountValues, accVal)

	prop = "MesosIP"
	accVal = &proto.AccountValue{
		Key: &prop,
		StringValue: &handler.clientConf.MesosIP,
	}
	accountValues = append(accountValues, accVal)

	prop = "MesosPort"
	accVal = &proto.AccountValue{
		Key: &prop,
		StringValue: &handler.clientConf.MesosPort,
	}
	accountValues = append(accountValues, accVal)

	prop = "ActionIP"
	accVal = &proto.AccountValue{
		Key: &prop,
		StringValue: &handler.clientConf.ActionIP,
	}
	accountValues = append(accountValues, accVal)

	prop = "ActionPort"
	accVal = &proto.AccountValue{
		Key: &prop,
		StringValue: &handler.clientConf.ActionPort,
	}
	accountValues = append(accountValues, accVal)

	prop = "SlavePort"
	accVal = &proto.AccountValue{
		Key: &prop,
		StringValue: &handler.clientConf.SlavePort,
	}
	accountValues = append(accountValues, accVal)

	targetInfo := &probe.TurboTarget{
		AccountValues: accountValues,
	}

	targetInfo.SetUser("defaultUsername")
	targetInfo.SetPassword("defaultPassword")
	return targetInfo
}

// Validate the Target
func (handler *MesosDiscoveryClient) Validate(accountValues[] *proto.AccountValue) *proto.ValidationResponse {
	fmt.Printf("[MesosDiscoveryClient] BEGIN Validation for MesosDiscoveryClient  %s", accountValues)
	// TODO: connect to the client and get validation response
	validationResponse := &proto.ValidationResponse{}

	fmt.Printf("[MesosDiscoveryClient] validation response %s\n", validationResponse)
	return validationResponse
}

// Discover the Target Topology
func (handler *MesosDiscoveryClient) Discover(accountValues[] *proto.AccountValue) *proto.DiscoveryResponse {
	fmt.Printf("[MesosDiscoveryClient] BEGIN Discovery for MesosDiscoveryClient %s", accountValues)
	//Discover the Mesos topology
	glog.V(3).Infof("[MesosDiscoveryClient] Discover topology request from server.")
	// 1. Get message ID
	var stopCh chan struct{} = make(chan struct{})
	defer close(stopCh)

	// 2. Build discoverResponse
	// Get the discovery client for the given target identifier from the account values map
	//handler := myProbe.GetTurboDiscoveryClient(accountValues)

	mesosProbe, err := handler.NewMesosProbe()
	if err != nil && err.Error() == "update leader" {
		mesosProbe, err = handler.NewMesosProbe()
		if err != nil {
			glog.Errorf("Error, need to update leader")
			return nil
		}
	}

	if err != nil {
		glog.Errorf("Error getting state from master : %s", err)
		return nil
	}

	nodeEntityDtos, err := handler.ParseNode(mesosProbe, handler.slaveUseMap)
	if err != nil {
		glog.Errorf("Error parsing nodes: %s. Will return.", err)
		fmt.Println("[MesosDiscoveryClient] Error parsing nodes: %s. Will return.", err)
		return nil
	}
	containerEntityDtos, err := handler.ParseTask(mesosProbe, handler.taskUseMap)
	if err != nil {
		// TODO, should here still send out msg to server? Or set errorDTO?
		glog.Errorf("Error parsing pods: %s. Will return.", err)
		fmt.Println("[MesosDiscoveryClient] Error parsing pods: %s. Will return.", err)
		return nil
	}

	entityDtos := nodeEntityDtos
	entityDtos = append(entityDtos, containerEntityDtos...)
	//	entityDtos = append(entityDtos, serviceEntityDtos...)
	discoveryResponse := &proto.DiscoveryResponse{
		EntityDTO: entityDtos,
	}

	currtime := time.Now()
	handler.lastDiscoveryTime = &currtime
	fmt.Printf("[MesosDiscoveryClient] END Discovery for MesosDiscoveryClient %s", accountValues)
	return discoveryResponse
}

// ===================== Detailed Discovery ===========================================
// TODO: TO REFACTOR
func (handler *MesosDiscoveryClient) NewMesosProbe() (*util.MesosAPIResponse, error) {
	var fullUrl string
	if handler.clientConf.MesosPort == "" {
		fullUrl = "http://" + handler.clientConf.MesosIP + "/mesos/state"
	} else {
		fullUrl = "http://" + handler.clientConf.MesosIP + ":" + handler.clientConf.MesosPort + "/state"
	}
	fmt.Println("[MesosDiscoveryClient] The full Url is ", fullUrl)

	req, err := http.NewRequest("GET", fullUrl, nil)

	if err != nil {
		glog.Errorf("Error in GET request: %s\n", err)
		return nil, err
	}

	// DCOS mode only
	if handler.clientConf.DCOS {
		req.Header.Add("content-type", "application/json")
		req.Header.Add("authorization", "token="+handler.clientConf.Token)
	}

	fmt.Println("%+v", req)
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		glog.Errorf("Error in GET request to mesos master: %s\n", err)
		return nil, err
	}

	defer resp.Body.Close()

	// Get token if response if OK
	if resp.Status == "" {
		glog.Errorf("Empty response status\n")
		return nil, errors.New("Empty response status\n")
	}

	if resp.StatusCode != 200 {
		mesosdcosCli := &mesoshttp.MesosHTTPClient{
			MesosMasterBase: fullUrl,
		}
		errormsg := mesosdcosCli.DCOSLoginRequest(handler.clientConf, handler.clientConf.Token)
		if errormsg != nil {
			glog.Errorf("Please check DCOS credentials and start mesosturbo again.\n")
			return nil, errormsg
		}
		glog.V(3).Infof("Current token has expired, updated DCOS token.\n")
	}

	respContent, err := handler.parseAPIStateResponse(resp)

	if err != nil {
		glog.Errorf("Error parsing Mesos master state response")
		return nil, err
	}

	currentLeader := respContent.Leader
	respContent.Leader = currentLeader[7 : len(currentLeader)-5]

	if respContent.Leader != handler.clientConf.MesosIP {
		// not good, update leader
		handler.clientConf.MesosIP = respContent.Leader
		glog.V(3).Infof("The mesos master IP has been updated to : %s \n", handler.clientConf.MesosIP)
		return nil, fmt.Errorf("update leader")
	}

	if respContent.SlaveIdIpMap == nil {
		respContent.SlaveIdIpMap = make(map[string]string)
	}

	// UPDATE RESOURCE UNITS AFTER HTTP REQUEST
	for idx := range respContent.Slaves {
		glog.V(3).Infof("Number of slaves %d \n", len(respContent.Slaves))
		s := &respContent.Slaves[idx]
		s.Resources.Mem = s.Resources.Mem * float64(1024)
		s.UsedResources.Mem = s.UsedResources.Mem * float64(1024)
		s.OfferedResources.Mem = s.OfferedResources.Mem * float64(1024)
		glog.V(3).Infof("=======> SLAVE idk: %d name: %s, mem: %.2f, cpu: %.2f, disk: %.2f \n", idx, s.Name, s.Resources.Mem, s.Resources.CPUs, s.Resources.Disk)
	}

	if err != nil {
		glog.Errorf("Error getting response: %s", err)
		return nil, err
	}

	glog.V(3).Infof("Get Succeed: %v\n", respContent)

	if respContent.Frameworks == nil {
		glog.Errorf("Error getting Frameworks response")
		return nil, errors.New("Error getting Frameworks response: %s")
	}
	/*
		configFile, err := os.Open("task.json")
		if err != nil {
			fmt.Println("opening config file", err.Error())
		}
		var jsonTasks = new(util.MasterTasks)

		jsonParser := json.NewDecoder(configFile)
		if err = jsonParser.Decode(jsonTasks); err != nil {
			fmt.Println("parsing config file", err.Error())
		}
		taskContent := jsonTasks
	*/

	//We pass the entire http response as the respContent object
	taskContent, err := handler.parseAPITasksResponse(respContent)
	if err != nil {
		glog.Errorf("Error getting response: %s", err)
		return nil, err
	}
	glog.V(3).Infof("Number of tasks \n", len(taskContent.Tasks))

	for j := range taskContent.Tasks {
		t := taskContent.Tasks[j]
		// MEM UNITS KB
		t.Resources.Mem = t.Resources.Mem * float64(1024)
		//	fmt.Printf("----> tasks from mesos: # %d, name : %s, state: %s\n", j, t.Name, t.State)
		//	glog.V(3).Infof("=======> TASK name: %s, mem: %.2f, cpu: %.2f, disk: %.2f \n", t.Name, t.Resources.Mem, t.Resources.CPUs, t.Resources.Disk)

	}
	respContent.TaskMasterAPI = *taskContent

	//Marathon
	fullUrlM := "http://" + handler.clientConf.MarathonIP + ":" + handler.clientConf.MarathonPort + "/v2/apps"
	glog.V(4).Infof("The full Url is ", fullUrlM)

	reqM, err := http.NewRequest("GET", fullUrlM, nil)

	if err != nil {
		glog.Errorf("GET request creation for url: %s failed \n", fullUrlM)
		return nil, err
	}

	if handler.clientConf.DCOS {
		reqM.Header.Add("content-type", "application/json")
		reqM.Header.Add("authorization", "token="+handler.clientConf.Token)
	}

	glog.V(4).Infof("%+v", reqM)
	clientM := &http.Client{}
	respM, err := clientM.Do(reqM)
	if err != nil {
		glog.Errorf("Error getting response from Marathon: %s \n", err)
		return nil, err
	}
	defer respM.Body.Close()

	marathonRespContent, err := handler.parseMarathonResponse(respM)

	if err != nil {
		glog.Errorf("Error parsing response form Marathon: %s \n", err)
		return nil, err
	}

	glog.V(3).Infof("Marathon Get Succeed: %v\n", marathonRespContent)

	respContent.MApps = marathonRespContent

	// STATS
	var mapTaskRes map[string]util.Statistics
	mapTaskRes = make(map[string]util.Statistics)
	var mapSlaveUse map[string]*util.CalculatedUse
	mapSlaveUse = make(map[string]*util.CalculatedUse)
	var mapTaskUse map[string]*util.CalculatedUse
	mapTaskUse = make(map[string]*util.CalculatedUse)
	var ports_slaves = []string{}
	var allports []string
	allports = make([]string, 1)

	for i := range respContent.Slaves {
		s := respContent.Slaves[i]
		someports, err := handler.monitorSlaveStatistics(s, handler.taskUseMap, mapTaskRes, mapSlaveUse, mapTaskUse, ports_slaves)
		if err != nil {
			glog.Errorf("Error getting use data for slave %s \n", s.Name)
			continue
		}
		for _, someport := range someports {
			allports = append(allports, someport)
		}
	} // slave loop

	glog.V(3).Info("--------------=======> ALL PORTS: %+v\n\n", allports)
	respContent.AllPorts = allports
	// map task to resources
	handler.taskUseMap = mapTaskUse
	handler.slaveUseMap = mapSlaveUse
	respContent.MapTaskStatistics = mapTaskRes
	respContent.SlaveUseMap = mapSlaveUse
	respContent.Cluster.MasterIP = handler.clientConf.MesosIP
	respContent.Cluster.ClusterName = respContent.ClusterName

	return respContent, nil
}

func (handler *MesosDiscoveryClient) monitorSlaveStatistics(s util.Slave, previousUseMap map[string]*util.CalculatedUse, mapTaskRes map[string]util.Statistics, mapSlaveUse map[string]*util.CalculatedUse, mapTaskUse map[string]*util.CalculatedUse, ports_slaves []string) ([]string, error) {
	fullUrl := "http://" + util.GetSlaveIP(s) + ":" + handler.clientConf.SlavePort + "/monitor/statistics.json"

	req, err := http.NewRequest("GET", fullUrl, nil)

	if err != nil {
		return nil, err
	}

	if handler.clientConf.DCOS {
		req.Header.Add("content-type", "application/json")
		req.Header.Add("authorization", "token="+handler.clientConf.Token)
	}

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
	glog.V(4).Infof("----> in parseAPICallResponse")
	if resp == nil {
		return nil, errors.New("Task information response received is nil")
	}
	glog.V(3).Infof(" Number of frameworks is %d\n", len(resp.Frameworks))

	allTasks := make([]util.Task, 0)
	for i := range resp.Frameworks {
		if resp.Frameworks[i].Tasks != nil {
			ftasks := resp.Frameworks[i].Tasks
			for j := range ftasks {
				allTasks = append(allTasks, ftasks[j])
			}
			glog.V(3).Infof(" Number of tasks is %d\n", len(resp.Frameworks[i].Tasks))
		}
	}
	tasksObj := &util.MasterTasks{
		Tasks: allTasks,
	}
	return tasksObj, nil
}

func (handler *MesosDiscoveryClient) parseAPIStateResponse(resp *http.Response) (*util.MesosAPIResponse, error) {
	glog.V(4).Infof("----> in parseAPICallResponse")
	if resp == nil {
		return nil, errors.New("Response sent from mesos/DCOS master is nil")
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("Error after ioutil.ReadAll: %s", err)
		return nil, err
	}

	glog.V(4).Infof("response content is %s", string(content))
	byteContent := []byte(content)
	var jsonMesosMaster = new(util.MesosAPIResponse)
	err = json.Unmarshal(byteContent, &jsonMesosMaster)
	if err != nil {
		glog.Errorf("error in json unmarshal : %s", err)
		return nil, errors.New("Error in json unmarshal")
	}
	return jsonMesosMaster, nil
}

func (handler *MesosDiscoveryClient) parseMarathonResponse(resp *http.Response) (*util.MarathonApps, error) {
	glog.V(4).Infof("----> in parseAPICallResponse")
	if resp == nil {
		return nil, fmt.Errorf("response sent in is nil")
	}
	glog.V(3).Infof(" from glog response body is %s", resp.Body)

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("Error after ioutil.ReadAll: %s", err)
		return nil, err
	}

	glog.V(4).Infof("response content is %s", string(content))
	byteContent := []byte(content)
	var jsonMarathonMaster = new(util.MarathonApps)
	err = json.Unmarshal(byteContent, &jsonMarathonMaster)
	if err != nil {
		glog.Errorf("error in json unmarshal : %s", err)
		return nil, err
	}
	for i, app := range jsonMarathonMaster.Apps {
		newN := app.Name[1:len(app.Name)]
		jsonMarathonMaster.Apps[i].Name = newN
	}
	glog.V(3).Infof(" MARATHON resp %+v", jsonMarathonMaster)
	return jsonMarathonMaster, nil
}

func (handler *MesosDiscoveryClient) ParseNode(m *util.MesosAPIResponse, slaveUseMap map[string]*util.CalculatedUse) ([]*proto.EntityDTO, error) {
	glog.V(4).Infof("in ParseNode\n")
	result := []*proto.EntityDTO{}
	for i := range m.Slaves {
		s := m.Slaves[i]
		// build sold commodities
		slaveProbe := &NodeProbe{
			MasterState:   m,
			Cluster:       &m.Cluster,
			AllSlavePorts: m.AllPorts,
		}
		commoditiesSold, err := slaveProbe.CreateCommoditySold(&s, slaveUseMap)
		if err != nil {
			glog.Errorf("error is : %s\n", err)
			return result, err
		}
		slaveIP := util.GetSlaveIP(s)
		m.SlaveIdIpMap[s.Id] = slaveIP
		entityDTO := buildVMEntityDTO(slaveIP, s.Id, s.Name, commoditiesSold)
		result = append(result, entityDTO)
	}
	glog.V(4).Infof(" entity DTOs : %d\n", len(result))
	return result, nil
}

func buildVMEntityDTO(slaveIP, nodeID, displayName string, commoditiesSold []*proto.CommodityDTO) *proto.EntityDTO {
	entityDTOBuilder := builder.NewEntityDTOBuilder(proto.EntityDTO_VIRTUAL_MACHINE, nodeID)
	entityDTOBuilder.DisplayName(displayName)
	entityDTOBuilder.SellsCommodities(commoditiesSold)
	// TODO stitch
	ipAddress := slaveIP //nodeProbe.getIPForStitching(displayName)
	ipPropName := "IP"
	ipProp := &proto.EntityDTO_EntityProperty{
		Namespace: &DEFAULT_NAMESPACE,
		Name: &ipPropName,
		Value: &ipAddress,
	}
	entityDTOBuilder = entityDTOBuilder.WithProperty(ipProp)	//"IP", ipAddress)
	glog.V(4).Infof("Parse node: The ip of vm to be reconcile with is %s", ipAddress)
	metaData := generateReconcilationMetaData()

	entityDTOBuilder = entityDTOBuilder.ReplacedBy(metaData)
	entityDto, _ := entityDTOBuilder.Create()
	return entityDto
}

func generateReconcilationMetaData() *proto.EntityDTO_ReplacementEntityMetaData {
	replacementEntityMetaDataBuilder := builder.NewReplacementEntityMetaDataBuilder()
	replacementEntityMetaDataBuilder.Matching("IP")
	replacementEntityMetaDataBuilder.PatchSelling(proto.CommodityDTO_CPU_ALLOCATION)
	replacementEntityMetaDataBuilder.PatchSelling(proto.CommodityDTO_MEM_ALLOCATION)
	replacementEntityMetaDataBuilder.PatchSelling(proto.CommodityDTO_STORAGE_ALLOCATION)
	replacementEntityMetaDataBuilder.PatchSelling(proto.CommodityDTO_CLUSTER)
	replacementEntityMetaDataBuilder.PatchSelling(proto.CommodityDTO_VCPU)
	replacementEntityMetaDataBuilder.PatchSelling(proto.CommodityDTO_VMEM)
	replacementEntityMetaDataBuilder.PatchSelling(proto.CommodityDTO_APPLICATION)
	replacementEntityMetaDataBuilder.PatchSelling(proto.CommodityDTO_VMPM_ACCESS)
	metaData := replacementEntityMetaDataBuilder.Build()
	return metaData
}

func (handler *MesosDiscoveryClient) ParseTask(m *util.MesosAPIResponse, taskUseMap map[string]*util.CalculatedUse) ([]*proto.EntityDTO, error) {
	result := []*proto.EntityDTO{}
	taskList := m.TaskMasterAPI.Tasks

	builder := &TaskBuilder{
		// map
	}
	builder.BuildConstraintMap(m.MApps.Apps)
	for i := range taskList {
		glog.V(3).Infof("entire Task ====================> %+v", taskList[i])
		if _, ok := taskUseMap[taskList[i].Id]; !ok {
			continue
		}
		taskProbe := &TaskProbe{
			Task:    &taskList[i],
			Cluster: &m.Cluster,
		}
		if taskProbe.Task.State != "TASK_RUNNING" {
			glog.V(4).Infof("=====> not running task is %s and state %s\n", taskProbe.Task.Name, taskProbe.Task.State)
			continue
		}
		glog.V(4).Infof("=====> task is %s and state %s\n", taskProbe.Task.Name, taskProbe.Task.State)

		builder.SetTaskConstraints(taskProbe)
		//ipAddress := slaveIdIpMap[taskProbe.Task.SlaveId]
		//usedResources := taskProbe.GetUsedResourcesForTask(ipAddress)
		taskResource, err := taskProbe.GetTaskResourceStat(m.MapTaskStatistics, taskProbe.Task, taskUseMap)
		if err != nil {
			glog.Errorf("error is : %s", err)
		}
		commoditiesSoldContainer := taskProbe.GetCommoditiesSoldByContainer(taskProbe.Task, taskResource)
		commoditiesBoughtContainer := taskProbe.GetCommoditiesBoughtByContainer(taskProbe.Task, taskResource)

		entityDTO, _ := buildTaskContainerEntityDTO(m.SlaveIdIpMap, taskProbe.Task, commoditiesSoldContainer, commoditiesBoughtContainer)

		result = append(result, entityDTO)

		commoditiesSoldApp := taskProbe.GetCommoditiesSoldByApp(taskProbe.Task, taskResource)
		commoditiesBoughtApp := taskProbe.GetCommoditiesBoughtByApp(taskProbe.Task, taskResource)

		entityDTO = buildTaskAppEntityDTO(m.SlaveIdIpMap, taskProbe.Task, commoditiesSoldApp, commoditiesBoughtApp)
		result = append(result, entityDTO)
	}
	glog.V(4).Infof("Task entity DTOs : %d", len(result))
	return result, nil
}

func buildTaskAppEntityDTO(slaveIdIp map[string]string, task *util.Task, commoditiesSold []*proto.CommodityDTO, commoditiesBoughtMap map[*builder.ProviderDTO][]*proto.CommodityDTO) *proto.EntityDTO {
	appEntityType := proto.EntityDTO_APPLICATION
	id := task.Name + "::" + "APP:" + task.Id
	dispName := "APP:" + task.Name
	entityDTOBuilder := builder.NewEntityDTOBuilder(appEntityType, id+"foo")
	entityDTOBuilder = entityDTOBuilder.DisplayName(dispName)

	entityDTOBuilder.SellsCommodities(commoditiesSold)

	for provider, commodities := range commoditiesBoughtMap {
		entityDTOBuilder.Provider(provider)
		entityDTOBuilder.BuysCommodities(commodities)
	}

	entityDto, _ := entityDTOBuilder.Create()

	//appType := task.Name
	//
	//ipAddress := slaveIdIp[task.SlaveId] //this.getIPAddress(host, nodeName)

	//appData := &proto.EntityDTO_ApplicationData{
	//	Type:      &appType,
	//	IpAddress: &ipAddress,
	//}
	//// entityDto.ApplicationData = appData // TODO:
	return entityDto

}

// Build entityDTO that contains all the necessary info of a pod.
func buildTaskContainerEntityDTO(slaveIdIpMap map[string]string, task *util.Task, commoditiesSold, commoditiesBought []*proto.CommodityDTO) (*proto.EntityDTO, error) {
	taskName := task.Name
	id := task.Id
	dispName := task.Name

	entityDTOBuilder := builder.NewEntityDTOBuilder(proto.EntityDTO_CONTAINER, id)
	entityDTOBuilder.DisplayName(dispName)

	slaveId := task.SlaveId
	if slaveId == "" {
		return nil, fmt.Errorf("Cannot find the hosting slave ID for task %s", taskName)
	}
	glog.V(4).Infof("Pod %s is hosted on %s", dispName, slaveId)

	entityDTOBuilder.SellsCommodities(commoditiesSold)
	//	providerUid := nodeUidTranslationMap[slaveId]
	providerDto := builder.CreateProvider(proto.EntityDTO_VIRTUAL_MACHINE, slaveId)
	entityDTOBuilder = entityDTOBuilder.Provider(providerDto)
	entityDTOBuilder.BuysCommodities(commoditiesBought)
	ipAddress := slaveIdIpMap[task.SlaveId]
	ipPropName := "IP"
	ipProp := &proto.EntityDTO_EntityProperty{
		Namespace: &DEFAULT_NAMESPACE,
		Name: &ipPropName,
		Value: &ipAddress,
	}
	entityDTOBuilder = entityDTOBuilder.WithProperty(ipProp)		//"ipAddress", ipAddress)
	glog.V(3).Infof("Pod %s will be stitched to VM with IP %s", dispName, ipAddress)

	entityDto, _ := entityDTOBuilder.Create()
	return entityDto, nil
}
