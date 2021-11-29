package addon

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/service/eks"
)

type AddonVersionPrinter struct {
	addonInfos []*eks.AddonInfo
}

func NewVersionPrinter(addonInfos []*eks.AddonInfo) *AddonVersionPrinter {
	return &AddonVersionPrinter{addonInfos}
}

func (p *AddonVersionPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Name", "Version", "Restrictions"})

	for _, addonInfo := range p.addonInfos {
		name := aws.StringValue(addonInfo.AddonName)

		for _, av := range addonInfo.AddonVersions {
			isDefault := ""
			restrictions := "-"

			if len(av.Compatibilities) > 0 {
				if *av.Compatibilities[0].DefaultVersion {
					isDefault = "*"
				}

				if len(av.Compatibilities[0].PlatformVersions) > 0 && aws.StringValue(av.Compatibilities[0].PlatformVersions[0]) != "*" {
					restrictions = strings.Join(aws.StringValueSlice(av.Compatibilities[0].PlatformVersions), ",")
				}
			}

			table.AppendRow([]string{
				name,
				aws.StringValue(av.AddonVersion) + isDefault,
				restrictions,
			})
		}

	}

	table.Print(writer)
	if len(p.addonInfos) > 0 {
		fmt.Println("* Indicates default version")
	}

	return nil
}

func (p *AddonVersionPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.addonInfos)
}

func (p *AddonVersionPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.addonInfos)
}
