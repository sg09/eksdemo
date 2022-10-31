package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type IAMClient struct {
	*iam.Client
}

func NewIAMClient() *IAMClient {
	return &IAMClient{iam.NewFromConfig(GetConfig())}
}

func (c *IAMClient) CreateRole(assumeRolePolicy, name, path string) (*types.Role, error) {
	if path == "" {
		path = "/"
	}

	result, err := c.Client.CreateRole(context.Background(), &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(assumeRolePolicy),
		RoleName:                 aws.String(name),
		Path:                     aws.String(path),
	})

	if err != nil {
		return nil, err
	}

	return result.Role, nil
}

func (c *IAMClient) DeleteRole(name string) error {
	_, err := c.Client.DeleteRole(context.Background(), &iam.DeleteRoleInput{
		RoleName: aws.String(name),
	})

	return err
}

// Deletes the specified inline policy that is embedded in the specified IAM role.
func (c *IAMClient) DeleteRolePolicy(roleName, policyName string) error {
	_, err := c.Client.DeleteRolePolicy(context.Background(), &iam.DeleteRolePolicyInput{
		PolicyName: aws.String(policyName),
		RoleName:   aws.String(roleName),
	})

	return err
}

func (c *IAMClient) DetachRolePolicy(roleName, policyArn string) error {
	_, err := c.Client.DetachRolePolicy(context.Background(), &iam.DetachRolePolicyInput{
		PolicyArn: aws.String(policyArn),
		RoleName:  aws.String(roleName),
	})

	return err
}

func (c *IAMClient) GetOpenIDConnectProvider(arn string) (*iam.GetOpenIDConnectProviderOutput, error) {
	return c.Client.GetOpenIDConnectProvider(context.Background(), &iam.GetOpenIDConnectProviderInput{
		OpenIDConnectProviderArn: aws.String(arn),
	})
}

func (c *IAMClient) GetRole(name string) (*types.Role, error) {
	result, err := c.Client.GetRole(context.Background(), &iam.GetRoleInput{
		RoleName: aws.String(name),
	})

	if err != nil {
		return nil, err
	}

	return result.Role, nil
}

// Lists all managed policies that are attached to the specified IAM role.
func (c *IAMClient) ListAttachedRolePolicies(roleName string) ([]types.AttachedPolicy, error) {
	policies := []types.AttachedPolicy{}
	pageNum := 0

	paginator := iam.NewListAttachedRolePoliciesPaginator(c.Client, &iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(roleName),
	})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		policies = append(policies, out.AttachedPolicies...)
		pageNum++
	}

	return policies, nil
}

func (c *IAMClient) ListOpenIDConnectProviders() ([]types.OpenIDConnectProviderListEntry, error) {
	result, err := c.Client.ListOpenIDConnectProviders(context.Background(), &iam.ListOpenIDConnectProvidersInput{})
	if err != nil {
		return nil, err
	}

	return result.OpenIDConnectProviderList, nil
}

// Lists the names of the inline policies that are embedded in the specified IAM role.
func (c *IAMClient) ListRolePolicies(roleName string) ([]string, error) {
	policyNames := []string{}
	pageNum := 0

	paginator := iam.NewListRolePoliciesPaginator(c.Client, &iam.ListRolePoliciesInput{
		RoleName: aws.String(roleName),
	})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		policyNames = append(policyNames, out.PolicyNames...)
		pageNum++
	}

	return policyNames, nil
}

func (c *IAMClient) ListRoles() ([]types.Role, error) {
	roles := []types.Role{}
	pageNum := 0

	paginator := iam.NewListRolesPaginator(c.Client, &iam.ListRolesInput{})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		roles = append(roles, out.Roles...)
		pageNum++
	}

	return roles, nil
}

func (c *IAMClient) PutRolePolicy(roleName, policyName, policyDoc string) error {
	_, err := c.Client.PutRolePolicy(context.Background(), &iam.PutRolePolicyInput{
		PolicyDocument: aws.String(policyDoc),
		PolicyName:     aws.String(policyName),
		RoleName:       aws.String(roleName),
	})

	return err
}
