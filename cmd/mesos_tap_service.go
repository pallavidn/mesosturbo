package main

import (
	goflag "flag"

	probe "github.com/turbonomic/turbo-go-sdk/pkg/probe"
	service "github.com/turbonomic/turbo-go-sdk/pkg/service"

	mesos "github.com/turbonomic/mesosturbo/pkg/probe"
	"github.com/turbonomic/mesosturbo/pkg/conf"
	"github.com/turbonomic/mesosturbo/pkg/factory"
	"fmt"
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

	probeCategory := "CloudNative"
	targetType := "MesosProbe"

	targetConf := 	 "src/github.com/turbonomic/mesosturbo/cmd/dcos-target-conf.json"
	target1 := "dcosentpr-elasticl-1y19vdkfxkx0s-1027774118.us-west-1.elb.amazonaws.com"

	//target1 := "10.10.174.91"
	//targetConf :=  "src/github.com/turbonomic/mesosturbo/cmd/apache-mesos-target-conf.json"

	turboCommConf := "src/github.com/turbonomic/mesosturbo/cmd/container-conf.json"

	mesosTargetConf, _ := conf.NewMesosTargetConf(targetConf)
	mesosMasterType := mesosTargetConf.Master
	client := factory.GetMasterRestClient(mesosTargetConf.Master, mesosTargetConf.MasterIP, mesosTargetConf.MasterPort, "", "")

	if client == nil {
		fmt.Println("Cannot find RestClient for Mesos : ", mesosTargetConf.Master)
	} else {
		fmt.Println(">>>>>>>>>>>> Found RestClient for Mesos : %s", client)
	}


	// Mesos Probe Registration Client
	registrationClient := mesos.NewRegistrationClient(mesosMasterType)
	// Mesos Probe Discovery Client
	discoveryClient := mesos.NewDiscoveryClient(mesosMasterType, target1, targetConf)

	tapService :=
		service.NewTAPServiceBuilder().
			WithTurboCommunicator(turboCommConf).
			WithTurboProbe(probe.NewProbeBuilder(targetType, probeCategory).
			RegisteredBy(registrationClient).
			DiscoversTarget(target1, discoveryClient)).
			Create()

	//// Connect to the Turbo server
	tapService.ConnectToTurbo()

	//
	//
	//// Connect to Mesos Master
	//
	select {}

} //end main

