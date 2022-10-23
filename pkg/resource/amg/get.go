package amg

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/managedgrafana"
)

type Getter struct {
	resource.EmptyInit
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	var workspaces []*managedgrafana.WorkspaceDescription
	var err error

	if name == "" {
		workspaces, err = g.GetAll()
	} else {
		workspaces, err = g.GetAllAmgByName(name)
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

func (g *Getter) GetAllAmgByName(name string) ([]*managedgrafana.WorkspaceDescription, error) {
	summaries, err := aws.AmgListWorkspaces()
	if err != nil {
		return nil, err
	}

	workspaces := make([]*managedgrafana.WorkspaceDescription, 0, len(summaries))

	for _, s := range summaries {
		if aws.StringValue(s.Name) != name {
			continue
		}

		result, err := aws.AmgDescribeWorkspace(aws.StringValue(s.Id))
		if err != nil {
			return nil, err
		}
		workspaces = append(workspaces, result)
	}

	if len(workspaces) == 0 {
		return nil, resource.NotFoundError(fmt.Sprintf("workspace name %q not found", name))
	}

	return workspaces, nil
}

func (g *Getter) GetAmgByName(name string) (*managedgrafana.WorkspaceDescription, error) {
	workspaces, err := g.GetAllAmgByName(name)
	if err != nil {
		return nil, err
	}

	found := []*managedgrafana.WorkspaceDescription{}

	for _, w := range workspaces {
		if aws.StringValue(w.Status) != "DELETING" {
			found = append(found, w)
		}
	}

	if len(found) == 0 {
		return nil, resource.NotFoundError(fmt.Sprintf("workspace name %q not found", name))
	}

	if len(found) > 1 {
		return nil, fmt.Errorf("multiple workspaces found with name: %s", name)
	}

	return found[0], nil
}
