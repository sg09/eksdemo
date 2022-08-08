package crossplane

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/aws"
	"eksdemo/pkg/cloudformation"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/template"
	"fmt"
)

func crossplaneIrsa(options application.Options) *resource.Resource {
	o := options.Common()

	return &resource.Resource{
		Options: &irsa.IrsaOptions{
			CommonOptions: resource.CommonOptions{
				ClusterName:    o.ClusterName,
				Name:           fmt.Sprintf("crossplane-%s-irsa", o.Namespace),
				Namespace:      o.Namespace,
				ServiceAccount: o.ServiceAccount,
			},
		},

		Manager: &cloudformation.ResourceManager{
			Capabilities: []aws.Capability{aws.CapabilityCapabilityNamedIam},
			Template: &template.TextTemplate{
				Template: cloudFormationTemplate,
			},
		},
	}
}

const cloudFormationTemplate = `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  CrossplaneIRSA:
    Type: "AWS::IAM::Role"
    Properties:
      RoleName: {{ .RoleName }}
      Path: /
      AssumeRolePolicyDocument:
        Version: "2012-10-17"
        Statement:
        - Effect: Allow
          Principal:
            Federated:
              !Sub "arn:${AWS::Partition}:iam::{{ .Account }}:oidc-provider/{{ .ClusterOIDCProvider }}"
          Action:
          - sts:AssumeRoleWithWebIdentity
          Condition:
            StringLike:
              "{{ .ClusterOIDCProvider }}:sub": "system:serviceaccount:{{ .Namespace }}:provider-aws-*"
      ManagedPolicyArns:
      - !Sub "arn:${AWS::Partition}:iam::aws:policy/AdministratorAccess"
`
