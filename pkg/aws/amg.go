package aws

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"

	"github.com/aws/aws-sdk-go/service/managedgrafana"
)

func AmgCreateWorkspace(name string, auth []string, roleArn string) (*managedgrafana.WorkspaceDescription, error) {
	sess := GetSession()
	svc := managedgrafana.New(sess)

	result, err := svc.CreateWorkspace(&managedgrafana.CreateWorkspaceInput{
		AccountAccessType:       aws.String(managedgrafana.AccountAccessTypeCurrentAccount),
		AuthenticationProviders: aws.StringSlice(auth),
		PermissionType:          aws.String(managedgrafana.PermissionTypeServiceManaged),
		WorkspaceDataSources:    aws.StringSlice([]string{managedgrafana.DataSourceTypePrometheus}),
		WorkspaceName:           aws.String(name),
		WorkspaceRoleArn:        aws.String(roleArn),
	})

	if err != nil {
		return nil, FormatError(err)
	}

	err = waitUntilWorkspaceActive(svc, &managedgrafana.DescribeWorkspaceInput{
		WorkspaceId: result.Workspace.Id,
	})

	return result.Workspace, err
}

func AmgDeleteWorkspace(id string) error {
	sess := GetSession()
	svc := managedgrafana.New(sess)

	_, err := svc.DeleteWorkspace(&managedgrafana.DeleteWorkspaceInput{
		WorkspaceId: aws.String(id),
	})

	return FormatError(err)
}

func AmgDescribeWorkspace(id string) (*managedgrafana.WorkspaceDescription, error) {
	sess := GetSession()
	svc := managedgrafana.New(sess)

	result, err := svc.DescribeWorkspace(&managedgrafana.DescribeWorkspaceInput{
		WorkspaceId: aws.String(id),
	})

	if err != nil {
		return nil, err
	}

	return result.Workspace, nil
}

func AmgListWorkspaces() ([]*managedgrafana.WorkspaceSummary, error) {
	sess := GetSession()
	svc := managedgrafana.New(sess)

	workspaces := []*managedgrafana.WorkspaceSummary{}
	pageNum := 0

	err := svc.ListWorkspacesPages(&managedgrafana.ListWorkspacesInput{},
		func(page *managedgrafana.ListWorkspacesOutput, lastPage bool) bool {
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

func AmgUpdateWorkspace(id string) (*managedgrafana.WorkspaceDescription, error) {
	sess := GetSession()
	svc := managedgrafana.New(sess)

	result, err := svc.UpdateWorkspace(&managedgrafana.UpdateWorkspaceInput{
		WorkspaceDataSources: aws.StringSlice([]string{managedgrafana.DataSourceTypePrometheus}),
		WorkspaceId:          aws.String(id),
	})

	if err != nil {
		return nil, FormatError(err)
	}

	return result.Workspace, nil
}

func waitUntilWorkspaceActive(svc *managedgrafana.ManagedGrafana, input *managedgrafana.DescribeWorkspaceInput, opts ...request.WaiterOption) error {
	ctx := aws.BackgroundContext()

	w := request.Waiter{
		Name:        "WaitUntilWorkspaceActive",
		MaxAttempts: 60,
		Delay:       request.ConstantWaiterDelay(2 * time.Second),
		Acceptors: []request.WaiterAcceptor{
			{
				State:   request.SuccessWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "workspace.status",
				Expected: "ACTIVE",
			},
			{
				State:   request.RetryWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "workspace.status",
				Expected: "UPDATING",
			},
			{
				State:   request.RetryWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "workspace.status",
				Expected: "CREATING",
			},
		},
		Logger: svc.Config.Logger,
		NewRequest: func(opts []request.Option) (*request.Request, error) {
			var inCpy *managedgrafana.DescribeWorkspaceInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req, _ := svc.DescribeWorkspaceRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}
