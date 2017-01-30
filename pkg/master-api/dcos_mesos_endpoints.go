package master


type DCOSEndpointPath string
const (
	DCOS_StatePath DCOSEndpointPath = "/mesos" +  "/state"	//string(Apache_StatePath)
	DCOS_FrameworksPath DCOSEndpointPath = "/mesos"+ "/frameworks" // Apache_FrameworksPath
	DCOS_TasksPath DCOSEndpointPath = "/mesos" + "/tasks"	//Apache_TasksPath
)


type DCOSMesosEndpointStore struct {
	EndpointMap map[MesosEndpointName]DCOSEndpointPath
}

func NewDCOSMesosEndpointStore() *DCOSMesosEndpointStore {
	store := &DCOSMesosEndpointStore{
		EndpointMap: make(map[MesosEndpointName]DCOSEndpointPath),
	}

	epMap := store.EndpointMap

	epMap[State] = DCOS_StatePath
	epMap[Frameworks] = DCOS_FrameworksPath
	epMap[Tasks] = DCOS_TasksPath
	return store
}

