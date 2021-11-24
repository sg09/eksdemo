package keycloak

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/helm"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/amg"
	"eksdemo/pkg/template"
)

// Docs:    https://www.keycloak.org/documentation
// GitHub:  https://github.com/keycloak/keycloak
// GitHub:  https://github.com/bitnami/bitnami-docker-keycloak
// Helm:    https://github.com/bitnami/charts/tree/master/bitnami/keycloak
// Repo:    https://hub.docker.com/r/bitnami/keycloak
// Version: Latest is 15.0.2 (as of 11/24/21)

func NewApp() *application.Application {
	options, flags := NewOptions()

	options.AmgOptions = &amg.AmgOptions{
		CommonOptions: resource.CommonOptions{
			Name: "amazon-managed-grafana",
		},
		Auth: []string{"SAML"},
	}

	app := &application.Application{
		Command: cmd.Command{
			Name:        "keycloak-amg",
			Description: "Keycloak SAML iDP for Amazon Managed Grafana",
			Aliases:     []string{"keycloak"},
		},

		Dependencies: []*resource.Resource{
			amg.NewResourceWithOptions(options.AmgOptions),
		},

		Installer: &helm.HelmInstaller{
			ChartName:     "keycloak",
			ReleaseName:   keycloakReleasName,
			RepositoryURL: "https://charts.bitnami.com/bitnami",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
			PVCLabels: map[string]string{
				"app.kubernetes.io/instance": keycloakReleasName,
			},
		},
	}

	app.Options = options
	app.Flags = flags

	return app
}

const keycloakReleasName = `keycloak`

const valuesTemplate = `
auth:
  adminUser: admin
  adminPassword: {{ .AdminPassword }}
image:
  tag: {{ .Version }}
proxyAddressForwarding: true
keycloakConfigCli:
  enabled: true
  command:
  - java
  - -jar
  - /opt/bitnami/keycloak-config-cli/keycloak-config-cli-{{ .Version }}.jar
  configuration:
    eksdemo.json: |
      {
        "realm": "eksdemo",
        "enabled": true,
        "roles": {
          "realm": [
            {
              "name": "admin"
            }
          ]
        },
{{- if not .TLSHost }}
        "sslRequired": "none",
{{- end }}
        "users": [
          {
            "username": "admin",
            "email": "admin@eksdemo",
            "enabled": true,
            "firstName": "Admin",
            "realmRoles": [
              "admin"
            ],
            "credentials": [
              {
                "type": "password",
                "value": "{{ .AdminPassword }}"
              }
            ]
          }
        ],
        "clients": [
          {
            "clientId": "https://{{ .AmgWorkspaceUrl }}/saml/metadata",
            "name": "amazon-managed-grafana",
            "enabled": true,
            "protocol": "saml",
            "adminUrl": "https://{{ .AmgWorkspaceUrl }}/login/saml",
            "redirectUris": [
              "https://{{ .AmgWorkspaceUrl }}/saml/acs"
            ],
            "attributes": {
              "saml.authnstatement": "true",
              "saml.server.signature": "true",
              "saml_name_id_format": "email",
              "saml_force_name_id_format": "true",
              "saml.assertion.signature": "true",
              "saml.client.signature": "false"
            },
            "defaultClientScopes": [],
            "protocolMappers": [
              {
                "name": "name",
                "protocol": "saml",
                "protocolMapper": "saml-user-property-mapper",
                "consentRequired": false,
                "config": {
                  "attribute.nameformat": "Unspecified",
                  "user.attribute": "firstName",
                  "attribute.name": "displayName"
                }
              },
              {
                "name": "email",
                "protocol": "saml",
                "protocolMapper": "saml-user-property-mapper",
                "consentRequired": false,
                "config": {
                  "attribute.nameformat": "Unspecified",
                  "user.attribute": "email",
                  "attribute.name": "mail"
                }
              },
              {
                "name": "role list",
                "protocol": "saml",
                "protocolMapper": "saml-role-list-mapper",
                "config": {
                  "single": "true",
                  "attribute.nameformat": "Unspecified",
                  "attribute.name": "role"
                }
              }
            ]
          }
        ]
      }
service:
  type: ClusterIP
ingress:
  enabled: true
  pathType: Prefix
  annotations:
    kubernetes.io/ingress.class: alb
    alb.ingress.kubernetes.io/scheme: internet-facing
    alb.ingress.kubernetes.io/target-type: 'ip'
{{- if .TLSHost }}
    alb.ingress.kubernetes.io/listen-ports: '[{"HTTPS":443}]'
  hostname: {{ .TLSHost }}
{{- else }}
  hostname: null
{{- end }}
`
