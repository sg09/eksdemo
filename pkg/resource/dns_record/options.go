package dns_record

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
	"fmt"

	"github.com/spf13/cobra"
)

type DnsRecordOptions struct {
	resource.CommonOptions

	AllRecords bool
	AllTypes   bool
	Filter     []string
	ZoneName   string
}

func newOptions() (options *DnsRecordOptions, deleteFlags cmd.Flags, getFlags cmd.Flags) {
	options = &DnsRecordOptions{
		CommonOptions: resource.CommonOptions{
			ArgumentOptional:    true,
			ClusterFlagDisabled: true,
		},
	}

	commonFlags := cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "zone",
				Description: "hosted zone name",
				Shorthand:   "z",
				Required:    true,
			},
			Option: &options.ZoneName,
		},
	}

	deleteFlags = append(cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "all",
				Description: "delete all records types for the record name",
				Shorthand:   "A",
			},
			Option: &options.AllTypes,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "all-records",
				Description: "delete all records in the zone (excluding root records)",
				Validate: func(cmd *cobra.Command, args []string) error {
					if !options.AllRecords && len(args) == 0 {
						return fmt.Errorf("must include either %q argument or %q flag", "NAME", "--all-records")
					}

					if options.AllRecords && len(args) > 0 {
						return fmt.Errorf("%q flag cannot be used with a record name", "--all-records")
					}
					return nil
				},
			},
			Option: &options.AllRecords,
		},
	}, commonFlags...)

	getFlags = append(cmd.Flags{
		&cmd.StringSliceFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "filter",
				Description: "filter records by types",
			},
			Option: &options.Filter,
		},
	}, commonFlags...)

	return
}
