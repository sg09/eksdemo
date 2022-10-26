package availability_zone

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"fmt"
	"os"
)

type Getter struct {
	ec2Client *aws.EC2Client
}

func NewGetter(ec2Client *aws.EC2Client) *Getter {
	return &Getter{ec2Client}
}

func (g *Getter) Init() {
	if g.ec2Client == nil {
		g.ec2Client = aws.NewEC2Client()
	}
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	azOptions, ok := options.(*AvailabilityZoneOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to AvailabilityZoneOptions")
	}

	// If getting by name, lookup all zones
	allZones := azOptions.AllZones || name != ""

	zones, err := g.ec2Client.DescribeAvailabilityZones(name, allZones)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(zones))
}
