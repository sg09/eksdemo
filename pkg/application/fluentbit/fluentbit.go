package fluentbit

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/helm"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/template"
)

// Docs:    https://docs.fluentbit.io/manual/
// GitHub:  https://github.com/aws/aws-for-fluent-bit
// Helm:    https://github.com/fluent/helm-charts/tree/main/charts/fluent-bit
// Repo:    https://gallery.ecr.aws/aws-observability/aws-for-fluent-bit
// Version: Latest is 2.19.1 aka Fluent-bit v1.8.6 (as of 9/23/21)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "fluent-bit",
			Description: "Fluent Bit Logging",
			Aliases:     []string{"fluentbit", "fb"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "fluent-bit-irsa",
				},
				PolicyType: irsa.PolicyDocument,
				Policy:     []string{policyDocument},
			}),
		},

		Installer: &helm.HelmInstaller{
			ChartName:     "fluent-bit",
			ReleaseName:   "fluent-bit",
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
