# awsops
A Go based CLI and Serverless functions for common AWS operations

This library is a bit scrappy and very opinionated for executing operational processes on AWS the way we like to.

## Installation
If you don't care to modify the source, you can grab a prebuilt binary from [Releases](https://github.com/silinternational/awsops/releases).

You can also clone this repo and use `go build` to build a binary or `go run` to run it directly. Or simply run `make` to build for multiple
platforms.

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
  ecs         ECS related actions, run 'awsops ecs' to view list of subcommands
  help        Help about any command

Flags:
      --config string    config file (default is $HOME/.awsops.yaml)
  -h, --help             help for awsops
  -p, --profile string   AWS shared credentials profile to use
  -r, --region string    AWS shared credentials profile to use (default "us-east-1")
  -t, --toggle           Help message for toggle

Use "awsops [command] --help" for more information about a command.
```

```
$ awsops ecs
ECS related actions, run 'awsops ecs' to view list of subcommands

Usage:
  awsops ecs [flags]
  awsops ecs [command]

Available Commands:
  listInstanceIPs  List Instance IPs for ECS Cluster
  replaceInstances Gracefully replace EC2 instances for given ECS cluster
  rightSizeCluster Scale ASG for ECS cluster to minimum needed servers

Flags:
  -c, --cluster string   ECS cluster name
  -h, --help             help for ecs

Global Flags:
      --config string    config file (default is $HOME/.awsops.yaml)
  -p, --profile string   AWS shared credentials profile to use
  -r, --region string    AWS shared credentials profile to use (default "us-east-1")

Use "awsops ecs [command] --help" for more information about a command.
```

```
$ awsops ecs listInstanceIPs --help
Command returns a space separated list of IP addresses for instances in an ECS cluster

Usage:
  awsops ecs listInstanceIPs [flags]

Flags:
  -h, --help   help for listInstanceIPs

Global Flags:
  -c, --cluster string   ECS cluster name
      --config string    config file (default is $HOME/.awsops.yaml)
  -p, --profile string   AWS shared credentials profile to use
  -r, --region string    AWS shared credentials profile to use (default "us-east-1")
```

```
$ awsops ecs replaceInstances --help
Gracefully replace EC2 instances for given ECS cluster

Usage:
  awsops ecs replaceInstances [flags]

Flags:
  -h, --help   help for replaceInstances

Global Flags:
  -c, --cluster string   ECS cluster name
      --config string    config file (default is $HOME/.awsops.yaml)
  -p, --profile string   AWS shared credentials profile to use
  -r, --region string    AWS shared credentials profile to use (default "us-east-1")
```

```
$ awsops ecs rightSizeCluster --help
This command calculates total memory and CPU needed
for all services in the given ECS cluster and then adjusts
instance count in the ASG based on instance type/size to
support running all tasks with as few servers as is needed.

This function may scale a cluster up or down depending on services.

Usage:
  awsops ecs rightSizeCluster [flags]

Flags:
  -h, --help   help for rightSizeCluster

Global Flags:
  -c, --cluster string   ECS cluster name
      --config string    config file (default is $HOME/.awsops.yaml)
  -p, --profile string   AWS shared credentials profile to use
  -r, --region string    AWS shared credentials profile to use (default "us-east-1")
```
