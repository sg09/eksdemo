package security_group_rule

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/network_interface"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type Getter struct {
	eniGetter network_interface.Getter
}

func (g *Getter) Get(id string, output printer.Output, options resource.Options) error {
	sgrOptions, ok := options.(*SecurityGroupRuleOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to SecurityGroupRuleOptions")
	}

	var err error
	var securityGroupRules []*ec2.SecurityGroupRule

	if sgrOptions.SecurityGroupId != "" {
		securityGroupRules, err = g.GetSecurityGroupRulesBySecurityGroupId(sgrOptions.SecurityGroupId)
	} else if sgrOptions.NetworkInterfaceId != "" {
		securityGroupRules, err = g.GetSecurityGroupRulesByNetworkInterfaceId(sgrOptions.NetworkInterfaceId)
	} else {
		securityGroupRules, err = g.GetSecurityGroupRulesById(id)
	}

	if err != nil {
		return err
	}

	if sgrOptions.Ingress {
		securityGroupRules = filterResults(securityGroupRules, false)
	} else if sgrOptions.Egress {
		securityGroupRules = filterResults(securityGroupRules, true)
	}

	return output.Print(os.Stdout, NewPrinter(securityGroupRules, options.Common().ClusterName))
}

func (g *Getter) GetSecurityGroupRulesById(id string) ([]*ec2.SecurityGroupRule, error) {
	return aws.EC2DescribeSecurityGroupRules(id, "")
}

func (g *Getter) GetSecurityGroupRulesByNetworkInterfaceId(eniId string) ([]*ec2.SecurityGroupRule, error) {
	networkInterface, err := g.eniGetter.GetNetworkInterfaceById(eniId)
	if err != nil {
		return nil, err
	}
	eniRules := []*ec2.SecurityGroupRule{}

	for _, groupIdentifier := range networkInterface.Groups {
		sgr, err := aws.EC2DescribeSecurityGroupRules("", aws.StringValue(groupIdentifier.GroupId))
		if err != nil {
			return nil, err
		}
		eniRules = append(eniRules, sgr...)
	}

	return eniRules, nil
}

func (g *Getter) GetSecurityGroupRulesBySecurityGroupId(securityGroupId string) ([]*ec2.SecurityGroupRule, error) {
	return aws.EC2DescribeSecurityGroupRules("", securityGroupId)
}

func filterResults(rules []*ec2.SecurityGroupRule, egress bool) []*ec2.SecurityGroupRule {
	filtered := make([]*ec2.SecurityGroupRule, 0, len(rules))

	if egress {
		for _, rule := range rules {
			if aws.BoolValue(rule.IsEgress) {
				filtered = append(filtered, rule)
			}
		}
	} else {
		for _, rule := range rules {
			if !aws.BoolValue(rule.IsEgress) {
				filtered = append(filtered, rule)
			}
		}
	}

	return filtered
}
