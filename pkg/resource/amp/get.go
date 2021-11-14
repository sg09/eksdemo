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
	ampSummaries, err := aws.AmpListWorkspaces(alias)
	workspaces := make([]*prometheusservice.WorkspaceDescription, 0, len(ampSummaries))

	if err != nil {
		return err
	}

	for _, summary := range ampSummaries {
		result, err := aws.AmpDescribeWorkspace(aws.StringValue(summary.WorkspaceId))
		if err != nil {
			return err
		}
		workspaces = append(workspaces, result)
	}

	return output.Print(os.Stdout, NewPrinter(workspaces))
}

func (g *Getter) GetAmpByAlias(alias string) (*prometheusservice.WorkspaceDescription, error) {
	found, err := aws.AmpListWorkspaces(alias)
	if err != nil {
		return nil, err
	}

	if len(found) == 0 {
		return nil, fmt.Errorf("workspace alias %q not found", alias)
	}

	if len(found) > 1 {
		return nil, fmt.Errorf("multiple workspaces found with alias: %s", alias)
	}

	workspace, err := aws.AmpDescribeWorkspace(aws.StringValue(found[0].WorkspaceId))
	if err != nil {
		return nil, err
	}

	return workspace, nil
}
