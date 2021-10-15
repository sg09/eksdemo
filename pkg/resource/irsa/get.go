package irsa

import (
	"eksdemo/pkg/aws"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/iam"
)

func Get(cluster *eks.Cluster) ([]*iam.Role, error) {
	allRoles, err := aws.IamListRoles()
	if err != nil {
		return nil, err
	}

	oidcProviders, err := aws.IamListOpenIDConnectProviders()
	if err != nil {
		return nil, err
	}

	providerARN := ""

	for _, p := range oidcProviders {
		provider, err := aws.IamGetOpenIDConnectProviders(aws.StringValue(p.Arn))
		if err != nil {
			return nil, err
		}

		if strings.Contains(aws.StringValue(cluster.Identity.Oidc.Issuer), aws.StringValue(provider.Url)) {
			providerARN = aws.StringValue(p.Arn)
		}
	}

	if providerARN == "" {
		return nil, fmt.Errorf("cluster %q has no IAM OIDC identity provider configured", aws.StringValue(cluster.Name))
	}

	roles := make([]*iam.Role, 0, len(allRoles))

	for _, r := range allRoles {

		doc, err := url.QueryUnescape(*r.AssumeRolePolicyDocument)
		if err != nil {
			return nil, fmt.Errorf("cannot unescape role policy document: %s", err)
		}

		if strings.Contains(doc, providerARN) {
			role, err := aws.IamGetRole(*r.RoleName)
			if err != nil {
				return nil, err
			}
			roles = append(roles, role)
		}
	}

	return roles, nil
}
