package templates

import (
	"errors"
	"fmt"
	"github.com/kube-sailmaker/template-gen/model"
	"log"
	"os"
	"path/filepath"
	"strings"
)

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
	kind := ""
	requiredTemplates := make([]string, 0)
	requiredTemplates = append(requiredTemplates, "ServiceAccountTemplate")
	if len(application.Kind) == 0 || application.Kind == "Deployment" {
		requiredTemplates = append(requiredTemplates, "DeploymentTemplate")
		kind = "deployment"
	} else if strings.EqualFold(application.Kind, "Job") {
		requiredTemplates = append(requiredTemplates, "JobTemplate")
		kind = "job"
	}

	if application.ServiceEnabled {
		requiredTemplates = append(requiredTemplates, "ServiceTemplate")
	}
	return requiredTemplates, kind
}
