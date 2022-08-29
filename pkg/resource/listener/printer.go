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
	table.SetHeader([]string{"Id", "Prot:Port", "Default Certificate Id", "Default Action"})

	resourceId := regexp.MustCompile(`[^:/]*$`)

	for _, l := range p.listeners {
		// DescribeListeners API documentation states that only the default certificate is included
		// https://docs.aws.amazon.com/elasticloadbalancing/latest/APIReference/API_DescribeListeners.html
		defaultCert := "-"
		if len(l.Certificates) > 0 {
			defaultCert = resourceId.FindString(aws.StringValue(l.Certificates[0].CertificateArn))
		}

		table.AppendRow([]string{
			resourceId.FindString(aws.StringValue(l.ListenerArn)),
			aws.StringValue(l.Protocol) + ":" + strconv.FormatInt(aws.Int64Value(l.Port), 10),
			defaultCert,
			strings.Join(PrintActions(l.DefaultActions), "\n"),
		})
	}

	table.SeparateRows()
	table.Print(writer)

	return nil
}

func (p *ListenerPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.listeners)
}

func (p *ListenerPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.listeners)
}

func PrintActions(elbActions []*elbv2.Action) (actions []string) {
	for _, a := range elbActions {
		switch {
		case a.AuthenticateCognitoConfig != nil || a.AuthenticateOidcConfig != nil:
			actions = append(actions, "TODO: authenticate action")

		case a.FixedResponseConfig != nil:
			actions = append(actions, "return fixed response "+aws.StringValue(a.FixedResponseConfig.StatusCode))

		case a.ForwardConfig != nil:
			tgNames := []string{}
			for _, tg := range a.ForwardConfig.TargetGroups {
				tgNames = append(tgNames, tgName(aws.StringValue(tg.TargetGroupArn)))
			}
			actions = append(actions, "forward to "+strings.Join(tgNames, "\n"))

		case a.RedirectConfig != nil:
			prot := aws.StringValue(a.RedirectConfig.Protocol)
			host := aws.StringValue(a.RedirectConfig.Host)
			port := aws.StringValue(a.RedirectConfig.Port)
			path := aws.StringValue(a.RedirectConfig.Path)
			query := aws.StringValue(a.RedirectConfig.Query)
			actions = append(actions, "redirect to "+prot+"://"+host+":"+port+path+"?"+query)
		}
	}
	return
}

func tgName(tgArn string) string {
	parts := strings.Split(tgArn, "/")
	if len(parts) < 2 {
		return "failed to parse arn"
	}
	return parts[1]
}
