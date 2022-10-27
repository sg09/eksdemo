package ecr_repository

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"os"
)

type Getter struct {
	ecrClient *aws.ECRClient
}

func NewGetter(ecrClient *aws.ECRClient) *Getter {
	return &Getter{ecrClient}
}

func (g *Getter) Init() {
	if g.ecrClient == nil {
		g.ecrClient = aws.NewECRClient()
	}
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	repos, err := g.ecrClient.DescribeRepositories(name)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(repos))
}
