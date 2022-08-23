package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var region string
var sess *session.Session

func GetSession() *session.Session {
	if sess != nil {
		return sess
	}
	var err error

	sess, err = session.NewSessionWithOptions(session.Options{
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
