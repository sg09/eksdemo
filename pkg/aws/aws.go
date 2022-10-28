package aws

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/smithy-go"
)

const maxPages = 3

func FormatErrorSDKv1(err error) error {
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

// Return cleaner error message for service API errors
func FormatError(err error) error {
	var ae smithy.APIError
	if err != nil && errors.As(err, &ae) {
		return ae
	}
	return err
}

// Return cleaner error message for service API errors
func FormatErrorAsMessageOnly(err error) error {
	var ae smithy.APIError
	if err != nil && errors.As(err, &ae) {
		return fmt.Errorf(ae.ErrorMessage())
	}
	return err
}

// BoolValue returns the value of the bool pointer passed in or
// false if the pointer is nil.
func BoolValue(v *bool) bool {
	if v != nil {
		return *v
	}
	return false
}

// Int64Value returns the value of the int64 pointer passed in or
// 0 if the pointer is nil.
func Int64Value(v *int64) int64 {
	if v != nil {
		return *v
	}
	return 0
}

// StringValue returns the value of the string pointer passed in or
// "" if the pointer is nil.
func StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

// StringValueSlice converts a slice of string pointers into a slice of
// string values
func StringValueSlice(src []*string) []string {
	dst := make([]string, len(src))
	for i := 0; i < len(src); i++ {
		if src[i] != nil {
			dst[i] = *(src[i])
		}
	}
	return dst
}

// TimeValue returns the value of the time.Time pointer passed in or
// time.Time{} if the pointer is nil.
func TimeValue(v *time.Time) time.Time {
	if v != nil {
		return *v
	}
	return time.Time{}
}
