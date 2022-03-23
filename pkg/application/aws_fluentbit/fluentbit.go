package aws_fluentbit

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/template"
)

// Docs:    https://docs.fluentbit.io/manual/
// GitHub:  https://github.com/aws/aws-for-fluent-bit
// Helm:    https://github.com/fluent/helm-charts/tree/main/charts/fluent-bit
// Repo:    https://gallery.ecr.aws/aws-observability/aws-for-fluent-bit
// Version: Latest is 2.23.1 aka Fluent-bit v1.8.13 (as of 03/23/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "aws-fluent-bit",
			Description: "AWS Fluent Bit",
			Aliases:     []string{"aws-fluentbit", "aws-fb", "awsfb"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "aws-fluent-bit-irsa",
				},
				PolicyType: irsa.PolicyDocument,
				PolicyDocTemplate: &template.TextTemplate{
					Template: policyDocument,
				},
			}),
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "fluent-bit",
			ReleaseName:   "aws-fluent-bit",
			RepositoryURL: "https://fluent.github.io/helm-charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	return addOptions(app)
}

// TODO: Limit, similar to Resource: "arn:aws:logs:*:*:log-group:/aws/eks/<cluster-name>/*"
const policyDocument = `
Version: "2012-10-17"
Statement:
- Effect: Allow
  Action:
  - "logs:CreateLogGroup"
  - "logs:CreateLogStream"
  - "logs:DescribeLogStreams"
  - "logs:PutLogEvents"
  - "logs:PutRetentionPolicy"
  Resource: '*'
`
