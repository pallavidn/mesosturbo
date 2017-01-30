package conf


type MasterRestClient interface {
	Login() (string, error)
	GetState()	(*MesosState, error)
	GetNodes()
	GetFrameworks()
}


type FrameworkRestClient interface {
	getFrameworkApps()
}

type MesosMasterType string
type MesosFrameworkType string

const (
	Apache MesosMasterType = "Apache Mesos"
	DCOS MesosMasterType = "Mesosphere DCOS"
)

const (
	Marathon MesosFrameworkType = "Marathon"
	Chronos MesosFrameworkType = "Chronos"
	Hadoop MesosFrameworkType = "Hadoop"
)


