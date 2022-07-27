package karpenter

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/cloudformation"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/template"
)

func karpenterNodeRole() *resource.Resource {
	res := &resource.Resource{
		Options: &resource.CommonOptions{
			Name: "karpenter-node-role",
		},

		Manager: &cloudformation.ResourceManager{
			Capabilities: []aws.Capability{aws.CapabilityCapabilityNamedIam},
			Template: &template.TextTemplate{
				Template: cloudFormationTemplate,
			},
		},
	}
	return res
}

const cloudFormationTemplate = `
AWSTemplateFormatVersion: "2010-09-09"
Description: Resources used by https://github.com/aws/karpenter
Parameters:
  ClusterName:
    Type: String
    Description: "EKS cluster name"
    Default: "{{ .ClusterName }}"
Resources:
  KarpenterNodeInstanceProfile:
    Type: "AWS::IAM::InstanceProfile"
    Properties:
      InstanceProfileName: !Sub "KarpenterNodeInstanceProfile-${ClusterName}"
      Path: "/"
      Roles:
        - Ref: "KarpenterNodeRole"
  KarpenterNodeRole:
    Type: "AWS::IAM::Role"
    Properties:
      RoleName: !Sub "KarpenterNodeRole-${ClusterName}"
      Path: /
      AssumeRolePolicyDocument:
        Version: "2012-10-17"
        Statement:
          - Effect: Allow
            Principal:
              Service:
                !Sub "ec2.${AWS::URLSuffix}"
            Action:
              - "sts:AssumeRole"
      ManagedPolicyArns:
        - !Sub "arn:${AWS::Partition}:iam::aws:policy/AmazonEKS_CNI_Policy"
        - !Sub "arn:${AWS::Partition}:iam::aws:policy/AmazonEKSWorkerNodePolicy"
        - !Sub "arn:${AWS::Partition}:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
        - !Sub "arn:${AWS::Partition}:iam::aws:policy/AmazonSSMManagedInstanceCore"
`
