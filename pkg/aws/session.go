package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
)

var sess *session.Session

func GetSession() *session.Session {
	if sess != nil {
		return sess
	}

	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create AWS session: %s", err))
	}
	return sess
}
