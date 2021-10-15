package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
)

var awsAccount string

func AccountId() string {
	if awsAccount != "" {
		return awsAccount
	}

	sess := GetSession()
	svc := sts.New(sess)
	input := &sts.GetCallerIdentityInput{}

	result, err := svc.GetCallerIdentity(input)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get AWS identity: %s", err))
	}
	awsAccount = aws.StringValue(result.Account)

	return awsAccount
}

func Region() string {
	return aws.StringValue(GetSession().Config.Region)
}
