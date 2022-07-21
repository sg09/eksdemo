package kubecost

import (
	"eksdemo/pkg/kubernetes"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/template"
)

type NginxSecretOptions struct {
	resource.CommonOptions
	AdminPassword string
}

func nginxSecret(password string) *resource.Resource {
	res := &resource.Resource{
		Options: &NginxSecretOptions{
			CommonOptions: resource.CommonOptions{
				Name: "kubecost-nginx-basic-auth-secret",
			},
			AdminPassword: password,
		},

		Manager: &kubernetes.ResourceManager{
			Template: &template.TextTemplate{
				Template: secretTemplate,
			},
		},
	}
	return res
}

const secretTemplate = `---
apiVersion: v1
kind: Secret
metadata:
  name: basic-auth
  namespace: {{ .Namespace }}
type: Opaque
data:
  auth: {{ htpasswd "admin" .AdminPassword | b64enc }}
`
