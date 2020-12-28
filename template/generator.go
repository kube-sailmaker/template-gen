package templates

import (
	"errors"
	"fmt"
	"github.com/kube-sailmaker/template-gen/model"
	"log"
	"os"
	"path/filepath"
)

const (
	DEPLOY         = "DeploymentTemplate"
	JOB            = "JobTemplate"
	CONFIGMAP      = "ConfigMapTemplate"
	SERVICE        = "ServiceTemplate"
	SERVICEACCOUNT = "ServiceAccountTemplate"
)

var templateMap = map[string]string{"deployment": DEPLOY, "job": JOB}

func createDirSafely(fileName string) error {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			return merr
		}
	}
	return nil
}

func Run(releaseTemplate *ReleaseTemplate, outputDir string) (*model.DeploymentItemSummary, error) {

	itemSummary := model.DeploymentItemSummary{}
	items := make([]model.DeploymentItem, 0)
	for _, application := range releaseTemplate.Application {
		appWorkDir := fmt.Sprintf("%s/%s/", outputDir, application.Name)
		cerr := createDirSafely(appWorkDir)
		if cerr != nil {
			return nil, cerr
		}
		log.Println("Generating template for: ", application.Name)

		requiredTemplates, kind := GetRequiredTemplates(&application)
		for _, tName := range requiredTemplates {
			tmpl, err := LoadTemplates(tName, &application)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("[app]: %s, [error]: %v", application.Name, err))
			}
			file, er := os.Create(fmt.Sprintf("%s/%s", appWorkDir, tmpl.Name()))
			if er != nil {
				return nil, er
			}

			exErr := tmpl.Execute(file, &application)
			if exErr != nil {
				return nil, errors.New(fmt.Sprintf("[app]: %s, [error]: %v", application.Name, err))
			}
		}
		items = append(items, model.DeploymentItem{
			Name: application.Name,
			Kind: kind,
			Path: appWorkDir,
		})
	}
	itemSummary = model.DeploymentItemSummary{
		Namespace: releaseTemplate.Namespace,
		Items:     items,
	}
	return &itemSummary, nil

}

func GetRequiredTemplates(application *Application) ([]string, string) {
	kind := application.Kind
	requiredTemplates := make([]string, 0)
	requiredTemplates = append(requiredTemplates, SERVICEACCOUNT)
	if application.Kind == "" {
		requiredTemplates = append(requiredTemplates, DEPLOY)
		kind = "deployment"
	} else {
		requiredTemplates = append(requiredTemplates, templateMap[kind])
		kind = application.Kind
	}
	if len(application.ConfigMaps) > 0 {
		requiredTemplates = append(requiredTemplates, CONFIGMAP)
	}
	if application.Service.Enabled {
		requiredTemplates = append(requiredTemplates, SERVICE)
	}
	return requiredTemplates, kind
}
