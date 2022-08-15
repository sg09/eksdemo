package ecr_repository

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/hako/durafmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type RepositoryPrinter struct {
	repos []*ecr.Repository
}

func NewPrinter(repos []*ecr.Repository) *RepositoryPrinter {
	return &RepositoryPrinter{repos}
}

func (p *RepositoryPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Name", "Tags", "Scan", "Encryption"})

	caser := cases.Title(language.English)

	for _, repo := range p.repos {
		age := durafmt.ParseShort(time.Since(aws.TimeValue(repo.CreatedAt)))
		scan := "Manual"
		if aws.BoolValue(repo.ImageScanningConfiguration.ScanOnPush) {
			scan = "Scan on push"
		}

		table.AppendRow([]string{
			age.String(),
			aws.StringValue(repo.RepositoryName),
			caser.String(aws.StringValue(repo.ImageTagMutability)),
			scan,
			aws.StringValue(repo.EncryptionConfiguration.EncryptionType),
		})
	}

	table.Print(writer)

	return nil
}

func (p *RepositoryPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.repos)
}

func (p *RepositoryPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.repos)
}
