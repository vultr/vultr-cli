// Package loadbalancer provides CLI command to control load balancers
package loadbalancer

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
)

var (
	long    = `Get commands available to Load Balancers`
	example = `
	# Full example
	vultr-cli load-balancer
	`

	listLong    = `Get all load balancers on your Vultr account`
	listExample = `
	# Full example
	vultr-cli load-balancer list

	# Full example with paging
	vultr-cli load-balancer list --per-page=1 --cursor="bmV4dF9fQU1T"

	# Shortened with alias commands
	vultr-cli lb l

	# Summarized view
	vultr-cli load-balancer list --summarize
	`

	createLong    = `Create a new Load Balancer with the desired settings`
	createExample = `
	# Full example
	vultr-cli load-balancer create --region="lax" --balancing-algorithm="roundrobin" --label="Example Load Balancer" \
		--port=80 --check-interval=10 --healthy-threshold=15

	You must pass --region; other arguments are optional

	#Shortened example with aliases
	vultr-cli lb c -r="lax" -b="roundrobin" -l="Example Load Balancer" -p=80 -c=10

	#Full example with attached VPC
	vultr-cli load-balancer create --region="lax"  --label="Example Load Balancer with VPC" \
		--vpc="e951822b-10b2-4c5e-b333-bf38033e7175" --balancing-algorithm="leastconn"
	`
	updateLong    = `Update a Load Balancer with the desired settings`
	updateExample = `
	# Full example
	vultr-cli load-balancer update 57539f6f-66a2-4580-936b-d0af934bce5d --label="Updated Load Balancer Label" \
		--balancing-algorithm="leastconn" --unhealthy-threshold=20

	#Shortened example with aliases
	vultr-cli lb u 57539f6f-66a2-4580-936b-d0af934bce5d -l="Updated Load Balancer Label" -b="leastconn" -u=20

	#Full example with attached VPC
	vultr-cli load-balancer update 57539f6f-66a2-4580-936b-d0af934bce5d --vpc="bff36707-977e-4357-8f30-bef3339155cc"
	`
)

// NewCmdLoadBalancer provides the CLI command for load balancers
func NewCmdLoadBalancer(base *cli.Base) *cobra.Command { //nolint:funlen,gocyclo
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "load-balancer",
		Short:   "Commands to managed load balancers",
		Long:    long,
		Example: example,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			utils.SetOptions(o.Base, cmd, args)
			if !o.Base.HasAuth {
				return errors.New(utils.APIKeyError)
			}
			return nil
		},
	}

	// List
	list := &cobra.Command{
		Use:     "list <Load Balancer ID>",
		Short:   "List load balancers",
		Aliases: []string{"l"},
		Long:    listLong,
		Example: listExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)

			summarize, errSu := cmd.Flags().GetBool("summarize")
			if errSu != nil {
				return fmt.Errorf("error parsing flag 'summarize' for load balancer list : %v", errSu)
			}

			lbs, meta, err := o.list()
			if err != nil {
				return fmt.Errorf("error getting load balancer : %v", err)
			}

			var data printer.ResourceOutput
			if summarize {
				data = &LBsSummaryPrinter{LBs: lbs, Meta: meta}
			} else {
				data = &LBsPrinter{LBs: lbs, Meta: meta}
			}

			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	list.Flags().StringP("cursor", "c", "", "(optional) cursor for paging.")
	list.Flags().IntP(
		"per-page",
		"p",
		utils.PerPageDefault,
		fmt.Sprintf(
			"(optional) Number of items requested per page. Default is %d and Max is 500.",
			utils.PerPageDefault,
		),
	)
	list.Flags().BoolP("summarize", "", false, "(optional) Summarize the list output. One line per load balancer.")

	// Get
	get := &cobra.Command{
		Use:     "get <Load Balancer ID>",
		Short:   "Retrieve a load balancer",
		Aliases: []string{"g"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a load balancer ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			lb, err := o.get()
			if err != nil {
				return fmt.Errorf("error getting load balancer : %v", err)
			}

			o.Base.Printer.Display(&LBPrinter{LB: lb}, nil)

			return nil
		},
	}

	// Create
	create := &cobra.Command{
		Use:     "create",
		Short:   "Create a load balancer",
		Aliases: []string{"c"},
		Long:    createLong,
		Example: createExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			region, errRg := cmd.Flags().GetString("region")
			if errRg != nil {
				return fmt.Errorf("error parsing flag 'region' for load balancer create : %v", errRg)
			}

			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				return fmt.Errorf("error parsing flag 'label' for load balancer create : %v", errLa)
			}

			algorithm, errAl := cmd.Flags().GetString("balancing-algorithm")
			if errAl != nil {
				return fmt.Errorf("error parsing flag 'balancing-algorithm' for load balancer create : %v", errAl)
			}

			sslRedirect, errSs := cmd.Flags().GetBool("ssl-redirect")
			if errSs != nil {
				return fmt.Errorf("error parsing flag 'ssl-redirect' for load balancer create : %v", errSs)
			}

			proxyProtocol, errPr := cmd.Flags().GetBool("proxy-protocol")
			if errPr != nil {
				return fmt.Errorf("error parsing flag 'proxy-protocol' for load balancer create : %v", errPr)
			}

			cookieName, errCo := cmd.Flags().GetString("cookie-name")
			if errCo != nil {
				return fmt.Errorf("error parsing flag 'cookie-name' for load balancer create : %v", errCo)
			}

			vpc, errVp := cmd.Flags().GetString("vpc")
			if errVp != nil {
				return fmt.Errorf("error parsing flag 'vpc' for load balancer create : %v", errVp)
			}

			rulesInForward, errFw := cmd.Flags().GetStringArray("forwarding-rules")
			if errFw != nil {
				return fmt.Errorf("error parsing flag 'forwarding-rules' for load balancer create : %v", errFw)
			}

			rulesInFire, errFi := cmd.Flags().GetStringArray("firewall-rules")
			if errFi != nil {
				return fmt.Errorf("error parsing flag 'firewall-rules' for load balancer create : %v", errFi)
			}

			protocol, errPo := cmd.Flags().GetString("protocol")
			if errPo != nil {
				return fmt.Errorf("error parsing flag 'protocol' for load balancer create : %v", errPo)
			}

			port, errPo := cmd.Flags().GetInt("port")
			if errPo != nil {
				return fmt.Errorf("error parsing flag 'port' for load balancer create : %v", errPo)
			}

			path, errPa := cmd.Flags().GetString("path")
			if errPa != nil {
				return fmt.Errorf("error parsing flag 'path' for load balancer create : %v", errPa)
			}

			checkInterval, errCh := cmd.Flags().GetInt("check-interval")
			if errCh != nil {
				return fmt.Errorf("error parsing flag 'check-interval' for load balancer create : %v", errCh)
			}

			responseTimeout, errRe := cmd.Flags().GetInt("response-timeout")
			if errRe != nil {
				return fmt.Errorf("error parsing flag 'response-timeout' for load balancer create : %v", errRe)
			}

			unhealthyThreshold, errUn := cmd.Flags().GetInt("unhealthy-threshold")
			if errUn != nil {
				return fmt.Errorf("error parsing flag 'unhealthy-threshold' for load balancer create : %v", errUn)
			}

			healthyThreshold, errHe := cmd.Flags().GetInt("healthy-threshold")
			if errHe != nil {
				return fmt.Errorf("error parsing flag 'healthy-threshold' for load balancer create : %v", errHe)
			}

			privateKey, errPi := cmd.Flags().GetString("private-key")
			if errPi != nil {
				return fmt.Errorf("error parsing flag 'private-key' for load balancer create : %v", errPi)
			}

			certificate, errCe := cmd.Flags().GetString("certificate")
			if errCe != nil {
				return fmt.Errorf("error parsing flag 'certificate' for load balancer create : %v", errCe)
			}

			certificateChain, errCr := cmd.Flags().GetString("certificate-chain")
			if errCr != nil {
				return fmt.Errorf("error parsing flag 'certificate-chain' for load balancer create : %v", errCr)
			}

			instances, errIn := cmd.Flags().GetStringSlice("instances")
			if errIn != nil {
				return fmt.Errorf("error parsing flag 'instances' for load balancer create : %v", errIn)
			}

			nodes, errNo := cmd.Flags().GetInt("nodes")
			if errNo != nil {
				return fmt.Errorf("error parsing flag 'nodes' for load balancer create : %v", errNo)
			}

			o.CreateReq = &govultr.LoadBalancerReq{
				Region:             region,
				Label:              label,
				VPC:                &vpc,
				ProxyProtocol:      &proxyProtocol,
				SSLRedirect:        &sslRedirect,
				BalancingAlgorithm: algorithm,
				Nodes:              nodes,
				Instances:          instances,
				HealthCheck: &govultr.HealthCheck{
					Port:               port,
					Protocol:           protocol,
					Path:               path,
					CheckInterval:      checkInterval,
					ResponseTimeout:    responseTimeout,
					UnhealthyThreshold: unhealthyThreshold,
					HealthyThreshold:   healthyThreshold,
				},
				SSL: &govultr.SSL{
					PrivateKey:  privateKey,
					Certificate: certificate,
					Chain:       certificateChain,
				},
			}

			if cmd.Flags().Changed("cookie-name") {
				o.CreateReq.StickySessions = &govultr.StickySessions{CookieName: cookieName}
			}

			if len(rulesInForward) > 0 {
				rulesFo, errFo := formatForwardingRules(rulesInForward)
				if errFo != nil {
					return fmt.Errorf("error creating load balancer : %v", errFo)
				}

				if len(rulesFo) > 0 {
					o.CreateReq.ForwardingRules = rulesFo
				}
			}

			if len(rulesInFire) > 0 {
				rulesFi, errFi := formatFirewallRules(rulesInFire)
				if errFi != nil {
					return fmt.Errorf("error creating load balancer : %v", errFi)
				}

				if len(rulesFi) > 0 {
					o.CreateReq.FirewallRules = rulesFi
				}
			}

			lb, err := o.create()
			if err != nil {
				return fmt.Errorf("error creating load balancer : %v", err)
			}

			o.Base.Printer.Display(&LBPrinter{LB: lb}, nil)

			return nil
		},
	}

	create.Flags().StringP("region", "r", "", "region id you wish to have the load balancer created in")
	if err := create.MarkFlagRequired("region"); err != nil {
		fmt.Printf("error marking load-balancer create 'region' flag required: %v", err)
		os.Exit(1)
	}

	create.Flags().StringP(
		"balancing-algorithm",
		"b",
		"roundrobin",
		"(optional) balancing algorithm that determines server selection | roundrobin or leastconn",
	)
	create.Flags().BoolP(
		"ssl-redirect",
		"s",
		false,
		`(optional) if true, this will redirect HTTP traffic to HTTPS.
		You must have an HTTPS rule and SSL certificate installed on the load balancer to enable this option.`,
	)
	create.Flags().BoolP(
		"proxy-protocol",
		"p",
		false,
		"(optional) if true, you must configure backend nodes to accept Proxy protocol.",
	)
	create.Flags().StringArrayP(
		"forwarding-rules",
		"f",
		[]string{},
		`(optional) a comma-separated, key-value pair list of forwarding rules. Use - between each new rule.
		E.g: "frontend_port:80,frontend_protocol:http,backend_port:80,backend_protocol:http-frontend_port:81,
		frontend_protocol:http,backend_port:81,backend_protocol:http"`,
	)
	create.Flags().StringP(
		"vpc",
		"v",
		"",
		`(optional) the VPC ID to attach to your load balancer. 
When not provided, load balancer defaults to public network.`,
	)

	create.Flags().StringArrayP(
		"firewall-rules",
		"",
		[]string{},
		`(optional) a comma-separated, key-value pair list of firewall rules. Use - between each new rule.
		E.g: "port:80,ip_type:v4,source:0.0.0.0/0-port:8080,ip_type:v4,source:1.1.1.1/4"`,
	)

	create.Flags().String("protocol", "http", "(optional) the protocol to use for health checks. | https, http, tcp")
	create.Flags().Int("port", 80, "(optional) the port to use for health checks.") //nolint:mnd
	create.Flags().String("path", "/", "(optional) HTTP Path to check. only applies if protocol is HTTP or HTTPS.")
	create.Flags().IntP("check-interval", "c", 15, "(optional) interval between health checks.")      //nolint:mnd
	create.Flags().IntP("response-timeout", "t", 15, "(optional) timeout before health check fails.") //nolint:mnd

	create.Flags().IntP(
		"unhealthy-threshold",
		"u",
		15, //nolint:mnd
		"(optional) number times a check must fail before becoming unhealthy.",
	)

	create.Flags().Int(
		"healthy-threshold",
		15, //nolint:mnd
		"(optional) number times a check must succeed before returning to healthy status.",
	)

	create.Flags().String("cookie-name", "", "(optional) the cookie name to make sticky.")

	create.Flags().String("private-key", "", "(optional) the private key component for a ssl certificate.")
	create.Flags().String("certificate", "", "(optional) the SSL certificate.")
	create.Flags().String("certificate-chain", "", "(optional) the certificate chain for a ssl certificate.")

	create.Flags().StringP("label", "l", "", "(optional) the label for your load balancer.")
	create.Flags().StringSliceP(
		"instances",
		"i",
		[]string{},
		"(optional) an array of instances IDs that you want attached to the load balancer.",
	)

	create.Flags().Int(
		"nodes",
		1,
		"(optional) The number of nodes to add to the load balancer (1-99), must be an odd number",
	)

	// Update
	update := &cobra.Command{
		Use:     "update <Load Balancer ID>",
		Short:   "Update a load balancer",
		Aliases: []string{"u"},
		Long:    updateLong,
		Example: updateExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a load balancer ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			label, errLa := cmd.Flags().GetString("label")
			if errLa != nil {
				return fmt.Errorf("error parsing flag 'label' for load balancer update : %v", errLa)
			}

			algorithm, errAl := cmd.Flags().GetString("balancing-algorithm")
			if errAl != nil {
				return fmt.Errorf("error parsing flag 'balancing-algorithm' for load balancer update : %v", errAl)
			}

			sslRedirect, errSs := cmd.Flags().GetBool("ssl-redirect")
			if errSs != nil {
				return fmt.Errorf("error parsing flag 'ssl-redirect' for load balancer update : %v", errSs)
			}

			proxyProtocol, errPr := cmd.Flags().GetBool("proxy-protocol")
			if errPr != nil {
				return fmt.Errorf("error parsing flag 'proxy-protocol' for load balancer update : %v", errPr)
			}

			cookieName, errCo := cmd.Flags().GetString("cookie-name")
			if errCo != nil {
				return fmt.Errorf("error parsing flag 'cookie-name' for load balancer update : %v", errCo)
			}

			vpc, errVp := cmd.Flags().GetString("vpc")
			if errVp != nil {
				return fmt.Errorf("error parsing flag 'vpc' for load balancer update : %v", errVp)
			}

			rulesInForward, errFw := cmd.Flags().GetStringArray("forwarding-rules")
			if errFw != nil {
				return fmt.Errorf("error parsing flag 'forwarding-rules' for load balancer update : %v", errFw)
			}

			rulesInFire, errFi := cmd.Flags().GetStringArray("firewall-rules")
			if errFi != nil {
				return fmt.Errorf("error parsing flag 'firewall-rules' for load balancer update : %v", errFi)
			}

			protocol, errPo := cmd.Flags().GetString("protocol")
			if errPo != nil {
				return fmt.Errorf("error parsing flag 'protocol' for load balancer update : %v", errPo)
			}

			port, errPo := cmd.Flags().GetInt("port")
			if errPo != nil {
				return fmt.Errorf("error parsing flag 'port' for load balancer update : %v", errPo)
			}

			path, errPa := cmd.Flags().GetString("path")
			if errPa != nil {
				return fmt.Errorf("error parsing flag 'path' for load balancer update : %v", errPa)
			}

			checkInterval, errCh := cmd.Flags().GetInt("check-interval")
			if errCh != nil {
				return fmt.Errorf("error parsing flag 'check-interval' for load balancer update : %v", errCh)
			}

			responseTimeout, errRe := cmd.Flags().GetInt("response-timeout")
			if errRe != nil {
				return fmt.Errorf("error parsing flag 'response-timeout' for load balancer update : %v", errRe)
			}

			unhealthyThreshold, errUn := cmd.Flags().GetInt("unhealthy-threshold")
			if errUn != nil {
				return fmt.Errorf("error parsing flag 'unhealthy-threshold' for load balancer update : %v", errUn)
			}

			healthyThreshold, errHe := cmd.Flags().GetInt("healthy-threshold")
			if errHe != nil {
				return fmt.Errorf("error parsing flag 'healthy-threshold' for load balancer update : %v", errHe)
			}

			privateKey, errPi := cmd.Flags().GetString("private-key")
			if errPi != nil {
				return fmt.Errorf("error parsing flag 'private-key' for load balancer update : %v", errPi)
			}

			certificate, errCe := cmd.Flags().GetString("certificate")
			if errCe != nil {
				return fmt.Errorf("error parsing flag 'certificate' for load balancer update : %v", errCe)
			}

			certificateChain, errCr := cmd.Flags().GetString("certificate-chain")
			if errCr != nil {
				return fmt.Errorf("error parsing flag 'certificate-chain' for load balancer update : %v", errCr)
			}

			instances, errIn := cmd.Flags().GetStringSlice("instances")
			if errIn != nil {
				return fmt.Errorf("error parsing flag 'instances' for load balancer update : %v", errIn)
			}

			nodes, errNo := cmd.Flags().GetInt("nodes")
			if errNo != nil {
				return fmt.Errorf("error parsing flag 'nodes' for load balancer update : %v", errNo)
			}

			o.UpdateReq = &govultr.LoadBalancerReq{}

			if len(rulesInForward) > 0 {
				rulesFo, errFo := formatForwardingRules(rulesInForward)
				if errFo != nil {
					return fmt.Errorf("error updating load balancer : %v", errFo)
				}

				if len(rulesFo) > 0 {
					o.UpdateReq.ForwardingRules = rulesFo
				}
			}

			if len(rulesInFire) > 0 {
				rulesFi, errFi := formatFirewallRules(rulesInFire)
				if errFi != nil {
					return fmt.Errorf("error updating load balancer : %v", errFi)
				}

				if len(rulesFi) > 0 {
					o.UpdateReq.FirewallRules = rulesFi
				}
			}

			// Health
			if port != 0 || protocol != "" || path != "" || checkInterval != 0 || responseTimeout != 0 || unhealthyThreshold != 0 || healthyThreshold != 0 { //nolint: lll
				o.UpdateReq.HealthCheck = &govultr.HealthCheck{}
			}

			if port != 0 {
				o.UpdateReq.HealthCheck.Port = port
			}

			if protocol != "" {
				o.UpdateReq.HealthCheck.Protocol = protocol
			}

			if path != "" {
				o.UpdateReq.HealthCheck.Path = path
			}

			if checkInterval != 0 {
				o.UpdateReq.HealthCheck.CheckInterval = checkInterval
			}

			if responseTimeout != 0 {
				o.UpdateReq.HealthCheck.ResponseTimeout = responseTimeout
			}

			if unhealthyThreshold != 0 {
				o.UpdateReq.HealthCheck.UnhealthyThreshold = unhealthyThreshold
			}

			if healthyThreshold != 0 {
				o.UpdateReq.HealthCheck.HealthyThreshold = healthyThreshold
			}

			// SSL
			if privateKey != "" && certificate != "" {
				o.UpdateReq.SSL = &govultr.SSL{
					PrivateKey:  privateKey,
					Certificate: certificate,
					Chain:       certificateChain,
				}
			}

			// Generic Info
			if cmd.Flags().Changed("label") {
				o.UpdateReq.Label = label
			}

			if cmd.Flags().Changed("vpc") {
				o.UpdateReq.VPC = govultr.StringToStringPtr(vpc)
			}

			if cmd.Flags().Changed("proxy-protocol") {
				o.UpdateReq.ProxyProtocol = &proxyProtocol
			}

			if cmd.Flags().Changed("ssl-redirect") {
				o.UpdateReq.SSLRedirect = &sslRedirect
			}

			if cmd.Flags().Changed("cookie-name") {
				o.UpdateReq.StickySessions = &govultr.StickySessions{
					CookieName: cookieName,
				}
			}

			if cmd.Flags().Changed("balancing-algorithm") {
				o.UpdateReq.BalancingAlgorithm = algorithm
			}

			if len(instances) > 0 {
				o.UpdateReq.Instances = instances
			}

			if cmd.Flags().Changed("nodes") {
				o.UpdateReq.Nodes = nodes
			}

			if err := o.update(); err != nil {
				return fmt.Errorf("error updating load balancer : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Load balancer has been updated"), nil)

			return nil
		},
	}

	update.Flags().StringP(
		"balancing-algorithm",
		"b",
		"roundrobin",
		"(optional) balancing algorithm that determines server selection | roundrobin or leastconn",
	)
	update.Flags().BoolP(
		"ssl-redirect",
		"s",
		false,
		`(optional) if true, this will redirect HTTP traffic to HTTPS. You must have an HTTPS rule
		and SSL certificate installed on the load balancer to enable this option.`,
	)
	update.Flags().BoolP(
		"proxy-protocol",
		"p",
		false,
		"(optional) if true, you must configure backend nodes to accept Proxy protocol.",
	)
	update.Flags().StringArrayP(
		"forwarding-rules",
		"f",
		[]string{},
		`(optional) a comma-separated, key-value pair list of forwarding rules. Use - between each new rule.
		E.g: "frontend_port:80,frontend_protocol:http,backend_port:80,backend_protocol:http-frontend_port:81,
		frontend_protocol:http,backend_port:81,backend_protocol:http"`,
	)
	update.Flags().StringArrayP(
		"firewall-rules",
		"",
		[]string{},
		`(optional) a comma-separated, key-value pair list of firewall rules. Use - between each new rule.
		E.g: "port:80,ip_type:v4,source:0.0.0.0/0-port:8080,ip_type:v4,source:1.1.1.1/4"`,
	)
	update.Flags().StringP("vpc", "v", "", "(optional) the VPC ID to attach to your load balancer.")

	update.Flags().String("protocol", "", "(optional) the protocol to use for health checks. | https, http, tcp")
	update.Flags().Int("port", 0, "(optional) the port to use for health checks.")
	update.Flags().String("path", "", "(optional) HTTP Path to check. only applies if protocol is HTTP or HTTPS.")
	update.Flags().IntP("check-interval", "c", 0, "(optional) interval between health checks.")
	update.Flags().IntP("response-timeout", "t", 0, "(optional) timeout before health check fails.")
	update.Flags().IntP(
		"unhealthy-threshold",
		"u",
		0,
		"(optional) number times a check must fail before becoming unhealthy.",
	)
	update.Flags().Int(
		"healthy-threshold",
		0,
		"(optional) number times a check must succeed before returning to healthy status.",
	)

	update.Flags().String("cookie-name", "", "(optional) the cookie name to make sticky.")

	update.Flags().String("private-key", "", "(optional) the private key component for a ssl certificate.")
	update.Flags().String("certificate", "", "(optional) the SSL certificate.")
	update.Flags().String("certificate-chain", "", "(optional) the certificate chain for a ssl certificate.")

	update.Flags().StringP("label", "l", "", "(optional) the label for your load balancer.")
	update.Flags().StringSliceP(
		"instances",
		"i",
		[]string{},
		"(optional) an array of instances IDs that you want attached to the load balancer.",
	)

	update.Flags().Int(
		"nodes",
		1,
		"(optional) The number of nodes to add to the load balancer (1-99), must be an odd number",
	)

	// Delete
	del := &cobra.Command{
		Use:     "delete <Load Balancer ID>",
		Short:   "Delete a load balancer",
		Aliases: []string{"destroy", "d"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a load balancer ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.del(); err != nil {
				return fmt.Errorf("error deleting load balancer : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Load balancer has been deleted"), nil)

			return nil
		},
	}

	// Forwarding Rules
	forwarding := &cobra.Command{
		Use:   "forwarding",
		Short: "Commands to manage forwarding rules on a load balancer",
	}

	// List Forwarding Rules
	listForwardingRules := &cobra.Command{
		Use:   "list <Load Balancer ID>",
		Short: "List all forwarding rules",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a load balancer ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			rules, meta, err := o.listForwardingRules()
			if err != nil {
				return fmt.Errorf("error listing load balancer forwarding rules : %v", err)
			}

			data := &LBRulesPrinter{LBRules: rules, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	listForwardingRules.Flags().StringP("cursor", "c", "", "(optional) cursor for paging.")
	listForwardingRules.Flags().IntP(
		"per-page",
		"p",
		utils.PerPageDefault,
		fmt.Sprintf(
			"(optional) Number of items requested per page. Default is %d and Max is 500.",
			utils.PerPageDefault,
		),
	)

	// Get Forwarding Rule
	getForwardingRule := &cobra.Command{
		Use:   "get <Load Balancer ID> <Rule ID>",
		Short: "Get a forwarding rule",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide a load balancer ID and a rule ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			rule, err := o.getForwardingRule()
			if err != nil {
				return fmt.Errorf("error getting load balancer forwarding rule : %v", err)
			}

			data := &LBRulePrinter{LBRule: *rule}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Create Forwarding Rule
	createForwardingRule := &cobra.Command{
		Use:   "create <Load Balancer ID>",
		Short: "Create a forwarding rule",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a load balancer ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			frontProtocol, errFr := cmd.Flags().GetString("frontend-protocol")
			if errFr != nil {
				return fmt.Errorf("error parsing flag 'frontend-protocol' for forwarding rule create : %v", errFr)
			}

			backProtocol, errBa := cmd.Flags().GetString("backend-protocol")
			if errBa != nil {
				return fmt.Errorf("error parsing flag 'backend-protocol' for forwarding rule create : %v", errBa)
			}

			frontPort, errFp := cmd.Flags().GetInt("frontend-port")
			if errFp != nil {
				return fmt.Errorf("error parsing flag 'frontend-port' for forwarding rule create : %v", errFp)
			}

			backPort, errBp := cmd.Flags().GetInt("backend-port")
			if errBp != nil {
				return fmt.Errorf("error parsing flag 'backend-port' for forwarding rule create : %v", errBp)
			}

			o.RuleCreateReq = &govultr.ForwardingRule{
				FrontendProtocol: frontProtocol,
				FrontendPort:     frontPort,
				BackendProtocol:  backProtocol,
				BackendPort:      backPort,
			}

			rule, err := o.createForwardingRule()
			if err != nil {
				return fmt.Errorf("error creating load balancer forwarding rule : %v", err)
			}

			data := &LBRulePrinter{LBRule: *rule}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	createForwardingRule.Flags().String(
		"frontend-protocol",
		"http",
		"the protocol on the Load Balancer to forward to the backend. | HTTP, HTTPS, TCP",
	)
	if err := createForwardingRule.MarkFlagRequired("frontend-protocol"); err != nil {
		fmt.Printf("error marking load-balancer rule create 'frontend-protocol' flag required: %v", err)
		os.Exit(1)
	}

	createForwardingRule.Flags().String(
		"backend-protocol",
		"http",
		"the protocol destination on the backend server. | HTTP, HTTPS, TCP",
	)
	if err := createForwardingRule.MarkFlagRequired("backend-protocol"); err != nil {
		fmt.Printf("error marking load-balancer rule create 'backend-protocol' flag required: %v", err)
		os.Exit(1)
	}

	createForwardingRule.Flags().Int(
		"frontend-port",
		80, //nolint:mnd
		"the port number on the Load Balancer to forward to the backend.",
	)
	if err := createForwardingRule.MarkFlagRequired("frontend-port"); err != nil {
		fmt.Printf("error marking load-balancer rule create 'frontend-port' flag required: %v", err)
		os.Exit(1)
	}

	createForwardingRule.Flags().Int(
		"backend-port",
		80, //nolint:mnd
		"the port number destination on the backend server.",
	)
	if err := createForwardingRule.MarkFlagRequired("backend-port"); err != nil {
		fmt.Printf("error marking load-balancer rule create 'backend-port' flag required: %v", err)
		os.Exit(1)
	}

	// Delete Forwarding Rule
	deleteForwardingRule := &cobra.Command{
		Use:   "delete <Load Balancer ID> <Rule ID>",
		Short: "Delete a forwarding rule",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide a load balancer ID and a rule ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.deleteForwardingRule(); err != nil {
				return fmt.Errorf("error deleting load balancer forwarding rule : %v", err)
			}

			o.Base.Printer.Display(printer.Info("Forwarding rule has been deleted"), nil)

			return nil
		},
	}

	forwarding.AddCommand(
		listForwardingRules,
		getForwardingRule,
		createForwardingRule,
		deleteForwardingRule,
	)

	// Firewall
	firewall := &cobra.Command{
		Use:   "firewall",
		Short: "Commands to retrieve firewall rules on a load balancer",
	}

	// List Firewall Rules
	listFirewallRules := &cobra.Command{
		Use:   "list <Load Balancer ID>",
		Short: "List all firewall rules",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a load balancer ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			rules, meta, err := o.listFirewallRules()
			if err != nil {
				return fmt.Errorf("error listing load balancer firewall rules : %v", err)
			}

			data := &FWRulesPrinter{Rules: rules, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	listFirewallRules.Flags().StringP("cursor", "c", "", "(optional) cursor for paging.")
	listFirewallRules.Flags().IntP(
		"per-page",
		"p",
		utils.PerPageDefault,
		fmt.Sprintf(
			"(optional) Number of items requested per page. Default is %d and Max is 500.",
			utils.PerPageDefault,
		),
	)

	// Get Firewall Rule
	getFirewallRule := &cobra.Command{
		Use:   "get <Load Balancer ID> <Rule ID>",
		Short: "Get a firewall rule",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide a load balancer ID and a rule ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			rule, err := o.getFirewallRule()
			if err != nil {
				return fmt.Errorf("error getting load balancer firewall rule : %v", err)
			}

			data := &FWRulePrinter{Rule: *rule}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	firewall.AddCommand(
		listFirewallRules,
		getFirewallRule,
	)

	cmd.AddCommand(
		list,
		get,
		create,
		update,
		del,
		forwarding,
		firewall,
	)

	return cmd
}

type options struct {
	Base          *cli.Base
	CreateReq     *govultr.LoadBalancerReq
	UpdateReq     *govultr.LoadBalancerReq
	RuleCreateReq *govultr.ForwardingRule
}

func (o *options) list() ([]govultr.LoadBalancer, *govultr.Meta, error) {
	lbs, meta, _, err := o.Base.Client.LoadBalancer.List(o.Base.Context, o.Base.Options)
	return lbs, meta, err
}

func (o *options) get() (*govultr.LoadBalancer, error) {
	lb, _, err := o.Base.Client.LoadBalancer.Get(o.Base.Context, o.Base.Args[0])
	return lb, err
}

func (o *options) create() (*govultr.LoadBalancer, error) {
	lb, _, err := o.Base.Client.LoadBalancer.Create(o.Base.Context, o.CreateReq)
	return lb, err
}

func (o *options) update() error {
	return o.Base.Client.LoadBalancer.Update(o.Base.Context, o.Base.Args[0], o.UpdateReq)
}

func (o *options) del() error {
	return o.Base.Client.LoadBalancer.Delete(o.Base.Context, o.Base.Args[0])
}

func (o *options) listForwardingRules() ([]govultr.ForwardingRule, *govultr.Meta, error) {
	rs, meta, _, err := o.Base.Client.LoadBalancer.ListForwardingRules(o.Base.Context, o.Base.Args[0], o.Base.Options)
	return rs, meta, err
}

func (o *options) getForwardingRule() (*govultr.ForwardingRule, error) {
	r, _, err := o.Base.Client.LoadBalancer.GetForwardingRule(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
	return r, err
}

func (o *options) createForwardingRule() (*govultr.ForwardingRule, error) {
	r, _, err := o.Base.Client.LoadBalancer.CreateForwardingRule(o.Base.Context, o.Base.Args[0], o.RuleCreateReq)
	return r, err
}

func (o *options) deleteForwardingRule() error {
	return o.Base.Client.LoadBalancer.DeleteForwardingRule(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
}

func (o *options) listFirewallRules() ([]govultr.LBFirewallRule, *govultr.Meta, error) {
	rs, meta, _, err := o.Base.Client.LoadBalancer.ListFirewallRules(o.Base.Context, o.Base.Args[0], o.Base.Options)
	return rs, meta, err
}

func (o *options) getFirewallRule() (*govultr.LBFirewallRule, error) {
	r, _, err := o.Base.Client.LoadBalancer.GetFirewallRule(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
	return r, err
}

// ======================================

// formatFirewallRules parses forwarding rules into proper format
func formatFirewallRules(rules []string) ([]govultr.LBFirewallRule, error) {
	var formattedList []govultr.LBFirewallRule
	rulesList := strings.Split(rules[0], "-")

	for i := range rulesList {
		rule := govultr.LBFirewallRule{}
		fwRule := strings.Split(rulesList[i], ",")

		if len(fwRule) != 3 {
			return nil, fmt.Errorf("unable to format firewall rules. each rule must include ip_type, source, and port")
		}

		for j := range fwRule {
			ruleKeyVal := strings.Split(fwRule[j], ":")

			if len(ruleKeyVal) != 2 {
				return nil, fmt.Errorf("invalid firewall rule format")
			}

			field := ruleKeyVal[0]
			val := ruleKeyVal[1]

			switch {
			case field == "ip_type":
				rule.IPType = val
			case field == "port":
				port, errCon := strconv.Atoi(val)
				if errCon != nil {
					return nil, fmt.Errorf("unable to parse firewall rule port value")
				}
				rule.Port = port
			case field == "source":
				rule.Source = val
			}
		}

		formattedList = append(formattedList, rule)
	}

	return formattedList, nil
}

// formatForwardingRules parses forwarding rules into proper format
func formatForwardingRules(rules []string) ([]govultr.ForwardingRule, error) {
	var formattedList []govultr.ForwardingRule
	var rulePartNum = 4
	rulesList := strings.Split(rules[0], "-")

	for i := range rulesList {
		rule := govultr.ForwardingRule{}
		fwRule := strings.Split(rulesList[i], ",")

		if len(fwRule) != rulePartNum {
			return nil, fmt.Errorf(
				"unable to format forwarding rules. each rule must include frontend and backend ports and protocols",
			)
		}

		for j := range fwRule {
			ruleKeyVal := strings.Split(fwRule[j], ":")

			if len(ruleKeyVal) != 2 {
				return nil, fmt.Errorf("invalid forwarding rule format")
			}

			field := ruleKeyVal[0]
			val := ruleKeyVal[1]

			switch {
			case field == "frontend_protocol":
				rule.FrontendProtocol = val
			case field == "frontend_port":
				port, errCon := strconv.Atoi(val)
				if errCon != nil {
					return nil, fmt.Errorf("unable to parse fowarding rule frontend port value")
				}
				rule.FrontendPort = port
			case field == "backend_protocol":
				rule.BackendProtocol = val
			case field == "backend_port":
				port, errCon := strconv.Atoi(val)
				if errCon != nil {
					return nil, fmt.Errorf("unable to parse fowarding rule backend port value")
				}
				rule.BackendPort = port
			}
		}

		formattedList = append(formattedList, rule)
	}

	return formattedList, nil
}
