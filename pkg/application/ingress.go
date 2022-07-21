package application

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/template"
	"fmt"
	"strings"
)

type IngressOptions struct {
	IngressClass string
	IngressHost  string
	NLB          bool
	TargetType   string

	ingressTemplate template.Template
	serviceTemplate template.Template
}

func NewIngressOptions() IngressOptions {
	return IngressOptions{
		IngressClass: "alb",
		TargetType:   "ip",
		ingressTemplate: &template.TextTemplate{
			Template: ingressAnnotationsTemplate,
		},
		serviceTemplate: &template.TextTemplate{
			Template: strings.TrimPrefix(serviceAnnotationsTemplate, "\n"),
		},
	}
}

func (o *IngressOptions) IngressAnnotations() (string, error) {
	return o.ingressTemplate.Render(o)
}

func (o *IngressOptions) ServiceAnnotations() (string, error) {
	return o.serviceTemplate.Render(o)
}

func (o *IngressOptions) ServiceType() string {
	if o.IngressHost != "" {
		return `ClusterIP`
	} else {
		return `LoadBalancer`
	}
}

func (o *IngressOptions) NewIngressFlags() cmd.Flags {
	return cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "ingress-class",
				Description: "name of IngressClass",
			},
			Option: &o.IngressClass,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "ingress-host",
				Description: "hostname for Ingress with TLS (default is Service of type LoadBalancer)",
				Shorthand:   "I",
			},
			Option: &o.IngressHost,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "nlb",
				Description: "use NLB instead of CLB (when not using Ingress)",
				Validate: func() error {
					if o.NLB && o.IngressHost != "" {
						return fmt.Errorf("%q flag cannot be used with %q flag", "nlb", "ingress-host")
					}

					return nil
				},
			},
			Option: &o.NLB,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "target-type",
				Description: "target type when deploying NLB or ALB Ingress",
				Validate: func() error {
					if o.TargetType == "instance" && !o.NLB && o.IngressHost == "" {
						return fmt.Errorf("%q flag can only be used with %q or %q flags", "target-type", "nlb", "ingress-host")
					}

					return nil
				},
			},
			Choices: []string{"instance", "ip"},
			Option:  &o.TargetType,
		},
	}
}

const ingressAnnotationsTemplate = `
{{- if eq .IngressClass "alb" -}}
alb.ingress.kubernetes.io/listen-ports: '[{"HTTPS":443}]'
alb.ingress.kubernetes.io/scheme: internet-facing
alb.ingress.kubernetes.io/target-type: {{ .TargetType }}
{{- else -}}
cert-manager.io/cluster-issuer: letsencrypt-prod
{{- end -}}
`

const serviceAnnotationsTemplate = `
service.beta.kubernetes.io/aws-load-balancer-cross-zone-load-balancing-enabled: "true"
{{- if .NLB }}
service.beta.kubernetes.io/aws-load-balancer-nlb-target-type: {{ .TargetType }}
service.beta.kubernetes.io/aws-load-balancer-scheme: internet-facing
service.beta.kubernetes.io/aws-load-balancer-type: external
{{- end -}}
`
