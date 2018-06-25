package lib

import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"fmt"
	"os"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func GetInstanceListForCluster(awsSess *session.Session, clusterName string) []*ecs.ContainerInstance {
	svc := ecs.New(awsSess)
	listResult, err := svc.ListContainerInstances(&ecs.ListContainerInstancesInput{
		Cluster: aws.String(clusterName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeClientException:
				fmt.Println(ecs.ErrCodeClientException, aerr.Error())
			case ecs.ErrCodeInvalidParameterException:
				fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
			case ecs.ErrCodeClusterNotFoundException:
				fmt.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}

	descResult, err := svc.DescribeContainerInstances(&ecs.DescribeContainerInstancesInput{
		Cluster:            aws.String(clusterName),
		ContainerInstances: listResult.ContainerInstanceArns,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeClientException:
				fmt.Println(ecs.ErrCodeClientException, aerr.Error())
			case ecs.ErrCodeInvalidParameterException:
				fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
			case ecs.ErrCodeClusterNotFoundException:
				fmt.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}

	return descResult.ContainerInstances
}

func GetInstanceIDsForCluster(awsSess *session.Session, clusterName string) []*string {
	instances := GetInstanceListForCluster(awsSess, clusterName)
	instanceIDs := []*string{}

	for _, instance := range instances {
		instanceIDs = append(instanceIDs, instance.Ec2InstanceId)
	}

	return instanceIDs
}

func GetInstanceIPsForCluster(awsSess *session.Session, clusterName string) []string {
	instanceIDs := GetInstanceIDsForCluster(awsSess, clusterName)

	svc := ec2.New(awsSess)
	instanceDetails, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: instanceIDs,
	})
	if err != nil {
		fmt.Println("Unable to get instance details", err)
		os.Exit(1)
	}

	var instanceIPs []string

	for _, r := range instanceDetails.Reservations {
		for _, i := range r.Instances {
			instanceIPs = append(instanceIPs, *i.PrivateIpAddress)
		}
	}


	return instanceIPs
}

func GetPendingTasksCount(awsSess *session.Session, cluster string) int64 {
	svc := ecs.New(awsSess)

	services, err := svc.ListServices(&ecs.ListServicesInput{
		Cluster: aws.String(cluster),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeClientException:
				fmt.Println(ecs.ErrCodeClientException, aerr.Error())
			case ecs.ErrCodeInvalidParameterException:
				fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
			case ecs.ErrCodeClusterNotFoundException:
				fmt.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}

	descResult, err := svc.DescribeServices(&ecs.DescribeServicesInput{
		Cluster:  aws.String(cluster),
		Services: services.ServiceArns,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeClientException:
				fmt.Println(ecs.ErrCodeClientException, aerr.Error())
			case ecs.ErrCodeInvalidParameterException:
				fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
			case ecs.ErrCodeClusterNotFoundException:
				fmt.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}

	var pendingTasks int64

	for _, service := range descResult.Services {
		pendingTasks += *service.PendingCount
	}

	return pendingTasks
}