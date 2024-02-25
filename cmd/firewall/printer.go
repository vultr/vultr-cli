package firewall

import (
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
)

// FirewallGroupsPrinter ...
type FirewallGroupsPrinter struct {
	Groups []govultr.FirewallGroup `json:"firewall_groups"`
	Meta   *govultr.Meta           `json:"meta"`
}

// JSON ...
func (f *FirewallGroupsPrinter) JSON() []byte {
	return printer.MarshalObject(f, "json")
}

// YAML ...
func (f *FirewallGroupsPrinter) YAML() []byte {
	return printer.MarshalObject(f, "yaml")
}

// Columns ...
func (f *FirewallGroupsPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"DATE CREATED",
		"DATE MODIFIED",
		"INSTANCE COUNT",
		"RULE COUNT",
		"MAX RULE COUNT",
		"DESCRIPTION",
	}}
}

// Data ...
func (f *FirewallGroupsPrinter) Data() [][]string {
	if len(f.Groups) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range f.Groups {
		data = append(data, []string{
			f.Groups[i].ID,
			f.Groups[i].DateCreated,
			f.Groups[i].DateModified,
			strconv.Itoa(f.Groups[i].InstanceCount),
			strconv.Itoa(f.Groups[i].RuleCount),
			strconv.Itoa(f.Groups[i].MaxRuleCount),
			f.Groups[i].Description,
		})
	}

	return data
}

// Paging ...
func (f *FirewallGroupsPrinter) Paging() [][]string {
	return printer.NewPagingFromMeta(f.Meta).Compose()
}

// ======================================

// FirewallGroupPrinter ...
type FirewallGroupPrinter struct {
	Group govultr.FirewallGroup `json:"firewall_group"`
}

// JSON ...
func (f *FirewallGroupPrinter) JSON() []byte {
	return printer.MarshalObject(f, "json")
}

// YAML ...
func (f *FirewallGroupPrinter) YAML() []byte {
	return printer.MarshalObject(f, "yaml")
}

// Columns ...
func (f *FirewallGroupPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"DATE CREATED",
		"DATE MODIFIED",
		"INSTANCE COUNT",
		"RULE COUNT",
		"MAX RULE COUNT",
		"DESCRIPTION",
	}}
}

// Data ...
func (f *FirewallGroupPrinter) Data() [][]string {
	return [][]string{0: {
		f.Group.ID,
		f.Group.DateCreated,
		f.Group.DateModified,
		strconv.Itoa(f.Group.InstanceCount),
		strconv.Itoa(f.Group.RuleCount),
		strconv.Itoa(f.Group.MaxRuleCount),
		f.Group.Description,
	}}
}

// Paging ...
func (f *FirewallGroupPrinter) Paging() [][]string {
	return nil
}

// ======================================

// FirewallRulesPrinter ...
type FirewallRulesPrinter struct {
	Rules []govultr.FirewallRule `json:"firewall_rules"`
	Meta  *govultr.Meta          `json:"meta"`
}

// JSON ...
func (f *FirewallRulesPrinter) JSON() []byte {
	return printer.MarshalObject(f, "json")
}

// YAML ...
func (f *FirewallRulesPrinter) YAML() []byte {
	return printer.MarshalObject(f, "yaml")
}

// Columns ...
func (f *FirewallRulesPrinter) Columns() [][]string {
	return [][]string{0: {
		"RULE NUMBER",
		"ACTION",
		"TYPE",
		"PROTOCOL",
		"PORT",
		"NETWORK",
		"SOURCE",
		"NOTES",
	}}
}

// Data ...
func (f *FirewallRulesPrinter) Data() [][]string {
	if len(f.Rules) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range f.Rules {
		data = append(data, []string{
			strconv.Itoa(f.Rules[i].ID),
			f.Rules[i].Action,
			f.Rules[i].IPType,
			f.Rules[i].Protocol,
			f.Rules[i].Port,
			utils.FormatFirewallNetwork(f.Rules[i].Subnet, f.Rules[i].SubnetSize),
			utils.GetFirewallSource(f.Rules[i].Source),
			f.Rules[i].Notes,
		})
	}

	return data
}

// Paging ...
func (f *FirewallRulesPrinter) Paging() [][]string {
	return printer.NewPagingFromMeta(f.Meta).Compose()
}

// ======================================

// FirewallRulePrinter ...
type FirewallRulePrinter struct {
	Rule govultr.FirewallRule `json:"firewall_rule"`
}

// JSON ...
func (f *FirewallRulePrinter) JSON() []byte {
	return printer.MarshalObject(f, "json")
}

// YAML ...
func (f *FirewallRulePrinter) YAML() []byte {
	return printer.MarshalObject(f, "yaml")
}

// Columns ...
func (f *FirewallRulePrinter) Columns() [][]string {
	return [][]string{0: {
		"RULE NUMBER",
		"ACTION",
		"TYPE",
		"PROTOCOL",
		"PORT",
		"NETWORK",
		"SOURCE",
		"NOTES",
	}}
}

// Data ...
func (f *FirewallRulePrinter) Data() [][]string {
	return [][]string{0: {
		strconv.Itoa(f.Rule.ID),
		f.Rule.Action,
		f.Rule.IPType,
		f.Rule.Protocol,
		f.Rule.Port,
		utils.FormatFirewallNetwork(f.Rule.Subnet, f.Rule.SubnetSize),
		utils.GetFirewallSource(f.Rule.Source),
		f.Rule.Notes,
	}}
}

// Paging ...
func (f *FirewallRulePrinter) Paging() [][]string {
	return nil
}
