package main

import (
	"fmt"
	"github.com/kube-sailmaker/template-gen/entry"
	"github.com/kube-sailmaker/template-gen/model"
	"os"
)

func main() {
	busybox := model.App{
		Name:    "busybox",
		Version: "latest",
	}
	eodJob := model.App{
		Name:    "eod-job",
		Version: "latest",
	}
	nginx := model.App{
		Name:    "nginx",
		Version: "latest",
	}

	appList := make([]model.App, 0)
	appList = append(appList, busybox, nginx, eodJob)

	appSpec := model.AppSpec{
		Namespace:   "apps",
		ReleaseName: "Release-2",
		Environment: "test",
		Apps:        appList,
	}
	path, _ := os.Getwd()
	appDir := path + "/sample-manifest/user/apps"
	resourceDir := path + "/sample-manifest/provider"
	outputDir := path + "/tmp"

	data, err := entry.TemplateGenerator(&appSpec, appDir, resourceDir, outputDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Generation Summary:")
	for _, d := range data.Items {
		fmt.Printf("Name: %s, Kind: %s, Path: %s\n", d.Name, d.Kind, d.Path)
	}
}
