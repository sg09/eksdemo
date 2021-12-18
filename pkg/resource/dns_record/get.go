package dns_record

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/route53"
)

type Getter struct{}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	dnsOptions, ok := options.(*DnsRecordOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to DnsRecordOptions")
	}

	zone, err := aws.Route53ListHostedZonesByName(dnsOptions.ZoneName)
	if err != nil {
		return err
	}

	if len(zone) == 0 {
		return fmt.Errorf("zone not found")
	}

	z := zone[0]

	if !strings.HasPrefix(strings.ToLower(aws.StringValue(z.Name)), dnsOptions.ZoneName) {
		return fmt.Errorf("zone not found")
	}

	recordSets, err := aws.Route53ListResourceRecordSets(aws.StringValue(z.Id))
	if err != nil {
		return err
	}

	if name != "" {
		n := strings.ToLower(name)
		filtered := []*route53.ResourceRecordSet{}
		for _, z := range recordSets {
			if strings.HasPrefix(strings.ToLower(aws.StringValue(z.Name)), n) {
				filtered = append(filtered, z)
			}
		}
		recordSets = filtered
	}

	return output.Print(os.Stdout, NewPrinter(recordSets))
}
