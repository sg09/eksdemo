package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/amp"
	"github.com/aws/aws-sdk-go-v2/service/amp/types"
)

type AMPClient struct {
	*amp.Client
}

func NewAMPClient() *AMPClient {
	return &AMPClient{amp.NewFromConfig(GetConfig())}
}

func (c *AMPClient) CreateWorkspace(alias string) (*amp.CreateWorkspaceOutput, error) {
	input := amp.CreateWorkspaceInput{}

	if alias != "" {
		input.Alias = aws.String(alias)
	}

	result, err := c.Client.CreateWorkspace(context.Background(), &input)
	if err != nil {
		return nil, err
	}

	err = amp.NewWorkspaceActiveWaiter(c.Client).Wait(context.Background(),
		&amp.DescribeWorkspaceInput{WorkspaceId: result.WorkspaceId},
		1*time.Minute,
	)

	return result, err
}

func (c *AMPClient) DeleteWorkspace(id string) error {
	_, err := c.Client.DeleteWorkspace(context.Background(), &amp.DeleteWorkspaceInput{
		WorkspaceId: aws.String(id),
	})

	return err
}

func (c *AMPClient) DescribeWorkspace(workspaceId string) (*types.WorkspaceDescription, error) {
	out, err := c.Client.DescribeWorkspace(context.Background(), &amp.DescribeWorkspaceInput{
		WorkspaceId: aws.String(workspaceId),
	})

	if err != nil {
		return nil, err
	}

	return out.Workspace, nil
}

func (c *AMPClient) ListWorkspaces(alias string) ([]types.WorkspaceSummary, error) {
	workspaces := []types.WorkspaceSummary{}
	pageNum := 0

	input := amp.ListWorkspacesInput{}
	if alias != "" {
		input.Alias = aws.String(alias)
	}

	paginator := amp.NewListWorkspacesPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		workspaces = append(workspaces, out.Workspaces...)
		pageNum++
	}

	return workspaces, nil
}
