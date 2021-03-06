### Template Generator
Given a manifest folder path, it generate chart to deploy into K8s.

#### Running Example App
Load a configuration file with details from config.yml
```
go run example.go
```

### Quickstart
```
import (
	"github.com/kube-sailmaker/template-gen/entry"
	"github.com/kube-sailmaker/template-gen/model"
	"os"
)

func main() {
	appList := make([]model.App, 0)
	
	appSpec := model.AppSpec{
		Namespace:   "apps",
		ReleaseName: "Release-2",
		Environment: "test",
		App:         appList,
	}
	path, _ := os.Getwd()
	appDir := path + "/sample-manifest/user/apps"
	resourceDir := path + "/sample-manifest/provider"
	outputDir := path + "/tmp"

	entry.TemplateGenerator(&appSpec, appDir, resourceDir, outputDir)
}
```

Based on the release manifest it generates the following for the each application listed in the release manifest

```
- service
- service-account
- deployment
```