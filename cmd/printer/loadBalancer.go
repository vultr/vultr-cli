package printer

import (
	"github.com/vultr/govultr/v3"
)

func LoadBalancerList(loadbalancer []govultr.LoadBalancer, meta *govultr.Meta) {
	for _, lb := range loadbalancer {
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
		display(columns{lb.HealthCheck.Protocol, lb.HealthCheck.Port, lb.HealthCheck.Path, lb.HealthCheck.CheckInterval, lb.HealthCheck.ResponseTimeout, lb.HealthCheck.UnhealthyThreshold, lb.HealthCheck.HealthyThreshold})

		display(columns{" "})
		display(columns{"GENERIC INFO"})
		display(columns{"BALANCING ALGORITHM", "SSL REDIRECT", "COOKIE NAME", "PROXY PROTOCOL", "VPC"})
		display(columns{lb.GenericInfo.BalancingAlgorithm, *lb.GenericInfo.SSLRedirect, lb.GenericInfo.StickySessions.CookieName, *lb.GenericInfo.ProxyProtocol, lb.GenericInfo.VPC})

		display(columns{" "})
		display(columns{"FORWARDING RULES"})
		display(columns{"RULEID", "FRONTEND PROTOCOL", "FRONTEND PORT", "BACKEND PROTOCOL", "BACKEND PORT"})
		for _, r := range lb.ForwardingRules {
			display(columns{r.RuleID, r.FrontendProtocol, r.FrontendPort, r.BackendProtocol, r.BackendPort})
		}

		display(columns{" "})
		display(columns{"FIREWALL RULES"})
		display(columns{"RULEID", "PORT", "SOURCE", "IP_TYPE"})
		for _, r := range lb.FirewallRules {
			display(columns{r.RuleID, r.Port, r.Source, r.IPType})
		}
		if len(lb.FirewallRules) < 1 {
			display(columns{"-", "-", "-"})
		}

		display(columns{"---------------------------"})
	}

	Meta(meta)
	flush()
}

func LoadBalancer(lb *govultr.LoadBalancer) {
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
	display(columns{lb.HealthCheck.Protocol, lb.HealthCheck.Port, lb.HealthCheck.Path, lb.HealthCheck.CheckInterval, lb.HealthCheck.ResponseTimeout, lb.HealthCheck.UnhealthyThreshold, lb.HealthCheck.HealthyThreshold})

	display(columns{" "})
	display(columns{"GENERIC INFO"})
	display(columns{"BALANCING ALGORITHM", "SSL REDIRECT", "COOKIE NAME", "PROXY PROTOCOL", "VPC"})
	display(columns{lb.GenericInfo.BalancingAlgorithm, *lb.GenericInfo.SSLRedirect, lb.GenericInfo.StickySessions.CookieName, *lb.GenericInfo.ProxyProtocol, lb.GenericInfo.VPC})

	display(columns{" "})
	display(columns{"FORWARDING RULES"})
	display(columns{"RULEID", "FRONTEND PROTOCOL", "FRONTEND PORT", "BACKEND PROTOCOL", "BACKEND PORT"})
	for _, r := range lb.ForwardingRules {
		display(columns{r.RuleID, r.FrontendProtocol, r.FrontendPort, r.BackendProtocol, r.BackendPort})
	}

	display(columns{" "})
	display(columns{"FIREWALL RULES"})
	display(columns{"RULEID", "PORT", "SOURCE", "IP_TYPE"})
	for _, r := range lb.FirewallRules {
		display(columns{r.RuleID, r.Port, r.Source, r.IPType})
	}
	if len(lb.FirewallRules) < 1 {
		display(columns{"-", "-", "-"})
	}
	flush()
}

func LoadBalancerRuleList(rules []govultr.ForwardingRule, meta *govultr.Meta) {
	display(columns{"RULEID", "FRONTEND PROTOCOL", "FRONTEND PORT", "BACKEND PROTOCOL", "BACKEND PORT"})

	for _, r := range rules {
		display(columns{r.RuleID, r.FrontendProtocol, r.FrontendPort, r.BackendProtocol, r.BackendPort})
	}

	Meta(meta)
	flush()
}

func LoadBalancerRule(rule *govultr.ForwardingRule) {
	display(columns{"RULEID", "FRONTEND PROTOCOL", "FRONTEND PORT", "BACKEND PROTOCOL", "BACKEND PORT"})
	display(columns{rule.RuleID, rule.FrontendProtocol, rule.FrontendPort, rule.BackendProtocol, rule.BackendPort})

	flush()
}

func LoadBalancerFWRuleList(rules []govultr.LBFirewallRule, meta *govultr.Meta) {
	display(columns{"RULEID", "PORT", "SOURCE", "IP_TYPE"})

	for _, r := range rules {
		display(columns{r.RuleID, r.Port, r.Source, r.IPType})
	}

	Meta(meta)
	flush()
}

func LoadBalancerFWRule(rule *govultr.LBFirewallRule) {
	display(columns{"RULEID", "PORT", "SOURCE", "IP_TYPE"})
	display(columns{rule.RuleID, rule.Port, rule.Source, rule.IPType})

	flush()
}

func LoadBalancerListSummary(loadbalancer []govultr.LoadBalancer, meta *govultr.Meta) {
	display(columns{"ID", "LABEL", "STATUS", "REGION", "FORWARD#", "FIREWALL#"})
	for _, lb := range loadbalancer {
		forwardRuleCount := len(lb.ForwardingRules)
		firewallRuleCount := len(lb.FirewallRules)

		display(columns{lb.ID, lb.Label, lb.Status, lb.Region, forwardRuleCount, firewallRuleCount})
	}

	Meta(meta)
	flush()
}
