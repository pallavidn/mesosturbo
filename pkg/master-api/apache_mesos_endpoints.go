package master

// Endpoint paths for Apache Mesos Master
type ApacheMesosEndpointPath string

const (
	Apache_StatePath      ApacheMesosEndpointPath = "/state"
	Apache_FrameworksPath ApacheMesosEndpointPath = "/frameworks"
	Apache_TasksPath      ApacheMesosEndpointPath = "/tasks"
)

// Endpoint store containing endpoint and parsers for Apache Mesos Master
func NewApacheMesosEndpointStore() *MasterEndpointStore {
	store := &MasterEndpointStore{
		EndpointMap: make(map[MasterEndpointName]*MasterEndpoint),
	}

	epMap := store.EndpointMap

	epMap[State] = &MasterEndpoint{
		EndpointName: string(State),
		EndpointPath: string(Apache_StatePath),
		Parser:       &GenericMasterStateParser{},
	}
	epMap[Frameworks] = &MasterEndpoint{
		EndpointName: string(Frameworks),
		EndpointPath: string(Apache_FrameworksPath),
		Parser:       &GenericMasterStateParser{},
	}
	epMap[Tasks] = &MasterEndpoint{
		EndpointName: string(Tasks),
		EndpointPath: string(Apache_TasksPath),
		Parser:       &GenericMasterStateParser{},
	}

	return store
}
