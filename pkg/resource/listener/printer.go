package listener

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/service/elbv2"
)

type ListenerPrinter struct {
	listeners []*elbv2.Listener
}

func NewPrinter(listeners []*elbv2.Listener) *ListenerPrinter {
	return &ListenerPrinter{listeners}
}

func (p *ListenerPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Proto:Port", "Target Group Name", "Default Certificate Id"})

	resourceId := regexp.MustCompile(`[^:/]*$`)

	for _, l := range p.listeners {
		tgNames := []string{}
		for _, da := range l.DefaultActions {
			tgNames = append(tgNames, tgName(aws.StringValue(da.TargetGroupArn)))
		}
		if len(tgNames) > 1 {
			table.SeparateRows()
		}

		// DescribeListeners API documentation states that only the default certificate is included
		// https://docs.aws.amazon.com/elasticloadbalancing/latest/APIReference/API_DescribeListeners.html
		defaultCert := "-"
		if len(l.Certificates) > 0 {
			defaultCert = resourceId.FindString(aws.StringValue(l.Certificates[0].CertificateArn))
		}

		table.AppendRow([]string{
			aws.StringValue(l.Protocol) + ":" + strconv.FormatInt(aws.Int64Value(l.Port), 10),
			strings.Join(tgNames, "\n"),
			defaultCert,
		})
	}

	table.Print(writer)

	return nil
}

func (p *ListenerPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.listeners)
}

func (p *ListenerPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.listeners)
}

func tgName(tgArn string) string {
	parts := strings.Split(tgArn, "/")
	if len(parts) < 2 {
		return "failed to parse arn"
	}
	return parts[1]
}
