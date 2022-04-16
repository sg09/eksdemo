package certificate

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/hosted_zone"
	"fmt"
	"strings"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/spf13/cobra"
)

type Manager struct {
	Getter

	arn        string
	zoneGetter hosted_zone.Getter
}

func (m *Manager) Create(options resource.Options) error {
	certOptions, ok := options.(*CertificateOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to CertificateOptions")
	}

	name := certOptions.Name
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
	m.arn, err = aws.AcmRequestCertificate(name, certOptions.sans)
	if err != nil {
		return err
	}
	fmt.Printf("done\nCreated ACM Certificate Id: %s\n", m.arn[strings.LastIndex(m.arn, "/")+1:])

	if certOptions.skipValidation {
		return nil
	}

	return m.validate()
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

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}

func (m *Manager) validate() error {
	cert, err := m.Getter.GetCert(m.arn)
	if err != nil {
		return fmt.Errorf("failed during valication to describe the cert: %w", err)
	}

	zones, err := m.zoneGetter.GetAllZones()
	if err != nil {
		return fmt.Errorf("failed during validation to list hosted zones: %w", err)
	}

	for _, z := range zones {
		changes := []*route53.Change{}
		zoneName := strings.TrimSuffix(aws.StringValue(z.Name), ".")

		for _, dv := range cert.DomainValidationOptions {
			if strings.HasSuffix(aws.StringValue(dv.DomainName), zoneName) {
				fmt.Printf("Validating domain %q using hosted zone %q\n", aws.StringValue(dv.DomainName), zoneName)
				rr := dv.ResourceRecord
				changes = append(changes, createChange(rr.Name, rr.Value, rr.Type, z.Id))
			}
		}

		if len(changes) == 0 {
			continue
		}

		changeBatch := &route53.ChangeBatch{
			Changes: changes,
			Comment: awssdk.String("certificate validation"),
		}

		if err := aws.Route53ChangeResourceRecordSets(changeBatch, aws.StringValue(z.Id)); err != nil {
			return err
		}
	}

	fmt.Printf("Waiting for certificate to be issued...")
	if err = aws.AcmWaitUntilCertificateValidated(m.arn); err != nil {
		return err
	}
	fmt.Println("done")

	return nil
}

func createChange(name, value, recType, zoneId *string) *route53.Change {
	return &route53.Change{
		Action: awssdk.String("UPSERT"),
		ResourceRecordSet: &route53.ResourceRecordSet{
			Name: name,
			ResourceRecords: []*route53.ResourceRecord{
				{
					Value: value,
				},
			},
			TTL:  awssdk.Int64(300),
			Type: recType,
		},
	}
}
