# awsops
A Go based CLI for common AWS operations

This CLI is a bit scrappy and very opinionated for executing operational processes on AWS the way we like to.

## Installation
If you don't care to modify the source, you can grab a prebuilt binary from the `dist/` folder for your 
platform and run it directly. 

You can also clone this repo and use `go build` or `go run` to run it. 

## Configuration
This app makes use of the AWS Go SDK - https://docs.aws.amazon.com/sdk-for-go/api/

This library makes use of standard AWS CLI configurations such as the `~/.aws/credentials` file. If you don't already use the AWS CLI, the easiest way to configure your system is to create a `.aws` folder in your home directory and a `credentials` file inside that directory. The format of the file should be:

```
[nameofprofile]
aws_access_key_id = value
aws_secret_access_key = value
```

Then when using `awsops` you can use the `-p` flag followed by whatever profile from the `credentials` file you want to use. For example: `awsops -p default ...`

## Usage

```
$ awsops
Utility app for common operational tasks for AWS

Usage:
  awsops [command]

Available Commands:
  ecsReplaceInstances Gracefully replace EC2 instances for given ECS cluster
  help                Help about any command

Flags:
      --config string    config file (default is $HOME/.awsops.yaml)
  -h, --help             help for awsops
  -p, --profile string   AWS shared credentials profile to use
  -r, --region string    AWS shared credentials profile to use (default "us-east-1")
  -t, --toggle           Help message for toggle

Use "awsops [command] --help" for more information about a command.
```

```
$ awsops ecsReplaceInstances -h
Gracefully replace EC2 instances for given ECS cluster

Usage:
  awsops ecsReplaceInstances [flags]

Flags:
  -c, --cluster string   ECS cluster name
  -h, --help             help for ecsReplaceInstances

Global Flags:
      --config string    config file (default is $HOME/.awsops.yaml)
  -p, --profile string   AWS shared credentials profile to use
  -r, --region string    AWS shared credentials profile to use (default "us-east-1")
```