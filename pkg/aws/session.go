package aws

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsv1 "github.com/aws/aws-sdk-go/aws"
	sessionv1 "github.com/aws/aws-sdk-go/aws/session"
)

var awsConfig *aws.Config
var profile string
var region string
var sess *sessionv1.Session

func Init(awsProfile, awsRegion string) {
	profile = awsProfile
	region = awsRegion
}

func GetConfig() aws.Config {
	if awsConfig != nil {
		return *awsConfig
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithSharedConfigProfile(profile),
		config.WithRegion(region),
	)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create AWS config: %w", err))
	}
	region = cfg.Region
	awsConfig = &cfg

	return cfg
}

func GetSession() *sessionv1.Session {
	if sess != nil {
		return sess
	}
	var err error

	sess, err = sessionv1.NewSessionWithOptions(sessionv1.Options{
		Config: awsv1.Config{
			CredentialsChainVerboseErrors: aws.Bool(true),
			Region:                        aws.String(region),
		},
		Profile:           profile,
		SharedConfigState: sessionv1.SharedConfigEnable,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create AWS session: %w", err))
	}
	region = aws.ToString(sess.Config.Region)

	return sess
}

func Region() string {
	if region == "" {
		GetConfig()
	}
	return region
}
