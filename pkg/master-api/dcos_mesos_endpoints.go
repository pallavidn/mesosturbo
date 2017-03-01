package master

// Endpoint paths for Apache Mesos Master
type DCOSEndpointPath string

const (
	DCOS_StatePath      DCOSEndpointPath = "/mesos" + "/state"
	DCOS_FrameworksPath DCOSEndpointPath = "/mesos" + "/frameworks"
	DCOS_TasksPath      DCOSEndpointPath = "/mesos" + "/tasks"
	DCOS_LoginPath      DCOSEndpointPath = "/acs/api/v1/auth/login"
)

// Endpoint store containing endpoint and parsers for DCOS Mesos Master
func NewDCOSMesosEndpointStore() *MasterEndpointStore {
	store := &MasterEndpointStore{
		EndpointMap: make(map[MasterEndpointName]*MasterEndpoint),
	}

	epMap := store.EndpointMap

	epMap[Login] = &MasterEndpoint{
		EndpointName: string(Login),
		EndpointPath: string(DCOS_LoginPath),
		Parser:       &DCOSLoginParser{},
	}

	epMap[State] = &MasterEndpoint{
		EndpointName: string(State),
		EndpointPath: string(DCOS_StatePath),
		Parser:       &GenericMasterStateParser{},
	}
	epMap[Frameworks] = &MasterEndpoint{
		EndpointName: string(Frameworks),
		EndpointPath: string(DCOS_FrameworksPath),
	}
	epMap[Tasks] = &MasterEndpoint{
		EndpointName: string(Tasks),
		EndpointPath: string(DCOS_TasksPath),
	}

	return store
}
