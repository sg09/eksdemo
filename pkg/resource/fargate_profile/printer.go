package fargate_profile

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"strings"
	"time"

	"github.com/hako/durafmt"

	"github.com/aws/aws-sdk-go/service/eks"
)

type FargateProfilePrinter struct {
	profiles []*eks.FargateProfile
}

func NewPrinter(profiles []*eks.FargateProfile) *FargateProfilePrinter {
	return &FargateProfilePrinter{profiles}
}

func (p *FargateProfilePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Status", "Name", "Selectors"})

	for _, profile := range p.profiles {
		age := durafmt.ParseShort(time.Since(aws.TimeValue(profile.CreatedAt)))
		name := aws.StringValue(profile.FargateProfileName)

		selectors := make([]string, 0, len(profile.Selectors))
		for _, s := range profile.Selectors {
			selectors = append(selectors, s.String())
		}

		table.AppendRow([]string{
			age.String(),
			aws.StringValue(profile.Status),
			name,
			strings.Join(selectors, ","),
		})
	}

	table.Print(writer)

	return nil
}

func (p *FargateProfilePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.profiles)
}

func (p *FargateProfilePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.profiles)
}
