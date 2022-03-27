package irsa

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hako/durafmt"
)

type IrsaPrinter struct {
	Roles []*iam.Role
}

func NewPrinter(Roles []*iam.Role) *IrsaPrinter {
	return &IrsaPrinter{Roles}
}

func (p *IrsaPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Role", "Last Used"})

	lastUsed := ""

	for _, r := range p.Roles {
		age := durafmt.ParseShort(time.Since(aws.TimeValue(r.CreateDate)))
		rlu := r.RoleLastUsed

		if rlu.LastUsedDate != nil {
			lastUsed = durafmt.ParseShort(time.Since(aws.TimeValue(rlu.LastUsedDate))).String()
		} else {
			lastUsed = "-"
		}

		table.AppendRow([]string{
			age.String(),
			aws.StringValue(r.RoleName),
			lastUsed,
		})
	}

	table.Print(writer)

	return nil
}

func (p *IrsaPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.Roles)
}

func (p *IrsaPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.Roles)
}
