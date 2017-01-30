package probe

type ProbeAcctDefEntryName string
const (
	MasterIP ProbeAcctDefEntryName = "MasterIP"
	MasterPort ProbeAcctDefEntryName = "MasterPort"
	Username ProbeAcctDefEntryName = "Username"
	Password ProbeAcctDefEntryName = "Password"

	FrameworkIP ProbeAcctDefEntryName = "FrameworkIP"
	FrameworkPort ProbeAcctDefEntryName = "FrameworkPort"
	FrameworkUsername ProbeAcctDefEntryName = "FrameworkUsername"
	FrameworkPassword ProbeAcctDefEntryName = "FrameworkPassword"

	ActionIP ProbeAcctDefEntryName = "ActionIP"
	ActionPort ProbeAcctDefEntryName = "ActionPort"
)
