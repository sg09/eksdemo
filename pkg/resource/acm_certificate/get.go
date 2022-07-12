package acm_certificate

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/acm"
)

type Getter struct{}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	var certs []*acm.CertificateDetail
	var err error

	if name != "" {
		certs, err = g.GetAllCertsStartingWithName(name)
	} else {
		certs, err = g.GetAllCerts()
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(certs))
}

func (g *Getter) GetAllCerts() ([]*acm.CertificateDetail, error) {
	certSummaries, err := aws.AcmListCertificates()
	certs := make([]*acm.CertificateDetail, 0, len(certSummaries))

	if err != nil {
		return nil, err
	}

	for _, summary := range certSummaries {
		cert, err := aws.AcmDescribeCertificate(aws.StringValue(summary.CertificateArn))
		if err != nil {
			return nil, err
		}
		certs = append(certs, cert)
	}

	return certs, nil
}

func (g *Getter) GetCert(arn string) (*acm.CertificateDetail, error) {
	return aws.AcmDescribeCertificate(arn)
}

func (g *Getter) GetOneCertStartingWithName(name string) (*acm.CertificateDetail, error) {
	certs, err := g.GetAllCertsStartingWithName(name)
	if err != nil {
		return nil, err
	}

	if len(certs) == 0 {
		return nil, resource.NotFoundError(fmt.Sprintf("certificate name %q not found", name))
	}

	if len(certs) > 1 {
		return nil, fmt.Errorf("multiple certificates found starting with: %s", name)
	}

	return certs[0], nil
}

func (g *Getter) GetAllCertsStartingWithName(name string) ([]*acm.CertificateDetail, error) {
	certSummaries, err := aws.AcmListCertificates()
	if err != nil {
		return nil, err
	}

	n := strings.ToLower(name)
	certs := []*acm.CertificateDetail{}

	for _, summary := range certSummaries {
		if strings.HasPrefix(strings.ToLower(aws.StringValue(summary.DomainName)), n) {
			cert, err := aws.AcmDescribeCertificate(aws.StringValue(summary.CertificateArn))
			if err != nil {
				return nil, err
			}
			certs = append(certs, cert)
		}
	}

	return certs, nil
}
