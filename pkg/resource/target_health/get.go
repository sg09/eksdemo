package target_health

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/target_group"
	"fmt"
	"os"
)

type Getter struct {
	tgGetter target_group.Getter
}

func (g *Getter) Get(id string, output printer.Output, options resource.Options) error {
	thOptions, ok := options.(*TargetHealthOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to TargetHealthOptions")
	}

	targetGroup, err := g.tgGetter.GetTargetGroupByName(thOptions.TargetGroupName)
	if err != nil {
		return err
	}

	targets, err := aws.ELBDescribeTargetHealth(aws.StringValue(targetGroup.TargetGroupArn), id)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(targets))
}
