package hosted_zone

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/route53"
)

type Getter struct{}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	zones, err := g.GetAllZonesStartingWithName(name)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(zones))
}

func (g *Getter) GetAllZones() ([]*route53.HostedZone, error) {
	return aws.Route53ListHostedZones()
}

func (g *Getter) GetAllZonesStartingWithName(name string) ([]*route53.HostedZone, error) {
	zones, err := g.GetAllZones()
	if err != nil {
		return nil, err
	}

	if name != "" {
		n := strings.ToLower(name)
		filtered := []*route53.HostedZone{}
		for _, z := range zones {
			if strings.HasPrefix(strings.ToLower(aws.StringValue(z.Name)), n) {
				filtered = append(filtered, z)
			}
		}
		zones = filtered
	}

	return zones, nil
}
