package iam_oidc

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"fmt"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/iam"
)

type Getter struct {
	resource.EmptyInit
}

func (g *Getter) Get(providerUrl string, output printer.Output, options resource.Options) error {
	var err error
	var oidcProvider *iam.GetOpenIDConnectProviderOutput
	var oidcProviders []*iam.GetOpenIDConnectProviderOutput

	if providerUrl != "" {
		oidcProvider, err = g.GetOidcProviderByUrl(providerUrl)
		oidcProviders = []*iam.GetOpenIDConnectProviderOutput{oidcProvider}
	} else if options.Common().Cluster != nil {
		oidcProvider, err = g.GetOidcProviderByCluster(options.Common().Cluster)
		oidcProviders = []*iam.GetOpenIDConnectProviderOutput{oidcProvider}
	} else {
		oidcProviders, err = g.GetAllOidcProviders()
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(oidcProviders))
}

func (g *Getter) GetAllOidcProviders() ([]*iam.GetOpenIDConnectProviderOutput, error) {
	providerList, err := aws.IamListOpenIDConnectProviders()
	if err != nil {
		return nil, err
	}
	oidcProviders := make([]*iam.GetOpenIDConnectProviderOutput, 0, len(providerList))

	for _, p := range providerList {
		provider, err := aws.IamGetOpenIDConnectProvider(aws.StringValue(p.Arn))
		if err != nil {
			return nil, err
		}
		oidcProviders = append(oidcProviders, provider)
	}
	return oidcProviders, nil
}

func (g *Getter) GetOidcProviderByCluster(cluster *eks.Cluster) (*iam.GetOpenIDConnectProviderOutput, error) {
	u, err := url.Parse(aws.StringValue(cluster.Identity.Oidc.Issuer))
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL when validating options: %w", err)
	}
	return g.GetOidcProviderByUrl(u.Hostname() + u.Path)
}

func (g *Getter) GetOidcProviderByUrl(url string) (*iam.GetOpenIDConnectProviderOutput, error) {
	arn := fmt.Sprintf("arn:%s:iam::%s:oidc-provider/%s", aws.Partition(), aws.AccountId(), url)

	provider, err := aws.IamGetOpenIDConnectProvider(arn)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				return nil, resource.NotFoundError(fmt.Sprintf("oidc-provider %q not found", url))
			default:
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return provider, nil
}
