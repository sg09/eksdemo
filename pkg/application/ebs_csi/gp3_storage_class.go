package ebs_csi

import (
	"eksdemo/pkg/kubernetes"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/template"
)

func gp3StorageClass() *resource.Resource {
	res := &resource.Resource{
		Options: &resource.CommonOptions{
			Name: "ebs-csi-gp3-storage-class",
		},

		Manager: &kubernetes.ResourceManager{
			Template: &template.TextTemplate{
				Template: yamlTemplate,
			},
		},
	}
	return res
}

const yamlTemplate = `
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: gp3
parameters:
  csi.storage.k8s.io/fstype: ext4
  type: gp3
provisioner: ebs.csi.aws.com
volumeBindingMode: WaitForFirstConsumer
`
