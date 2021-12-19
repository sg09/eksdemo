package certificate

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/resource"
	"fmt"
	"strings"
)

type Manager struct {
	Getter
}

func (m *Manager) Create(options resource.Options) error {
	name := options.Common().Name

	cert, err := m.Getter.GetOneCertStartingWithName(name)
	if err != nil {
		if _, ok := err.(resource.NotFoundError); !ok {
			// Return an error if it's anything other than resource not found
			return err
		}
	}

	if cert != nil && strings.EqualFold(aws.StringValue(cert.DomainName), name) {
		fmt.Printf("Certificate %q already exists\n", name)
		return nil
	}

	fmt.Printf("Creating ACM Certificate request for: %s...", name)
	arn, err := aws.AcmRequestCertificate(name)
	if err != nil {
		return err
	}
	fmt.Printf("done\nCreated ACM Certificate Id: %s\n", arn[strings.LastIndex(arn, "/")+1:])

	return nil
}

func (m *Manager) Delete(options resource.Options) error {
	name := options.Common().Name

	cert, err := m.Getter.GetOneCertStartingWithName(name)
	if err != nil {
		return err
	}

	err = aws.AcmDeleteCertificate(aws.StringValue(cert.CertificateArn))
	if err != nil {
		return err
	}
	fmt.Printf("ACM Certificate Domain name %q deleted\n", aws.StringValue(cert.DomainName))

	return nil
}

func (m *Manager) SetDryRun() {}
