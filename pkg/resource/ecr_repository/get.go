package ecr_repository

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"os"
)

type Getter struct{}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	repos, err := aws.ECRDescribeRepositories(name)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(repos))
}
