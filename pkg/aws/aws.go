package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/eks"
)

func FormatError(err error) error {
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case eks.ErrCodeResourceNotFoundException:
				return fmt.Errorf(awsErr.Message())
			default:
				return err
			}
		}
	}
	return err
}

// StringValue returns the value of the string pointer passed in or
// "" if the pointer is nil.
func StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}
