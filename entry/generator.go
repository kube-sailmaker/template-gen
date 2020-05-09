package entry

import (
	"github.com/kube-sailmaker/template-gen/model"
	"github.com/kube-sailmaker/template-gen/task"
	templates "github.com/kube-sailmaker/template-gen/template"
)

func TemplateGenerator(appSpec *model.AppSpec, appDir string, resourceDir string, outputDir string) (*model.DeploymentItemSummary, error) {
	appTemplate := make([]templates.Application, 0)

	validationErr := appSpec.Validate()
	if validationErr != nil {
		return nil, validationErr
	}
	appSpec.Normalise()

	for _, app := range appSpec.Apps {
		application, err := task.ProcessApplication(&app, appSpec.ReleaseName, appSpec.Namespace, appSpec.Environment, appDir, resourceDir)
		if err != nil {
			return nil, err
		}
		appTemplate = append(appTemplate, *application)
	}

	releaseTemplate := templates.ReleaseTemplate{
		Namespace:   appSpec.Namespace,
		Application: appTemplate,
	}
	return templates.Run(&releaseTemplate, outputDir)
}
