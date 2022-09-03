package ssm_session

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"fmt"
	"os"
)

type Getter struct{}

func (g *Getter) Get(id string, output printer.Output, options resource.Options) error {
	sessOptions, ok := options.(*SessionOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to SessionOptions")
	}

	state := "Active"
	if sessOptions.History {
		state = "History"
	}

	sessions, err := aws.NewSSMClient().DescribeSessions(id, state)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(sessions))
}
