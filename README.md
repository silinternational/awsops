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

mQINBFt7Ku8BEAC5BEsVHEhXyISZc9QqGpJz2vbJ3NxMP3QsdLGtJo3nT8VClO7a
VWusTPy7eIgLfPX53OsFt++TiV/OFgjQHthvaQWtbOfvViwHz1Okn8/cZhQbhHkw
Miimv3w/zqADhE+TAkk8RIc9EC9nSIPgMjFP5w/K6BjIGTYXtgexDvik0/g1u63a
eRauZUvksABAvZn+i/O8fwsrfeDHTvhU8HpERWv4y6iDocAU8AZdB7s66sHNeJzE
usU0u1I3OLwo7ds5OlGiecltcDqVr/Nx1W1wuHbLDOPLLhlbK0WCiV38YTWRh3Ii
U5Dh2EYFo+2ujpiI2ftA8B1SwXIDgq1pisrY6H9pVWWsqEY+yv5EwqmcSPRYJ/2C
ykTG+iwbjt181m2suLof9n3Q9/fXYnbZ/x3Wsz3zT0NQIU3zCRrH5buD/N+GErIq
JZUk1A1Kx2nDfEQAvzSznFZmlHdc/86LloCZGiaFJ7OCRUCD2VHz1fqFTrladPIG
0EFdASgEjslq+8+WHLuihwgJFjTGw2anmFXyyIoqJJbJFXFQr+9WVah2ZXX85CY2
0ZQ7QJ/RsvsaR1PHqH+wKLOUvocFzec/JPvQUhkS5rzbf5XSSGhkGvEPjPRe2y2E
17kLjvQAdPKREEnkGKW4yDSpCmUr/1a5oN+GK39uzRs++E08F/PJ8s9nNQARAQAB
tCFHVElTIEF3c09wcyA8Z3Rpcy1hd3NvcHNAc2lsLm9yZz6JAk4EEwEKADgWIQQd
q/nfvDBuANwiZh3Z933aYLKQtgUCW3sq7wIbAwULCQgHAwUVCgkICwUWAgMBAAIe
AQIXgAAKCRDZ933aYLKQtpeDEACJHWbXfG7GI42117o6hiWsKcRwGMMq/71Hq0v+
jkSDo8MFzqGTg8Aodf6i1S3FfDb56ZvQ2FA6aNy33pMOXyWwpijjwBFkeH+EaHjO
KEHi0MFnplT6fny9NLiTOlvI/ou1bJ8AFJVXYnonB2EsSKit9hk0hh48vwtOVkww
yfrIedyvo/0bk1BqiX8GwPVaspb9F1WrjLy7tW7LeqIOjs0lSoI2MHaTou5GST9z
n/UwvUkUQuJhjpivMSbS67TjvRry1TAST4O0XMN2Cj8Ri+Ubi5vAWg6Q6XTo6CrM
jEJTLlISHk4kjG5uXBAB+ikF+nMIJroNTussAS/9zCOx7FIil3hVHDjew08X3QDR
SAc4RBRenmZk9Td1AQFHgu7bgtLNe12qM45oB+Jj0goXAipkXDrbVbbrKT+Ta+31
MA5D6WXn2dvUbaQM0PDhN8zUap+dENlptPkzbKrWpSdLF+wePAmPRaXVsjDlsLPC
y1NfxrCeqbJYqDTlakR5x2f1hBzYk+UWxcbRpQJQzoV+Zxo/WPVcILr7njO6NI/D
h2tnjj6GebdrCZKl51eb7bCqYT8yyv4cvEW5sehzMv3Ppb3SAP+NJcd0k4CJFMX3
7n2bOE3jVt7mNomrLZhe43ZmumI8nktLvwWlWH2PVudifOjlg9c6lJSFKPU1RmAk
Jnn5mrkCDQRbeyrvARAAwf5pM5hRm362iljV7q63uG+8M+ZN6Zy9sFzeH/roMmpo
DBddS/NSqelymc9xHumNdYzHeewjzJyMc7F8WOcay7Znzom9BPB65Sm145SBTqFQ
tNGaKMwGsY3G8CPOyFZoKRy17GkbKwT0YLCCza5uFukIHFfn7zBCCYXF+Jdu40ik
1ppC2VG5S1Z6PL4ioAyFmHIdGCGimiwJd1d974SXeOeqdUvk598HpceeaK66pM6q
fW5HDe/POwkaArWrg2MUDWkv+pm9azwQ6O3YTcxy44YDX4sYSpv+sl/3AU61lcRI
j4kllXGNFrc0/Xt2kB6yZIrIWqOmIwX9tHq3CdS2qVclu7HVFzCwI9YeEer5OEO8
a/7Kbs9dWgzkYa67fbXfvUz+MDTZQzadgfnzWHQMUwfzR8ut6i2x0f0sa6VK2YnF
VxJpQJXOCdDXfna1DM+eCPx0MZXgW0UjAs0yW3ZgcCiGfHdj7Ej3LwCmoBhay3AK
mA9wEfhoK2YT1+dcfvAfFwOmiebD2wJEu/g6YXF47SGGgwI9B8mMcF2SAaRkEXCt
nhkcB8UfXfUhhWX/DHyiUueKXsHBe6c9vP9uFfmBovR2K+XjylKMjaygQ4cec6kO
ymjWRGfsx7+sSx+3dx++f9pKEWs2ccdyonNw9g8+IrIqKOUY6R6608ETQZiRflcA
EQEAAYkCNgQYAQoAIBYhBB2r+d+8MG4A3CJmHdn3fdpgspC2BQJbeyrvAhsMAAoJ
ENn3fdpgspC2RAUQAIRDQimbfrP0n1FNHV1Lie+VA8gsp1QFn/69xVmqrhyXlsXC
zzH2IGc5cjTj5s8R5tUB99Yiv+ZZ2PmCBGsCydnfVW2g9jnmDQMDltcU81IXRd50
km6USoSsQ/su1YwYGyQnJCzgmQqwoZFXf9cDRZWKD/bAeDRsp4Io0gnRBUj16wI5
7lYyc+oqK4Zlr6vC7ZoERPrcaJGgntxJrSFzXeDsYmE+yY8acxdZWAg5H4wSh5fH
X4s0I9FboaRmNW1E4OhzYunQGpzan5hRstdc41WGl4Eu7xI3yb2UtXPdW0+mgkhu
EHjIG5jMjOhgg2fu16gfHdxFJ3bIx8dsCXhB2JpN+8NWoXVwj8aKE0P7E1n/O9PI
Z4WwfzsNCVUMeCspWzfgB44Og8vmJD9hamQp26bDCtxiCS8ai9HVApnTJiTa6eaF
uTlNr1jcl0iTddRW4+nKkICTSTijEgtr98iaHNW51GG8lIGFZov/RPQmnZ1xnqYu
hFbRYFRZh1Bor0+cAsa3X3okugScPf4uG0o8XUQV7PCGg6Ie+4Y1Dk+hYnCQtxjR
7B1oBLdfcf1GR88RYXpeyCShr5dl5V9BudzcJjMcbvISAAklmJJz/xVNPsbWyouD
x/2vwNgKlFWr68smFLupuf4qCMmVU516RUfmg/JDbKovsd6DRHVizfVp8jes
=b8aO
-----END PGP PUBLIC KEY BLOCK-----
```