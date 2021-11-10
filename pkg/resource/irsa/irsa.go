package irsa

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/eksctl"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/template"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "irsa",
			Description: "IAM Role for a Service Account",
			Args:        []string{"SERVICEACCOUNT"},
		},

		Getter: &Getter{},

		Manager: &eksctl.ResourceManager{
			Resource: "iamserviceaccount",
			Template: &template.TextTemplate{
				Template: eksctl.EksctlHeader + eksctlIamHeader + EksctlTemplate,
			},
			ApproveCreate: true,
			ApproveDelete: true,
		},
	}
	return addOptions(res)
}

func NewResourceWithOptions(options *IrsaOptions) *resource.Resource {
	res := NewResource()
	res.Options = options
	return res
}

const eksctlIamHeader = `
iam:
  withOIDC: true
  serviceAccounts:
`

const EksctlTemplate = `
  - metadata:
      name: {{ .ServiceAccount }}
      namespace: {{ .Namespace }}
    roleName: eksdemo.{{ .ClusterName }}.{{ .Namespace }}.{{ .ServiceAccount }}
    roleOnly: true
{{- if .PolicyType | .IsPolicyDocument }}
    attachPolicy:
{{- first .Policy | indent 6 }}
{{- end }}
{{- if .PolicyType | .IsPolicyARN }}
    attachPolicyARNs:
  {{- range .Policy }}
    - {{ . }}
  {{- end }}
{{- end }}
{{- if .PolicyType | .IsWellKnownPolicy }}
    wellKnownPolicies:
    {{- range .Policy }}
      {{ . }}: true
    {{- end }}
{{- end -}}
`
