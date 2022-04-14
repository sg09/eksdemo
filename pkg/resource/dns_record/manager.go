package dns_record

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/hosted_zone"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/route53"

	awssdk "github.com/aws/aws-sdk-go/aws"
)

type Manager struct {
	DryRun bool

	dnsRecordGetter Getter
	zoneGetter      hosted_zone.Getter
}

// Record types to skip when deleting
var deleteFilterTypes = map[string]bool{
	"NS":  true,
	"SOA": true,
	"TXT": true,
}

// Don't delete NS or SOA record types when deleting with --all flag
var deleteFilterTypesAll = map[string]bool{
	"NS":  true,
	"SOA": true,
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

	filterTypes := deleteFilterTypes
	if dnsOptions.All {
		filterTypes = deleteFilterTypesAll
	}

	records, err := m.dnsRecordGetter.GetRecords(dnsOptions.Name, aws.StringValue(zone.Id), filterTypes)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return fmt.Errorf("no records found with name %q", dnsOptions.Name)
	}

	if len(records) > 1 && !dnsOptions.All {
		return fmt.Errorf("multiple records found starting with name %q, use --all flag to delete all records", dnsOptions.Name)
	}

	changes := make([]*route53.Change, 0, len(records))

	for _, rec := range records {
		change := &route53.Change{
			Action:            awssdk.String("DELETE"),
			ResourceRecordSet: rec,
		}
		changes = append(changes, change)
		fmt.Printf("Deleting %s record %q...\n", aws.StringValue(rec.Type), strings.TrimSuffix(aws.StringValue(rec.Name), "."))
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

func (m *Manager) Update(options resource.Options) error {
	return fmt.Errorf("feature not supported")
}
