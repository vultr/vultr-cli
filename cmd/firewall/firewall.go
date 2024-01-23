// Package firewall provides the functionality for firewall commands in the CLI
package firewall

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
)

var (
	ruleLong    = `Show commands available for firewall rules`
	ruleExample = `
	# Full example
	vultr-cli firewall rule

	# Shortened example with aliases
	vultr-cli fw r
	`
	ruleCreateLong = `
	Create a new firewall rule in the provided firewall group

	If protocol is TCP or UDP, port must be provided.

	An ip-type of v4 or v6 must be supplied for all rules.
	`
	ruleCreateExample = `
	# Full examples
	vultr-cli firewall rule create --id=f04ae5aa-ff6a-4078-900d-78cc17dca2d5 --ip-type=v4 --protocol=tcp --size=24 \
		--subnet=127.0.0.0 --port=30000

	vultr-cli firewall rule create --id=f04ae5aa-ff6a-4078-900d-78cc17dca2d5 --ip-type=v4 --protocol=icmp --size=24 --subnet=127.0.0.0

	# Shortened example with aliases
	vultr-cli fw r c -i=f04ae5aa-ff6a-4078-900d-78cc17dca2d5 -t=v4 -p=tcp -z=24 -s=127.0.0.0 -r=30000
	`
	ruleDeleteLong    = `Delete a firewall rule in the provided firewall group`
	ruleDeleteExample = `
	# Full example
	vultr-cli firewall rule delete 704ac064-4ff2-49ca-a6e6-88262cca8f8a f31ade4f-2308-4a58-82c6-2d1bae0837b3

	# Shortened example with aliases
	vultr-cli fw r d 704ac064-4ff2-49ca-a6e6-88262cca8f8a f31ade4f-2308-4a58-82c6-2d1bae0837b3
	`
	ruleGetLong    = `Get a firewall rule in the provided firewall group`
	ruleGetExample = `
	# Full example
	vultr-cli firewall rule get 704ac064-4ff2-49ca-a6e6-88262cca8f8a f31ade4f-2308-4a58-82c6-2d1bae0837b3

	# Shortened example with aliases
	vultr-cli fw r get 704ac064-4ff2-49ca-a6e6-88262cca8f8a f31ade4f-2308-4a58-82c6-2d1bae0837b3
	`
	ruleListLong    = `List all firewall rules in the provided firewall group`
	ruleListExample = `
	# Full example
	vultr-cli firewall rule list 704ac064-4ff2-49ca-a6e6-88262cca8f8a

	# Shortened example with aliases
	vultr-cli fw r l 704ac064-4ff2-49ca-a6e6-88262cca8f8a
	`
)

type Options struct {
	Base     *cli.Base
	GroupReq *govultr.FirewallGroupReq
	RuleReq  *govultr.FirewallRuleReq
}

// NewCmdFirewall provides the CLI command functionality for Firewall
func NewCmdFirewall(base *cli.Base) *cobra.Command {
	o := &Options{Base: base}

	cmd := &cobra.Command{
		Use:     "firewall",
		Short:   "firewall is used to access firewall commands",
		Long:    ``,
		Aliases: []string{"fw"},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			utils.SetOptions(o.Base, cmd, args)
			if !o.Base.HasAuth {
				return errors.New(utils.APIKeyError)
			}
			return nil
		},
	}

	// Group
	group := &cobra.Command{
		Use:     "group",
		Short:   "group is used to access firewall group commands",
		Long:    ``,
		Aliases: []string{"g"},
	}

	// Group List
	groupList := &cobra.Command{
		Use:     "list",
		Short:   "List all firewall groups",
		Aliases: []string{"l"},
		Run: func(cmd *cobra.Command, args []string) {
			o.Base.Options = utils.GetPaging(cmd)

			groups, meta, err := o.GroupList()
			if err != nil {
				printer.Error(fmt.Errorf("error retrieving firewall group list : %v", err))
				os.Exit(1)
			}

			data := &FirewallGroupsPrinter{Groups: groups, Meta: meta}
			o.Base.Printer.Display(data, nil)
		},
	}

	groupList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	groupList.Flags().IntP(
		"per-page",
		"p",
		utils.PerPageDefault,
		"(optional) Number of items requested per page. Default is 100 and Max is 500.",
	)

	// Group Get
	groupGet := &cobra.Command{
		Use:   "get <Firewall Group ID>",
		Short: "get firewall group",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a firewall group ID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			group, err := o.GroupGet()
			if err != nil {
				printer.Error(fmt.Errorf("error getting firewall group : %v", err))
				os.Exit(1)
			}

			data := &FirewallGroupPrinter{Group: *group}
			o.Base.Printer.Display(data, nil)
		},
	}

	// Group Create
	groupCreate := &cobra.Command{
		Use:     "create",
		Short:   "create a firewall group",
		Aliases: []string{"c"},
		Run: func(cmd *cobra.Command, args []string) {
			description, errDe := cmd.Flags().GetString("description")
			if errDe != nil {
				printer.Error(fmt.Errorf("error parsing 'description' flag for firewall group create: %v", errDe))
				os.Exit(1)
			}

			o.GroupReq = &govultr.FirewallGroupReq{
				Description: description,
			}

			grp, err := o.GroupCreate()
			if err != nil {
				printer.Error(fmt.Errorf("error creating firewall group : %v", err))
				os.Exit(1)
			}

			data := &FirewallGroupPrinter{Group: *grp}
			o.Base.Printer.Display(data, nil)
		},
	}

	groupCreate.Flags().StringP("description", "d", "", "(optional) Description of firewall group.")

	// Group Update
	groupUpdate := &cobra.Command{
		Use:     "update <Firewall Group ID>",
		Short:   "Update firewall group description",
		Aliases: []string{"u"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a firewall group ID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			description, errDe := cmd.Flags().GetString("description")
			if errDe != nil {
				printer.Error(fmt.Errorf("error parsing 'description' flag for firewall group update : %v", errDe))
				os.Exit(1)
			}

			o.GroupReq = &govultr.FirewallGroupReq{
				Description: description,
			}

			if err := o.GroupUpdate(); err != nil {
				printer.Error(fmt.Errorf("error updating firewall group : %v", err))
				os.Exit(1)
			}

			o.Base.Printer.Display(printer.Info("firewall group has been updated"), nil)
		},
	}

	groupUpdate.Flags().StringP("description", "d", "", "Description of firewall group.")
	if err := groupUpdate.MarkFlagRequired("description"); err != nil {
		printer.Error(fmt.Errorf("error marking firewall group 'description' flag required: %v", err))
		os.Exit(1)
	}

	// Group Delete
	groupDelete := &cobra.Command{
		Use:     "delete <Firewall Group ID>",
		Short:   "Delete a firewall group",
		Aliases: []string{"d", "destroy"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a firewall group ID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := o.GroupDelete(); err != nil {
				printer.Error(fmt.Errorf("error deleting firewall group : %v", err))
				os.Exit(1)
			}

			o.Base.Printer.Display(printer.Info("firewall group has been deleted"), nil)
		},
	}

	group.AddCommand(groupList, groupGet, groupCreate, groupUpdate, groupDelete)

	// Rule
	rule := &cobra.Command{
		Use:     "rule",
		Short:   "rule is used to access firewall rule commands",
		Long:    ruleLong,
		Example: ruleExample,
		Aliases: []string{"r"},
	}

	// Rule List
	ruleList := &cobra.Command{
		Use:     "list <firewall group ID>",
		Short:   "List all firewall rules",
		Long:    ruleListLong,
		Example: ruleListExample,
		Aliases: []string{"l"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a firewall group ID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			o.Base.Options = utils.GetPaging(cmd)

			rules, meta, err := o.RuleList()
			if err != nil {
				printer.Error(fmt.Errorf("error retrieving firewall rule list : %v", err))
				os.Exit(1)
			}

			data := &FirewallRulesPrinter{Rules: rules, Meta: meta}
			o.Base.Printer.Display(data, nil)
		},
	}

	ruleList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	ruleList.Flags().IntP(
		"per-page",
		"p",
		utils.PerPageDefault,
		"(optional) Number of items requested per page. Default is 100 and Max is 500.",
	)

	// Rule Get
	ruleGet := &cobra.Command{
		Use:     "get <Firewall Group ID> <Firewall Rule Number>",
		Short:   "Get firewall rule",
		Long:    ruleGetLong,
		Example: ruleGetExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide a firewall group ID and firewall rule number")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			rule, err := o.RuleGet()
			if err != nil {
				printer.Error(fmt.Errorf("error getting firewall rule : %v", err))
				os.Exit(1)
			}

			data := &FirewallRulePrinter{Rule: *rule}
			o.Base.Printer.Display(data, nil)
		},
	}

	// Rule Create
	ruleCreate := &cobra.Command{
		Use:     "create",
		Short:   "create a firewall rule",
		Long:    ruleCreateLong,
		Example: ruleCreateExample,
		Aliases: []string{"c"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a firewall group ID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			protocol, errPr := cmd.Flags().GetString("protocol")
			if errPr != nil {
				printer.Error(fmt.Errorf("error parsing 'protocol' flag for firewall group create : %v", errPr))
				os.Exit(1)
			}

			subnet, errSu := cmd.Flags().GetString("subnet")
			if errSu != nil {
				printer.Error(fmt.Errorf("error parsing 'subnet' flag for firewall group create : %v", errSu))
				os.Exit(1)
			}

			ipType, errIp := cmd.Flags().GetString("ip-type")
			if errIp != nil {
				printer.Error(fmt.Errorf("error parsing 'ip-type' flag for firewall group create : %v", errIp))
				os.Exit(1)
			}

			size, errSi := cmd.Flags().GetInt("size")
			if errSi != nil {
				printer.Error(fmt.Errorf("error parsing 'size' flag for firewall group create : %v", errSi))
				os.Exit(1)
			}

			source, errSo := cmd.Flags().GetString("source")
			if errSo != nil {
				printer.Error(fmt.Errorf("error parsing 'source' flag for firewall group create : %v", errSo))
				os.Exit(1)
			}

			port, errPo := cmd.Flags().GetString("port")
			if errPo != nil {
				printer.Error(fmt.Errorf("error parsing 'port' flag for firewall group create : %v", errPo))
				os.Exit(1)
			}

			notes, errNo := cmd.Flags().GetString("notes")
			if errNo != nil {
				printer.Error(fmt.Errorf("error parsing 'notes' flag for firewall group create : %v", errNo))
				os.Exit(1)
			}

			o.RuleReq = &govultr.FirewallRuleReq{
				Protocol:   protocol,
				Subnet:     subnet,
				SubnetSize: size,
				Notes:      notes,
			}

			if port != "" {
				o.RuleReq.Port = port
			}

			if source != "" {
				o.RuleReq.Source = source
			}

			if ipType == "" {
				printer.Error(fmt.Errorf("a firewall rule requires an IP type. Pass an --ip-type value of v4 or v6"))
				os.Exit(1)
			}

			if ipType != "" {
				o.RuleReq.IPType = ipType
			}

			rule, err := o.RuleCreate()
			if err != nil {
				printer.Error(fmt.Errorf("error creating firewall rule : %v", err))
				os.Exit(1)
			}

			data := &FirewallRulePrinter{Rule: *rule}
			o.Base.Printer.Display(data, nil)
		},
	}

	ruleCreate.Flags().StringP("protocol", "p", "", "Protocol type. Possible values: 'icmp', 'tcp', 'udp', 'gre'.")
	if err := ruleCreate.MarkFlagRequired("protocol"); err != nil {
		printer.Error(fmt.Errorf("error marking firewall rule create 'protocol' flag required : %v", err))
		os.Exit(1)
	}

	ruleCreate.Flags().StringP("subnet", "s", "", "The IPv4 network in CIDR notation.")
	if err := ruleCreate.MarkFlagRequired("subnet"); err != nil {
		printer.Error(fmt.Errorf("error marking firewall rule create 'subnet' flag required : %v", err))
		os.Exit(1)
	}

	ruleCreate.Flags().IntP("size", "z", 0, "The number of bits for the netmask in CIDR notation.")
	if err := ruleCreate.MarkFlagRequired("size"); err != nil {
		printer.Error(fmt.Errorf("error marking firewall rule create 'size' flag required : %v", err))
		os.Exit(1)
	}

	ruleCreate.Flags().StringP(
		"source",
		"",
		"",
		"(optional) When empty, uses value from subnet and size. If \"cloudflare\", allows all Cloudflare IP space through firewall.",
	)

	ruleCreate.Flags().StringP("ip-type", "t", "", "The type of IP rule - v4 or v6.")
	if err := ruleCreate.MarkFlagRequired("ip-type"); err != nil {
		printer.Error(fmt.Errorf("error marking firewall rule create 'ip-type' flag required : %v", err))
		os.Exit(1)
	}

	ruleCreate.Flags().StringP(
		"port",
		"r",
		"",
		"(optional) TCP/UDP only. This field can be an integer value specifying a port or a colon separated port range.",
	)

	ruleCreate.Flags().StringP("notes", "n", "", "(optional) This field supports notes up to 255 characters.")

	// Rule Delete
	ruleDelete := &cobra.Command{
		Use:     "delete <Firewall Group ID> <Firewall Rule Number>",
		Short:   "delete a firewall rule",
		Long:    ruleDeleteLong,
		Example: ruleDeleteExample,
		Aliases: []string{"d", "destroy"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide a firewall group ID and firewall rule number")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := o.RuleDelete(); err != nil {
				printer.Error(fmt.Errorf("error deleting firewall rule : %v", err))
				os.Exit(1)
			}

			o.Base.Printer.Display(printer.Info("firewall rule deleted"), nil)
		},
	}

	rule.AddCommand(ruleList, ruleGet, ruleCreate, ruleDelete)

	cmd.AddCommand(group, rule)

	return cmd
}

// GroupList ...
func (o *Options) GroupList() ([]govultr.FirewallGroup, *govultr.Meta, error) {
	groups, meta, _, err := o.Base.Client.FirewallGroup.List(o.Base.Context, o.Base.Options)
	return groups, meta, err
}

// GroupGet ...
func (o *Options) GroupGet() (*govultr.FirewallGroup, error) {
	group, _, err := o.Base.Client.FirewallGroup.Get(o.Base.Context, o.Base.Args[0])
	return group, err
}

// GroupCreate ...
func (o *Options) GroupCreate() (*govultr.FirewallGroup, error) {
	group, _, err := o.Base.Client.FirewallGroup.Create(o.Base.Context, o.GroupReq)
	return group, err
}

// GroupUpdate ...
func (o *Options) GroupUpdate() error {
	return o.Base.Client.FirewallGroup.Update(o.Base.Context, o.Base.Args[0], o.GroupReq)
}

// GroupDelete ...
func (o *Options) GroupDelete() error {
	return o.Base.Client.FirewallGroup.Delete(o.Base.Context, o.Base.Args[0])
}

// RuleList ...
func (o *Options) RuleList() ([]govultr.FirewallRule, *govultr.Meta, error) {
	rules, meta, _, err := o.Base.Client.FirewallRule.List(o.Base.Context, o.Base.Args[0], o.Base.Options)
	return rules, meta, err
}

// RuleGet ...
func (o *Options) RuleGet() (*govultr.FirewallRule, error) {
	id, errIn := strconv.Atoi(o.Base.Args[1])
	if errIn != nil {
		return nil, fmt.Errorf("unable to convert int to string : %v", errIn)
	}

	rule, _, err := o.Base.Client.FirewallRule.Get(o.Base.Context, o.Base.Args[0], id)
	return rule, err
}

// RuleCreate ...
func (o *Options) RuleCreate() (*govultr.FirewallRule, error) {
	rule, _, err := o.Base.Client.FirewallRule.Create(o.Base.Context, o.Base.Args[0], o.RuleReq)
	return rule, err
}

// RuleDelete ...
func (o *Options) RuleDelete() error {
	id, err := strconv.Atoi(o.Base.Args[1])
	if err != nil {
		return fmt.Errorf("unable to convert int to string : %v", err)
	}
	return o.Base.Client.FirewallRule.Delete(o.Base.Context, o.Base.Args[0], id)
}
