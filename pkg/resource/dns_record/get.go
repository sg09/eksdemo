package dns_record

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/hosted_zone"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/route53"
)

type Getter struct {
	zoneGetter hosted_zone.Getter
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	dnsOptions, ok := options.(*DnsRecordOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to DnsRecordOptions")
	}

	zone, err := g.zoneGetter.GetZoneByName(dnsOptions.ZoneName)
	if err != nil {
		return err
	}

	filterTypes := map[string]bool{}
	for _, f := range dnsOptions.Filter {
		filterTypes[f] = true
	}

	recordSets, err := g.GetRecordsWithFilter(name, aws.StringValue(zone.Id), filterTypes)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(recordSets))
}

func (g *Getter) GetRecords(name, zoneId string) ([]*route53.ResourceRecordSet, error) {
	return g.GetRecordsWithFilter(name, zoneId, map[string]bool{})
}

func (g *Getter) GetRecordsWithFilter(name, zoneId string, filterTypes map[string]bool) ([]*route53.ResourceRecordSet, error) {
	recordSets, err := aws.Route53ListResourceRecordSets(zoneId)
	if err != nil {
		return nil, err
	}

	if name != "" {
		n := strings.ToLower(name) + "."
		filtered := []*route53.ResourceRecordSet{}
		for _, rs := range recordSets {
			if n == strings.ToLower(aws.StringValue(rs.Name)) {
				filtered = append(filtered, rs)
			}
		}
		recordSets = filtered
	}

	if len(filterTypes) > 0 {
		filtered := make([]*route53.ResourceRecordSet, 0, len(recordSets))
		for _, rs := range recordSets {
			if filterTypes[aws.StringValue(rs.Type)] {
				filtered = append(filtered, rs)
			}
		}
		recordSets = filtered
	}

	return recordSets, nil
}
