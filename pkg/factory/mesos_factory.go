package factory


import (
	conf "github.com/turbonomic/mesosturbo/pkg/conf"
	master "github.com/turbonomic/mesosturbo/pkg/master-api"
	// scheduler "github.com/turbonomic/mesosturbo/pkg/framework-api"
	"fmt"
)

func GetMasterRestClient(mesosType conf.MesosMasterType, masterIP, masterPort, username, password string) conf.MasterRestClient {
	if mesosType == conf.Apache {
		fmt.Println("[GetMasterRestClient] Creating Apache Mesos Master Client")
		return master.NewApacheMesosRestClient(masterIP, masterPort, username, password)
	} else if mesosType == conf.DCOS {
		fmt.Println("[GetMasterRestClient] Creating DCOS Mesos Master Client")
		return master.NewDCOSMesosRestClient(masterIP, masterPort, username, password)
	}
	return nil
}

//
//func GetFrameworkRestClient(framework conf.MesosFramework, frameworkIP, frameworkPort string) conf.FrameworkRestClient {
//	if framework == conf.Marathon {
//		return scheduler.NewMarathonRestClient(frameworkIP, frameworkPort)
//	}
//	return nil
//}
