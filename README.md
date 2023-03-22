# vultr-cli

The Vultr Command Line Interface

```sh
vultr-cli is a command line interface for the Vultr API

Usage:
  vultr-cli [command]

Available Commands:
  account        Retrieve information about your account
  apps           Display all available applications
  backups        Display backups
  bare-metal     bare-metal is used to access bare metal server commands
  billing        Display billing information
  block-storage  block storage commands
  completion     Generate the autocompletion script for the specified shell
  database       Commands to interact with managed databases on vultr
  dns            dns is used to access dns commands
  firewall       firewall is used to access firewall commands
  help           Help about any command
  instance       commands to interact with instances on vultr
  iso            iso is used to access iso commands
  kubernetes     kubernetes is used to access kubernetes commands
  load-balancer  load balancer commands
  object-storage object storage commands
  os             os is used to access os commands
  plans          get information about Vultr plans
  regions        get regions
  reserved-ip    reserved-ip lets you interact with reserved-ip 
  script         startup script commands
  snapshot       snapshot commands
  ssh-key        ssh-key commands
  user           user commands
  version        Display current version of Vultr-cli
  vpc            Interact with VPCs

Flags:
  --config string   config file (default is $HOME/.vultr-cli.yaml) (default "#HOME/.vultr-cli.yaml")
  -h, --help            help for vultr-cli
  -t, --toggle          Help message for toggle

Use "vultr-cli [command] --help" for more information about a command.
```

## Installation

There are three ways to install `vultr-cli`:
1. Download a release from GitHub
2. From source
3. Package Manager
  - Arch Linux
  - Brew
  - OpenBSD (-current)
  - Snap (Coming soon)
  - Chocolatey (Coming soon)
4. [Docker Hub](https://hub.docker.com/repository/docker/vultr/vultr-cli)
  
### GitHub Release
If you are to visit the `vultr-cli` [releases](https://github.com/vultr/vultr-cli/releases) page. You can download a compiled version of `vultr-cli` for you Linux/MacOS/Windows in 64bit.

### Building from source 

You will need Go installed on your machine in order to work with the source (and make if you decide to pull the repo down).

`go get -u github.com/vultr/vultr-cli/v2`

Another way to build from source is to 

```
git clone git@github.com:vultr/vultr-cli.git or git clone https://github.com/vultr/vultr-cli.git
cd vultr-cli
make builds/vultr-cli_(pass name of os + arch, as shown below)
```

The available make build options are
- make builds/vultr-cli_darwin_amd64
- make builds/vultr-cli_darwin_arm64
- make builds/vultr-cli_linux_386
- make builds/vultr-cli_linux_amd64
- make builds/vultr-cli_linux_arm64
- make builds/vultr-cli_windows_386.exe
- make builds/vultr-cli_windows_amd64.exe
- make builds/vultr-cli_linux_arm

Note that the latter method will install the `vultr-cli` executable in `builds/vultr-cli_(name of os + arch)`.

### Installing on Arch Linux

```sh
pacman -S vultr-cli
```

### Installing via Brew

You will need to tap for formula
``` sh
brew tap vultr/vultr-cli
```

Then install the formula

```sh 
brew install vultr-cli
```

### Installing on Fedora

```sh
dnf install vultr-cli
```

### Installing on OpenBSD

```sh
pkg_add vultr-cli
```

## Using Vultr-cli

### Authentication

In order to use `vultr-cli` you will need to export your [Vultr API KEY](https://my.vultr.com/settings/#settingsapi) 

`export VULTR_API_KEY=your_api_key`

### Examples

`vultr-cli` can interact with all of your Vultr resources. Here are some basic examples to get you started:

##### List all available instances
`vultr-cli instance list`

##### Create an instance
`vultr-cli instance create --region <region-id> --plan <plan-id> --os <os-id> --host <hostname>`

##### Create a DNS Domain
`vultr-cli dns domain create --domain <domain-name> --ip <ip-address>`

##### Utilizing a boolean flag
You should use = when using a boolean flag.

`vultr-cli instance create --region <region-id> --plan <plan-id> --os <os-id> --host <hostname> --notify=true`

##### Utilizing the config flag
The config flag can be used to specify the vultr-cli.yaml file path when it's outside the default location. If the file has the `api-key` defined, the CLI will use the vultr-cli.yaml config, otherwise it will default to reading the environment variable for the api key.

`vultr-cli instance list --config /Users/myuser/vultr-cli.yaml`

### Example vultr-cli.yaml config file

Currently the only available field that you can use with a config file is `api-key`. Your yaml file will have a single entry which would be:

`api-key: MYKEY`

### CLI Autocompletion 
`vultr-cli completion` will return autocompletions, but this feature requires setup. 

Some guides:

<pre>
<h4>Bash:</h4>
  $ source <(yourprogram completion bash)

  <b>To load completions for each session, execute once:</b>
  <b>Linux:</b>
  $ yourprogram completion bash > /etc/bash_completion.d/yourprogram
  
  <b>macOS:</b>
  $ yourprogram completion bash > /usr/local/etc/bash_completion.d/yourprogram

<h4>Zsh:</h4>
  <b>If shell completion is not already enabled in your environment,
  you will need to enable it.  You can execute the following once:</b>

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  <b>To load completions for each session, execute once:</b>
  $ yourprogram completion zsh > "${fpath[1]}/_yourprogram"

  You will need to start a new shell for this setup to take effect.

<h4>fish:</h4>
  $ yourprogram completion fish | source

  <b>To load completions for each session, execute once:</b>
  $ yourprogram completion fish > ~/.config/fish/completions/yourprogram.fish

<h4>PowerShell:</h4>
  PS> yourprogram completion powershell | Out-String | Invoke-Expression

  <b>To load completions for every new session, run:</b>
  PS> yourprogram completion powershell > yourprogram.ps1
  <b>and source this file from your PowerShell profile.</b>
</pre>

## Contributing
Feel free to send pull requests our way! Please see the [contributing guidelines](CONTRIBUTING.md).
