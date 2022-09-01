package ssm_session

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hako/durafmt"
)

type SessionPrinter struct {
	sessions []*ssm.Session
}

func NewPrinter(sessions []*ssm.Session) *SessionPrinter {
	return &SessionPrinter{sessions}
}

func (p *SessionPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Status", "Id", "Instance"})

	for _, s := range p.sessions {
		age := durafmt.ParseShort(time.Since(aws.TimeValue(s.StartDate)))

		table.AppendRow([]string{
			age.String(),
			aws.StringValue(s.Status),
			aws.StringValue(s.SessionId),
			aws.StringValue(s.Target),
		})
	}

	table.Print(writer)

	return nil
}

func (p *SessionPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.sessions)
}

func (p *SessionPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.sessions)
}
