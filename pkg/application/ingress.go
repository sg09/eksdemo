package application

import (
	"eksdemo/pkg/application/ingress"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/template"
	"fmt"
	"strings"
)

type IngressOptions struct {
	IngressClass       string
	IngressHost        string
	IngressOnly        bool
	NginxBasicAuthPass string
	NLB                bool
	TargetType         string

	ingressTemplate template.Template
	serviceTemplate template.Template
}

func NewIngressOptions(ingressOnly bool) IngressOptions {
	return IngressOptions{
		IngressClass: "alb",
		IngressOnly:  ingressOnly,
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

func (o *IngressOptions) NewIngressFlags() cmd.Flags {
	var ingressHostDesc, targetTypeDesc string

	if o.IngressOnly {
		ingressHostDesc += "hostname for Ingress with TLS"
		targetTypeDesc += "target type when deploying ALB Ingress"
	} else {
		ingressHostDesc += "hostname for Ingress with TLS (default is Service of type LoadBalancer)"
		targetTypeDesc += "target type when deploying NLB or ALB Ingress"
	}

	ingressFlags := cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "ingress-class",
				Description: "name of IngressClass",
				Validate: func() error {
					if o.IngressClass != "alb" && o.IngressHost == "" {
						return fmt.Errorf("%q flag can only be used with %q flag", "ingress-class", "ingress-host")
					}
					return nil
				},
			},
			Option: &o.IngressClass,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "ingress-host",
				Description: ingressHostDesc,
				Shorthand:   "I",
				Required:    o.IngressOnly,
			},
			Option: &o.IngressHost,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "nginx-pass",
				Description: "basic auth password for admin user (only valid with --ingress-class=nginx)",
				Shorthand:   "X",
				Validate: func() error {
					if o.NginxBasicAuthPass != "" && o.IngressClass != "nginx" {
						return fmt.Errorf("%q flag can only be used with %q)", "nginx-pass", "--ingress-class=nginx")
					}
					return nil
				},
			},
			Option: &o.NginxBasicAuthPass,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "target-type",
				Description: targetTypeDesc,
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

	if !o.IngressOnly {
		serviceFlag := &cmd.BoolFlag{
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
		}

		ingressFlags = append(ingressFlags, serviceFlag)
	}

	return ingressFlags
}

func (o *IngressOptions) PostInstallResources(name string) []*resource.Resource {
	if len(o.NginxBasicAuthPass) > 0 {
		return []*resource.Resource{ingress.BasicAuthSecret(name, o.NginxBasicAuthPass)}
	}
	return nil
}

func (o *IngressOptions) ServiceAnnotations() (string, error) {
	if o.IngressHost != "" {
		return "{}", nil
	}

	return o.serviceTemplate.Render(o)
}

func (o *IngressOptions) ServiceType() string {
	if o.IngressHost != "" {
		return `ClusterIP`
	} else {
		return `LoadBalancer`
	}
}

const ingressAnnotationsTemplate = `
{{- if eq .IngressClass "alb" -}}
alb.ingress.kubernetes.io/listen-ports: '[{"HTTP": 80}, {"HTTPS":443}]'
alb.ingress.kubernetes.io/scheme: internet-facing
alb.ingress.kubernetes.io/ssl-redirect: '443'
alb.ingress.kubernetes.io/target-type: {{ .TargetType }}
{{- else -}}
cert-manager.io/cluster-issuer: letsencrypt-prod
{{- end -}}
{{- if .NginxBasicAuthPass }}
nginx.ingress.kubernetes.io/auth-type: basic
nginx.ingress.kubernetes.io/auth-secret: basic-auth
nginx.ingress.kubernetes.io/auth-realm: "Authentication Required"
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
