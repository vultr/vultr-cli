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

// NewCmdFirewall provides the CLI command functionality for Firewall
func NewCmdFirewall(base *cli.Base) *cobra.Command {
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "firewall",
		Short:   "Access firewall commands",
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
		Short:   "Commands to access firewall group functions",
		Aliases: []string{"g"},
	}

	// Group List
	groupList := &cobra.Command{
		Use:     "list",
		Short:   "List all firewall groups",
		Aliases: []string{"l"},
		Run: func(cmd *cobra.Command, args []string) {
			o.Base.Options = utils.GetPaging(cmd)

			groups, meta, err := o.listGroups()
			if err != nil {
				printer.Error(fmt.Errorf("error retrieving firewall group list : %v", err))
				os.Exit(1)
			}

			data := &FirewallGroupsPrinter{Groups: groups, Meta: meta}
			o.Base.Printer.Display(data, nil)
		},
	}

	groupList.Flags().StringP("cursor", "c", "", "(optional) cursor for paging.")
	groupList.Flags().IntP(
		"per-page",
		"p",
		utils.PerPageDefault,
		fmt.Sprintf("(optional) Number of items requested per page. Default is %d and Max is 500.", utils.PerPageDefault),
	)

	// Group Get
	groupGet := &cobra.Command{
		Use:   "get <Firewall Group ID>",
		Short: "Get firewall group",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a firewall group ID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			group, err := o.getGroup()
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

			grp, err := o.createGroup()
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

			if err := o.updateGroup(); err != nil {
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
			if err := o.deleteGroup(); err != nil {
				printer.Error(fmt.Errorf("error deleting firewall group : %v", err))
				os.Exit(1)
			}

			o.Base.Printer.Display(printer.Info("firewall group has been deleted"), nil)
		},
	}

	group.AddCommand(
		groupList,
		groupGet,
		groupCreate,
		groupUpdate,
		groupDelete,
	)

	// Rule
	rule := &cobra.Command{
		Use:     "rule",
		Short:   "Commands to access firewall rules",
		Aliases: []string{"r"},
		Long:    ruleLong,
		Example: ruleExample,
	}

	// Rule List
	ruleList := &cobra.Command{
		Use:     "list <firewall group ID>",
		Short:   "List all firewall rules",
		Aliases: []string{"l"},
		Long:    ruleListLong,
		Example: ruleListExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a firewall group ID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			o.Base.Options = utils.GetPaging(cmd)

			rules, meta, err := o.listRules()
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
			rule, err := o.getRule()
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
		Short:   "Create a firewall rule",
		Aliases: []string{"c"},
		Long:    ruleCreateLong,
		Example: ruleCreateExample,
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

			rule, err := o.createRule()
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
		Short:   "Delete a firewall rule",
		Aliases: []string{"d", "destroy"},
		Long:    ruleDeleteLong,
		Example: ruleDeleteExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide a firewall group ID and firewall rule number")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := o.deleteRule(); err != nil {
				printer.Error(fmt.Errorf("error deleting firewall rule : %v", err))
				os.Exit(1)
			}

			o.Base.Printer.Display(printer.Info("firewall rule deleted"), nil)
		},
	}

	rule.AddCommand(
		ruleList,
		ruleGet,
		ruleCreate,
		ruleDelete,
	)

	cmd.AddCommand(
		group,
		rule,
	)

	return cmd
}

type options struct {
	Base     *cli.Base
	GroupReq *govultr.FirewallGroupReq
	RuleReq  *govultr.FirewallRuleReq
}

// listGroups ...
func (o *options) listGroups() ([]govultr.FirewallGroup, *govultr.Meta, error) {
	groups, meta, _, err := o.Base.Client.FirewallGroup.List(o.Base.Context, o.Base.Options)
	return groups, meta, err
}

// getGroup ...
func (o *options) getGroup() (*govultr.FirewallGroup, error) {
	group, _, err := o.Base.Client.FirewallGroup.Get(o.Base.Context, o.Base.Args[0])
	return group, err
}

// createGroup ...
func (o *options) createGroup() (*govultr.FirewallGroup, error) {
	group, _, err := o.Base.Client.FirewallGroup.Create(o.Base.Context, o.GroupReq)
	return group, err
}

// updateGroup ...
func (o *options) updateGroup() error {
	return o.Base.Client.FirewallGroup.Update(o.Base.Context, o.Base.Args[0], o.GroupReq)
}

// deleteGroup ...
func (o *options) deleteGroup() error {
	return o.Base.Client.FirewallGroup.Delete(o.Base.Context, o.Base.Args[0])
}

// listRules ...
func (o *options) listRules() ([]govultr.FirewallRule, *govultr.Meta, error) {
	rules, meta, _, err := o.Base.Client.FirewallRule.List(o.Base.Context, o.Base.Args[0], o.Base.Options)
	return rules, meta, err
}

// getRule ...
func (o *options) getRule() (*govultr.FirewallRule, error) {
	id, errIn := strconv.Atoi(o.Base.Args[1])
	if errIn != nil {
		return nil, fmt.Errorf("unable to convert int to string : %v", errIn)
	}

	rule, _, err := o.Base.Client.FirewallRule.Get(o.Base.Context, o.Base.Args[0], id)
	return rule, err
}

// createRule ...
func (o *options) createRule() (*govultr.FirewallRule, error) {
	rule, _, err := o.Base.Client.FirewallRule.Create(o.Base.Context, o.Base.Args[0], o.RuleReq)
	return rule, err
}

// deleteRule ...
func (o *options) deleteRule() error {
	id, err := strconv.Atoi(o.Base.Args[1])
	if err != nil {
		return fmt.Errorf("unable to convert int to string : %v", err)
	}
	return o.Base.Client.FirewallRule.Delete(o.Base.Context, o.Base.Args[0], id)
}
