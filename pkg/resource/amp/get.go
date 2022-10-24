package amp

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"fmt"
	"os"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/amp/types"
)

type Getter struct {
	prometheusClient *aws.AMPClient
}

func NewGetter(prometheusClient *aws.AMPClient) *Getter {
	return &Getter{prometheusClient}
}

func (g *Getter) Init() {
	if g.prometheusClient == nil {
		g.prometheusClient = aws.NewAMPClient()
	}
}

func (g *Getter) Get(alias string, output printer.Output, options resource.Options) error {
	workspaces, err := g.GetAll(alias)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(workspaces))
}

func (g *Getter) GetAll(alias string) ([]*types.WorkspaceDescription, error) {
	ampSummaries, err := g.prometheusClient.ListWorkspaces(alias)
	if err != nil {
		return nil, err
	}

	workspaces := make([]*types.WorkspaceDescription, 0, len(ampSummaries))

	for _, summary := range ampSummaries {
		// ListWorkspaces API will return workspaces that begin with alias, so drop those that don't match exactly
		if alias != "" && awssdk.ToString(summary.Alias) != alias {
			continue
		}

		result, err := g.prometheusClient.DescribeWorkspace(awssdk.ToString(summary.WorkspaceId))
		if err != nil {
			return nil, err
		}
		workspaces = append(workspaces, result)
	}

	if alias != "" && len(workspaces) == 0 {
		return nil, resource.NotFoundError(fmt.Sprintf("workspace alias %q not found", alias))
	}

	return workspaces, nil
}

func (g *Getter) GetAmpByAlias(alias string) (*types.WorkspaceDescription, error) {
	workspaces, err := g.GetAll(alias)
	if err != nil {
		return nil, err
	}

	found := []*types.WorkspaceDescription{}

	for _, w := range workspaces {
		if w.Status.StatusCode != types.WorkspaceStatusCodeDeleting {
			found = append(found, w)
		}
	}

	if len(found) == 0 {
		return nil, resource.NotFoundError(fmt.Sprintf("workspace alias %q not found", alias))
	}

	if len(found) > 1 {
		return nil, fmt.Errorf("multiple workspaces found with alias: %s", alias)
	}

	return found[0], nil
}
