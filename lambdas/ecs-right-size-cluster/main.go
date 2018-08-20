package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
	"github.com/silinternational/awsops/lib"
	"strings"
)

type EcsRightSizeClusterConfig struct {
	ClusterNamesCSV string
	Region          string
}

func main() {
	lambda.Start(handler)
}

func handler(config EcsRightSizeClusterConfig) error {
	if config.Region == "" {
		config.Region = "us-east-1"
	}

	clusters := strings.Split(config.ClusterNamesCSV, ",")
	if len(clusters) == 0 {
		err := errors.New("error: EcsRightSizeConfig.ClusterNamesCSV is empty")
		fmt.Println(err.Error())
		return err
	}

	AwsSess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(config.Region),
	}))

	for _, cluster := range clusters {
		err := lib.RightSizeAsgForEcsCluster(AwsSess, cluster)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}

	return nil
}
