package nodegroup

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/eksctl"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/template"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "nodegroup",
			Description: "EKS Managed Nodegroup",
			Aliases:     []string{"nodegroups", "mng", "ng"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},

		Manager: &Manager{
			Eksctl: &eksctl.ResourceManager{
				Resource: "nodegroup",
				Template: &template.TextTemplate{
					Template: eksctl.EksctlHeader + EksctlTemplate,
				},
				ApproveDelete: true,
			},
		},
	}

	res.Options, res.CreateFlags, res.UpdateFlags = NewOptions()

	return res
}

const EksctlTemplate = `
managedNodeGroups:
- name: {{ .NodegroupName }}
{{- if .AMI }}
  ami: {{ .AMI }}
{{- end }}
  amiFamily: {{ .OperatingSystem }}
  iam:
    attachPolicyARNs:
    - arn:{{ .Partition }}:iam::aws:policy/AmazonEKSWorkerNodePolicy
    - arn:{{ .Partition }}:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly
    - arn:{{ .Partition }}:iam::aws:policy/AmazonSSMManagedInstanceCore
{{- if .Spot }}
  instanceSelector:
    vCPUs: {{ .SpotvCPUs }}
    memory: {{ .SpotMemory | toString | printf "%q" }}
{{- else }}
  instanceType: {{ .InstanceType }}
{{- end }}
  minSize: {{ .MinSize }}
  desiredCapacity: {{ .DesiredCapacity }}
  maxSize: {{ .MaxSize }}
{{- if .Containerd }}
  overrideBootstrapCommand: |
    #!/bin/bash
    /etc/eks/bootstrap.sh {{ .ClusterName }} --container-runtime containerd
{{- end }}
  privateNetworking: true
  spot: {{ .Spot }}
`
