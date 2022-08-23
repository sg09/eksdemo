package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/sts"
)

var accountId, partition string

func AccountId() string {
	if accountId == "" {
		getCallerIdentity()
	}
	return accountId
}

func Partition() string {
	if partition == "" {
		getCallerIdentity()
	}
	return partition
}

func getCallerIdentity() {
	sess := GetSession()
	svc := sts.New(sess)

	result, err := svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get AWS identity: %w", err))
	}

	arn, err := arn.Parse(aws.StringValue(result.Arn))
	if err != nil {
		log.Fatal(fmt.Errorf("failed to parse ARN looking up AWS identity: %w", err))
	}

	accountId = arn.AccountID
	partition = arn.Partition
}
