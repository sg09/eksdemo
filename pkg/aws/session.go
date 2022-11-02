package aws

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var awsConfig *aws.Config
var profile string
var region string

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

func Region() string {
	if region == "" {
		GetConfig()
	}
	return region
}
