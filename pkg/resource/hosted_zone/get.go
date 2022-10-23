package hosted_zone

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/route53"
)

type Getter struct {
	resource.EmptyInit
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	var err error
	var zone *route53.HostedZone
	var zones []*route53.HostedZone

	if name != "" {
		zone, err = g.GetZoneByName(name)
		zones = []*route53.HostedZone{zone}
	} else {
		zones, err = g.GetAllZones()
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(zones))
}

func (g *Getter) GetAllZones() ([]*route53.HostedZone, error) {
	return aws.Route53ListHostedZones()
}

func (g *Getter) GetZoneByName(name string) (*route53.HostedZone, error) {
	zone, err := aws.Route53ListHostedZonesByName(name)
	if err != nil {
		return nil, err
	}

	if len(zone) == 0 {
		return nil, fmt.Errorf("hosted-zone %q not found", name)
	}

	z := zone[0]

	if strings.ToLower(aws.StringValue(z.Name)) != name+"." {
		return nil, fmt.Errorf("hosted-zone %q not found", name)
	}

	return z, nil
}
