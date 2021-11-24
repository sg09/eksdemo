package keycloak

import (
	"context"
	"eksdemo/pkg/application"
	"eksdemo/pkg/aws"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/kubernetes"
	"eksdemo/pkg/resource/amg"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const AmgAliasSuffix = `keycloak-amg`
const samlMetadataPath = `auth/realms/eksdemo/protocol/saml/descriptor`

type KeycloakOptions struct {
	application.ApplicationOptions

	AdminPassword   string
	AmgWorkspaceUrl string
	TLSHost         string
	amgWorkspaceId  string
	*amg.AmgOptions
}

func NewOptions() (options *KeycloakOptions, flags cmd.Flags) {
	options = &KeycloakOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:      "keycloak",
			ServiceAccount: "keycloak",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "15.0.2",
				Previous: "15.0.2",
			},
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "admin-pass",
				Description: "Keycloak admin password (required)",
				Required:    true,
			},
			Option: &options.AdminPassword,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "tls-host",
				Description: "FQDN of host to secure with TLS (requires ExternalDNS for cert discovery) ",
			},
			Option: &options.TLSHost,
		},
	}
	return
}

func (o *KeycloakOptions) PreDependencies(application.Action) error {
	o.AmgOptions.WorkspaceName = fmt.Sprintf("%s-%s", o.ClusterName, AmgAliasSuffix)
	return nil
}

func (o *KeycloakOptions) PreInstall() error {
	amgGetter := amg.Getter{}

	workspace, err := amgGetter.GetAmgByName(o.AmgOptions.WorkspaceName)
	if err != nil {
		return fmt.Errorf("failed to lookup AMG URL to use in Helm chart: %w", err)
	}

	o.amgWorkspaceId = aws.StringValue(workspace.Id)
	o.AmgWorkspaceUrl = aws.StringValue(workspace.Endpoint)

	return nil
}

func (o *KeycloakOptions) PostInstall() error {
	fmt.Print("Waiting for Keycloak SAML metadata URL to become active...")

	k8sclient, err := kubernetes.Client(o.KubeContext())
	if err != nil {
		return err
	}

	svc, err := k8sclient.NetworkingV1().Ingresses(o.Namespace).Get(context.Background(), keycloakReleasName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if len(svc.Status.LoadBalancer.Ingress) == 0 {
		return fmt.Errorf("failed to get Ingress load balancer address")
	}

	var metadataUrl string
	if o.TLSHost == "" {
		metadataUrl = fmt.Sprintf("http://%s/%s", svc.Status.LoadBalancer.Ingress[0].Hostname, samlMetadataPath)
	} else {
		metadataUrl = fmt.Sprintf("https://%s/%s", o.TLSHost, samlMetadataPath)
	}

	logger := logrus.New()
	logger.Out = ioutil.Discard

	_, err = resty.New().
		SetLogger(logger).
		SetRetryCount(10).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second).
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				if err != nil {
					// retry on "no such host" error when using TLS as we wait for DNS
					return true
				}
				// retry as we wait for the ALB healthcheck for the target group
				return r.StatusCode() == http.StatusServiceUnavailable
			},
		).
		R().Get(metadataUrl)

	if err != nil {
		fmt.Println()
		return fmt.Errorf("%w\n\nTo finish configuration, update AMG with the SAML metadata URL: %s", err, metadataUrl)
	}
	fmt.Println("done")
	fmt.Printf("Updating AMG with Keyclock SAML Metadata URL to complete SAML configuration\n")

	err = aws.AmgUpdateWorkspaceAuthentication(o.amgWorkspaceId, metadataUrl)
	if err != nil {
		return err
	}
	fmt.Printf("Amazon Managed Grafana available at: https://%s\n", o.AmgWorkspaceUrl)

	return nil
}
