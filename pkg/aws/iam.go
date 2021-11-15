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

func IamGetOpenIDConnectProviders(arn string) (*iam.GetOpenIDConnectProviderOutput, error) {
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

func IamListOpenIDConnectProviders() ([]*iam.OpenIDConnectProviderListEntry, error) {
	sess := GetSession()
	svc := iam.New(sess)

	result, err := svc.ListOpenIDConnectProviders(&iam.ListOpenIDConnectProvidersInput{})
	if err != nil {
		return nil, FormatError(err)
	}

	return result.OpenIDConnectProviderList, nil
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
			return pageNum <= 10
		},
	)

	if err != nil {
		return nil, FormatError(err)
	}

	return roles, nil
}
