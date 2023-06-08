// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"log"

	"github.com/spf13/cobra"

	"github.com/silinternational/awsops/lib"
)

var atLeastServiceDesiredCount bool

// rightSizeClusterCmd represents the scaleCluster command
var rightSizeClusterCmd = &cobra.Command{
	Use:   "rightSizeCluster",
	Short: "Scale ASG for ECS cluster to minimum needed servers",
	Long: `This command calculates total memory and CPU needed
for all services in the given ECS cluster and then adjusts 
instance count in the ASG based on instance type/size to 
support running all tasks with as few servers as is needed.

This function may scale a cluster up or down depending on services.`,
	Run: func(cmd *cobra.Command, args []string) {
		initAwsSess()
		err := lib.RightSizeAsgForEcsCluster(AwsSess, cluster, atLeastServiceDesiredCount)
		if err != nil {
			log.Fatalln(err.Error())
		}
	},
}

func init() {
	ecsCmd.AddCommand(rightSizeClusterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rightSizeClusterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	rightSizeClusterCmd.Flags().BoolVar(&atLeastServiceDesiredCount, "atLeastServiceDesiredCount", false, "Ensure at least as many EC2 instances as largest ECS service desired count.")
}
