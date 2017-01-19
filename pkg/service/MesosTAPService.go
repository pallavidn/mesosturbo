package probe

import (
	"fmt"

	// turbo go sdk
	"github.com/turbonomic/turbo-go-sdk/pkg/probe"

	mesos "github.com/turbonomic/mesosturbo/pkg/probe"
)

type MesosTAPService struct {
	// Interface to the Turbo Server
	Probe *probe.TurboProbe
	// Interface to the Mesos world
}

// Create the TurboProbe responsible for communicating with the Turbo server
func NewMesosTAPService(probeConf *probe.ProbeConfig, vmtApiConf *probe.TurboAPIConfig) *MesosTAPService {
	mesosTapService := &MesosTAPService{
	}

	// Create the Probe for the turbo server
	mesosTapService.Probe = probe.NewTurboProbe(probeConf)

	// Turbo Rest API Handler
	turboApiHandler := probe.NewTurboAPIHandler(vmtApiConf)
	mesosTapService.Probe.SetTurboAPIHandler(turboApiHandler)

	// Mesos Probe Registration Client
	registrationClient := &mesos.MesosRegistrationClient{}
	mesosTapService.Probe.SetProbeRegistrationClient(registrationClient)

	fmt.Printf("[MesosTAPService] : Created MesosTAPService %s\n", mesosTapService)
	return mesosTapService
}

// Create the Target that will execute the discovery for the Turbo server
func (mesosTAPService *MesosTAPService) CreateMesosTAPServiceAgent(targetIdentifier string, confFile string)  {
	if mesosTAPService.Probe == nil {
		fmt.Println("[MesosTAPService] : Mesos Turbo Probe is null")
		return
	}
	// Discovery client for Mesos Target
	discoveryClient := mesos.NewDiscoveryClient(targetIdentifier, confFile)
	mesosTAPService.Probe.SetDiscoveryClient(targetIdentifier, discoveryClient)

	fmt.Printf("[MesosTAPService] : Created MesosDiscoveryClient %s ", discoveryClient)
}


