package listener_rule

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/load_balancer"
	"fmt"
	"os"
	"strings"
)

type Getter struct {
	elbGetter load_balancer.Getter
}

func (g *Getter) Get(id string, output printer.Output, options resource.Options) (err error) {
	lrOptions, ok := options.(*ListenerRuleOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to ListenerRuleOptions")
	}

	elbs, err := g.elbGetter.GetLoadBalancers(lrOptions.LoadBalancerName)
	if err != nil {
		return err
	}
	if len(elbs.V1) > 0 {
		return fmt.Errorf("%q is a classic load balancer", lrOptions.LoadBalancerName)
	}

	lbArn := aws.StringValue(elbs.V2[0].LoadBalancerArn)
	listernArn := ""
	ruleArn := ""

	if id == "" {
		listernArn = strings.Replace(lbArn, ":loadbalancer/", ":listener/", 1) + "/" + lrOptions.ListenerId
	} else {
		ruleArn = strings.Replace(lbArn, ":loadbalancer/", ":listener-rule/", 1) + "/" + lrOptions.ListenerId + "/" + id
	}

	rules, err := aws.ELBDescribeRules(listernArn, []string{ruleArn})
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(rules))
}
