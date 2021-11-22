package aws

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"

	amg "github.com/aws/aws-sdk-go/service/managedgrafana"
)

func AmgCreateWorkspace(name string, auth []string, roleArn string) (*amg.WorkspaceDescription, error) {
	sess := GetSession()
	svc := amg.New(sess)

	result, err := svc.CreateWorkspace(&amg.CreateWorkspaceInput{
		AccountAccessType:       aws.String(amg.AccountAccessTypeCurrentAccount),
		AuthenticationProviders: aws.StringSlice(auth),
		PermissionType:          aws.String(amg.PermissionTypeServiceManaged),
		WorkspaceDataSources:    aws.StringSlice([]string{amg.DataSourceTypePrometheus}),
		WorkspaceName:           aws.String(name),
		WorkspaceRoleArn:        aws.String(roleArn),
	})

	if err != nil {
		return nil, FormatError(err)
	}

	err = waitUntilWorkspaceActive(svc, &amg.DescribeWorkspaceInput{
		WorkspaceId: result.Workspace.Id,
	})

	return result.Workspace, err
}

func AmgDeleteWorkspace(id string) error {
	sess := GetSession()
	svc := amg.New(sess)

	_, err := svc.DeleteWorkspace(&amg.DeleteWorkspaceInput{
		WorkspaceId: aws.String(id),
	})

	return FormatError(err)
}

func AmgDescribeWorkspace(id string) (*amg.WorkspaceDescription, error) {
	sess := GetSession()
	svc := amg.New(sess)

	result, err := svc.DescribeWorkspace(&amg.DescribeWorkspaceInput{
		WorkspaceId: aws.String(id),
	})

	if err != nil {
		return nil, err
	}

	return result.Workspace, nil
}

func AmgListWorkspaces() ([]*amg.WorkspaceSummary, error) {
	sess := GetSession()
	svc := amg.New(sess)

	workspaces := []*amg.WorkspaceSummary{}
	pageNum := 0

	err := svc.ListWorkspacesPages(&amg.ListWorkspacesInput{},
		func(page *amg.ListWorkspacesOutput, lastPage bool) bool {
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

func AmgUpdateWorkspaceAuthentication(id, samlMetadataUrl string) error {
	sess := GetSession()
	svc := amg.New(sess)

	result, err := svc.DescribeWorkspace(&amg.DescribeWorkspaceInput{
		WorkspaceId: aws.String(id),
	})
	if err != nil {
		return err
	}

	_, err = svc.UpdateWorkspaceAuthentication(&amg.UpdateWorkspaceAuthenticationInput{
		AuthenticationProviders: result.Workspace.Authentication.Providers,
		SamlConfiguration: &amg.SamlConfiguration{
			IdpMetadata: &amg.IdpMetadata{
				Url: aws.String(samlMetadataUrl),
			},
			AssertionAttributes: &amg.AssertionAttributes{
				Role: aws.String("role"),
			},
			RoleValues: &amg.RoleValues{
				Admin: aws.StringSlice([]string{"admin"}),
			},
		},
		WorkspaceId: aws.String(id),
	})

	return err
}

func waitUntilWorkspaceActive(svc *amg.ManagedGrafana, input *amg.DescribeWorkspaceInput, opts ...request.WaiterOption) error {
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
			var inCpy *amg.DescribeWorkspaceInput
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
