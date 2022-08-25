package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var profile string
var region string
var sess *session.Session

func Init(awsProfile, awsRegion string) {
	profile = awsProfile
	region = awsRegion
}

func GetSession() *session.Session {
	if sess != nil {
		return sess
	}
	var err error

	sess, err = session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			CredentialsChainVerboseErrors: aws.Bool(true),
			Region:                        aws.String(region),
		},
		Profile:           profile,
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create AWS session: %s", err))
	}
	region = aws.StringValue(sess.Config.Region)

	return sess
}

func Region() string {
	if region == "" {
		GetSession()
	}
	return region
}
