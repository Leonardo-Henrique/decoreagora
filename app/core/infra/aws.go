package infra

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
)

func NewAWSConfig() *aws.Config {
	cfg, err := awscfg.LoadDefaultConfig(context.TODO(), awscfg.WithRegion("us-east-1"))
	if err != nil {
		log.Fatal(err)
	}
	return &cfg
}
