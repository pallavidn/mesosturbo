package conf

import "github.com/turbonomic/mesosturbo/pkg/util"

// Interface for the client to handle Rest API communication with the Mesos Master
type MasterRestClient interface {
	Login() (string, error)
	GetState() (*util.MesosAPIResponse, error)	//(*MesosState, error)
	GetNodes()
	GetFrameworks()
}


// Interface for the client to handle Rest API communication with the Mesos Framework
type FrameworkRestClient interface {
	GetFrameworkApps() (*util.FrameworkApps, error)
}

