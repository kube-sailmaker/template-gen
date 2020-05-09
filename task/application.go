package task

import (
	"errors"
	"fmt"
	"github.com/kube-sailmaker/template-gen/functions"
	"github.com/kube-sailmaker/template-gen/model"
	templates "github.com/kube-sailmaker/template-gen/template"
	"log"
	"strings"
)

const (
	cpu      = "cpu"
	memory   = "memory"
	replicas = "replicas"
	sep      = "/"
)

//CPU value mapping
var CPU = map[string]string{
	"c05":     "0.5",
	"c1":      "1",
	"c2":      "2",
	"c3":      "3",
	"default": "0.5",
}

//Memory value mapping
var MEMORY = map[string]string{
	"m05":     "0.5Gi",
	"m1":      "1Gi",
	"m2":      "2Gi",
	"m3":      "3Gi",
	"default": "256Mi",
}

func ProcessApplication(app *model.App, releaseName string, namespace string, env string, appDir string, resourceDir string) (*templates.Application, error) {
	if app == nil {
		return nil, errors.New("app specification cannot be nil")
	}
	appFile := fmt.Sprintf("%s/%s.yaml", appDir, app.Name)
	application := &model.Application{}
	err := functions.UnmarshalFile(appFile, application)
	if err != nil {
		return nil, err
	}

	if application.Service == nil {
		application.Service = &model.ServiceSpec{
			Port:    8080,
			Enabled: true,
		}
	}
	if application.Service.Enabled {
		if application.Service.Port == 0 {
			application.Service.Port = 8080
		}
		if application.Service.Port < 0 || application.Service.Port > 65535 {
			return nil, fmt.Errorf("invalid port %d specified for %s", application.Service.Port, application.Name)
		}
	}

	appValues := templates.Application{
		Name:           app.Name,
		Tag:            app.Version,
		Kind:           application.Kind,
		Namespace:      namespace,
		ReleaseName:    releaseName,
		Annotations:    application.Annotations,
		LivenessProbe:  application.LivenessProbe,
		ReadinessProbe: application.ReadinessProbe,
		ServiceEnabled: application.Service.Enabled,
		ContainerPort:  application.Service.Port,
	}

	err = GenerateResourceLimit(application, env, &appValues)
	if err != nil {
		return nil, err
	}
	err = GenerateEnvVars(application, resourceDir, &appValues)
	if err != nil {
		return nil, err
	}
	err = GenerateMixins(application, resourceDir, &appValues)
	if err != nil {
		return nil, err
	}

	return &appValues, nil
}

//Function to set the mixins
func GenerateMixins(application *model.Application, resourceDir string, appValues *templates.Application) error {
	appValues.Command = make([]string, 0)
	appValues.Entrypoint = make([]string, 0)
	for _, mxin := range application.Mixins {
		mixinType := strings.Split(mxin, sep)
		if len(mixinType) < 2 {
			eMsg := fmt.Sprintf("application mixin %s has missing value, eg: java/java-default", mixinType)
			return errors.New(eMsg)
		}
		name := mixinType[0]
		mType := mixinType[1]
		mixinList := model.MixinList{}
		err := GetMixin(name, &mixinList, resourceDir)
		if err != nil {
			return err
		}
		match := false
		for _, m := range mixinList.Mixin {
			if mType == m.Name {
				for k, v := range m.Env {
					appValues.EnvVars[k] = v
				}
				appValues.Command = m.Cmd
				appValues.Entrypoint = m.Entrypoint
				match = true
				break
			}
		}
		if match == false {
			log.Print(fmt.Sprintf("[WARN] could not find matching mixin %s of app %s", mType, application.Name))
		}
	}
	return nil
}

//Function to set resource limit, request and replicas
func GenerateResourceLimit(application *model.Application, environment string, appValues *templates.Application) error {
	appValues.Limits = make(map[string]string, 0)
	//apply default values if missing
	if len(application.Template) == 0 {
		appValues.Limits[cpu] = CPU["default"]
		appValues.Limits[memory] = MEMORY["default"]
		appValues.Replicas = "1"
		log.Println("[WARN] missing resource, applying default values for application", application.Name)
		return nil
	}
	//Process app template
	found := false
	for _, tmpl := range application.Template {
		if tmpl.Name == environment {
			found = true
			configs := tmpl.Config
			cpuLimit := CPU["default"]
			if val, ok := configs[cpu]; ok && val != "" {
				cpuLimit = CPU[val]
			}
			memLimit := MEMORY["default"]
			if val, ok := configs[memory]; ok && val != "" {
				memLimit = MEMORY[val]
			}

			replicaLimit := "1"
			if val, ok := configs[replicas]; ok && val != "" {
				replicaLimit = val
			}
			appValues.Limits[cpu] = cpuLimit
			appValues.Limits[memory] = memLimit
			appValues.Replicas = replicaLimit
			break
		}
	}
	if !found {
		return fmt.Errorf("unknown environment %s", environment)
	}
	return nil
}

//Function to set environment variable from resources
func GenerateEnvVars(application *model.Application, resourceDir string, appValues *templates.Application) error {
	appEnvVars := make(map[string]string, 0)

	for _, appRes := range application.Resources {
		//elasticsearch-user:sit
		resDetails := strings.Split(appRes, sep)
		if len(resDetails) < 2 {
			eMsg := fmt.Sprintf("application resource %s has missing template type, eg: cassandra/test1", resDetails)
			return errors.New(eMsg)
		}
		name := resDetails[0]
		envType := resDetails[1]
		resource := &model.Resource{}
		err := GetResource(name, &resource, resourceDir)
		if err != nil {
			return err
		}
		matchEnvType := false
		for _, resTemplate := range resource.Spec.ResourceTemplate {
			//Only using the context
			if resTemplate.Name == envType {
				addToEnvVars(name, appEnvVars, resTemplate.Element)

				if len(resTemplate.Infra) > 0 {
					infra := strings.Split(resTemplate.Infra, sep)
					if len(infra) < 2 {
						eMsg := fmt.Sprintf("resource infrastructure %s has missing template type, eg: cassandra-a/test", infra)
						return errors.New(eMsg)
					}
					infraName := infra[0]
					infraEnv := infra[1]
					infrastructure := &model.Infrastructure{}
					GetInfrastructure(infraName, &infrastructure, resourceDir)
					matchInfra := false
					for _, infraTemplate := range infrastructure.Spec.Template {
						if infraEnv == infraTemplate.Name {
							addToEnvVars(name, appEnvVars, infraTemplate.Attributes)
							matchInfra = true
							break
						}
					}
					if matchInfra == false {
						log.Print(fmt.Sprintf("[WARN] could not find matching infra for env type %s of %s", infraEnv, resTemplate.Infra))
					}
				}
				matchEnvType = true
				break
			}
		}
		if matchEnvType == false {
			log.Print(fmt.Sprintf("[WARN] could not find matching env type %s of app %s", envType, application.Name))
		}
	}
	appValues.EnvVars = appEnvVars
	return nil
}

func addToEnvVars(name string, appEnvVars map[string]string, items map[string]string) {
	infraName := strings.ReplaceAll(name, "-", "_")
	for k, v := range items {
		key := strings.ToUpper(fmt.Sprintf("%s_%s", infraName, k))
		appEnvVars[key] = v
	}
}
