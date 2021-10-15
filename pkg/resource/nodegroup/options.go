package nodegroup

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
	"fmt"
)

type NodegroupOptions struct {
	resource.CommonOptions

	AMI             string
	InstanceType    string
	Containerd      bool
	DesiredCapacity int
	MinSize         int
	MaxSize         int
	NodegroupName   string
	Spot            bool
	SpotvCPUs       int
	SpotMemory      int
}

func NewOptions() (options *NodegroupOptions, flags cmd.Flags) {
	options = &NodegroupOptions{
		InstanceType:    "t3.large",
		DesiredCapacity: 1,
		MinSize:         1,
		MaxSize:         5,
		SpotvCPUs:       2,
		SpotMemory:      4,
	}

	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "containerd",
				Description: "use containerd runtime",
			},
			Option: &options.Containerd,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "instance",
				Description: "instance type",
				Shorthand:   "i",
			},
			Option: &options.InstanceType,
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "max",
				Description: "max nodes",
			},
			Option: &options.MaxSize,
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "min",
				Description: "min nodes",
				Validate: func() error {
					if options.MinSize >= options.MaxSize {
						return fmt.Errorf("min nodes must be less than max nodes")
					}
					return nil
				},
			},
			Option: &options.MinSize,
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "nodes",
				Description: "initial nodes",
				Shorthand:   "N",
				Validate: func() error {
					if options.DesiredCapacity > options.MaxSize {
						options.MaxSize = options.DesiredCapacity
					}
					if options.DesiredCapacity < options.MinSize {
						options.MinSize = options.DesiredCapacity
					}
					return nil
				},
			},
			Option: &options.DesiredCapacity,
		},
	}
	return
}

func (o *NodegroupOptions) PreCreate() error {
	if !o.Containerd {
		return nil
	}

	ami, err := aws.EksOptimizedAmi(o.KubernetesVersion)
	if err != nil {
		return err
	}

	o.AMI = ami

	return nil
}

func (o *NodegroupOptions) SetName(name string) {
	o.NodegroupName = name
}
