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
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	AwsSess *session.Session
	cfgFile string
	Profile string
	Region  string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "awsops",
	Short: "Utility app for common operational tasks for AWS",
	Long:  `Utility app for common operational tasks for AWS`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.awsops.yaml)")
	rootCmd.PersistentFlags().StringVarP(&Profile, "profile", "p", "", "AWS shared credentials profile to use")
	rootCmd.PersistentFlags().StringVarP(&Region, "region", "r", "us-east-1", "AWS shared credentials profile to use")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalln(err)
		}

		// Search config in home directory with name ".awsops" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".awsops")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initAwsSess() {
	// If profile is provided, use shared creds file and specific profile,
	// otherwise use default credential identification order
	if Profile != "" {
		AwsSess = session.Must(session.NewSession(&aws.Config{
			Region:      aws.String(Region),
			Credentials: credentials.NewSharedCredentials("", Profile),
		}))
	} else {
		AwsSess = session.Must(session.NewSession(&aws.Config{
			Region: aws.String(Region),
		}))
	}
}
