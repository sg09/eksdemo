package servicelb

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/kubernetes"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/template"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "servicelb",
			Description: "Kubernetes Service of type LoadBalancer",
			Aliases:     []string{"ng", "mng"},
			Args:        []string{"NAME"},
		},

		Manager: &kubernetes.ResourceManager{
			Template: &template.TextTemplate{
				Template: yamlTemplate,
			},
		},

		Options: &resource.CommonOptions{
			Name: "service-loadbalancer",
		},
	}

	return res
}

const yamlTemplate = `
apiVersion: v1
kind: Service
metadata:
  annotations:
    service.beta.kubernetes.io/aws-load-balancer-type: "external"
    service.beta.kubernetes.io/aws-load-balancer-nlb-target-type: "ip"
  labels:
    app: foo
  name: {{ .Name }}
spec:
  ports:
  - name: "80"
    port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: foo
  type: LoadBalancer
`
