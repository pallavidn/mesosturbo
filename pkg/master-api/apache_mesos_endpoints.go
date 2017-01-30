package master


type ApacheMesosEndpointPath string
const (
	Apache_StatePath ApacheMesosEndpointPath = "/state"
	Apache_FrameworksPath ApacheMesosEndpointPath = "/frameworks"
	Apache_TasksPath ApacheMesosEndpointPath = "/tasks"
)

type ApacheMesosEndpointStore struct {
	EndpointMap map[MesosEndpointName]ApacheMesosEndpointPath
}

func NewApacheMesosEndpointStore() *ApacheMesosEndpointStore {
	store := &ApacheMesosEndpointStore{
			EndpointMap: make(map[MesosEndpointName]ApacheMesosEndpointPath),
		}

	epMap := store.EndpointMap

	epMap[State] = Apache_StatePath
	epMap[Frameworks] = Apache_FrameworksPath
	epMap[Tasks] = Apache_TasksPath
	return store
}



type MesosEndpoint struct {
	Path string
	Parser EndpointParser
}
//
//type EndpointRequest interface {
//	GetRequest() *http.Request
//}

