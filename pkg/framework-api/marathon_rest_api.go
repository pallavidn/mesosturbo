package framework

import "github.com/turbonomic/mesosturbo/pkg/conf"

type MarathonRestClient struct {
	FrameworkIP 	string
	FrameworkPort 	string
	//Username 	string
	//Password 	string
}

func NewMarathonRestClient (frameworkIP, frameworkPort string) conf.FrameworkRestClient {
	return &MarathonRestClient {
		FrameworkIP: frameworkIP,
		FrameworkPort: frameworkPort,
		//Username: username,
		//Password: password,
	}
}


func (frameworkClient *MarathonRestClient) getFrameworkApps() {

}
