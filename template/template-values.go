package templates

type ReleaseTemplate struct {
	Namespace   string
	Environment string
	Application []Application
}

type Application struct {
	ReleaseName    string
	Name           string
	Tag            string
	Annotations    map[string]string
	Replicas       string
	LivenessProbe  string
	ReadinessProbe string
	EnvVars        map[string]string
	Limits         map[string]string
	Command        []string
	Entrypoint     []string
	ContainerPort  int
	ServiceEnabled bool
}
