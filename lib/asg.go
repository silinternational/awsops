package lib

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
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

func GetInstanceTypeFromLaunchConfiguration(awsSess *session.Session, launchConfigurationName string) string {
	input := &autoscaling.DescribeLaunchConfigurationsInput{
		LaunchConfigurationNames: []*string{
			aws.String(launchConfigurationName),
		},
	}

	lc, err := autoscaling.New(awsSess).DescribeLaunchConfigurations(input)
	if err != nil {
		fmt.Println("Unable to describe launch configuration: ", err.Error())
		os.Exit(1)
	}

	if len(lc.LaunchConfigurations) != 1 {
		fmt.Println("Expected one Launch Configuration, received ", len(lc.LaunchConfigurations))
		os.Exit(1)
	}

	return *lc.LaunchConfigurations[0].InstanceType
}

func GetInstanceTypeFromLaunchTemplate(awsSess *session.Session, launchTemplateName string) string {
	input := &ec2.DescribeLaunchTemplatesInput{
		LaunchTemplateNames: []*string{
			aws.String(launchTemplateName),
		},
	}

	ec2Client := ec2.New(awsSess)

	lt, err := ec2Client.DescribeLaunchTemplates(input)
	if err != nil {
		fmt.Println("Unable to describe Launch Template, err: ", err.Error())
		os.Exit(1)
	}

	if len(lt.LaunchTemplates) != 1 {
		fmt.Println("Expected one Launch Template, found ", len(lt.LaunchTemplates))
		os.Exit(1)
	}

	ltvInput := ec2.DescribeLaunchTemplateVersionsInput{
		LaunchTemplateId: lt.LaunchTemplates[0].LaunchTemplateId,
		Versions:         []*string{aws.String("$Latest")},
	}
	ltv, err := ec2Client.DescribeLaunchTemplateVersions(&ltvInput)
	if err != nil {
		fmt.Println("Unable to describe Launch Template version, error: ", err.Error())
		os.Exit(1)
	}

	if len(ltv.LaunchTemplateVersions) != 1 {
		fmt.Println(`Expected one "$Latest" Launch Template version, received `, len(lt.LaunchTemplates))
		os.Exit(1)
	}

	return *ltv.LaunchTemplateVersions[0].LaunchTemplateData.InstanceType
}

func GetInstanceTypeForAsg(awsSess *session.Session, asgName string) string {
	asg := GetAsg(awsSess, asgName)

	if asg.LaunchConfigurationName != nil {
		return GetInstanceTypeFromLaunchConfiguration(awsSess, *asg.LaunchConfigurationName)
	}

	if asg.LaunchTemplate != nil {
		return GetInstanceTypeFromLaunchTemplate(awsSess, *asg.LaunchTemplate.LaunchTemplateName)
	}

	fmt.Println("Unable to determine the ASG instance type. No Launch Template nor Launch Configuration is defined.")
	os.Exit(1)
	return ""
}

// HowManyServersNeededForAsg computes the theoretical number of servers needed based on the total resources needed,
// assuming near-perfect utilization of server resources. It does not take into account the "wasted" resources on an
// individual server when the free resources are not sufficient to place any of the desired containers.
func HowManyServersNeededForAsg(serverType string, resourcesNeeded ResourceSizes) int64 {
	instanceSpecs, valid := InstanceTypes[serverType]
	if !valid {
		fmt.Println("Invalid server type provided: ", serverType)
		os.Exit(1)
	}

	if resourcesNeeded.LargestMemory > instanceSpecs.MemoryMb {
		fmt.Printf("Configured instance type is not large enough. Available memory is %d, but largest task needs %d",
			instanceSpecs.MemoryMb, resourcesNeeded.LargestMemory)
		os.Exit(1)
	}

	if resourcesNeeded.LargestCPU > instanceSpecs.CPUUnits {
		fmt.Printf("Configured instance type is not large enough. Available CPU is %d, but largest task needs %d",
			instanceSpecs.CPUUnits, resourcesNeeded.LargestCPU)
		os.Exit(1)
	}

	// Some memory in each instance cannot be used because no container can be placed in the last portion available.
	// This assumes the best-case container placement.
	usableMemory := max(1, instanceSpecs.MemoryMb-resourcesNeeded.SmallestMemory)
	usableCPU := max(1, instanceSpecs.CPUUnits-resourcesNeeded.SmallestCPU)

	neededForMem := divideAndRoundUp(resourcesNeeded.TotalMemory, usableMemory)
	neededForCPU := divideAndRoundUp(resourcesNeeded.TotalCPU, usableCPU)

	serversNeeded := max(neededForCPU, neededForMem)
	if serversNeeded > 100 {
		fmt.Printf("Calculated need of %d instances, which is over the predefined threshold. Exiting.", serversNeeded)
		os.Exit(1)
	}

	return serversNeeded
}

func divideAndRoundUp(numerator, divisor int64) int64 {
	return int64(math.Ceil(float64(numerator) / float64(divisor)))
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
