package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/prometheusservice"
)

func AmpCreateWorkspace(alias string) (*prometheusservice.CreateWorkspaceOutput, error) {
	sess := GetSession()
	svc := prometheusservice.New(sess)

	input := prometheusservice.CreateWorkspaceInput{}

	if alias != "" {
		input.Alias = aws.String(alias)
	}

	result, err := svc.CreateWorkspace(&input)
	if err != nil {
		return nil, err
	}

	err = svc.WaitUntilWorkspaceActive(&prometheusservice.DescribeWorkspaceInput{
		WorkspaceId: result.WorkspaceId,
	})

	return result, err
}

func AmpDeleteWorkspace(id string) error {
	sess := GetSession()
	svc := prometheusservice.New(sess)

	_, err := svc.DeleteWorkspace(&prometheusservice.DeleteWorkspaceInput{
		WorkspaceId: aws.String(id),
	})

	return FormatError(err)
}

func AmpDescribeWorkspace(workspaceId string) (*prometheusservice.WorkspaceDescription, error) {
	sess := GetSession()
	svc := prometheusservice.New(sess)

	result, err := svc.DescribeWorkspace(&prometheusservice.DescribeWorkspaceInput{
		WorkspaceId: aws.String(workspaceId),
	})

	if err != nil {
		return nil, err
	}

	return result.Workspace, nil
}

func AmpListWorkspaces(alias string) ([]*prometheusservice.WorkspaceSummary, error) {
	sess := GetSession()
	svc := prometheusservice.New(sess)

	workspaces := []*prometheusservice.WorkspaceSummary{}
	pageNum := 0

	input := prometheusservice.ListWorkspacesInput{}
	if alias != "" {
		input.Alias = aws.String(alias)
	}

	err := svc.ListWorkspacesPages(&input,
		func(page *prometheusservice.ListWorkspacesOutput, lastPage bool) bool {
			pageNum++
			workspaces = append(workspaces, page.Workspaces...)
			return pageNum <= maxPages
		},
	)

	if err != nil {
		return nil, err
	}

	return workspaces, nil
}
