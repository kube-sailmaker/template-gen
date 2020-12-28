package model

import (
	"errors"
)

type AppSpec struct {
	Namespace   string `json:"namespace"`
	ReleaseName string `json:"release-name"`
	Environment string `json:"environment"`
	Apps        []App  `json:"apps"`
}

type App struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (as *AppSpec) Validate() error {
	if as.Namespace == "" {
		return errors.New("namespace is required")
	}
	if as.Environment == "" {
		return errors.New("environment is required")
	}
	if len(as.Apps) == 0 {
		return errors.New("app to deploy cannot be empty")
	}
	return nil
}
func (as *AppSpec) Normalise() {
	if as.ReleaseName == "" {
		as.ReleaseName = as.Namespace
	}
}

type DeploymentItem struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
	Path string `json:"path"`
}

type DeploymentItemSummary struct {
	Namespace string           `json:"namespace"`
	Items     []DeploymentItem `json:"items"`
}
