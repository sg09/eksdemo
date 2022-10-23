package irsa

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/iam"
)

type Getter struct {
	resource.EmptyInit
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	cluster := options.Common().Cluster

	allRoles, err := aws.IamListRoles()
	if err != nil {
		return err
	}

	oidcProviders, err := aws.IamListOpenIDConnectProviders()
	if err != nil {
		return err
	}

	providerARN := ""

	for _, p := range oidcProviders {
		provider, err := aws.IamGetOpenIDConnectProvider(aws.StringValue(p.Arn))
		if err != nil {
			return err
		}

		if strings.Contains(aws.StringValue(cluster.Identity.Oidc.Issuer), aws.StringValue(provider.Url)) {
			providerARN = aws.StringValue(p.Arn)
		}
	}

	if providerARN == "" {
		return fmt.Errorf("cluster %q has no IAM OIDC identity provider configured", aws.StringValue(cluster.Name))
	}

	roles := make([]*iam.Role, 0, len(allRoles))

	for _, r := range allRoles {

		doc, err := url.QueryUnescape(*r.AssumeRolePolicyDocument)
		if err != nil {
			return fmt.Errorf("cannot unescape role policy document: %s", err)
		}

		if strings.Contains(doc, providerARN) {
			role, err := aws.IamGetRole(*r.RoleName)
			if err != nil {
				return err
			}
			if name != "" && name != *r.RoleName {
				continue
			}
			roles = append(roles, role)
		}
	}

	return output.Print(os.Stdout, NewPrinter(roles))
}
