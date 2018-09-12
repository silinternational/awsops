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
	"fmt"
	"github.com/silinternational/awsops/lib"
	"github.com/spf13/cobra"
	"os"
)

var functionName string
var payload string

// invokeCmd represents the invoke command
var invokeCmd = &cobra.Command{
	Use:   "invoke",
	Short: "Invoke a lambda function",
	Long:  "Invoke a lambda function",
	Run: func(cmd *cobra.Command, args []string) {
		initAwsSess()

		result, err := lib.LambdaInvoke(AwsSess, functionName, payload)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Printf("Response: [code: %v] %s", *result.StatusCode, result.Payload)
	},
}

func init() {
	lambdaCmd.AddCommand(invokeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// invokeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// invokeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	invokeCmd.Flags().StringVarP(&functionName, "function", "f", "", "Lambda function name")
	invokeCmd.Flags().StringVarP(&payload, "payload", "b", "", "Lambda function input payload as JSON string")
}
