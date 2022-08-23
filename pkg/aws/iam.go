package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
)

func IamCreateRole(assumeRolePolicy, name, path string) (*iam.Role, error) {
	sess := GetSession()
	svc := iam.New(sess)

	if path == "" {
		path = "/"
	}

	result, err := svc.CreateRole(&iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(assumeRolePolicy),
		RoleName:                 aws.String(name),
		Path:                     aws.String(path),
	})
	if err != nil {
		return nil, err
	}

	return result.Role, nil
}

func IamDeleteRole(name string) error {
	sess := GetSession()
	svc := iam.New(sess)

	_, err := svc.DeleteRole(&iam.DeleteRoleInput{
		RoleName: aws.String(name),
	})

	return err
}

// Deletes the specified inline policy that is embedded in the specified IAM role.
func IamDeleteRolePolicy(roleName, policyName string) error {
	sess := GetSession()
	svc := iam.New(sess)

	_, err := svc.DeleteRolePolicy(&iam.DeleteRolePolicyInput{
		PolicyName: aws.String(policyName),
		RoleName:   aws.String(roleName),
	})
	if err != nil {
		return err
	}

	return nil
}

func IamDetachRolePolicy(roleName, policyArn string) error {
	sess := GetSession()
	svc := iam.New(sess)

	_, err := svc.DetachRolePolicy(&iam.DetachRolePolicyInput{
		PolicyArn: aws.String(policyArn),
		RoleName:  aws.String(roleName),
	})

	if err != nil {
		return err
	}

	return nil
}

func IamGetOpenIDConnectProvider(arn string) (*iam.GetOpenIDConnectProviderOutput, error) {
	sess := GetSession()
	svc := iam.New(sess)

	result, err := svc.GetOpenIDConnectProvider(&iam.GetOpenIDConnectProviderInput{
		OpenIDConnectProviderArn: aws.String(arn),
	})
	if err != nil {
		return nil, FormatError(err)
	}

	return result, nil
}

func IamGetRole(name string) (*iam.Role, error) {
	sess := GetSession()
	svc := iam.New(sess)

	result, err := svc.GetRole(&iam.GetRoleInput{
		RoleName: aws.String(name),
	})
	if err != nil {
		return nil, FormatError(err)
	}

	return result.Role, nil
}

// Lists all managed policies that are attached to the specified IAM role.
func IamListAttachedRolePolicies(roleName string) ([]*iam.AttachedPolicy, error) {
	sess := GetSession()
	svc := iam.New(sess)

	policies := []*iam.AttachedPolicy{}
	pageNum := 0

	err := svc.ListAttachedRolePoliciesPages(&iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(roleName),
	},
		func(page *iam.ListAttachedRolePoliciesOutput, lastPage bool) bool {
			pageNum++
			policies = append(policies, page.AttachedPolicies...)
			return pageNum <= maxPages
		},
	)
	if err != nil {
		return nil, err
	}

	return policies, err
}

func IamListOpenIDConnectProviders() ([]*iam.OpenIDConnectProviderListEntry, error) {
	sess := GetSession()
	svc := iam.New(sess)

	result, err := svc.ListOpenIDConnectProviders(&iam.ListOpenIDConnectProvidersInput{})
	if err != nil {
		return nil, FormatError(err)
	}

	return result.OpenIDConnectProviderList, nil
}

// Lists the names of the inline policies that are embedded in the specified IAM role.
func IamListRolePolicies(roleName string) ([]string, error) {
	sess := GetSession()
	svc := iam.New(sess)

	policyNames := []string{}
	pageNum := 0

	err := svc.ListRolePoliciesPages(&iam.ListRolePoliciesInput{
		RoleName: aws.String(roleName),
	},
		func(page *iam.ListRolePoliciesOutput, lastPage bool) bool {
			pageNum++
			policyNames = append(policyNames, aws.StringValueSlice(page.PolicyNames)...)
			return pageNum <= maxPages
		},
	)
	if err != nil {
		return nil, err
	}

	return policyNames, err
}

func IamListRoles() ([]*iam.Role, error) {
	sess := GetSession()
	svc := iam.New(sess)

	roles := []*iam.Role{}
	pageNum := 0

	err := svc.ListRolesPages(&iam.ListRolesInput{},
		func(page *iam.ListRolesOutput, lastPage bool) bool {
			pageNum++
			roles = append(roles, page.Roles...)
			return pageNum <= maxPages
		},
	)

	if err != nil {
		return nil, FormatError(err)
	}

	return roles, nil
}

func IamPutRolePolicy(roleName, policyName, policyDoc string) error {
	sess := GetSession()
	svc := iam.New(sess)

	_, err := svc.PutRolePolicy(&iam.PutRolePolicyInput{
		PolicyDocument: aws.String(policyDoc),
		PolicyName:     aws.String(policyName),
		RoleName:       aws.String(roleName),
	})
	if err != nil {
		return err
	}
	return nil
}
