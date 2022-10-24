package amg

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/template"
	"fmt"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/grafana/types"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/spf13/cobra"
)

type Manager struct {
	AssumeRolePolicyTemplate template.TextTemplate
	DryRun                   bool
	grafanaClient            *aws.GrafanaClient
	grafanaGetter            *Getter
}

func (m *Manager) Init() {
	if m.grafanaClient == nil {
		m.grafanaClient = aws.NewGrafanaClient()
	}
	m.grafanaGetter = NewGetter(m.grafanaClient)
}

func (m *Manager) Create(options resource.Options) error {
	amgOptions, ok := options.(*AmgOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to AmpOptions")
	}

	workspace, err := m.grafanaGetter.GetAmgByName(amgOptions.WorkspaceName)
	if err != nil {
		if _, ok := err.(resource.NotFoundError); !ok {
			// Return an error if it's anything other than resource not found
			return err
		}
	}

	if workspace != nil {
		fmt.Printf("AMG Workspace %q already exists\n", amgOptions.WorkspaceName)
		return nil
	}

	if m.DryRun {
		return m.dryRun(amgOptions)
	}

	role, err := m.createIamRole(amgOptions)
	if err != nil {
		return err
	}

	err = aws.IamPutRolePolicy(awssdk.ToString(role.RoleName), rolePolicName, rolePolicyDoc)
	if err != nil {
		return err
	}

	fmt.Printf("Creating AMG Workspace Name: %s...", amgOptions.WorkspaceName)
	result, err := m.grafanaClient.CreateWorkspace(amgOptions.WorkspaceName, amgOptions.Auth, awssdk.ToString(role.Arn))
	if err != nil {
		fmt.Println()
		return err
	}

	fmt.Printf("done\nCreated AMG Workspace Id: %s\n", awssdk.ToString(result.Id))

	return nil
}

func (m *Manager) Delete(options resource.Options) error {
	amgOptions, ok := options.(*AmgOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to AmgOptions")
	}

	var amg *types.WorkspaceDescription
	var err error

	if options.Common().Id == "" {
		amg, err = m.grafanaGetter.GetAmgByName(amgOptions.WorkspaceName)
		if err != nil {
			if _, ok := err.(resource.NotFoundError); ok {
				fmt.Printf("AMG Workspace Name %q does not exist\n", amgOptions.WorkspaceName)
				return nil
			}
			return err
		}
	} else {
		amg, err = m.grafanaClient.DescribeWorkspace(options.Common().Id)
		if err != nil {
			return err
		}
	}

	err = m.deleteIamRole(awssdk.ToString(amg.WorkspaceRoleArn))
	if err != nil {
		return err
	}

	id := awssdk.ToString(amg.Id)

	err = m.grafanaClient.DeleteWorkspace(id)
	if err != nil {
		return err
	}
	fmt.Printf("AMG Workspace Id %q deleting...\n", id)

	return nil
}

func (m *Manager) SetDryRun() {
	m.DryRun = true
}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}

func (m *Manager) createIamRole(options *AmgOptions) (*iam.Role, error) {
	assumeRolePolicy, err := m.AssumeRolePolicyTemplate.Render(options)
	if err != nil {
		return nil, err
	}

	roleName := options.iamRoleName()

	role, err := aws.IamCreateRole(assumeRolePolicy, roleName, "/service-role/")
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == iam.ErrCodeEntityAlreadyExistsException {
				fmt.Printf("IAM Role %q already exists\n", roleName)
				return aws.IamGetRole(roleName)
			}
		}
		return nil, err
	}

	fmt.Printf("Created IAM Role: %s\n", awssdk.ToString(role.RoleName))

	return role, nil
}

func (m *Manager) deleteIamRole(roleArn string) error {
	roleName := roleArn[strings.LastIndex(roleArn, "/")+1:]

	// Delete inline policies before deleting role
	inlinePolicyNames, err := aws.IamListRolePolicies(roleName)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == iam.ErrCodeNoSuchEntityException {
				return nil
			}
		}
		return err
	}

	for _, policyName := range inlinePolicyNames {
		err := aws.IamDeleteRolePolicy(roleName, policyName)
		if err != nil {
			return err
		}
	}

	// Remove managed policies before deleting role
	mgdPolicies, err := aws.IamListAttachedRolePolicies(roleName)
	if err != nil {
		return err
	}

	for _, policy := range mgdPolicies {
		err := aws.IamDetachRolePolicy(roleName, awssdk.ToString(policy.PolicyArn))
		if err != nil {
			return err
		}
	}

	return aws.IamDeleteRole(roleName)
}

func (m *Manager) dryRun(options *AmgOptions) error {
	fmt.Println("\nAMG Resource Manager Dry Run:")

	fmt.Printf("Amazon Managed Grafana API Call %q with request parameters:\n", "CreateWorkspace")
	fmt.Printf("AccountAccessType: %q\n", types.AccountAccessTypeCurrentAccount)
	fmt.Printf("AuthenticationProviders: %q\n", options.Auth)

	fmt.Printf("PermissionType: %q\n", types.PermissionTypeServiceManaged)
	fmt.Printf("WorkspaceDataSources: %q\n", []types.DataSourceType{types.DataSourceTypePrometheus})
	fmt.Printf("WorkspaceName: %q\n", options.WorkspaceName)
	fmt.Printf("WorkspaceRoleArn: %q\n", "<role-to-be-created>")

	return nil
}
