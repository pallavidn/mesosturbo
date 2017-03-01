package framework

type ApacheMarathonEndpointPath string
const (
	Apache_Marathon_AppsPath ApacheMarathonEndpointPath = "/v2/apps"
)

func NewApacheMarathonEndpointStore() *FrameworkEndpointStore {
	store := &FrameworkEndpointStore {}
	store.EndpointMap = make(map[FrameworkEndpointName]*FrameworkEndpoint)

	epMap := store.EndpointMap

	epMap[Apps] = &FrameworkEndpoint {
				EndpointName: string(Apps),
				EndpointPath: string(Apache_Marathon_AppsPath),
				Parser: &FrameworkAppsParser{},
			}
	return store
}

