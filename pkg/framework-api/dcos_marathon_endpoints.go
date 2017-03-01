package framework

type DCOSMarathonEndpointPath string
const (
	DCOS_Marathon_AppsPath DCOSMarathonEndpointPath = "/service/marathon/v2/apps"
)

func NewDCOSMarathonEndpointStore() *FrameworkEndpointStore {
	store := &FrameworkEndpointStore{}
	store.EndpointMap = make(map[FrameworkEndpointName]*FrameworkEndpoint)

	epMap := store.EndpointMap

	endpoint := &FrameworkEndpoint{
		EndpointName: string(Apps),
		EndpointPath: string(DCOS_Marathon_AppsPath),
		Parser: &FrameworkAppsParser{},
	}

	epMap[Apps] = endpoint
	return store
}


