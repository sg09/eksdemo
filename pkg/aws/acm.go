package aws

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/acm"
)

func AcmDeleteCertificate(arn string) error {
	sess := GetSession()
	svc := acm.New(sess)

	_, err := svc.DeleteCertificate(&acm.DeleteCertificateInput{
		CertificateArn: aws.String(arn),
	})

	return err
}

func AcmDescribeCertificate(arn string) (*acm.CertificateDetail, error) {
	sess := GetSession()
	svc := acm.New(sess)

	cert, err := svc.DescribeCertificate(&acm.DescribeCertificateInput{
		CertificateArn: aws.String(arn),
	})

	if err != nil {
		return nil, err
	}

	return cert.Certificate, nil
}

func AcmListCertificates() ([]*acm.CertificateSummary, error) {
	sess := GetSession()
	svc := acm.New(sess)

	certs := []*acm.CertificateSummary{}
	pageNum := 0

	err := svc.ListCertificatesPages(&acm.ListCertificatesInput{},
		func(page *acm.ListCertificatesOutput, lastPage bool) bool {
			pageNum++
			certs = append(certs, page.CertificateSummaryList...)
			return pageNum <= maxPages
		},
	)

	if err != nil {
		return nil, err
	}

	return certs, nil
}

func AcmRequestCertificate(fqdn string) (string, error) {
	sess := GetSession()
	svc := acm.New(sess)

	result, err := svc.RequestCertificate(&acm.RequestCertificateInput{
		DomainName:       aws.String(fqdn),
		ValidationMethod: aws.String("DNS"),
	})

	if err != nil {
		return "", err
	}

	err = waitUntilCertificateValidationMetadata(svc, &acm.DescribeCertificateInput{
		CertificateArn: result.CertificateArn,
	})

	return aws.StringValue(result.CertificateArn), err
}

func waitUntilCertificateValidationMetadata(svc *acm.ACM, input *acm.DescribeCertificateInput, opts ...request.WaiterOption) error {
	ctx := aws.BackgroundContext()

	w := request.Waiter{
		Name:        "WaitUntilCertificateValidationMetadata",
		MaxAttempts: 40,
		Delay:       request.ConstantWaiterDelay(2 * time.Second),
		Acceptors: []request.WaiterAcceptor{
			{
				State:   request.SuccessWaiterState,
				Matcher: request.PathAllWaiterMatch, Argument: "Certificate.DomainValidationOptions[].ResourceRecord.Type",
				Expected: "CNAME",
			},
			{
				State:   request.RetryWaiterState,
				Matcher: request.PathAnyWaiterMatch, Argument: "Certificate.DomainValidationOptions[].ResourceRecord",
				Expected: nil,
			},
			{
				State:   request.FailureWaiterState,
				Matcher: request.PathWaiterMatch, Argument: "Certificate.Status",
				Expected: "FAILED",
			},
			{
				State:    request.FailureWaiterState,
				Matcher:  request.ErrorWaiterMatch,
				Expected: "ResourceNotFoundException",
			},
		},
		Logger: svc.Config.Logger,
		NewRequest: func(opts []request.Option) (*request.Request, error) {
			var inCpy *acm.DescribeCertificateInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req, _ := svc.DescribeCertificateRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}
