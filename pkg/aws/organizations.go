package aws

import (
	"github.com/aws/aws-sdk-go/service/organizations"
)

func OrgsCreateOrganization() (*organizations.Organization, error) {
	sess := GetSession()
	svc := organizations.New(sess)

	result, err := svc.CreateOrganization(&organizations.CreateOrganizationInput{})
	if err != nil {
		return nil, err
	}

	return result.Organization, nil
}

func OrgsDeleteOrganization() error {
	sess := GetSession()
	svc := organizations.New(sess)

	_, err := svc.DeleteOrganization(&organizations.DeleteOrganizationInput{})
	return err
}

func OrgsDescribeOrganization() (*organizations.Organization, error) {
	sess := GetSession()
	svc := organizations.New(sess)

	result, err := svc.DescribeOrganization(&organizations.DescribeOrganizationInput{})
	if err != nil {
		return nil, err
	}

	return result.Organization, nil
}
