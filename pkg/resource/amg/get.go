package amg

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/managedgrafana"
)

type Getter struct{}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	var workspaces []*managedgrafana.WorkspaceDescription
	var err error

	if name == "" {
		workspaces, err = g.GetAll()
	} else {
		workspaces, err = g.GetAmgByName(name)
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(workspaces))
}

func (g *Getter) GetAll() ([]*managedgrafana.WorkspaceDescription, error) {
	amgSummaries, err := aws.AmgListWorkspaces()
	workspaces := make([]*managedgrafana.WorkspaceDescription, 0, len(amgSummaries))

	if err != nil {
		return nil, err
	}

	for _, summary := range amgSummaries {
		result, err := aws.AmgDescribeWorkspace(aws.StringValue(summary.Id))
		if err != nil {
			return nil, err
		}
		workspaces = append(workspaces, result)
	}

	return workspaces, nil
}

func (g *Getter) GetAmgByName(name string) ([]*managedgrafana.WorkspaceDescription, error) {
	all, err := aws.AmgListWorkspaces()
	if err != nil {
		return nil, err
	}

	ampIds := []string{}

	for _, workspace := range all {
		if aws.StringValue(workspace.Name) == name && aws.StringValue(workspace.Status) != "DELETING" {
			ampIds = append(ampIds, aws.StringValue(workspace.Id))
		}
	}

	if len(ampIds) == 0 {
		return nil, fmt.Errorf("workspace name %q not found", name)
	}

	workspaces := make([]*managedgrafana.WorkspaceDescription, 0, len(ampIds))

	for _, id := range ampIds {
		result, err := aws.AmgDescribeWorkspace(id)
		if err != nil {
			return nil, err
		}
		workspaces = append(workspaces, result)
	}

	return workspaces, nil
}
