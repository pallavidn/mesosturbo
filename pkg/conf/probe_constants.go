package conf

// Represents the Mesos Master Vendor
type MesosMasterType string
const (

	Apache MesosMasterType = "Apache Mesos"
	DCOS   MesosMasterType = "Mesosphere DCOS"
)

//
//type MesosMaster int
//var  MesosMasters []string
//
//func (master MesosMaster) String() string {
//	return MesosMasters[int(master)]
//}
//func ConvertToMesosMaster(s string) MesosMaster {
//	MesosMasters = append(MesosMasters, s)
//	fmt.Printf("mesosMasters=%s\n", MesosMasters)
//	r := MesosMaster(len(MesosMasters) -1)
//	fmt.Printf("r=%s\n", r)
//	return r
//}
//
//var (
//	Apache_Master  = ConvertToMesosMaster("Apache Mesos")	//"Apache Mesos"
//	DCOS_Master    = ConvertToMesosMaster("Mesosphere DCOS") //"Mesosphere DCOS"
//)


// ==========================================================================

// Represents the Framework used by the Mesos Master
type MesosFrameworkType string

const (
	Marathon      MesosFrameworkType = "Marathon"
	DCOS_Marathon MesosFrameworkType = "DCOS Marathon"
	Chronos       MesosFrameworkType = "Chronos"
	Hadoop        MesosFrameworkType = "Hadoop"
)

// Constants for the Mesos probe account definition fields in the Turbo server
type ProbeAcctDefEntryName string

const (
	MasterIP       ProbeAcctDefEntryName = "MasterIP"
	MasterPort     ProbeAcctDefEntryName = "MasterPort"
	MasterUsername ProbeAcctDefEntryName = "Username"
	MasterPassword ProbeAcctDefEntryName = "Password"

	FrameworkIP       ProbeAcctDefEntryName = "FrameworkIP"
	FrameworkPort     ProbeAcctDefEntryName = "FrameworkPort"
	FrameworkUsername ProbeAcctDefEntryName = "FrameworkUsername"
	FrameworkPassword ProbeAcctDefEntryName = "FrameworkPassword"

	ActionIP   ProbeAcctDefEntryName = "ActionIP"
	ActionPort ProbeAcctDefEntryName = "ActionPort"
)
