// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"

	"github.com/silinternational/awsops/lib"
)

// ecsReplaceInstancesCmd represents the ecsReplaceInstances command
var replaceInstancesCmd = &cobra.Command{
	Use:   "replaceInstances",
	Short: "Gracefully replace EC2 instances for given ECS cluster",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		initAwsSess()

		asgName := lib.GetAsgNameForEcsCluster(AwsSess, cluster)
		if asgName == "" {
			log.Fatalln("Unable to find ASG name for ECS cluster", cluster)
		}

		instancesToTerminate := lib.GetInstanceListForAsg(AwsSess, asgName)

		fmt.Println("Replacing EC2 instances one at a time for ECS cluster: ", cluster)
		fmt.Println("ASG: ", asgName)

		lib.DetachAndReplaceAsgInstances(AwsSess, asgName, instancesToTerminate)

		fmt.Printf("Terminating %v instances...\n", len(instancesToTerminate))
		for _, instanceID := range instancesToTerminate {
			_, err := terminateInstance(*instanceID)
			if err != nil {
				log.Fatalln("Unable to terminate instance:", err)
			}
			waitForZeroPendingTasks(cluster)
		}
		fmt.Println("Finished terminating instances")

		instances := lib.GetInstanceListForEcsCluster(AwsSess, cluster)
		fmt.Println("Final instances in cluster: ", len(instances))
		fmt.Println("All done. Be sure to tip your waiter and thank AppsDev for making your life better.")
	},
}

func init() {
	ecsCmd.AddCommand(replaceInstancesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ecsReplaceInstancesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ecsReplaceInstancesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func terminateInstance(id string) (bool, error) {
	svc := ec2.New(AwsSess)
	instanceStatus, err := svc.DescribeInstanceStatus(&ec2.DescribeInstanceStatusInput{
		InstanceIds: []*string{&id},
	})
	if err != nil {
		return false, err
	}

	if *instanceStatus.InstanceStatuses[0].InstanceState.Name != "terminated" {
		fmt.Println("Terminating instance: ", id)
		_, err := svc.TerminateInstances(&ec2.TerminateInstancesInput{
			InstanceIds: []*string{&id},
		})
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func waitForZeroPendingTasks(cluster string) {
	var pendingTasks int64

	time.Sleep(120 * time.Second)
	for pendingTasks = 1000; pendingTasks > 0; {
		time.Sleep(30 * time.Second)
		pendingTasks = lib.GetPendingEcsTasksCount(AwsSess, cluster)
		fmt.Printf("\rPending tasks: %v", pendingTasks)
	}
	fmt.Println()
}
