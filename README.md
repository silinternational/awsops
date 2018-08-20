# awsops
A Go based CLI and Serverless functions for common AWS operations

This library is a bit scrappy and very opinionated for executing operational processes on AWS the way we like to.

## Installation
If you don't care to modify the source, you can grab a prebuilt binary from the `dist/` folder for your 
platform and run it directly. 

You can also clone this repo and use `go build` or `go run` from the `cli/` folder to run it. 

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

## GPG Public Key
Binaries for `awsops` are also signed for you to verify it is from us. Our public GPG key is:

```
-----BEGIN PGP PUBLIC KEY BLOCK-----

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
=NK8G
-----END PGP PUBLIC KEY BLOCK-----
```