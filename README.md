# vultr-cli

The Vultr Command Line Interface
```
vultr-cli is a command line interface for the Vultr API

Usage:
  vultr-cli [command]

Available Commands:
  account       Retrieve information about your account
  api-key       retrieve information about the current API key
  apps          Display all available applications
  backups       display all available backups
  bare-metal    bare-metal is used to access bare metal server commands
  block-storage block storage commands
  dns           dns is used to access dns commands
  firewall      firewall is used to access firewall commands
  help          Help about any command
  iso           iso is used to access iso commands
  network       network interacts with network actions
  os            grab all available operating systems
  plans         get information about Vultr plans
  regions       get regions
  reserved-ip   reserved-ip lets you interact with reserved-ip
  script        startup script commands
  server        commands to interact with servers on vultr
  snapshot      snapshot commands
  ssh-key       ssh-key commands
  user          user commands
  version       Display current version of Vultr-cli

Flags:
      --config string   config file (default is $HOME/.vultr-cli.yaml)
  -h, --help            help for vultr-cli
  -t, --toggle          Help message for toggle

Use "vultr-cli [command] --help" for more information about a command.
```

## Installation

There are two ways to install `vultr-cli`:
1. From source
2. Package Manager (Coming soon)
  - Brew
  - Snap
  - Chocolatey

### Building from source 

You will need Go installed on your machine in order to work with the source (and make if you decide to pull the repo down).

`go get -u github.com/vultr/vultr-cli`

Another way to build from source is to 

```
git clone git@github.com:vultr/vultr-cli.git or git clone https://github.com/vultr/vultr-cli.git
cd vultr-cli
make build_(pass name of os + arch)
```

The available make build options are
- build_mac
- build_linux_386
- build_linux_amd64
- build_windows_64
- build_windows_32

### Installing via Brew

You will need to tap for formula
``` sh
brew tap vultr/vultr-cli
```

Then install the formula

```sh 
brew install vultr-cli
```

## Using Vultr-cli

### Authentication

In order to use `vultr-cli` you will need to export your [Vultr API KEY](https://my.vultr.com/settings/#settingsapi) 

`export VULTR_API_KEY=your_api_key`

### Examples

`vultr-cli` can interact with all of your Vultr resources. Here are some basic examples to get you started:

##### List all available servers
`vultr-cli server list`

##### Create a server
`vultr-cli server create --region <region-id> --plan <plan-id> --os <os-id> --hostname <hostname>` 

##### Create a DNS Domain
`vultr-cli dns domain create --domain <domain-name> --ip <ip-address>`

## Contributing
Feel free to send pull requests our way! Please see the [contributing guidelines](CONTRIBUTING.md).
