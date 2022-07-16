package s3_bucket

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hako/durafmt"
)

type BucketPrinter struct {
	buckets []*s3.Bucket
}

func NewPrinter(buckets []*s3.Bucket) *BucketPrinter {
	return &BucketPrinter{buckets}
}

func (p *BucketPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Name"})

	for _, b := range p.buckets {
		age := durafmt.ParseShort(time.Since(aws.TimeValue(b.CreationDate)))

		table.AppendRow([]string{
			age.String(),
			aws.StringValue(b.Name),
		})
	}
	table.Print(writer)

	return nil
}

func (p *BucketPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.buckets)
}

func (p *BucketPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.buckets)
}
