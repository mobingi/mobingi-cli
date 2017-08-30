[![logo](./logos/mobinginewlogo.png)](https://mobingi.co.jp/)

[![Build Status](https://travis-ci.org/mobingi/mobingi-cli.svg?branch=master)](https://travis-ci.org/mobingi/mobingi-cli)
[![Build status](https://ci.appveyor.com/api/projects/status/k7tmyr3l4dws4usd/branch/master?svg=true)](https://ci.appveyor.com/project/flowerinthenight/mobingi-cli/branch/master)

# mobingi-cli

```
Command line interface for Mobingi API and services.

Usage:
  mobingi-cli [command]

Available Commands:
  creds       manage your credentials
  help        help about any command
  login       login to Mobingi API
  registry    manage your Mobingi docker registry
  reset       reset config to defaults
  stack       manage your stack
  svrconf     manage your server config file
  template    manage your ALM templates
  version     print the version

Flags:
      --token string    access token
      --url string      base url for API
      --rurl string     base url for Docker Registry
      --apiver string   API version (default "v3")
  -f, --fmt string      output format (values depends on command)
  -o, --out string      full file path to write the output
      --indent int      indent padding when fmt is 'json' (default 2)
      --timeout int     timeout in seconds (default 120)
      --verbose         verbose output
      --debug           debug mode when error occurs
  -h, --help            help for mobingi-cli

Use "mobingi-cli [command] --help" for more information about a command.
```

# Getting started

### Getting mobingi-cli

The easiest way to get mobingi-cli is to use one of the [pre-built release binaries](https://github.com/mobingi/mobingi-cli/releases) which are available for OSX, Linux, and Windows.

If you want to try the latest version, you can build mobingi-cli from the master branch. You need to have [Go](https://golang.org/) installed (version 1.8+ is required). Note that the master branch may be in an unstable or even broken state during development.

### Building mobingi-cli

```
$ git clone https://github.com/mobingi/mobingi-cli
$ cd mobingi-cli
$ go build -v
$ ./mobingi-cli
```

You can also install the binary to your `$GOPATH/bin` folder (`$GOPATH/bin` should be added to your `$PATH` environment variable). 

```
$ go get -u github.com/mobingi/mobingi-cli
$ mobingi-cli
```

# Usage

## Login

This is the first command you need to run to use the other commands. To login, run

```
$ mobingi-cli login --client-id=foo --client-secret=bar
```

This will create a file `config.yml` under `$HOME/.mobingi-cli/` folder that will contain the access token to be used for your subsequent commands, alongside other configuration values.

## Stack operations

### List stacks

Examples:

```
$ mobingi-cli stack list
$ mobingi-cli stack list --fmt=json
$ mobingi-cli stack list --fmt=raw --out=`echo $HOME`/out.txt
```

Enclose with double quotes if absolute file path has whitespace(s) in it.

### Describe a stack

Examples:

```
$ mobingi-cli stack describe --id=foo
$ mobingi-cli stack describe --id=foo --fmt=json
$ mobingi-cli stack describe --id=foo --fmt=raw --out=/home/bar/out.txt
```

You can get the stack id from the `stack list` command.

### Create a stack

#### API v3

Starting in v3, we create stacks using ALM templates. Below is an example of a very simple template that creates a single EC2 instance:

```
{
  "version": "2017-03-03",
  "label": "template version label #1",
  "description": "This template creates a sample stack with EC2 instance on AWS",
  "vendor": {
    "aws": {
      "cred": "** Your AWS Security Key ID **",
      "secret": "** Your AWS Security Key Secret, remove line if you have a Mobingi account **",
      "region": "ap-northeast-1"
    }
  },
  "configurations": [
    {
      "role": "web",
      "flag": "Single1",
      "provision": {
        "instance_type": "t2.micro",
        "instance_count": 1,
        "keypair": false,
        "subnet": {
          "cidr": "10.0.1.0/24",
          "public": true,
          "auto_assign_public_ip": true
        },
        "availability_zone": "ap-northeast-1c"
      }
    }
  ]
}
```

After saving your template to a file, run the command:

```
$ mobingi-cli stack create --alm-template=<path_to_template-file>
```

Example:

```
$ mobingi-cli stack create --alm-template=/home/user/aws-single-ec2.json
```

#### API v2

You can run `$ mobingi-cli stack create -h` to see the defaults.

Examples:

```
$ mobingi-cli stack create --nickname=sample
$ mobingi-cli stack create --nickname=sample --min=2 --max=2
```

If the `--cred` option is not provided (just like in the examples above), cli will attempt to get your list of credentials and use the first one (if more than one). You can view your credentials list using the command:

```
$ mobingi-cli creds list
```

### Update stack

#### API v3

Similar to stack creation, you only need to update some parts of your ALM template to update your stack.

```
$ mobingi-cli stack update --id=foo --alm-template=/home/user/aws-single-ec2-updated.json
```

#### API v2

Examples:

```
$ mobingi-cli stack update --id=foo --min=5 --max=20
$ mobingi-cli stack update --id=foo --spot-range=25
```

### Delete a stack

Example:

```
$ mobingi-cli stack delete --id=foo
```

## ALM template operations

### List stack template versions

Examples:

```
$ mobingi-cli template versions --id=foo
$ mobingi-cli template versions --id=foo --fmt=json
```

You can get the stack id from the `stack list` command.

### Compare template versions

You can compare template versions from the same stack, versions from different stacks, or a local template file to a specific template version.

Example:

```
$ mobingi-cli stack list
STACK ID                          STACK NAME                   PLATFORM     STATUS              ...
mo-58c2297d25645-q38pTmeey-tk     small lunch behave           AWS          CREATE_COMPLETE     ...
mo-58c2297d25645-PxviFSJQV-tk     chronic leaflet flourish     AWS          CREATE_COMPLETE     ...

$ mobingi-cli template versions --id mo-58c2297d25645-PxviFSJQV-tk
VERSION ID                           LATEST     LAST MODIFIED                     SIZE
jbyW_PxMAauQmOS31dUhij4KIqHAtqW2     true       Wed, 30 Aug 2017 02:32:43 UTC     472
1xoPd.cg3juHK94vC8IdUh1bexx7sQ1T     false      Tue, 29 Aug 2017 09:47:50 UTC     453

$ mobingi-cli template compare --src-sid mo-58c2297d25645-PxviFSJQV-tk \
      --src-vid jbyW_PxMAauQmOS31dUhij4KIqHAtqW2 \
      --tgt-vid 1xoPd.cg3juHK94vC8IdUh1bexx7sQ1T
[mobingi-cli]: info: diff:
{
  "new": [],
  "removed": [],
  "edited": {
    "label": {
      "oldvalue": "template version label #1",
      "newvalue": "template version label #1 (update)"
    },
    "description": {
      "oldvalue": "This template creates a sample stack with EC2 instance on AWS",
      "newvalue": "This template creates a sample stack with EC2 instance on AWS (update)"
    },
    "configurations\/provision\/instance_type": {
      "oldvalue": "t2.micro",
      "newvalue": "m3.medium"
    },
    "configurations\/provision\/instance_count": {
      "oldvalue": 1,
      "newvalue": 2
    }
  }
}
```

To view help information, run the command:

```
$ mobingi-cli template compare -h
```

## Server config operations

### Show server config

Example:

```
$ mobingi-cli svrconf show --id=foo
```

You can get the stack id from the `stack list` command.

### Update server config

Examples:

```
$ mobingi-cli svrconf update --id=foo --env=KEY1:value1,KEY2:value2,KEYx:valuex
```

If you have whitespaces in the input, enclose it with double quotes

```
$ mobingi-cli svrconf update --id=foo --env="KEY1: value1, KEY2: value2, KEYx: valuex"
```

To clear all environment variables, set `--env=null` option

```
$ mobingi-cli svrconf update --id=foo --env=null
```

To update server config filepath, run

```
$ mobingi-cli svrconf update --id=foo --filepath=git://github.com/mobingilabs/default
```

Note that when you provide update options simultaneously (for example, you provide `--env=FOO:bar` and `--filepath=test` at the same time), the tool will send each option as a separate request.

## Vendor credentials

### View vendor credentials

Examples:

```
$ mobingi-cli creds list
$ mobingi-cli creds list --fmt=json
$ mobingi-cli creds list --fmt=raw
```

## Mobingi Docker registry

### Get token for Docker Registry API

To get token for Docker Registry API access, run

```
$ mobingi-cli registry token \
      --username=foo \
      --password=bar \
      --service="Mobingi Docker Registry" \
      --scope="repository:foo/container:*"
```

where `username` is a subuser under your Mobingi account. You can also remove `--service`, `--username` and/or `--password`.

```
$ mobingi-cli registry token --scope="repository:foo/container:*"
Username:
Password:
```

By default, it will only print the token value. To print the raw JSON output, append the `--fmt=raw` option.
 
This is useful when you want to access the registry directly using other tools. For example, you can use the token when using Docker Registry API via `curl`.

```
$ curl -H "Authorization: Bearer token" \
      -H "Accept application/vnd.docker.distribution.manifest.v2+json" \
      https://registry.mobingi.com/v2/foo/container/manifests/latest
```

### List image tags

To list image tags, run

```
$ mobingi-cli registry tags [--username=foo] [--password=bar] --image=image_name
```

### List registry catalog

To print the catalog, run

```
$ mobingi-cli registry catalog [--username=foo] [--password=bar]
```

Note that this command is inherently slow.

### Print tag manifest

To print a tag manifest, run

```
$ mobingi-cli registry manifest [--username=foo] [--password=bar] --image=hello:latest
```

You can also write the output to a file via the `--fmt=full_path_to_file` option.

### Delete a tag

To delete a tag, run

```
$ mobingi-cli registry delete [--username=foo] [--password=bar] --image=hello:latest
```

## Verbose output

You can use the global `--verbose` option if you want to see more information during the command execution.
