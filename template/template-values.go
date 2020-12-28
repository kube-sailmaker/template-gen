package templates

type ReleaseTemplate struct {
	Namespace   string
	Environment string
	Application []Application
}

type Application struct {
	Kind                    string
	ReleaseName             string
	Namespace               string
	Name                    string
	Tag                     string
	Annotations             map[string]string
	Replicas                string
	LivenessProbe           string
	ReadinessProbe          string
	EnvVars                 map[string]string
	Limits                  map[string]string
	Command                 []string
	Entrypoint              []string
	ContainerPort           int
	Service                 ServiceSpec
	Parallelism             int
	BackoffLimit            int
	ActiveDeadLine          int
	TTLSecondsAfterFinished int
	RestartPolicy           string
	VolumeMounts            []Mount
	Volumes                 []map[string]interface{}
	ConfigMaps              []map[string]interface{}
	NodeSelector            map[string]string
}

type ServiceSpec struct {
	Enabled    bool
	Name       string
	Type       string
	Port       int
	TargetPort int
}

type Mount struct {
	Name      string
	MountPath string
}
