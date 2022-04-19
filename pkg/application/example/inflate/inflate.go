package inflate

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/template"
)

// Manifest: https://karpenter.sh/v0.8.2/getting-started/getting-started-with-eksctl/#automatic-node-provisioning
// Repo:     https://public.ecr.aws/eks-distro/kubernetes/pause

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "inflate",
			Description: "Karpenter Example App to Demonstrate Autoscaling",
		},

		Installer: &installer.ManifestInstaller{
			AppName: "example-inflate",
			ResourceTemplate: &template.TextTemplate{
				Template: manifestTemplate,
			},
		},
	}

	app.Options, app.Flags = NewOptions()

	return app
}

const manifestTemplate = `
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: inflate
  namespace: {{ .Namespace }}
spec:
  replicas: {{ .Replicas }}
  selector:
    matchLabels:
      app: inflate
  template:
    metadata:
      labels:
        app: inflate
    spec:
      terminationGracePeriodSeconds: 0
      containers:
        - name: inflate
          image: public.ecr.aws/eks-distro/kubernetes/pause:3.2
          resources:
            requests:
              cpu: 1
...`
