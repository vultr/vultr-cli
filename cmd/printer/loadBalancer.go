package printer

import (
	"github.com/vultr/govultr"
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
		display(columns{"HAS SSL", lb.SSLInfo})
		display(columns{"INSTANCES", lb.Instances})

		display(columns{" "})
		display(columns{"HEALTH CHECKS", lb.HealthCheck})
		display(columns{"PROTOCOL", "PORT", "PATH", "CHECK INTERVAL", "RESPONSE TIMEOUT", "UNHEALTHY THRESHOLD", "HEALTHY THRESHOLD"})
		display(columns{lb.HealthCheck.Protocol, lb.HealthCheck.Port, lb.HealthCheck.Path, lb.HealthCheck.CheckInterval, lb.HealthCheck.ResponseTimeout, lb.HealthCheck.UnhealthyThreshold, lb.HealthCheck.HealthyThreshold})

		display(columns{" "})
		display(columns{"GENERIC INFO"})
		display(columns{"BALANCING ALGORITHM", "SSL REDIRECT", "COOKIE NAME", "PROXY PROTOCOL"})
		display(columns{lb.GenericInfo.BalancingAlgorithm, lb.GenericInfo.SSLRedirect, lb.GenericInfo.StickySessions.CookieName, lb.GenericInfo.ProxyProtocol})

		display(columns{" "})
		display(columns{"FORWARDING RULES"})
		display(columns{"RULEID", "FRONTEND PROTOCOL", "FRONTEND PORT", "BACKEND PROTOCOL", "BACKEND PORT"})
		for _, r := range lb.ForwardingRules {
			display(columns{r.RuleID, r.FrontendProtocol, r.FrontendPort, r.BackendProtocol, r.BackendPort})
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
	display(columns{"GENERIC INFO", lb.GenericInfo})
	display(columns{"HAS SSL", lb.SSLInfo})
	display(columns{"HEALTH CHECKS", lb.HealthCheck})
	display(columns{"INSTANCES", lb.Instances})

	display(columns{" "})
	display(columns{"HEALTH CHECKS", lb.HealthCheck})
	display(columns{"PROTOCOL", "PORT", "PATH", "CHECK INTERVAL", "RESPONSE TIMEOUT", "UNHEALTHY THRESHOLD", "HEALTHY THRESHOLD"})
	display(columns{lb.HealthCheck.Protocol, lb.HealthCheck.Port, lb.HealthCheck.Path, lb.HealthCheck.CheckInterval, lb.HealthCheck.ResponseTimeout, lb.HealthCheck.UnhealthyThreshold, lb.HealthCheck.HealthyThreshold})

	display(columns{" "})
	display(columns{"GENERIC INFO"})
	display(columns{"BALANCING ALGORITHM", "SSL REDIRECT", "COOKIE NAME", "PROXY PROTOCOL"})
	display(columns{lb.GenericInfo.BalancingAlgorithm, lb.GenericInfo.SSLRedirect, lb.GenericInfo.StickySessions.CookieName, lb.GenericInfo.ProxyProtocol})

	display(columns{" "})
	display(columns{"FORWARDING RULES"})
	display(columns{"RULEID", "FRONTEND PROTOCOL", "FRONTEND PORT", "BACKEND PROTOCOL", "BACKEND PORT"})
	for _, r := range lb.ForwardingRules {
		display(columns{r.RuleID, r.FrontendProtocol, r.FrontendPort, r.BackendProtocol, r.BackendPort})
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
