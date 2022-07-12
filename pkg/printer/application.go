package printer

import (
	"fmt"
	"io"

	"helm.sh/helm/v3/pkg/release"
)

type ApplicationPrinter struct {
	releases []*release.Release
}

func NewApplicationPrinter(releases []*release.Release) *ApplicationPrinter {
	return &ApplicationPrinter{releases}
}

func (p *ApplicationPrinter) PrintTable(writer io.Writer) error {
	if len(p.releases) == 0 {
		fmt.Fprint(writer, "No helm releases found.\n")
		return nil
	}

	table := NewTablePrinter()
	table.SetHeader([]string{"Name", "Namespace", "Version", "Status", "Chart"})

	for _, r := range p.releases {
		table.AppendRow([]string{
			r.Name,
			r.Namespace,
			r.Chart.Metadata.AppVersion,
			r.Info.Status.String(),
			r.Chart.Metadata.Version,
		})
	}
	table.Print(writer)

	return nil
}

func (p *ApplicationPrinter) PrintJSON(writer io.Writer) error {
	return EncodeJSON(writer, p.releases)
}

func (p *ApplicationPrinter) PrintYAML(writer io.Writer) error {
	return EncodeYAML(writer, p.releases)
}
