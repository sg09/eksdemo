package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

func S3CreateBucket(name, region string) error {
	sess := GetSession()
	svc := s3.New(sess)

	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(name),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(region),
		},
	})

	return err
}

func S3GetBucketLocation(name string) (exists bool, err error) {
	sess := GetSession()
	svc := s3.New(sess)

	_, err = svc.GetBucketLocation(&s3.GetBucketLocationInput{
		Bucket: aws.String(name),
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return false, nil
			case "AccessDenied":
				return false, fmt.Errorf(awsErr.Message())
			default:
				return false, awsErr
			}
		}
		return false, err
	}
	return true, nil
}

func S3ListBuckets() ([]*s3.Bucket, error) {
	sess := GetSession()
	svc := s3.New(sess)

	result, err := svc.ListBuckets(&s3.ListBucketsInput{})

	if err != nil {
		return nil, err
	}

	return result.Buckets, nil
}
