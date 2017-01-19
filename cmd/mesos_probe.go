package main

import (
	goflag "flag"

	container "github.com/turbonomic/turbo-go-sdk/pkg/communication"
	probe "github.com/turbonomic/turbo-go-sdk/pkg/probe"
	service "github.com/turbonomic/mesosturbo/pkg/service"

	//mesoshttp "github.com/turbonomic/mesosturbo/pkg/mesoshttp"
)

var dcosToken string
func init() {
	goflag.Set("logtostderr", "true")
	goflag.StringVar(&dcosToken, "token", "", "dcos token")
}

func main() {
	goflag.Parse()
	//switch {
	//case 1 : glog.V(3).Infof("Test")
	//}
	// Read one conf file for the application
	// Create TAPGoApplication

	probeCategory := "CloudNative"
	targetType := "MesosProbe"
	targetConf := 	 "src/github.com/turbonomic/mesosturbo/cmd/mesos-target-conf.json"
	containerConf := "src/github.com/turbonomic/mesosturbo/cmd/container-conf.json"
	target1 := "10.10.174.91"

	vmtServerAddress := "127.0.0.1:8080"
	vmtUser := "administrator"
	vmtPwd := "admin"

	probeConf := &probe.ProbeConfig {
		ProbeCategory: probeCategory,
		ProbeType: targetType,
	}

	apiConf := &probe.TurboAPIConfig {
		VmtServerAddress: vmtServerAddress,
		VmtUser: vmtUser,
		VmtPassword: vmtPwd,
	}

	// The Probe
	MesosProbe :=  service.NewMesosTAPService(probeConf, apiConf)
	// The Targets
	MesosProbe.CreateMesosTAPServiceAgent(target1, targetConf)

	// The Container
	theContainer := container.CreateMediationContainer(containerConf)

	// Load the probe in the container
	theContainer.LoadProbe(MesosProbe.Probe)
	theContainer.GetProbe(targetType)

	// Connect to the Turbo server
	theContainer.Init()

	//MesosProbe.Probe.AddTarget()

	//theContainer.CreateTarget(targetConfFile)
	//TAPGoApplication.start()

	//// --------- Http client for DCOS
	//mesosClientConfig, err := mesoshttp.NewConnectionClient("src/github.com/turbonomic/mesosturbo/cmd/mesos-target-conf.json")
	//if err != nil {
	//	fmt.Println("Error loading configuration for mesos client")
	//}
	//
	//if mesosClientConfig.DCOS {
	//	// check DCOS username and password work
	//	mesosAPIClient := &mesoshttp.MesosHTTPClient{}
	//	err := mesosAPIClient.DCOSLoginRequest(mesosClientConfig, dcosToken)
	//	if err != nil {
	//		return
	//	}
	//}

	select {}

} //end main

