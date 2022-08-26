package listener

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/load_balancer"
	"fmt"
	"os"
)

type Getter struct {
	elbGetter load_balancer.Getter
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) (err error) {
	listenerOptions, ok := options.(*ListenerOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to ListenerOptions")
	}

	elbs, err := g.elbGetter.GetLoadBalancers(listenerOptions.LoadBalancerName)
	if err != nil {
		return err
	}
	if len(elbs.V1) > 0 {
		return fmt.Errorf("%q is a classic load balancer", listenerOptions.LoadBalancerName)
	}

	lbArn := aws.StringValue(elbs.V2[0].LoadBalancerArn)

	listeners, err := aws.ELBDescribeListeners(lbArn)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(listeners))
}
