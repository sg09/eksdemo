package availability_zone

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"fmt"
	"os"
)

type Getter struct{}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	azOptions, ok := options.(*AvailabilityZoneOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to AvailabilityZoneOptions")
	}

	// If getting by name, lookup all zones
	allZones := azOptions.AllZones || name != ""

	zones, err := aws.EC2DescribeAvailabilityZones(name, allZones)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(zones))
}
