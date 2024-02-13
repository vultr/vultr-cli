package loadbalancer

import (
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// LBsPrinter ...
type LBsPrinter struct {
	LBs  []govultr.LoadBalancer `json:"load_balancers"`
	Meta *govultr.Meta          `json:"meta"`
}

// JSON ...
func (l *LBsPrinter) JSON() []byte {
	return printer.MarshalObject(l, "json")
}

// YAML ...
func (l *LBsPrinter) YAML() []byte {
	return printer.MarshalObject(l, "yaml")
}

// Columns ...
func (l *LBsPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (l *LBsPrinter) Data() [][]string {
	if len(l.LBs) == 0 {
		return [][]string{0: {"No active load balancers"}}

	}

	var data [][]string
	for i := range l.LBs {
		data = append(data,
			[]string{"---------------------------"},
			[]string{"ID", l.LBs[i].ID},
			[]string{"DATE CREATED", l.LBs[i].DateCreated},
			[]string{"REGION", l.LBs[i].Region},
			[]string{"LABEL", l.LBs[i].Label},
			[]string{"STATUS", l.LBs[i].Status},
			[]string{"IPV4", l.LBs[i].IPV4},
			[]string{"IPV6", l.LBs[i].IPV6},
			[]string{"HAS SSL", strconv.FormatBool(*l.LBs[i].SSLInfo)},
			[]string{"INSTANCES", printer.ArrayOfStringsToString(l.LBs[i].Instances)},

			[]string{" "},
			[]string{"HEALTH CHECKS"},
			[]string{"PROTOCOL", "PORT", "PATH", "CHECK INTERVAL", "RESPONSE TIMEOUT", "UNHEALTHY THRESHOLD", "HEALTHY THRESHOLD"},
			[]string{
				l.LBs[i].HealthCheck.Protocol,
				strconv.Itoa(l.LBs[i].HealthCheck.Port),
				l.LBs[i].HealthCheck.Path,
				strconv.Itoa(l.LBs[i].HealthCheck.CheckInterval),
				strconv.Itoa(l.LBs[i].HealthCheck.ResponseTimeout),
				strconv.Itoa(l.LBs[i].HealthCheck.UnhealthyThreshold),
				strconv.Itoa(l.LBs[i].HealthCheck.HealthyThreshold),
			},

			[]string{" "},
			[]string{"GENERIC INFO"},
			[]string{"BALANCING ALGORITHM", "SSL REDIRECT", "COOKIE NAME", "PROXY PROTOCOL", "VPC"},
			[]string{
				l.LBs[i].GenericInfo.BalancingAlgorithm,
				strconv.FormatBool(*l.LBs[i].GenericInfo.SSLRedirect),
				l.LBs[i].GenericInfo.StickySessions.CookieName,
				strconv.FormatBool(*l.LBs[i].GenericInfo.ProxyProtocol),
				l.LBs[i].GenericInfo.VPC,
			},

			[]string{" "},
			[]string{"FORWARDING RULES"},
			[]string{"RULEID", "FRONTEND PROTOCOL", "FRONTEND PORT", "BACKEND PROTOCOL", "BACKEND PORT"},
		)

		if len(l.LBs[i].ForwardingRules) == 0 {
			data = append(data, []string{"---", "---", "---", "---", "---"})

		} else {

			for j := range l.LBs[i].ForwardingRules {
				data = append(data,
					[]string{
						l.LBs[i].ForwardingRules[j].RuleID,
						l.LBs[i].ForwardingRules[j].FrontendProtocol,
						strconv.Itoa(l.LBs[i].ForwardingRules[j].FrontendPort),
						l.LBs[i].ForwardingRules[j].BackendProtocol,
						strconv.Itoa(l.LBs[i].ForwardingRules[j].BackendPort),
					},
				)
			}
		}

		data = append(data,
			[]string{" "},
			[]string{"FIREWALL RULES"},
			[]string{"RULEID", "PORT", "SOURCE", "IP_TYPE"},
		)

		if len(l.LBs[i].FirewallRules) == 0 {
			data = append(data, []string{"---", "---", "---", "---"})

		} else {

			for j := range l.LBs[i].FirewallRules {
				data = append(data,
					[]string{
						l.LBs[i].FirewallRules[j].RuleID,
						strconv.Itoa(l.LBs[i].FirewallRules[j].Port),
						l.LBs[i].FirewallRules[j].Source,
						l.LBs[i].FirewallRules[j].IPType,
					},
				)
			}
		}

		data = append(data, []string{" "})
	}

	return data
}

// Paging ...
func (l *LBsPrinter) Paging() [][]string {
	return printer.NewPaging(l.Meta.Total, &l.Meta.Links.Next, &l.Meta.Links.Prev).Compose()
}

// ======================================

// LBPrinter ...
type LBPrinter struct {
	LB *govultr.LoadBalancer `json:"load_balancer"`
}

// JSON ...
func (l *LBPrinter) JSON() []byte {
	return printer.MarshalObject(l, "json")
}

// YAML ...
func (l *LBPrinter) YAML() []byte {
	return printer.MarshalObject(l, "yaml")
}

// Columns ...
func (l *LBPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (l *LBPrinter) Data() [][]string {
	var data [][]string
	data = append(data,
		[]string{"ID", l.LB.ID},
		[]string{"DATE CREATED", l.LB.DateCreated},
		[]string{"REGION", l.LB.Region},
		[]string{"LABEL", l.LB.Label},
		[]string{"STATUS", l.LB.Status},
		[]string{"IPV4", l.LB.IPV4},
		[]string{"IPV6", l.LB.IPV6},
		[]string{"HAS SSL", strconv.FormatBool(*l.LB.SSLInfo)},
		[]string{"INSTANCES", printer.ArrayOfStringsToString(l.LB.Instances)},

		[]string{" "},
		[]string{"HEALTH CHECKS"},
		[]string{"PROTOCOL", "PORT", "PATH", "CHECK INTERVAL", "RESPONSE TIMEOUT", "UNHEALTHY THRESHOLD", "HEALTHY THRESHOLD"},
		[]string{
			l.LB.HealthCheck.Protocol,
			strconv.Itoa(l.LB.HealthCheck.Port),
			l.LB.HealthCheck.Path,
			strconv.Itoa(l.LB.HealthCheck.CheckInterval),
			strconv.Itoa(l.LB.HealthCheck.ResponseTimeout),
			strconv.Itoa(l.LB.HealthCheck.UnhealthyThreshold),
			strconv.Itoa(l.LB.HealthCheck.HealthyThreshold),
		},

		[]string{" "},
		[]string{"GENERIC INFO"},
		[]string{"BALANCING ALGORITHM", "SSL REDIRECT", "COOKIE NAME", "PROXY PROTOCOL", "VPC"},
		[]string{
			l.LB.GenericInfo.BalancingAlgorithm,
			strconv.FormatBool(*l.LB.GenericInfo.SSLRedirect),
			l.LB.GenericInfo.StickySessions.CookieName,
			strconv.FormatBool(*l.LB.GenericInfo.ProxyProtocol),
			l.LB.GenericInfo.VPC,
		},

		[]string{" "},
		[]string{"FORWARDING RULES"},
		[]string{"RULEID", "FRONTEND PROTOCOL", "FRONTEND PORT", "BACKEND PROTOCOL", "BACKEND PORT"},
	)

	if len(l.LB.ForwardingRules) == 0 {
		data = append(data, []string{"---", "---", "---", "---", "---"})

	} else {

		for i := range l.LB.ForwardingRules {
			data = append(data,
				[]string{
					l.LB.ForwardingRules[i].RuleID,
					l.LB.ForwardingRules[i].FrontendProtocol,
					strconv.Itoa(l.LB.ForwardingRules[i].FrontendPort),
					l.LB.ForwardingRules[i].BackendProtocol,
					strconv.Itoa(l.LB.ForwardingRules[i].BackendPort),
				},
			)
		}
	}

	data = append(data,
		[]string{" "},
		[]string{"FIREWALL RULES"},
		[]string{"RULEID", "PORT", "SOURCE", "IP_TYPE"},
	)

	if len(l.LB.FirewallRules) == 0 {
		data = append(data, []string{"---", "---", "---", "---"})

	} else {

		for i := range l.LB.FirewallRules {
			data = append(data,
				[]string{
					l.LB.FirewallRules[i].RuleID,
					strconv.Itoa(l.LB.FirewallRules[i].Port),
					l.LB.FirewallRules[i].Source,
					l.LB.FirewallRules[i].IPType,
				},
			)
		}
	}

	return data
}

// Paging ...
func (l *LBPrinter) Paging() [][]string {
	return nil
}

// ======================================

// LBsSummaryPrinter ...
type LBsSummaryPrinter struct {
	LBs  []govultr.LoadBalancer `json:"load_balancers"`
	Meta *govultr.Meta          `json:"meta"`
}

// JSON ...
func (l *LBsSummaryPrinter) JSON() []byte {
	return printer.MarshalObject(l, "json")
}

// YAML ...
func (l *LBsSummaryPrinter) YAML() []byte {
	return printer.MarshalObject(l, "yaml")
}

// Columns ...
func (l *LBsSummaryPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"LABEL",
		"STATUS",
		"REGION",
		"INSTANCE#",
		"FORWARD#",
		"FIREWALL#",
	}}
}

// Data ...
func (l *LBsSummaryPrinter) Data() [][]string {
	if len(l.LBs) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---", "---", "---"}}

	}

	var data [][]string
	for i := range l.LBs {

		forwardRuleCount := len(l.LBs[i].ForwardingRules)
		firewallRuleCount := len(l.LBs[i].FirewallRules)
		instanceCount := len(l.LBs[i].Instances)

		data = append(data, []string{

			l.LBs[i].ID,
			l.LBs[i].Label,
			l.LBs[i].Status,
			l.LBs[i].Region,
			strconv.Itoa(instanceCount),
			strconv.Itoa(forwardRuleCount),
			strconv.Itoa(firewallRuleCount),
		})
	}

	return data
}

// Paging ...
func (l *LBsSummaryPrinter) Paging() [][]string {
	return printer.NewPaging(l.Meta.Total, &l.Meta.Links.Next, &l.Meta.Links.Prev).Compose()
}

// ======================================

// LBRulesPrinter ...
type LBRulesPrinter struct {
	LBRules []govultr.ForwardingRule `json:"forwarding_rules"`
	Meta    *govultr.Meta            `json:"meta"`
}

// JSON ...
func (l *LBRulesPrinter) JSON() []byte {
	return printer.MarshalObject(l, "json")
}

// YAML ...
func (l *LBRulesPrinter) YAML() []byte {
	return printer.MarshalObject(l, "yaml")
}

// Columns ...
func (l *LBRulesPrinter) Columns() [][]string {
	return [][]string{0: {
		"RULEID",
		"FRONTEND PROTOCOL",
		"FRONTEND PORT",
		"BACKEND PROTOCOL",
		"BACKEND PORT",
	}}
}

// Data ...
func (l *LBRulesPrinter) Data() [][]string {
	if len(l.LBRules) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---"}}

	}

	var data [][]string
	for i := range l.LBRules {
		data = append(data, []string{
			l.LBRules[i].RuleID,
			l.LBRules[i].FrontendProtocol,
			strconv.Itoa(l.LBRules[i].FrontendPort),
			l.LBRules[i].BackendProtocol,
			strconv.Itoa(l.LBRules[i].BackendPort),
		})
	}

	return data
}

// Paging ...
func (l *LBRulesPrinter) Paging() [][]string {
	return printer.NewPaging(l.Meta.Total, &l.Meta.Links.Next, &l.Meta.Links.Prev).Compose()
}

// ======================================

// LBRulePrinter ...
type LBRulePrinter struct {
	LBRule govultr.ForwardingRule `json:"forwarding_rule"`
}

// JSON ...
func (l *LBRulePrinter) JSON() []byte {
	return printer.MarshalObject(l, "json")
}

// YAML ...
func (l *LBRulePrinter) YAML() []byte {
	return printer.MarshalObject(l, "yaml")
}

// Columns ...
func (l *LBRulePrinter) Columns() [][]string {
	return [][]string{0: {
		"RULEID",
		"FRONTEND PROTOCOL",
		"FRONTEND PORT",
		"BACKEND PROTOCOL",
		"BACKEND PORT",
	}}
}

// Data ...
func (l *LBRulePrinter) Data() [][]string {
	return [][]string{0: {
		l.LBRule.RuleID,
		l.LBRule.FrontendProtocol,
		strconv.Itoa(l.LBRule.FrontendPort),
		l.LBRule.BackendProtocol,
		strconv.Itoa(l.LBRule.BackendPort),
	}}
}

// Paging ...
func (l *LBRulePrinter) Paging() [][]string {
	return nil
}

// ======================================

// FWRulesPrinter ...
type FWRulesPrinter struct {
	Rules []govultr.LBFirewallRule `json:"firewall_rules"`
	Meta  *govultr.Meta            `json:"meta"`
}

// JSON ...
func (f *FWRulesPrinter) JSON() []byte {
	return printer.MarshalObject(f, "json")
}

// YAML ...
func (f *FWRulesPrinter) YAML() []byte {
	return printer.MarshalObject(f, "yaml")
}

// Columns ...
func (f *FWRulesPrinter) Columns() [][]string {
	return [][]string{0: {
		"RULEID",
		"PORT",
		"SOURCE",
		"IP_TYPE",
	}}
}

// Data ...
func (f *FWRulesPrinter) Data() [][]string {
	if len(f.Rules) == 0 {
		return [][]string{0: {"---", "---", "---", "---"}}

	}

	var data [][]string
	for i := range f.Rules {
		data = append(data, []string{
			f.Rules[i].RuleID,
			strconv.Itoa(f.Rules[i].Port),
			f.Rules[i].Source,
			f.Rules[i].IPType,
		})

	}

	return data
}

// Paging ...
func (f *FWRulesPrinter) Paging() [][]string {
	return printer.NewPaging(f.Meta.Total, &f.Meta.Links.Next, &f.Meta.Links.Prev).Compose()
}

// ======================================

// FWRulePrinter ...
type FWRulePrinter struct {
	Rule govultr.LBFirewallRule `json:"firewall_rule"`
}

// JSON ...
func (f *FWRulePrinter) JSON() []byte {
	return printer.MarshalObject(f, "json")
}

// YAML ...
func (f *FWRulePrinter) YAML() []byte {
	return printer.MarshalObject(f, "yaml")
}

// Columns ...
func (f *FWRulePrinter) Columns() [][]string {
	return [][]string{0: {
		"RULEID",
		"PORT",
		"SOURCE",
		"IP_TYPE",
	}}
}

// Data ...
func (f *FWRulePrinter) Data() [][]string {
	return [][]string{0: {
		f.Rule.RuleID,
		strconv.Itoa(f.Rule.Port),
		f.Rule.Source,
		f.Rule.IPType,
	}}
}

// Paging ...
func (f *FWRulePrinter) Paging() [][]string {
	return nil
}
