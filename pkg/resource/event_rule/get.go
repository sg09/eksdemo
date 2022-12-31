package event_rule

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"os"
)

type Getter struct {
	eventBridgeClient *aws.EventBridgeClient
}

func NewGetter(eventBridgeClient *aws.EventBridgeClient) *Getter {
	return &Getter{eventBridgeClient}
}

func (g *Getter) Init() {
	if g.eventBridgeClient == nil {
		g.eventBridgeClient = aws.NewEventBridgeClient()
	}
}

func (g *Getter) Get(namePrefix string, output printer.Output, options resource.Options) error {
	rules, err := g.eventBridgeClient.ListRules(namePrefix)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(rules))
}
