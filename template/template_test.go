package templates

import (
	"github.com/kube-sailmaker/template-gen/test"
	"testing"
)

func TestLoadTemplates(t *testing.T) {
	cmd := []string{"/bin/sh", "sleep"}

	application := Application{
		ReleaseName:    "apps",
		Name:           "busybox",
		Tag:            "latest",
		Annotations:    nil,
		Replicas:       "1",
		LivenessProbe:  "/health",
		ReadinessProbe: "/ready",
		EnvVars:        nil,
		Limits:         nil,
		Command:        cmd,
		Entrypoint:     nil,
	}

	template, err := LoadTemplates("DeploymentTemplate", &application)
	test.Null(t, err)
	test.NotNull(t, template)
	test.EqualTo(t, "busybox-deployment.yaml", template.Name())
}
