package templates

import (
	"errors"
	"fmt"
	"strings"
	"text/template"
)

var ChartTemplate = `apiVersion: v1
description: A Helm chart for Kubernetes {{ .ReleaseName }}
name: {{ .ReleaseName }}
version: 1.0
`

var ServiceAccountTemplate = `apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .ReleaseName | ToLower }}-{{ .Name  | ToLower }}
  namespace: {{ .Namespace }}
  labels:
    app: {{ .Name }}
    release: {{ .ReleaseName }}
    version: {{ .Tag }}
`

var ServiceTemplate = `apiVersion: v1
kind: Service
metadata:
  name: {{ .ReleaseName | ToLower }}-{{ .Name  | ToLower }}
  namespace: {{ .Namespace }}
  labels:
    app: {{ .Name }}
    release: {{ .ReleaseName }}
    version: {{ .Tag }}
spec:
  type: ClusterIP
  ports:
  - name: http
    port: 80
    targetPort: http
    protocol: TCP
  selector:
    app: {{ .Name }}
    release: {{ .ReleaseName }}
`

var DeploymentTemplate = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .ReleaseName | ToLower }}-{{ .Name  | ToLower }}
  namespace: {{ .Namespace }}
  labels:
    app: {{ .Name }}
    release: {{ .ReleaseName }}
    version: {{ .Tag }}
  annotations:{{ if .Annotations }}
    {{ range $key, $value := .Annotations }}{{ $key }}: {{ $value }}
    {{ end }}{{ end }}
spec:
  replicas: {{ .Replicas }}
  selector:
    matchLabels:
      app: {{ .Name }}
      release: {{ .ReleaseName }}
  template:
    metadata:
      labels:
        app: {{ .Name }}
        release: {{ .ReleaseName }}
    spec:
      serviceAccountName: {{ .ReleaseName | ToLower }}-{{ .Name  | ToLower }}
      containers:
       - name: {{ .Name }}
         image: {{ .Name}}:{{ .Tag}}
         imagePullPolicy: IfNotPresent
         {{ if .Entrypoint }}command: [{{ range $entry := .Entrypoint }}'{{$entry}}', {{ end }}]{{ end }}
         {{ if .Command }}args: [{{ range $cmd := .Command }}'{{$cmd}}', {{ end }}]{{ end }}

         {{ if .ServiceEnabled -}}ports:
         - name: http
           containerPort: {{ .ContainerPort }}
           protocol: TCP {{- end }}
         {{ if .LivenessProbe -}}livenessProbe:
           httpGet:
             path: {{ .LivenessProbe }}
             port: http
           initialDelaySeconds: 30
           timeoutSeconds: 100{{- end }}
         {{ if .ReadinessProbe -}}readinessProbe:
           httpGet:
             path: {{ .ReadinessProbe }}
             port: http
           initialDelaySeconds: 30
           timeoutSeconds: 100 {{- end }}
         resources:
           limits:
             cpu: "{{ index .Limits "cpu" }}"
             memory:  "{{ index .Limits "memory" }}"   
           requests:
             cpu:  "{{ index .Limits "cpu" }}"
             memory:  "{{ index .Limits "memory" }}"
         env:{{ range $key, $value := .EnvVars }}
          - name: "{{ $key | ToUpper }}"
            value: "{{ $value }}"{{end}}
      affinity:
      nodeSelector:
      tolerations:
`

var JobTemplate = `apiVersion: batch/v1
kind: Job
metadata:
  name: {{ .ReleaseName | ToLower }}-{{ .Name  | ToLower }}
  namespace: {{ .Namespace }}
  labels:
    app: {{ .Name }}
    release: {{ .ReleaseName }}
    version: {{ .Tag }}
  annotations:{{ if .Annotations }}
    {{ range $key, $value := .Annotations }}{{ $key }}: {{ $value }}
    {{ end }}{{ end }}
spec:
  {{ if .Replicas -}}completions: {{ .Replicas }}{{else}}1{{- end }}
  {{ if .Parallelism -}}parallelism: {{ .Parallelism }}{{- end }}
  {{ if .BackoffLimit -}}backoffLimit: {{ .BackoffLimit }}{{- end }}
  {{ if .ActiveDeadLine -}}activeDeadlineSeconds: {{ .ActiveDeadLine }}{{- end }}
  {{ if .TTLSecondsAfterFinished -}}ttlSecondsAfterFinished: {{ .TTLSecondsAfterFinished }}{{- end }}
  template:
    spec:
      serviceAccountName: {{ .ReleaseName | ToLower }}-{{ .Name  | ToLower }}
      containers:
       - name: {{ .Name }}
         image: {{ .Name}}:{{ .Tag}}
         imagePullPolicy: IfNotPresent
         {{ if .Entrypoint }}command: [{{ range $entry := .Entrypoint }}'{{$entry}}', {{ end }}]{{ end }}
         {{ if .Command }}args: [{{ range $cmd := .Command }}'{{$cmd}}', {{ end }}]{{ end }}

         resources:
           limits:
             cpu: "{{ index .Limits "cpu" }}"
             memory:  "{{ index .Limits "memory" }}"   
           requests:
             cpu:  "{{ index .Limits "cpu" }}"
             memory:  "{{ index .Limits "memory" }}"
         env:{{ range $key, $value := .EnvVars }}
          - name: "{{ $key | ToUpper }}"
            value: "{{ $value }}"{{end}}
      restartPolicy: {{ if .RestartPolicy -}}{{ .RestartPolicy }}{{ else }}Never{{end}} 
      affinity:
      nodeSelector:
      tolerations:
`

//LoadTemplates parse static template to helm chart
func LoadTemplates(tName string, app *Application) (*template.Template, error) {
	switch tName {
	case "ChartTemplate":
		return getTemplate("Chart.yaml", ChartTemplate)
	case "DeploymentTemplate":
		return getTemplate(fmt.Sprintf("%s-deployment.yaml", app.Name), DeploymentTemplate)
	case "ServiceTemplate":
		return getTemplate(fmt.Sprintf("%s-service.yaml", app.Name), ServiceTemplate)
	case "ServiceAccountTemplate":
		return getTemplate(fmt.Sprintf("%s-serviceaccount.yaml", app.Name), ServiceAccountTemplate)
	case "JobTemplate":
		return getTemplate(fmt.Sprintf("%s-job.yaml", app.Name), JobTemplate)
	}
	return nil, nil
}

func getTemplate(name string, templateType string) (*template.Template, error) {
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
	}

	tmpl, err := template.New(name).Funcs(funcMap).Parse(templateType)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error parsing %v ", err))
	}
	return tmpl, nil
}
