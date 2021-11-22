package amp

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/prometheusservice"
)

type Getter struct{}

func (g *Getter) Get(alias string, output printer.Output, options resource.Options) error {
	workspaces, err := g.GetAll(alias)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(workspaces))
}

func (g *Getter) GetAll(alias string) ([]*prometheusservice.WorkspaceDescription, error) {
	ampSummaries, err := aws.AmpListWorkspaces(alias)
	if err != nil {
		return nil, err
	}

	workspaces := make([]*prometheusservice.WorkspaceDescription, 0, len(ampSummaries))

	for _, summary := range ampSummaries {
		// ListWorkspaces API will return workspaces that begin with alias, so drop those that don't match exactly
		if alias != "" && aws.StringValue(summary.Alias) != alias {
			continue
		}

		result, err := aws.AmpDescribeWorkspace(aws.StringValue(summary.WorkspaceId))
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

func (g *Getter) GetAmpByAlias(alias string) (*prometheusservice.WorkspaceDescription, error) {
	workspaces, err := g.GetAll(alias)
	if err != nil {
		return nil, err
	}

	found := []*prometheusservice.WorkspaceDescription{}

	for _, w := range workspaces {
		if aws.StringValue(w.Status.StatusCode) != "DELETING" {
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
