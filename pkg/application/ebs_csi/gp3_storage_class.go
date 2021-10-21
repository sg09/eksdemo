package ebs_csi

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/kubernetes"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/template"
	"fmt"
)

type EbsCsiOptions struct {
	application.ApplicationOptions
}

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

func (o *EbsCsiOptions) PostInstall() error {
	res := gp3StorageClass()
	o.AssignCommonResourceOptions(res)
	fmt.Printf("Creating post-install resource: %s\n", res.Common().Name)

	return res.Create()
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
