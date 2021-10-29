package resource

import "eksdemo/pkg/printer"

type Getter interface {
	Get(string, printer.Output, Options) error
}
