package s3_bucket

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/resource"
	"fmt"

	"github.com/spf13/cobra"
)

type Manager struct {
	DryRun bool
}

func (m *Manager) Create(options resource.Options) error {
	bucketOptions, ok := options.(*BucketOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to BucketOptions")
	}

	if m.DryRun {
		return m.dryRun(bucketOptions)
	}

	exists, err := aws.S3GetBucketLocation(bucketOptions.BucketName)
	if err != nil {
		return fmt.Errorf("failed checking if bucket %q exists: %w", bucketOptions.BucketName, err)
	}

	if exists {
		fmt.Printf("Bucket %q already exists\n", bucketOptions.BucketName)
		return nil
	}

	fmt.Printf("Creating Bucket: %s...", bucketOptions.BucketName)

	err = aws.S3CreateBucket(bucketOptions.BucketName, options.Common().Region)
	if err != nil {
		return err
	}
	fmt.Println("done")

	return nil
}

func (m *Manager) Delete(options resource.Options) error {
	bucketOptions, ok := options.(*BucketOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to BucketOptions")
	}

	fmt.Printf("Deletion of Bucket %q not supported. Please delete manually.\n", bucketOptions.BucketName)

	return nil
}

func (m *Manager) SetDryRun() {
	m.DryRun = true
}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}

func (m *Manager) dryRun(options *BucketOptions) error {
	fmt.Printf("\nS3 Bucket Manager Dry Run:\n")
	fmt.Printf("S3 API Call %q with request parameters:\n", "CreateBucket")
	fmt.Printf("Bucket: %q\n", options.BucketName)
	fmt.Printf("CreateBucketConfiguration.LocationConstraint: %q\n", options.Region)
	return nil
}
