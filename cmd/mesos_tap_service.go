package main

import (
	goflag "flag"
	"fmt"
	"os"
	"github.com/golang/glog"
	"github.com/turbonomic/turbo-go-sdk/pkg/service"
	"github.com/turbonomic/turbo-go-sdk/pkg/probe"

	mesos "github.com/turbonomic/mesosturbo/pkg/probe"
	"github.com/turbonomic/mesosturbo/pkg/conf"
	"github.com/turbonomic/mesosturbo/pkg/factory"
	//"time"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

func init() {
	goflag.Set("logtostderr", "true")
}

func main() {
	goflag.Parse()

	probeCategory := "CloudNative"
	targetType := "Apache Mesos"	//"MesosProbe"

	targetConf := 	 "src/github.com/turbonomic/mesosturbo/cmd/apache-mesos-target-conf.json"
	//targetConf :=  "src/github.com/turbonomic/mesosturbo/cmd/dcos-target-conf.json"

	turboCommConf := "src/github.com/turbonomic/mesosturbo/cmd/container-conf.json"

	mesosTargetConf, err := conf.NewMesosTargetConf(targetConf)
	if err != nil {
		glog.Errorf("Cannot start Mesos TAP service, invalid config : %s\n", err.Error())
		os.Exit(1)
	}
	mesosMasterType := mesosTargetConf.Master
	target1 := mesosTargetConf.MasterIP

	// --------------------------------------------------------------------
	client := factory.GetMasterRestClient(mesosMasterType, mesosTargetConf)
	if client == nil {
		glog.Errorf("Cannot find RestClient for Mesos : ", mesosTargetConf.Master)
	} else {
		glog.Infof(">>>>>>>>>>>> Found RestClient for Mesos : %s", client)
	}

	if client != nil {
		token, err := client.Login()
		if err != nil {
			glog.Errorf("Error logging to " + string(mesosMasterType) + "::" + mesosTargetConf.MasterIP + "\n", err.Error())
			glog.Flush()
		}
		mesosTargetConf.Token = token
		mesosState, err := client.GetState()
		if err != nil {
			glog.Errorf("Error getting state from master : %s\n", err)
		}
		glog.V(2).Infof("MesosState %s", mesosState)

		//// Framework response
		//// Based on the Framework vendor, instantiate the Frameworks RestClient
		//frameworkClient := factory.GetFrameworkRestClient(mesosTargetConf.Framework, mesosTargetConf)
		//if frameworkClient == nil {
		//	glog.Errorf("Cannot find framework Client for Mesos : ", mesosTargetConf.Framework)
		//} else {
		//	glog.Infof(">>>>>>>>>>>> Found RestClient for Framework : %s", client)
		//}
		//
		//frameworkResp, err := frameworkClient.GetFrameworkApps()
		//
		//if err != nil {
		//	glog.Errorf("Error parsing response from the Framework: %s \n", err)
		//}
		//glog.V(2).Infof("Marathon Get Succeeded: %v\n", frameworkResp)

	}
	// --------------------------------------------------------------------

	fmt.Println("============================================================================================")

	turboCommConfigData, err := service.ParseTurboCommunicationConfig(turboCommConf)
	if turboCommConfigData == nil || err != nil {
		glog.Errorf("Cannot start Mesos TAP service, invalid turbo server communication config : %s\n", err.Error())
		os.Exit(1)
	}
	// Mesos Probe Registration Client
	registrationClient := mesos.NewRegistrationClient(mesosMasterType)

	// Mesos Probe Discovery Client
	discoveryClient, err := mesos.NewDiscoveryClient(mesosMasterType,  mesosTargetConf)

	if err != nil {
		glog.Errorf("Error creating discovery client for " + string(mesosMasterType) + "::" + mesosTargetConf.MasterIP +"\n", err.Error())
		os.Exit(1)
	}

	//var accountValues []*proto.AccountValue
	//accountValues = discoveryClient.GetAccountValues().AccountValues
	// Convert target conf to account values
	//discoveryClient.Discover(accountValues)
	targetType = string(conf.Apache)
	probeConfig := &probe.ProbeConfig{
		ProbeCategory: probeCategory,
		ProbeType: targetType,
		//NewConf: conf.CreateMesosTargetConf,
	}

	turboProbe, err := probe.NewTurboProbe(probeConfig)
	turboProbe.SetProbeRegistrationClient(registrationClient)

	glog.Infof("Created turbo probe : %s", turboProbe)

	accountValues := createApacheAccValues()
	//target := turboProbe.GetTurboTarget(accountValues)

	//accountValues = target.GetAccountValues()
	glog.Infof("Target account Values : %s", accountValues)

	targetType = string(conf.DCOS)
	probeConfig = &probe.ProbeConfig{
		ProbeCategory: probeCategory,
		ProbeType: targetType,
		//NewConf: conf.CreateMesosTargetConf,
	}
	glog.Infof("************************************************")
	turboProbe, err = probe.NewTurboProbe(probeConfig)
	turboProbe.SetProbeRegistrationClient(registrationClient)

	glog.Infof("Created turbo probe : %s", turboProbe)

	accountValues = createDCOSAccValues()
	//target = turboProbe.GetTurboTarget(accountValues)

	//accountValues = target.GetAccountValues()
	glog.Infof("Target account Values : %s", accountValues)
	// ============================================
	tapService, err :=
		service.NewTAPServiceBuilder().
			WithTurboCommunicator(turboCommConfigData).
			WithTurboProbe(probe.NewProbeBuilder(targetType, probeCategory).
						RegisteredBy(registrationClient).
						DiscoversTarget(target1, discoveryClient)).
			Create()

	if (err != nil) {
		glog.Errorf("Error creating TAP Service : ", err)

	}
	//
	// Connect to the Turbo server
	go tapService.ConnectToTurbo()
	glog.Infof("Connected to Turbo")

	// =============================================
	//t := time.NewTimer(time.Minute * 2)
	//for {
	//	select {
	//	case <-t.C:
	//		glog.Infof("============= END TAP SERVICE =============")
	//		tapService.DisconnectFromTurbo()
	//		break
	//	default:
	//	// do nothing
	//
	//	}
	//}
	//
	//glog.Infof("DisConnected to Turbo")
	select {}

} //end main


func createApacheAccValues() []*proto.AccountValue {
	var accountValues []*proto.AccountValue
	prop1 := string(conf.MasterIP)
	val1 := "10.10.10.10"
	accVal := &proto.AccountValue{
		Key: &prop1,
		StringValue: &val1,
	}
	accountValues = append(accountValues, accVal)

	prop2 := string(conf.MasterPort)
	val2 := "5050"
	accVal = &proto.AccountValue{
		Key: &prop2,
		StringValue: &val2,
	}
	accountValues = append(accountValues, accVal)

	prop3 := string(conf.MasterUsername)
	val3 := "pallavi.debnath"
	accVal = &proto.AccountValue{
		Key: &prop3,
		StringValue: &val3,
	}
	accountValues = append(accountValues, accVal)

	prop4 := string(conf.MasterPassword)
	val4 := "sysdreamworks"
	accVal = &proto.AccountValue{
		Key: &prop4,
		StringValue: &val4,
	}
	accountValues = append(accountValues, accVal)
	return accountValues
}

func createDCOSAccValues() []*proto.AccountValue {
	var accountValues []*proto.AccountValue
	prop1 := string(conf.MasterIP)
	val1 := "11.11.11.11"
	accVal := &proto.AccountValue{
		Key: &prop1,
		StringValue: &val1,
	}
	accountValues = append(accountValues, accVal)

	prop2 := string(conf.MasterPort)
	val2 := "8080"
	accVal = &proto.AccountValue{
		Key: &prop2,
		StringValue: &val2,
	}
	accountValues = append(accountValues, accVal)

	prop3 := string(conf.MasterUsername)
	val3 := "pallavi.debnath-2"
	accVal = &proto.AccountValue{
		Key: &prop3,
		StringValue: &val3,
	}
	accountValues = append(accountValues, accVal)

	prop4 := string(conf.MasterPassword)
	val4 := "sysdreamworks"
	accVal = &proto.AccountValue{
		Key: &prop4,
		StringValue: &val4,
	}
	accountValues = append(accountValues, accVal)
	return accountValues
}