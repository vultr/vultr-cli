// Copyright © 2020 The Vultr-cli Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v2"
	"github.com/vultr/vultr-cli/v2/cmd/printer"
)

var (
	lbLong    = `Get commands available to Load Balancers`
	lbExample = `
	# Full example
	vultr-cli load-balancer
	`
	lbCreateLong    = `Create a new Load Balancer with the desired settings`
	lbCreateExample = `
	# Full example
	vultr-cli load-balancer create --region="lax" --balancing-algorithm="roundrobin" --label="Example Load Balancer" --port=80 --check-interval=10 --healthy-threshold=15

	You must pass --region; other arguments are optional

	#Shortened example with aliases
	vultr-cli lb c -r="lax" -b="roundrobin" -l="Example Load Balancer" -p=80 -c=10 

	#Full example with attached VPC
	vultr-cli load-balancer create --region="lax"  --label="Example Load Balancer with VPC" --vpc="e951822b-10b2-4c5e-b333-bf38033e7175" --balancing-algorithm="leastconn"
	`
	lbUpdateLong    = `Update a Load Balancer with the desired settings`
	lbUpdateExample = `
	# Full example
	vultr-cli load-balancer update 57539f6f-66a2-4580-936b-d0af934bce5d --label="Updated Load Balancer Label" --balancing-algorithm="leastconn" --unhealthy-threshold=20

	#Shortened example with aliases
	vultr-cli lb u 57539f6f-66a2-4580-936b-d0af934bce5d -l="Updated Load Balancer Label" -b="leastconn" -u=20

	#Full example with attached VPC
	vultr-cli load-balancer update 57539f6f-66a2-4580-936b-d0af934bce5d --vpc="bff36707-977e-4357-8f30-bef3339155cc"
	`
)

// LoadBalancer represents the load-balancer command
func LoadBalancer() *cobra.Command {

	lbCmd := &cobra.Command{
		Use:     "load-balancer",
		Aliases: []string{"lb"},
		Short:   "load balancer commands",
		Long:    lbLong,
		Example: lbExample,
	}

	lbCmd.AddCommand(lbCreate, lbDelete, lbGet, lbList, lbUpdate)

	// Create
	lbCreate.Flags().StringP("region", "r", "", "region id you wish to have the load balancer created in")
	lbCreate.MarkFlagRequired("region")

	lbCreate.Flags().StringP("balancing-algorithm", "b", "roundrobin", "(optional) balancing algorithm that determines server selection | roundrobin or leastconn")
	lbCreate.Flags().StringP("ssl-redirect", "s", "", "(optional) if true, this will redirect HTTP traffic to HTTPS. You must have an HTTPS rule and SSL certificate installed on the load balancer to enable this option.")
	lbCreate.Flags().StringP("proxy-protocol", "p", "", "(optional) if true, you must configure backend nodes to accept Proxy protocol.")
	lbCreate.Flags().StringArrayP("forwarding-rules", "f", []string{}, "(optional) a comma-separated, key-value pair list of forwarding rules. Use - between each new rule. E.g: `frontend_port:80,frontend_protocol:http,backend_port:80,backend_protocol:http-frontend_port:81,frontend_protocol:http,backend_port:81,backend_protocol:http`")
	lbCreate.Flags().StringP("private-network", "", "", "(optional) Deprecated: use vpc instead. the private network for your load balancer. When not provided, load balancer defaults to public network.")
	lbCreate.Flags().StringP("vpc", "v", "", "(optional) the VPC ID to attach to your load balancer. When not provided, load balancer defaults to public network.")

	lbCreate.Flags().StringArrayP("firewall-rules", "", []string{}, "(optional) a comma-separated, key-value pair list of firewall rules. Use - between each new rule. E.g: `port:80,ip_type:v4,source:0.0.0.0/0-port:8080,ip_type:v4,source:1.1.1.1/4`")

	lbCreate.Flags().String("protocol", "http", "(optional) the protocol to use for health checks. | https, http, tcp")
	lbCreate.Flags().Int("port", 80, "(optional) the port to use for health checks.")
	lbCreate.Flags().String("path", "/", "(optional) HTTP Path to check. only applies if protocol is HTTP or HTTPS.")
	lbCreate.Flags().IntP("check-interval", "c", 15, "(optional) interval between health checks.")
	lbCreate.Flags().IntP("response-timeout", "t", 15, "(optional) timeout before health check fails.")
	lbCreate.Flags().IntP("unhealthy-threshold", "u", 15, "(optional) number times a check must fail before becoming unhealthy.")
	lbCreate.Flags().Int("healthy-threshold", 15, "(optional) number times a check must succeed before returning to healthy status.")

	lbCreate.Flags().String("cookie-name", "", "(optional) the cookie name to make sticky.")

	lbCreate.Flags().String("private-key", "", "(optional) the private key component for a ssl certificate.")
	lbCreate.Flags().String("certificate", "", "(optional) the SSL certificate.")
	lbCreate.Flags().String("certificate-chain", "", "(optional) the certificate chain for a ssl certificate.")

	lbCreate.Flags().StringP("label", "l", "", "(optional) the label for your load balancer.")
	lbCreate.Flags().StringSliceP("instances", "i", []string{}, "(optional) an array of instances IDs that you want attached to the load balancer.")

	// List
	lbList.Flags().StringP("cursor", "c", "", "(optional) cursor for paging.")
	lbList.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	// Update
	lbUpdate.Flags().StringP("balancing-algorithm", "b", "roundrobin", "(optional) balancing algorithm that determines server selection | roundrobin or leastconn")
	lbUpdate.Flags().StringP("ssl-redirect", "s", "", "(optional) if true, this will redirect HTTP traffic to HTTPS. You must have an HTTPS rule and SSL certificate installed on the load balancer to enable this option.")
	lbUpdate.Flags().StringP("proxy-protocol", "p", "", "(optional) if true, you must configure backend nodes to accept Proxy protocol.")
	lbUpdate.Flags().StringArrayP("forwarding-rules", "f", []string{}, "(optional) a comma-separated, key-value pair list of forwarding rules. Use - between each new rule. E.g: `frontend_port:80,frontend_protocol:http,backend_port:80,backend_protocol:http-frontend_port:81,frontend_protocol:http,backend_port:81,backend_protocol:http`")
	lbUpdate.Flags().StringArrayP("firewall-rules", "", []string{}, "(optional) a comma-separated, key-value pair list of firewall rules. Use - between each new rule. E.g: `port:80,ip_type:v4,source:0.0.0.0/0-port:8080,ip_type:v4,source:1.1.1.1/4`")
	lbUpdate.Flags().StringP("vpc", "v", "", "(optional) the VPC ID to attach to your load balancer.")

	lbUpdate.Flags().String("protocol", "", "(optional) the protocol to use for health checks. | https, http, tcp")
	lbUpdate.Flags().Int("port", 0, "(optional) the port to use for health checks.")
	lbUpdate.Flags().String("path", "", "(optional) HTTP Path to check. only applies if protocol is HTTP or HTTPS.")
	lbUpdate.Flags().IntP("check-interval", "c", 0, "(optional) interval between health checks.")
	lbUpdate.Flags().IntP("response-timeout", "t", 0, "(optional) timeout before health check fails.")
	lbUpdate.Flags().IntP("unhealthy-threshold", "u", 0, "(optional) number times a check must fail before becoming unhealthy.")
	lbUpdate.Flags().Int("healthy-threshold", 0, "(optional) number times a check must succeed before returning to healthy status.")

	lbUpdate.Flags().String("cookie-name", "", "(optional) the cookie name to make sticky.")

	lbUpdate.Flags().String("private-key", "", "(optional) the private key component for a ssl certificate.")
	lbUpdate.Flags().String("certificate", "", "(optional) the SSL certificate.")
	lbUpdate.Flags().String("certificate-chain", "", "(optional) the certificate chain for a ssl certificate.")

	lbUpdate.Flags().StringP("label", "l", "", "(optional) the label for your load balancer.")
	lbUpdate.Flags().StringSliceP("instances", "i", []string{}, "(optional) an array of instances IDs that you want attached to the load balancer.")

	// Forwarding Rules SubCommands
	rulesCmd := &cobra.Command{
		Use:   "rule",
		Short: "create/delete/list forwarding rules for a load balancer",
		Long:  ``,
	}

	// rule list
	ruleList.Flags().StringP("cursor", "c", "", "(optional) cursor for paging.")
	ruleList.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	// Firewall Rules SubCommands
	fwrulesCmd := &cobra.Command{
		Use:   "firewall-rule",
		Short: "get/list firewall rules for a load balancer",
		Long:  ``,
	}

	// firewall rule list
	fwRuleList.Flags().StringP("cursor", "c", "", "(optional) cursor for paging.")
	fwRuleList.Flags().IntP("per-page", "p", 100, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	fwrulesCmd.AddCommand(fwRuleList, fwRuleGet)
	lbCmd.AddCommand(fwrulesCmd)

	// rule create
	ruleCreate.Flags().String("frontend-protocol", "http", "the protocol on the Load Balancer to forward to the backend. | HTTP, HTTPS, TCP")
	ruleCreate.Flags().String("backend-protocol", "http", "the protocol destination on the backend server. | HTTP, HTTPS, TCP")
	ruleCreate.Flags().Int("frontend-port", 80, "the port number on the Load Balancer to forward to the backend.")
	ruleCreate.Flags().Int("backend-port", 80, "the port number destination on the backend server.")

	ruleCreate.MarkFlagRequired("frontend-protocol")
	ruleCreate.MarkFlagRequired("backend-protocol")
	ruleCreate.MarkFlagRequired("frontend-port")
	ruleCreate.MarkFlagRequired("backend-port")

	rulesCmd.AddCommand(ruleCreate, ruleDelete, ruleGet, ruleList)
	lbCmd.AddCommand(rulesCmd)

	return lbCmd
}

var lbCreate = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "create a load balancer",
	Long:    lbCreateLong,
	Example: lbCreateExample,
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		label, _ := cmd.Flags().GetString("label")

		algorithm, _ := cmd.Flags().GetString("balancing-algorithm")
		sslRedirect, _ := cmd.Flags().GetString("ssl-redirect")
		proxyProtocol, _ := cmd.Flags().GetString("proxy-protocol")

		fwRules, _ := cmd.Flags().GetStringArray("forwarding-rules")
		firewallRules, _ := cmd.Flags().GetStringArray("firewall-rules")

		protocol, _ := cmd.Flags().GetString("protocol")
		port, _ := cmd.Flags().GetInt("port")
		path, _ := cmd.Flags().GetString("path")
		checkInterval, _ := cmd.Flags().GetInt("check-interval")
		responseTimeout, _ := cmd.Flags().GetInt("response-timeout")
		unhealthyThreshold, _ := cmd.Flags().GetInt("unhealthy-threshold")
		healthyThreshold, _ := cmd.Flags().GetInt("healthy-threshold")

		privateKey, _ := cmd.Flags().GetString("private-key")
		certificate, _ := cmd.Flags().GetString("certificate")
		certificateChain, _ := cmd.Flags().GetString("certificate-chain")

		cookieName, _ := cmd.Flags().GetString("cookie-name")
		instances, _ := cmd.Flags().GetStringSlice("instances")

		privateNetwork, _ := cmd.Flags().GetString("private-network")
		vpc, _ := cmd.Flags().GetString("vpc")

		healthCheck := &govultr.HealthCheck{
			Protocol:           protocol,
			Path:               path,
			Port:               port,
			CheckInterval:      checkInterval,
			ResponseTimeout:    responseTimeout,
			UnhealthyThreshold: unhealthyThreshold,
			HealthyThreshold:   healthyThreshold,
		}

		ssl := &govultr.SSL{
			PrivateKey:  privateKey,
			Certificate: certificate,
			Chain:       certificateChain,
		}

		options := &govultr.LoadBalancerReq{
			Region:             region,
			Label:              label,
			BalancingAlgorithm: algorithm,
			HealthCheck:        healthCheck,
			SSL:                ssl,
		}

		if cookieName != "" {
			options.StickySessions = &govultr.StickySessions{
				CookieName: cookieName,
			}
		}

		if len(fwRules) > 0 {
			rules, err := formatFWRules(fwRules)
			if err != nil {
				fmt.Printf("error creating load balancer : %v\n", err)
				os.Exit(1)
			}

			if len(rules) > 0 {
				options.ForwardingRules = rules
			}
		}

		if len(firewallRules) > 0 {
			frules, err := formatFirewallRules(firewallRules)
			if err != nil {
				fmt.Printf("error creating load balancer : %v\n", err)
				os.Exit(1)
			}

			if len(frules) > 0 {
				options.FirewallRules = frules
			}
		}

		if sslRedirect == "yes" {
			options.SSLRedirect = govultr.BoolToBoolPtr(true)
		}

		if proxyProtocol == "yes" {
			options.ProxyProtocol = govultr.BoolToBoolPtr(true)
		}

		if len(instances) > 0 {
			options.Instances = instances
		}

		if privateNetwork != "" && vpc == "" {
			options.VPC = govultr.StringToStringPtr(privateNetwork)
		} else {
			options.VPC = govultr.StringToStringPtr(vpc)
		}

		lb, err := client.LoadBalancer.Create(context.Background(), options)
		if err != nil {
			fmt.Printf("error creating load balancer : %v\n", err)
			os.Exit(1)
		}

		printer.LoadBalancer(lb)
	},
}

var lbDelete = &cobra.Command{
	Use:   "delete <loadBalancerID>",
	Short: "deletes a load balancer",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a loadBalancerID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		if err := client.LoadBalancer.Delete(context.Background(), id); err != nil {
			fmt.Printf("error deleting load balancer : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Deleted load balancer")
	},
}

var lbGet = &cobra.Command{
	Use:   "get <loadBalancerID>",
	Short: "retrieves a load balancer",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a loadBalancerID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		lb, err := client.LoadBalancer.Get(context.Background(), id)
		if err != nil {
			fmt.Printf("error getting load balancer : %v\n", err)
			os.Exit(1)
		}

		printer.LoadBalancer(lb)
	},
}

var lbList = &cobra.Command{
	Use:   "list",
	Short: "retrieves a list of active load balancers",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		options := getPaging(cmd)
		list, meta, err := client.LoadBalancer.List(context.Background(), options)
		if err != nil {
			fmt.Printf("error listing load balancers : %v\n", err)
			os.Exit(1)
		}

		printer.LoadBalancerList(list, meta)
	},
}

var lbUpdate = &cobra.Command{
	Use:     "update <loadBalancerID>",
	Aliases: []string{"u"},
	Short:   "updates a load balancer",
	Long:    lbUpdateLong,
	Example: lbUpdateExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a loadBalancerID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		label, _ := cmd.Flags().GetString("label")

		algorithm, _ := cmd.Flags().GetString("balancing-algorithm")
		sslRedirect, _ := cmd.Flags().GetString("ssl-redirect")
		proxyProtocol, _ := cmd.Flags().GetString("proxy-protocol")
		cookieName, _ := cmd.Flags().GetString("cookie-name")
		vpc, _ := cmd.Flags().GetString("vpc")

		fwRules, _ := cmd.Flags().GetStringArray("forwarding-rules")
		firewallRules, _ := cmd.Flags().GetStringArray("firewall-rules")

		protocol, _ := cmd.Flags().GetString("protocol")
		port, _ := cmd.Flags().GetInt("port")
		path, _ := cmd.Flags().GetString("path")
		checkInterval, _ := cmd.Flags().GetInt("check-interval")
		responseTimeout, _ := cmd.Flags().GetInt("response-timeout")
		unhealthyThreshold, _ := cmd.Flags().GetInt("unhealthy-threshold")
		healthyThreshold, _ := cmd.Flags().GetInt("healthy-threshold")

		privateKey, _ := cmd.Flags().GetString("private-key")
		certificate, _ := cmd.Flags().GetString("certificate")
		certificateChain, _ := cmd.Flags().GetString("certificate-chain")

		instances, _ := cmd.Flags().GetStringSlice("instances")

		options := &govultr.LoadBalancerReq{}

		if len(fwRules) > 0 {
			rules, err := formatFWRules(fwRules)
			if err != nil {
				fmt.Printf("error updating load balancer : %v\n", err)
				os.Exit(1)
			}

			if len(rules) > 0 {
				options.ForwardingRules = rules
			}
		}

		if len(firewallRules) > 0 {
			frules, err := formatFirewallRules(firewallRules)
			if err != nil {
				fmt.Printf("error updating load balancer : %v\n", err)
				os.Exit(1)
			}

			if len(frules) > 0 {
				options.FirewallRules = frules
			}
		}

		// Health
		if port != 0 || protocol != "" || path != "" || checkInterval != 0 || responseTimeout != 0 || unhealthyThreshold != 0 || healthyThreshold != 0 {
			options.HealthCheck = &govultr.HealthCheck{}
		}

		if port != 0 {
			options.HealthCheck.Port = port
		}

		if protocol != "" {
			options.HealthCheck.Protocol = protocol
		}

		if path != "" {
			options.HealthCheck.Path = path
		}

		if checkInterval != 0 {
			options.HealthCheck.CheckInterval = checkInterval
		}

		if responseTimeout != 0 {
			options.HealthCheck.ResponseTimeout = responseTimeout
		}

		if unhealthyThreshold != 0 {
			options.HealthCheck.UnhealthyThreshold = unhealthyThreshold
		}

		if healthyThreshold != 0 {
			options.HealthCheck.HealthyThreshold = healthyThreshold
		}

		// SSL
		if privateKey != "" && certificate != "" {
			options.SSL = &govultr.SSL{
				PrivateKey:  privateKey,
				Certificate: certificate,
				Chain:       certificateChain,
			}
		}

		// Generic Info
		if label != "" {
			options.Label = label
		}

		options.VPC = govultr.StringToStringPtr(vpc)

		if proxyProtocol == "yes" {
			options.ProxyProtocol = govultr.BoolToBoolPtr(true)
		} else if proxyProtocol == "no" {
			options.ProxyProtocol = govultr.BoolToBoolPtr(false)
		}

		if sslRedirect == "yes" {
			options.SSLRedirect = govultr.BoolToBoolPtr(true)
		} else if sslRedirect == "no" {
			options.SSLRedirect = govultr.BoolToBoolPtr(false)
		}

		if cookieName != "" {
			options.StickySessions = &govultr.StickySessions{
				CookieName: cookieName,
			}
		}

		if algorithm != "" {
			options.BalancingAlgorithm = algorithm
		}

		if len(instances) > 0 {
			options.Instances = instances
		}

		if err := client.LoadBalancer.Update(context.Background(), id, options); err != nil {
			fmt.Printf("error updating load balancer : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Updated load balancer")
	},
}

var ruleList = &cobra.Command{
	Use:   "list rule <loadBalancerID>",
	Short: "lists a load balancers forwarding rules",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a loadBalancerID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		options := getPaging(cmd)
		rules, meta, err := client.LoadBalancer.ListForwardingRules(context.Background(), id, options)
		if err != nil {
			fmt.Printf("error listing load balancer rules : %v\n", err)
			os.Exit(1)
		}

		printer.LoadBalancerRuleList(rules, meta)
	},
}

var ruleCreate = &cobra.Command{
	Use:   "create rule <loadBalancerID>",
	Short: "creates a load balancer forwarding rule",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a loadBalancerID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		options := &govultr.ForwardingRule{}
		rule, err := client.LoadBalancer.CreateForwardingRule(context.Background(), id, options)
		if err != nil {
			fmt.Printf("error listing load balancer rules : %v\n", err)
			os.Exit(1)
		}

		printer.LoadBalancerRule(rule)
	},
}

var ruleGet = &cobra.Command{
	Use:   "get rule <loadBalancerID> <ruleID>",
	Short: "Gets a load balancer forwarding rule",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a loadBalancerID and ruleID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		ruleID := args[1]
		rule, err := client.LoadBalancer.GetForwardingRule(context.Background(), id, ruleID)
		if err != nil {
			fmt.Printf("error getting load balancer rule : %v\n", err)
			os.Exit(1)
		}

		printer.LoadBalancerRule(rule)
	},
}

var ruleDelete = &cobra.Command{
	Use:     "delete rule <loadBalancerID> <ruleID>",
	Short:   "deletes a load balancer forwarding rule",
	Long:    ``,
	Aliases: []string{"destroy"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a loadBalancerID and ruleID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		ruleID := args[1]

		if err := client.LoadBalancer.DeleteForwardingRule(context.Background(), id, ruleID); err != nil {
			fmt.Printf("error deleting load balancer rule : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Deleted load balancer rule")
	},
}

var fwRuleList = &cobra.Command{
	Use:   "list firewall-rule <loadBalancerID>",
	Short: "lists a load balancers firewall rules",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a loadBalancerID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		options := getPaging(cmd)
		rules, meta, err := client.LoadBalancer.ListFirewallRules(context.Background(), id, options)
		if err != nil {
			fmt.Printf("error listing load balancer firewall rules : %v\n", err)
			os.Exit(1)
		}

		printer.LoadBalancerFWRuleList(rules, meta)
	},
}

var fwRuleGet = &cobra.Command{
	Use:   "get rule <loadBalancerID> <ruleID>",
	Short: "Gets a load balancer firewall rule",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please provide a loadBalancerID and ruleID")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		ruleID := args[1]
		rule, err := client.LoadBalancer.GetFirewallRule(context.Background(), id, ruleID)
		if err != nil {
			fmt.Printf("error getting load balancer rule : %v\n", err)
			os.Exit(1)
		}

		printer.LoadBalancerFWRule(rule)
	},
}

// formatFirewallRules parses forwarding rules into proper format
func formatFirewallRules(rules []string) ([]govultr.LBFirewallRule, error) {
	var formattedList []govultr.LBFirewallRule
	rulesList := strings.Split(rules[0], "-")

	for _, r := range rulesList {
		rule := govultr.LBFirewallRule{}
		fwRule := strings.Split(r, ",")

		if len(fwRule) != 3 {
			return nil, fmt.Errorf("unable to format firewall rules. each rule must include ip_type, source, and port")
		}

		for _, f := range fwRule {
			ruleKeyVal := strings.Split(f, ":")

			if len(ruleKeyVal) != 2 {
				return nil, fmt.Errorf("invalid firewall rule format")
			}

			field := ruleKeyVal[0]
			val := ruleKeyVal[1]

			switch true {
			case field == "ip_type":
				rule.IPType = val
			case field == "port":
				port, _ := strconv.Atoi(val)
				rule.Port = port
			case field == "source":
				rule.Source = val
			}
		}

		formattedList = append(formattedList, rule)
	}

	return formattedList, nil
}

// formatFWRules parses forwarding rules into proper format
func formatFWRules(rules []string) ([]govultr.ForwardingRule, error) {
	var formattedList []govultr.ForwardingRule
	rulesList := strings.Split(rules[0], "-")

	for _, r := range rulesList {
		rule := govultr.ForwardingRule{}
		fwRule := strings.Split(r, ",")

		if len(fwRule) != 4 {
			return nil, fmt.Errorf("unable to format forwarding rules. each rule must include frontend and backend ports and protocols")
		}

		for _, f := range fwRule {
			ruleKeyVal := strings.Split(f, ":")

			if len(ruleKeyVal) != 2 {
				return nil, fmt.Errorf("invalid forwarding rule format")
			}

			field := ruleKeyVal[0]
			val := ruleKeyVal[1]

			switch true {
			case field == "frontend_protocol":
				rule.FrontendProtocol = val
			case field == "frontend_port":
				port, _ := strconv.Atoi(val)
				rule.FrontendPort = port
			case field == "backend_protocol":
				rule.BackendProtocol = val
			case field == "backend_port":
				port, _ := strconv.Atoi(val)
				rule.BackendPort = port
			}
		}

		formattedList = append(formattedList, rule)
	}

	return formattedList, nil
}
