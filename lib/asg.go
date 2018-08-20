package lib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"math"
	"os"
	"time"
)

func GetAsgNameForEcsCluster(awsSess *session.Session, cluster string) string {
	instanceIDs := GetInstanceIDsForEcsCluster(awsSess, cluster)

	svc := ec2.New(awsSess)
	instanceDetails, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: instanceIDs,
	})
	if err != nil {
		fmt.Println("Unable to get asg name from instance: ", err)
		os.Exit(1)
	}

	for _, tag := range instanceDetails.Reservations[0].Instances[0].Tags {
		if *tag.Key == "aws:autoscaling:groupName" {
			return *tag.Value
		}
	}

	return ""
}

func DetachAndReplaceAsgInstances(awsSess *session.Session, asgName string, instancesToTerminate []*string) {
	svc := autoscaling.New(awsSess)

	decrement := false

	fmt.Printf("Detaching %v instances...", len(instancesToTerminate))
	_, err := svc.DetachInstances(&autoscaling.DetachInstancesInput{
		AutoScalingGroupName:           &asgName,
		InstanceIds:                    instancesToTerminate,
		ShouldDecrementDesiredCapacity: &decrement,
	})
	if err != nil {
		fmt.Println("Unable to detach instances: ", err)
		os.Exit(1)
	}

	fmt.Printf("done\n")

	for ready := false; ready != true; {
		time.Sleep(15 * time.Second)
		instances := GetInstanceListForAsg(awsSess, asgName)
		fmt.Printf("\rNew instances created: %v", len(instances))
		if len(instances) == len(instancesToTerminate) {
			ready = true
			fmt.Println()
			fmt.Println("Finished creating new instances")
		}
	}
}

func GetInstanceListForAsg(awsSess *session.Session, asgName string) []*string {
	asg := GetAsg(awsSess, asgName)

	var instanceIds []*string
	for _, ins := range asg.Instances {
		instanceIds = append(instanceIds, ins.InstanceId)
	}

	return instanceIds
}

func GetInstanceTypeForAsg(awsSess *session.Session, asgName string) string {
	svc := autoscaling.New(awsSess)

	asg := GetAsg(awsSess, asgName)

	input := &autoscaling.DescribeLaunchConfigurationsInput{
		LaunchConfigurationNames: []*string{
			aws.String(*asg.LaunchConfigurationName),
		},
	}

	lc, err := svc.DescribeLaunchConfigurations(input)
	if err != nil {
		fmt.Println("Unable to describe launch configuration, err: ", err.Error())
	}

	if len(lc.LaunchConfigurations) != 1 {
		fmt.Println("DescribeLaunchConfigurations did not return expected number of results. Expected: 1, Actual: ", len(lc.LaunchConfigurations))
		os.Exit(1)
	}

	return *lc.LaunchConfigurations[0].InstanceType
}

func HowManyServersNeededForAsg(serverType string, memory, cpu int64) int64 {
	instanceSpecs, valid := InstanceTypes[serverType]
	if !valid {
		fmt.Println("Invalid server type provided: ", serverType)
		os.Exit(1)
	}

	neededForMem := math.Ceil(float64(memory) / float64(instanceSpecs.MemoryMb))
	neededForCPU := math.Ceil(float64(cpu) / float64(instanceSpecs.CPUUnits))

	if neededForMem > neededForCPU {
		return int64(neededForMem)
	}

	return int64(neededForCPU)
}

func GetAsgServerCount(awsSess *session.Session, asgName string) (desired int64, min int64, max int64) {
	asg := GetAsg(awsSess, asgName)

	return *asg.DesiredCapacity, *asg.MinSize, *asg.MaxSize
}

func GetAsg(awsSess *session.Session, asgName string) *autoscaling.Group {
	svc := autoscaling.New(awsSess)

	groups, err := svc.DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{&asgName},
	})
	if err != nil {
		fmt.Println("Unable to get list of ASG groups: ", err)
		os.Exit(1)
	}

	if len(groups.AutoScalingGroups) != 1 {
		fmt.Println("DescribeAutoScalingGroups did not return expected number of results. Expected: 1, Actual: ", len(groups.AutoScalingGroups))
		os.Exit(1)
	}

	return groups.AutoScalingGroups[0]
}

func UpdateAsgServerCount(awsSess *session.Session, asgName string, serverCount int64) error {
	svc := autoscaling.New(awsSess)
	input := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(asgName),
		MaxSize:              aws.Int64(serverCount),
		MinSize:              aws.Int64(serverCount),
		DesiredCapacity:      aws.Int64(serverCount),
	}

	_, err := svc.UpdateAutoScalingGroup(input)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
