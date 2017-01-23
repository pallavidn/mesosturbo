package main

import (
	goflag "flag"

	probe "github.com/turbonomic/turbo-go-sdk/pkg/probe"
	service "github.com/turbonomic/turbo-go-sdk/pkg/service"

	mesos "github.com/turbonomic/mesosturbo/pkg/probe"
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
	turboCommConf := "src/github.com/turbonomic/mesosturbo/cmd/container-conf.json"
	target1 := "10.10.174.91"

	// Mesos Probe Registration Client
	registrationClient := &mesos.MesosRegistrationClient{}
	// Mesos Probe Discovery Client
	discoveryClient := mesos.NewDiscoveryClient(target1, targetConf)
	tapService := service.NewTAPServiceBuilder().
				WithTurboCommunicator(turboCommConf).
				WithTurboProbe(probe.NewProbeBuilder(targetType, probeCategory).
							RegisteredBy(registrationClient).
							DiscoversTarget(target1, discoveryClient)).
				Create()

	// Connect to the Turbo server
	tapService.ConnectToTurbo()

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

