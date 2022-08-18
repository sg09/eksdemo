package dns_record

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/hosted_zone"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/spf13/cobra"

	awssdk "github.com/aws/aws-sdk-go/aws"
)

type Manager struct {
	DryRun bool

	dnsRecordGetter Getter
	zoneGetter      hosted_zone.Getter
}

func (m *Manager) Create(options resource.Options) error {
	return fmt.Errorf("create dns-record not implemented")
}

func (m *Manager) Delete(options resource.Options) error {
	dnsOptions, ok := options.(*DnsRecordOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to DnsRecordOptions")
	}

	zone, err := m.zoneGetter.GetZoneByName(dnsOptions.ZoneName)
	if err != nil {
		return err
	}

	records, err := m.dnsRecordGetter.GetRecords(dnsOptions.Name, aws.StringValue(zone.Id))
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return fmt.Errorf("no records found with name %q", dnsOptions.Name)
	}

	if len(records) > 1 && !dnsOptions.AllTypes && !dnsOptions.AllRecords {
		return fmt.Errorf("multiple records found with name %q, use %q flag to delete all records", dnsOptions.Name, "--all")
	}

	changes := make([]*route53.Change, 0, len(records))

	for _, rec := range records {
		if aws.StringValue(rec.Type) == "NS" || aws.StringValue(rec.Type) == "SOA" {
			continue
		}

		change := &route53.Change{
			Action:            awssdk.String("DELETE"),
			ResourceRecordSet: rec,
		}
		changes = append(changes, change)
		fmt.Printf("Deleting %s record %q...\n", aws.StringValue(rec.Type), strings.TrimSuffix(aws.StringValue(rec.Name), "."))
	}

	if len(changes) == 0 {
		fmt.Println("No records to delete.")
		return nil
	}

	changeBatch := &route53.ChangeBatch{
		Changes: changes,
		Comment: awssdk.String("eksdemo delete dns-record"),
	}

	if err := aws.Route53ChangeResourceRecordSets(changeBatch, aws.StringValue(zone.Id)); err != nil {
		return err
	}
	fmt.Println("Record(s) deleted successfully")

	return nil
}

func (m *Manager) SetDryRun() {
	m.DryRun = true
}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}
