package util

type FrameworkApps struct {
	Apps []FrameworkApp `json:"apps"`
}

type FrameworkApp struct {
	Id                         string        `json:"id"`
	Cmd                        string        `json:"cmd"`
	Args                       []interface{} `json:"args"`
	User                       string        `json:"user"`
	Env                        interface{}   `json:"env"`
	Instances                  int           `json:"instances"`
	Cpus                       float64       `json:"cpus"`
	Mem                        int           `json:"mem"`
	Disk                       int           `json:"disk"`
	GPUS                       int           `json:"gpus"`
	Executor                   string        `json:"executor"`
	Constraints                []interface{} `json:"constraints"`
	URIs                       []interface{} `json:"uris"`
	Fetch                      []interface{} `json:"fetch"`
	StoreUrls                  []interface{} `json:"storeUrls"`
	BackoffSeconds             int           `json:"backoffSeconds"`
	BackoffFactor              float64       `json:"backoffFactor"`
	MaxLaunchDelaySeconds      int           `json:"maxLaunchDelaySeconds"`
	Container                  interface{}   `json:"container"`
	HealthChecks               []interface{} `json:"healthChecks"`
	ReadinessChecks            []interface{} `json:"readinessChecks"`
	Dependencies               []interface{} `json:"dependencies"`
	UpgradeStrategy            interface{}   `json:"upgradeStrategy"`
	Labels                     interface{}   `json:"labels"`
	AcceptedResourceRoles      string        `json:"acceptedResourceRoles"`
	IPAddress                  string        `json:"ipAddress"`
	Version                    string        `json:"version"`
	Residency                  string        `json:"residency"`
	Secrets                    string        `json:"secrets"`
	TaskKillGracePeriodSeconds string        `json:"taskKillGracePeriodSeconds"`
	Ports                      []interface{} `json:"ports"`
	PortDefinitions            []interface{} `json:"portDefinitions"`
	RequirePorts               bool          `json:"requirePorts"`
	VersionInfo                interface{}   `json:"versionInfo"`
	TasksStaged                int           `json:"tasksStaged"`
	TasksRunning               int           `json:"tasksRunning"`
	TasksHealthy               int           `json:"tasksHealthy"`
	TasksUnhealthy             int           `json:"tasksUnhealthy"`
	Deployments                []interface{} `json:"deployments"`
}
