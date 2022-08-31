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
			Parent:      "autoscaling",
			Name:        "inflate",
			Description: "Example App to Demonstrate Autoscaling",
		},

		Installer: &installer.ManifestInstaller{
			AppName: "autoscaling-inflate",
			ResourceTemplate: &template.TextTemplate{
				Template: manifestTemplate,
			},
		},
	}

	app.Options, app.Flags = NewOptions()

	return app
}

const manifestTemplate = `---
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
{{- if .OnDemand }}
      nodeSelector:
        karpenter.sh/capacity-type: on-demand
{{- end }}
{{- if .Spread }}
      topologySpreadConstraints:
      - maxSkew: 1
        topologyKey: topology.kubernetes.io/zone
        whenUnsatisfiable: ScheduleAnyway
        labelSelector:
          matchLabels:
            app: inflate
{{- end }}
`
