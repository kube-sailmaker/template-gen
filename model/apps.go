package model

type Args struct {
	ManifestDir string
	Namespace   string
	Env         string
	ReleaseName string
}

type Application struct {
	Name           string            `yaml:"name"`
	LivenessProbe  string            `yaml:"liveness_probe"`
	ReadinessProbe string            `yaml:"readiness_probe"`
	Annotations    map[string]string `yaml:"annotations"`
	Resources      []string          `yaml:"resources"`
	Capabilities   []string          `yaml:"capabilities"`
	Mixins         []string          `yaml:"mixins"`
	Template       []AppTemplate     `yaml:"template"`
	Service        *ServiceSpec       `yaml:"service"`
}

type ServiceSpec struct {
	Enabled bool `yaml:"enabled"`
	Port    int  `yaml:"port"`
}

type AppTemplate struct {
	Name    string            `yaml:"name"`
	Replica int               `yaml:"replica"`
	Config  map[string]string `yaml:"config"`
}
