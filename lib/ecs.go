package lib

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type ResourceSizes struct {
	TotalCPU      int64
	TotalMemory   int64
	LargestCPU    int64
	LargestMemory int64
}

func GetInstanceListForEcsCluster(awsSess *session.Session, clusterName string) []*ecs.ContainerInstance {
	svc := ecs.New(awsSess)
	listResult, err := svc.ListContainerInstances(&ecs.ListContainerInstancesInput{
		Cluster: aws.String(clusterName),
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	descResult, err := svc.DescribeContainerInstances(&ecs.DescribeContainerInstancesInput{
		Cluster:            aws.String(clusterName),
		ContainerInstances: listResult.ContainerInstanceArns,
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return descResult.ContainerInstances
}

func GetInstanceIDsForEcsCluster(awsSess *session.Session, clusterName string) []*string {
	instances := GetInstanceListForEcsCluster(awsSess, clusterName)
	instanceIDs := []*string{}

	for _, instance := range instances {
		instanceIDs = append(instanceIDs, instance.Ec2InstanceId)
	}

	return instanceIDs
}

func GetInstanceIPsForEcsCluster(awsSess *session.Session, clusterName string) []string {
	instanceIDs := GetInstanceIDsForEcsCluster(awsSess, clusterName)

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

func GetPendingEcsTasksCount(awsSess *session.Session, cluster string) int64 {
	ecsServices := ListServicesForEcsCluster(awsSess, cluster)

	var pendingTasks int64

	for _, service := range ecsServices {
		pendingTasks += *service.PendingCount
	}

	return pendingTasks
}

func ListServicesForEcsCluster(awsSess *session.Session, cluster string) []*ecs.Service {
	svc := ecs.New(awsSess)

	var allServices []*ecs.Service
	err := svc.ListServicesPages(&ecs.ListServicesInput{
		Cluster: aws.String(cluster),
	}, func(page *ecs.ListServicesOutput, lastPage bool) bool {
		services, err := DescribeEcsServicesForArns(awsSess, page.ServiceArns, cluster)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		for _, service := range services {
			allServices = append(allServices, service)
		}

		return !lastPage
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return allServices
}

func DescribeEcsServicesForArns(awsSess *session.Session, serviceArns []*string, cluster string) ([]*ecs.Service, error) {
	svc := ecs.New(awsSess)

	descResult, err := svc.DescribeServices(&ecs.DescribeServicesInput{
		Cluster:  aws.String(cluster),
		Services: serviceArns,
	})
	if err != nil {
		return []*ecs.Service{}, err
	}

	return descResult.Services, nil
}

func GetMemoryCpuNeededForEcsServices(awsSess *session.Session, ecsServices []*ecs.Service) ResourceSizes {
	var resourcesNeeded ResourceSizes

	svc := ecs.New(awsSess)

	for _, service := range ecsServices {
		if *service.DesiredCount == 0 {
			continue
		}

		// fmt.Printf("Looking at service %s, count = %v\n", *service.ServiceName, *service.DesiredCount)
		taskDef, err := svc.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
			TaskDefinition: service.TaskDefinition,
		})
		if err != nil {
			fmt.Printf("Unable to describe task definition %s\n", *service.TaskDefinition)
			os.Exit(1)
		}

		var serviceMemory int64 = 0
		var serviceCpu int64 = 0

		for _, c := range taskDef.TaskDefinition.ContainerDefinitions {
			// fmt.Printf("    Looking at container %s, needs %v mem and %v cpu\n", *c.Name, *c.Memory, *c.Cpu)
			serviceMemory += *c.Memory
			serviceCpu += *c.Cpu
		}

		resourcesNeeded.LargestMemory = max(serviceMemory, resourcesNeeded.LargestMemory)
		resourcesNeeded.LargestCPU = max(serviceCpu, resourcesNeeded.LargestCPU)
		resourcesNeeded.TotalMemory += serviceMemory * *service.DesiredCount
		resourcesNeeded.TotalCPU += serviceCpu * *service.DesiredCount
	}

	// Add back in the largest service memory and cpu needs to ensure there is enough extra capacity
	// to launch another instance of the largest service for rolling updates
	resourcesNeeded.TotalMemory += resourcesNeeded.LargestMemory
	resourcesNeeded.TotalCPU += resourcesNeeded.LargestCPU

	return resourcesNeeded
}

func RightSizeAsgForEcsCluster(awsSess *session.Session, cluster string, atLeastServiceDesiredCount bool) error {
	asgName := GetAsgNameForEcsCluster(awsSess, cluster)
	if asgName == "" {
		fmt.Println("Unable to find ASG name for ECS cluster ", cluster)
		os.Exit(1)
	}

	fmt.Println("ASG found: ", asgName)

	instanceType := GetInstanceTypeForAsg(awsSess, asgName)
	fmt.Println("ASG uses instance type: ", instanceType)

	ecsServices := ListServicesForEcsCluster(awsSess, cluster)
	resourcesNeeded := GetMemoryCpuNeededForEcsServices(awsSess, ecsServices)
	fmt.Printf("Resources needed for all services with desired count > 0: %+v\n", resourcesNeeded)

	serversNeeded := HowManyServersNeededForAsg(instanceType, resourcesNeeded)
	fmt.Printf("ASG should have %v servers to fit all tasks\n", serversNeeded)

	// If an ECS service has a desired count > serversNeeded, and atLeastServiceDesiredCount is true, set serversNeeded to
	// largest ecs service desired count value
	if atLeastServiceDesiredCount {
		largestDesiredCount := GetLargestDesiredCountFromEcsServices(ecsServices)
		serversNeeded = max(largestDesiredCount, serversNeeded)
	}

	asgDesired, asgMin, asgMax := GetAsgServerCount(awsSess, asgName)
	fmt.Printf("ASG server count currently set to: desired = %v, min = %v, max = %v\n", asgDesired, asgMin, asgMax)

	if asgMin < serversNeeded {
		fmt.Printf("ASG needs to be scaled up by %v servers\n", serversNeeded-asgMin)
		fmt.Printf("Scaling ASG to %v servers...", serversNeeded)
		err := UpdateAsgServerCount(awsSess, asgName, serversNeeded)
		if err != nil {
			return err
		}
		fmt.Printf("done.\n")
	} else if asgMin > serversNeeded {
		fmt.Printf("ASG can be scaled down by %v servers\n", asgMin-serversNeeded)
		fmt.Printf("Scaling ASG to %v servers (desired/min/max)...", serversNeeded)
		err := UpdateAsgServerCount(awsSess, asgName, serversNeeded)
		if err != nil {
			return err
		}
		fmt.Printf("done.\n")
	} else {
		fmt.Printf("Looks like this ASG is already right sized, good day sir.\n")
	}

	return nil
}

func GetLargestDesiredCountFromEcsServices(ecsServices []*ecs.Service) int64 {
	largestDesiredCount := int64(0)

	for _, service := range ecsServices {
		largestDesiredCount = max(*service.DesiredCount, largestDesiredCount)
	}

	return largestDesiredCount
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
