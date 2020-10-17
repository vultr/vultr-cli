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
	"github.com/vultr/vultr-cli/cmd/printer"
)

// LoadBalancer represents the load-balancer command
func LoadBalancer() *cobra.Command {

	lbCmd := &cobra.Command{
		Use:     "load-balancer",
		Aliases: []string{"lb"},
		Short:   "load balancer commands",
		Long:    `load-balancer is used to interact with the load-balancer api`,
	}

	lbCmd.AddCommand(lbCreate, lbDelete, lbGet, lbList, lbUpdate)

	// Create
	lbCreate.Flags().StringP("region", "r", "", "region id you wish to have the load balancer created in")
	lbCreate.MarkFlagRequired("region")

	lbCreate.Flags().StringP("balancing-algorithm", "b", "roundrobin", "(optional) balancing algorithm that determines server selection | roundrobin or leastconn")
	lbCreate.Flags().StringP("ssl-redirect", "s", "", "(optional) if true, this will redirect HTTP traffic to HTTPS. You must have an HTTPS rule and SSL certificate installed on the load balancer to enable this option.")
	lbCreate.Flags().StringP("proxy-protocol", "p", "", "(optional) if true, you must configure backend nodes to accept Proxy protocol.")
	lbCreate.Flags().StringArrayP("forwarding-rules", "f", []string{}, "(optional) a comma-separated, key-value pair list of forwarding rules. Use \"\" between each new rule. E.g: `frontend_port:80,frontend_protocol:http,backend_port:80,backend_protocol:http`")

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
	lbList.Flags().IntP("per-page", "p", 25, "(optional) number of items requested per page. default and max are 25.")

	// Update
	lbUpdate.Flags().StringP("balancing-algorithm", "b", "roundrobin", "(optional) balancing algorithm that determines server selection | roundrobin or leastconn")
	lbUpdate.Flags().StringP("ssl-redirect", "s", "", "(optional) if true, this will redirect HTTP traffic to HTTPS. You must have an HTTPS rule and SSL certificate installed on the load balancer to enable this option.")
	lbUpdate.Flags().StringP("proxy-protocol", "p", "", "(optional) if true, you must configure backend nodes to accept Proxy protocol.")
	lbUpdate.Flags().StringArrayP("forwarding-rules", "f", []string{}, "(optional) a comma-separated, key-value pair list of forwarding rules. Use \"\" between each new rule. E.g: `frontend_port:80,frontend_protocol:http,backend_port:80,backend_protocol:http`")

	lbUpdate.Flags().String("protocol", "http", "(optional) the protocol to use for health checks. | https, http, tcp")
	lbUpdate.Flags().Int("port", 80, "(optional) the port to use for health checks.")
	lbUpdate.Flags().String("path", "/", "(optional) HTTP Path to check. only applies if protocol is HTTP or HTTPS.")
	lbUpdate.Flags().IntP("check-interval", "c", 15, "(optional) interval between health checks.")
	lbUpdate.Flags().IntP("response-timeout", "t", 15, "(optional) timeout before health check fails.")
	lbUpdate.Flags().IntP("unhealthy-threshold", "u", 15, "(optional) number times a check must fail before becoming unhealthy.")
	lbUpdate.Flags().Int("healthy-threshold", 15, "(optional) number times a check must succeed before returning to healthy status.")

	lbUpdate.Flags().String("cookie-name", "", "(optional) the cookie name to make sticky.")

	lbUpdate.Flags().String("private-key", "", "(optional) the private key component for a ssl certificate.")
	lbUpdate.Flags().String("certificate", "", "(optional) the SSL certificate.")
	lbUpdate.Flags().String("certificate-chain", "", "(optional) the certificate chain for a ssl certificate.")

	lbUpdate.Flags().StringP("label", "l", "", "(optional) the label for your load balancer.")
	lbUpdate.Flags().StringSliceP("instances", "i", []string{}, "(optional) an array of instances IDs that you want attached to the load balancer.")

	// Rules SubCommands
	rulesCmd := &cobra.Command{
		Use:   "rule",
		Short: "create/delete/list rules for an load balancer",
		Long:  ``,
	}

	// rule list
	ruleList.Flags().StringP("cursor", "c", "", "(optional) cursor for paging.")
	ruleList.Flags().IntP("per-page", "p", 25, "(optional) number of items requested per page. default and max are 25.")

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
	Use:   "create",
	Short: "create a load balancer",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		label, _ := cmd.Flags().GetString("label")

		algorithm, _ := cmd.Flags().GetString("balancing-algorithm")
		sslRedirect, _ := cmd.Flags().GetString("ssl-redirect")
		proxyProtocol, _ := cmd.Flags().GetString("proxy-protocol")

		fwRules, _ := cmd.Flags().GetStringArray("forwarding-rules")

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
				CookieName:            cookieName,
				StickySessionsEnabled: "yes",
			}
		}

		rules, err := formatFWRules(fwRules)
		if err != nil {
			fmt.Printf("error creating load balancer : %v\n", err)
			os.Exit(1)
		}

		if len(rules) > 0 {
			options.ForwardingRules = rules
		}

		if sslRedirect == "yes" {
			options.SSLRedirect = true
		}

		if proxyProtocol == "yes" {
			options.ProxyProtocol = true
		}

		if len(instances) > 0 {
			options.Instances = instances
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
	Use:   "update <loadBalancerID>",
	Short: "updates a load balancer",
	Long:  ``,
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

		fwRules, _ := cmd.Flags().GetStringArray("forwarding-rules")

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

		rules, err := formatFWRules(fwRules)
		if err != nil {
			fmt.Printf("error updating load balancer : %v\n", err)
			os.Exit(1)
		}

		if len(rules) > 0 {
			options.ForwardingRules = rules
		}

		// Health
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

		if proxyProtocol == "yes" {
			options.ProxyProtocol = true
		} else if proxyProtocol == "no" {
			options.ProxyProtocol = false
		}

		if sslRedirect == "yes" {
			options.SSLRedirect = true
		} else if sslRedirect == "no" {
			options.SSLRedirect = false
		}

		if cookieName != "" {
			options.StickySessions.CookieName = cookieName
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

// formatFWRules parses forwarding rules into proper format
func formatFWRules(rules []string) ([]govultr.ForwardingRule, error) {
	var formattedList []govultr.ForwardingRule
	rulesList := strings.Split(rules[0], ",,")

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
