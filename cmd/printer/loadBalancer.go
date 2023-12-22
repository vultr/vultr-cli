package printer

import (
	"github.com/vultr/govultr/v3"
)

func LoadBalancerList(loadbalancer []govultr.LoadBalancer, meta *govultr.Meta) {
	defer flush()

	if len(loadbalancer) == 0 {
		displayString("No active load balancers")
		return
	}

	for i := range loadbalancer {
		display(columns{"ID", loadbalancer[i].ID})
		display(columns{"DATE CREATED", loadbalancer[i].DateCreated})
		display(columns{"REGION", loadbalancer[i].Region})
		display(columns{"LABEL", loadbalancer[i].Label})
		display(columns{"STATUS", loadbalancer[i].Status})
		display(columns{"IPV4", loadbalancer[i].IPV4})
		display(columns{"IPV6", loadbalancer[i].IPV6})
		display(columns{"HAS SSL", *loadbalancer[i].SSLInfo})
		display(columns{"INSTANCES", loadbalancer[i].Instances})

		display(columns{" "})
		display(columns{"HEALTH CHECKS"})
		display(columns{"PROTOCOL", "PORT", "PATH", "CHECK INTERVAL", "RESPONSE TIMEOUT", "UNHEALTHY THRESHOLD", "HEALTHY THRESHOLD"})
		display(columns{
			loadbalancer[i].HealthCheck.Protocol,
			loadbalancer[i].HealthCheck.Port,
			loadbalancer[i].HealthCheck.Path,
			loadbalancer[i].HealthCheck.CheckInterval,
			loadbalancer[i].HealthCheck.ResponseTimeout,
			loadbalancer[i].HealthCheck.UnhealthyThreshold,
			loadbalancer[i].HealthCheck.HealthyThreshold,
		})

		display(columns{" "})
		display(columns{"GENERIC INFO"})
		display(columns{"BALANCING ALGORITHM", "SSL REDIRECT", "COOKIE NAME", "PROXY PROTOCOL", "VPC"})
		display(columns{
			loadbalancer[i].GenericInfo.BalancingAlgorithm,
			*loadbalancer[i].GenericInfo.SSLRedirect,
			loadbalancer[i].GenericInfo.StickySessions.CookieName,
			*loadbalancer[i].GenericInfo.ProxyProtocol,
			loadbalancer[i].GenericInfo.VPC,
		})

		display(columns{" "})
		display(columns{"FORWARDING RULES"})
		display(columns{"RULEID", "FRONTEND PROTOCOL", "FRONTEND PORT", "BACKEND PROTOCOL", "BACKEND PORT"})
		for _, r := range loadbalancer[i].ForwardingRules {
			display(columns{r.RuleID, r.FrontendProtocol, r.FrontendPort, r.BackendProtocol, r.BackendPort})
		}

		display(columns{" "})
		display(columns{"FIREWALL RULES"})
		display(columns{"RULEID", "PORT", "SOURCE", "IP_TYPE"})
		for _, r := range loadbalancer[i].FirewallRules {
			display(columns{r.RuleID, r.Port, r.Source, r.IPType})
		}
		if len(loadbalancer[i].FirewallRules) < 1 {
			display(columns{"-", "-", "-"})
		}

		display(columns{"---------------------------"})
	}

	Meta(meta)
}

func LoadBalancer(lb *govultr.LoadBalancer) {
	defer flush()

	display(columns{"ID", lb.ID})
	display(columns{"DATE CREATED", lb.DateCreated})
	display(columns{"REGION", lb.Region})
	display(columns{"LABEL", lb.Label})
	display(columns{"STATUS", lb.Status})
	display(columns{"IPV4", lb.IPV4})
	display(columns{"IPV6", lb.IPV6})
	display(columns{"HAS SSL", *lb.SSLInfo})
	display(columns{"INSTANCES", lb.Instances})

	display(columns{" "})
	display(columns{"HEALTH CHECKS"})
	display(columns{"PROTOCOL", "PORT", "PATH", "CHECK INTERVAL", "RESPONSE TIMEOUT", "UNHEALTHY THRESHOLD", "HEALTHY THRESHOLD"})
	display(columns{
		lb.HealthCheck.Protocol,
		lb.HealthCheck.Port,
		lb.HealthCheck.Path,
		lb.HealthCheck.CheckInterval,
		lb.HealthCheck.ResponseTimeout,
		lb.HealthCheck.UnhealthyThreshold,
		lb.HealthCheck.HealthyThreshold,
	})

	display(columns{" "})
	display(columns{"GENERIC INFO"})
	display(columns{"BALANCING ALGORITHM", "SSL REDIRECT", "COOKIE NAME", "PROXY PROTOCOL", "VPC"})
	display(columns{
		lb.GenericInfo.BalancingAlgorithm,
		*lb.GenericInfo.SSLRedirect,
		lb.GenericInfo.StickySessions.CookieName,
		*lb.GenericInfo.ProxyProtocol,
		lb.GenericInfo.VPC,
	})

	display(columns{" "})
	display(columns{"FORWARDING RULES"})
	display(columns{"RULEID", "FRONTEND PROTOCOL", "FRONTEND PORT", "BACKEND PROTOCOL", "BACKEND PORT"})
	for i := range lb.ForwardingRules {
		display(columns{
			lb.ForwardingRules[i].RuleID,
			lb.ForwardingRules[i].FrontendProtocol,
			lb.ForwardingRules[i].FrontendPort,
			lb.ForwardingRules[i].BackendProtocol,
			lb.ForwardingRules[i].BackendPort,
		})
	}

	display(columns{" "})
	display(columns{"FIREWALL RULES"})
	display(columns{"RULEID", "PORT", "SOURCE", "IP_TYPE"})
	for i := range lb.FirewallRules {
		display(columns{
			lb.FirewallRules[i].RuleID,
			lb.FirewallRules[i].Port,
			lb.FirewallRules[i].Source,
			lb.FirewallRules[i].IPType,
		})
	}
	if len(lb.FirewallRules) < 1 {
		display(columns{"-", "-", "-"})
	}
}

func LoadBalancerRuleList(rules []govultr.ForwardingRule, meta *govultr.Meta) {
	defer flush()

	display(columns{"RULEID", "FRONTEND PROTOCOL", "FRONTEND PORT", "BACKEND PROTOCOL", "BACKEND PORT"})

	if len(rules) == 0 {
		display(columns{"---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range rules {
		display(columns{
			rules[i].RuleID,
			rules[i].FrontendProtocol,
			rules[i].FrontendPort,
			rules[i].BackendProtocol,
			rules[i].BackendPort,
		})
	}

	Meta(meta)
}

func LoadBalancerRule(rule *govultr.ForwardingRule) {
	defer flush()

	display(columns{"RULEID", "FRONTEND PROTOCOL", "FRONTEND PORT", "BACKEND PROTOCOL", "BACKEND PORT"})
	display(columns{rule.RuleID, rule.FrontendProtocol, rule.FrontendPort, rule.BackendProtocol, rule.BackendPort})
}

func LoadBalancerFWRuleList(rules []govultr.LBFirewallRule, meta *govultr.Meta) {
	defer flush()

	display(columns{"RULEID", "PORT", "SOURCE", "IP_TYPE"})

	if len(rules) == 0 {
		display(columns{"---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range rules {
		display(columns{
			rules[i].RuleID,
			rules[i].Port,
			rules[i].Source,
			rules[i].IPType,
		})
	}

	Meta(meta)
}

func LoadBalancerFWRule(rule *govultr.LBFirewallRule) {
	defer flush()

	display(columns{"RULEID", "PORT", "SOURCE", "IP_TYPE"})
	display(columns{rule.RuleID, rule.Port, rule.Source, rule.IPType})
}

func LoadBalancerListSummary(loadbalancer []govultr.LoadBalancer, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "LABEL", "STATUS", "REGION", "INSTANCE#", "FORWARD#", "FIREWALL#"})

	if len(loadbalancer) == 0 {
		display(columns{"---", "---", "---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range loadbalancer {
		forwardRuleCount := len(loadbalancer[i].ForwardingRules)
		firewallRuleCount := len(loadbalancer[i].FirewallRules)
		instanceCount := len(loadbalancer[i].Instances)

		display(columns{
			loadbalancer[i].ID,
			loadbalancer[i].Label,
			loadbalancer[i].Status,
			loadbalancer[i].Region,
			instanceCount,
			forwardRuleCount,
			firewallRuleCount,
		})
	}

	Meta(meta)
}
