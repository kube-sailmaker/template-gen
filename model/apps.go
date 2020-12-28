package model

type Args struct {
	ManifestDir string
	Namespace   string
	Env         string
	ReleaseName string
}

type Application struct {
	Name                    string                   `yaml:"name"`
	Kind                    string                   `yaml:"kind"`
	LivenessProbe           string                   `yaml:"liveness_probe"`
	ReadinessProbe          string                   `yaml:"readiness_probe"`
	Annotations             map[string]string        `yaml:"annotations"`
	Resources               []string                 `yaml:"resources"`
	Capabilities            []string                 `yaml:"capabilities"`
	Mixins                  []string                 `yaml:"mixins"`
	Template                []AppTemplate            `yaml:"template"`
	Service                 ServiceSpec              `yaml:"service"`
	BackoffLimit            int                      `yaml:"backoffLimit"`
	ActiveDeadLine          int                      `yaml:"activeDeadlineSeconds"`
	TTLSecondsAfterFinished int                      `yaml:"ttlSecondsAfterFinished"`
	RestartPolicy           string                   `yaml:"restartPolicy"`
	VolumeMounts            []Mount                  `yaml:"volumeMounts"`
	Volumes                 []map[string]interface{} `yaml:"volumes"`
	ConfigMaps              []map[string]string      `yaml:"configMaps"`
	NodeSelector            map[string]string        `yaml:"nodeSelector"`
}

type ServiceSpec struct {
	Enabled    bool   `yaml:"enabled", default: false`
	Name       string `yaml:"name"`
	Type       string `yaml:"type"`
	Port       int    `yaml:"port"`
	TargetPort int    `yaml:"targetPort"`
}

type AppTemplate struct {
	Name    string            `yaml:"name"`
	Replica int               `yaml:"replica"`
	Config  map[string]string `yaml:"config"`
}

type Mount struct {
	Name      string `yaml:"name"`
	MountPath string `yaml:"mountPath"`
}
