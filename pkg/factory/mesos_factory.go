package factory


import (
	"github.com/golang/glog"
	"github.com/turbonomic/mesosturbo/pkg/conf"
	master "github.com/turbonomic/mesosturbo/pkg/master-api"
	scheduler "github.com/turbonomic/mesosturbo/pkg/framework-api"
)

// Get the Rest API client to handle communication with the Mesos Master
func GetMasterRestClient(mesosType conf.MesosMasterType, mesosConf *conf.MesosTargetConf) conf.MasterRestClient {
	var endpointStore *master.MasterEndpointStore
	if mesosType == conf.Apache {
		glog.V(2).Infof("[GetMasterRestClient] Creating Apache Mesos Master Client")
		endpointStore = master.NewApacheMesosEndpointStore()
	} else if mesosType == conf.DCOS {
		glog.V(2).Infof("[GetMasterRestClient] Creating DCOS Mesos Master Client")
		endpointStore = master.NewDCOSMesosEndpointStore()
	}

	if endpointStore == nil {
		glog.Errorf("[GetMasterRestClient] Unsupported Mesos Master ", mesosType)
		return nil
	}

	return master.NewGenericMasterAPIClient(mesosConf, endpointStore)
}


// Get the Rest API client to handle communication with the given Mesos Framework
func GetFrameworkRestClient(framework conf.MesosFrameworkType, mesosConf *conf.MesosTargetConf) conf.FrameworkRestClient {
	var endpointStore *scheduler.FrameworkEndpointStore
	if framework == conf.Marathon {
		glog.V(2).Infof("[GetFrameworkRestClient] Creating Apache Marathon Client")
		endpointStore = scheduler.NewApacheMarathonEndpointStore()
	} else if framework == conf.DCOS_Marathon {
		glog.V(2).Infof("[GetFrameworkRestClient] Creating DCOS Marathon Client")
		endpointStore = scheduler.NewDCOSMarathonEndpointStore()
	}

	if endpointStore == nil {
		glog.Errorf("[GetFrameworkRestClient] Unsupported Mesos Framework ", framework)
		return nil
	}

	return scheduler.NewFrameworkAPIClient(mesosConf, endpointStore)
}
