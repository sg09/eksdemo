package iam_role

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hako/durafmt"
)

type IamRolePrinter struct {
	roles    []*iam.Role
	lastUsed bool
}

func NewPrinter(roles []*iam.Role, lastUsed bool) *IamRolePrinter {
	return &IamRolePrinter{roles, lastUsed}
}

func (p *IamRolePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()

	header := []string{"Age", "Role"}
	if p.lastUsed {
		header = append(header, "Last Used")
	}

	table.SetHeader(header)

	for _, r := range p.roles {
		age := durafmt.ParseShort(time.Since(aws.TimeValue(r.CreateDate)))
		rlu := r.RoleLastUsed

		row := []string{
			age.String(),
			aws.StringValue(r.RoleName),
		}

		if p.lastUsed {
			var lastUsed string

			if rlu != nil && rlu.LastUsedDate != nil {
				lastUsed = durafmt.ParseShort(time.Since(aws.TimeValue(rlu.LastUsedDate))).String()
			} else {
				lastUsed = "-"
			}
			row = append(row, lastUsed)
		}

		table.AppendRow(row)
	}

	table.Print(writer)

	return nil
}

func (p *IamRolePrinter) PrintJSON(writer io.Writer) error {
	if err := p.decodeAssumeRolePolicyDocuments(); err != nil {
		return err
	}
	return printer.EncodeJSON(writer, p.roles)
}

func (p *IamRolePrinter) PrintYAML(writer io.Writer) error {
	if err := p.decodeAssumeRolePolicyDocuments(); err != nil {
		return err
	}
	return printer.EncodeYAML(writer, p.roles)
}

func (p *IamRolePrinter) decodeAssumeRolePolicyDocuments() error {
	for i, r := range p.roles {
		decodedValue, err := url.QueryUnescape(aws.StringValue(r.AssumeRolePolicyDocument))
		if err != nil {
			return fmt.Errorf("unable to decode AssumeRolePolicyDocument for role %q: %w", aws.StringValue(r.RoleName), err) // TODO
		}
		p.roles[i] = r.SetAssumeRolePolicyDocument(decodedValue)
	}
	return nil
}
