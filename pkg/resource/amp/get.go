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
	ampSummary, err := aws.AmpListWorkspaces(alias)
	if err != nil {
		return nil, err
	}

	if len(ampSummary) == 0 {
		return nil, fmt.Errorf("workspace alias %q not found", alias)
	}

	workspace, err := aws.AmpDescribeWorkspace(aws.StringValue(ampSummary[0].WorkspaceId))
	if err != nil {
		return nil, err
	}

	return workspace, nil
}
