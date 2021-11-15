package amg

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/resource"
	"fmt"
	"strings"
)

type Manager struct {
	Getter
}

func (m *Manager) Create(options resource.Options) error {
	amgOptions, ok := options.(*AmgOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to AmpOptions")
	}

	all, err := aws.AmgListWorkspaces()
	if err != nil {
		return err
	}

	ampIds := []string{}

	for _, workspace := range all {
		if aws.StringValue(workspace.Name) == amgOptions.WorkspaceName && aws.StringValue(workspace.Status) != "DELETING" {
			ampIds = append(ampIds, aws.StringValue(workspace.Id))
		}
	}

	if len(ampIds) == 1 {
		fmt.Printf("AMG Workspace %q already exists\n", amgOptions.WorkspaceName)
		return nil
	}

	if len(ampIds) > 1 {
		return fmt.Errorf("multiple workspaces found with name: %s", amgOptions.WorkspaceName)
	}

	fmt.Printf("Creating AMG with Name: %s...", amgOptions.WorkspaceName)

	roleName := m.iamRoleName(amgOptions.WorkspaceName)
	role, err := aws.IamCreateRole(assumeRolePolicy, roleName, "/service-role/")
	if err != nil {
		return err
	}

	result, err := aws.AmgCreateWorkspace(amgOptions.WorkspaceName, amgOptions.Auth, aws.StringValue(role.Arn))
	if err != nil {
		return err
	}
	fmt.Printf("done\nCreated AMG Workspace Id: %s\n", aws.StringValue(result.Id))

	return nil
}

func (m *Manager) Delete(options resource.Options) error {
	amgOptions, ok := options.(*AmgOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to AmgOptions")
	}

	amg, err := aws.AmgDescribeWorkspace(amgOptions.WorkspaceName)
	if err != nil {
		return err
	}

	roleArn := aws.StringValue(amg.WorkspaceRoleArn)
	roleName := roleArn[strings.LastIndex(roleArn, "/")+1:]

	err = aws.IamDeleteRole(roleName)
	if err != nil {
		return err
	}

	err = aws.AmgDeleteWorkspace(amgOptions.WorkspaceName)
	if err != nil {
		return err
	}
	fmt.Printf("AMG %q deleting...\n", amgOptions.WorkspaceName)

	return nil
}

func (m *Manager) SetDryRun() {}

func (m *Manager) iamRoleName(name string) string {
	return fmt.Sprintf("eksdemo.amg.%s", name)
}

const assumeRolePolicy = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "grafana.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
`
