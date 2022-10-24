package auto_scaling_group

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"os"
)

type Getter struct {
	autoscalingClient *aws.AutoscalingClient
}

func NewGetter(autoscalingClient *aws.AutoscalingClient) *Getter {
	return &Getter{autoscalingClient}
}

func (g *Getter) Init() {
	if g.autoscalingClient == nil {
		g.autoscalingClient = aws.NewAutoscalingClient()
	}
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	autoScalingGroups, err := g.autoscalingClient.DescribeAutoScalingGroups(name)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(autoScalingGroups))
}
