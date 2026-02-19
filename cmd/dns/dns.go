// Package dns provides the functionality for dns commands in the CLI
package dns

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
)

var (
	dnsLong    = ``
	dnsExample = ``

	createLong    = ``
	createExample = ``

	domainLong    = ``
	domainExample = ``
)

// NewCmdDNS provides the CLI command functionality for DNS
func NewCmdDNS(base *cli.Base) *cobra.Command { //nolint:funlen,gocyclo
	o := &options{Base: base}

	cmd := &cobra.Command{
		Use:     "dns",
		Short:   "Commands to control DNS records",
		Long:    dnsLong,
		Example: dnsExample,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			utils.SetOptions(o.Base, cmd, args)
			if !o.Base.HasAuth {
				return errors.New(utils.APIKeyError)
			}
			return nil
		},
	}

	domain := &cobra.Command{
		Use:     "domain",
		Short:   "DNS domain commands",
		Long:    domainLong,
		Example: domainExample,
	}

	// Domain List
	domainList := &cobra.Command{
		Use:   "list",
		Short: "Get a list of domains",
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)

			dms, meta, err := o.domainList()
			if err != nil {
				return fmt.Errorf("error retrieving domain list : %v", err)
			}

			data := &DNSDomainsPrinter{Domains: dms, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	domainList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	domainList.Flags().IntP(
		"per-page",
		"p",
		utils.PerPageDefault,
		fmt.Sprintf(
			"(optional) Number of items requested per page. Default is %d and Max is 500.",
			utils.PerPageDefault,
		),
	)

	// Domain Get
	domainGet := &cobra.Command{
		Use:   "get <Domain Name>",
		Short: "Get a domain",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			dm, err := o.domainGet()
			if err != nil {
				return fmt.Errorf("error retrieving domain : %v", err)
			}

			data := &DNSDomainPrinter{Domain: *dm}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Domain Create
	domainCreate := &cobra.Command{
		Use:     "create",
		Short:   "Create a domain",
		Long:    createLong,
		Example: createExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			domain, errDo := cmd.Flags().GetString("domain")
			if errDo != nil {
				return fmt.Errorf("error parsing 'domain' flag for domain create : %v", errDo)
			}

			ip, errIP := cmd.Flags().GetString("ip")
			if errIP != nil {
				return fmt.Errorf("error parsing 'ip' flag for domain create : %v", errIP)
			}

			o.DomainCreateReq = &govultr.DomainReq{
				Domain: domain,
				IP:     ip,
			}

			dm, err := o.domainCreate()
			if err != nil {
				return fmt.Errorf("error creating dns domain : %v", err)
			}

			data := &DNSDomainPrinter{Domain: *dm}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	domainCreate.Flags().StringP("domain", "d", "", "name of the domain")
	if err := domainCreate.MarkFlagRequired("domain"); err != nil {
		fmt.Printf("error marking domain create 'domain' flag required: %v", err)
		os.Exit(1)
	}
	domainCreate.Flags().StringP("ip", "i", "", "instance ip you want to assign this domain to")

	// Domain Delete
	domainDelete := &cobra.Command{
		Use:     "delete <Domain Name>",
		Short:   "Delete a domain",
		Aliases: []string{"destroy"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a domain name")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.domainDelete(); err != nil {
				return fmt.Errorf("error delete dns domain : %v", err)
			}

			o.Base.Printer.Display(printer.Info("dns domain has been deleted"), nil)

			return nil
		},
	}

	// Domain DNSSEC Update
	domainDNSSEC := &cobra.Command{
		Use:   "dnssec <Domain Name>",
		Short: "Enable or disable DNSSEC",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a domain name")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			enabled, errEn := cmd.Flags().GetBool("enabled")
			if errEn != nil {
				return fmt.Errorf("error parsing 'enabled' flag for dnssec : %v", errEn)
			}

			disabled, errDi := cmd.Flags().GetBool("disabled")
			if errDi != nil {
				return fmt.Errorf("error parsing 'disabled' flag for dnssec : %v", errDi)
			}

			if cmd.Flags().Changed("enabled") {
				if enabled {
					o.DomainDNSSECEnabled = "enabled"
				} else {
					o.DomainDNSSECEnabled = "disabled"
				}
			}

			if cmd.Flags().Changed("disabled") {
				if disabled {
					o.DomainDNSSECEnabled = "disabled"
				} else {
					o.DomainDNSSECEnabled = "enabled"
				}
			}

			if err := o.domainUpdate(); err != nil {
				return fmt.Errorf("error toggling dnssec : %v", err)
			}

			o.Base.Printer.Display(printer.Info("dns domain DNSSEC has been updated"), nil)

			return nil
		},
	}

	domainDNSSEC.Flags().BoolP("enabled", "e", true, "enable dnssec")
	domainDNSSEC.Flags().BoolP("disabled", "d", true, "disable dnssec")
	domainDNSSEC.MarkFlagsOneRequired("enabled", "disabled")
	domainDNSSEC.MarkFlagsMutuallyExclusive("enabled", "disabled")

	// Domain DNSSEC Info
	domainDNSSECInfo := &cobra.Command{
		Use:   "dnssec-info <Domain Name>",
		Short: "Get DNSSEC info",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a domain name")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := o.domainDNSSECGet()
			if err != nil {
				return fmt.Errorf("error getting domain dnssec info : %v", err)
			}

			data := &DNSSECPrinter{SEC: info}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Domain SOA Info
	domainSOAInfo := &cobra.Command{
		Use:   "soa-info <Domain Name>",
		Short: "Get SOA info",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a domain name")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := o.domainSOAGet()
			if err != nil {
				return fmt.Errorf("error getting domain soa info : %v", err)
			}

			data := &DNSSOAPrinter{SOA: *info}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Domain SOA Update
	domainSOAUpdate := &cobra.Command{
		Use:   "soa-update <Domain Name>",
		Short: "Update SOA for a domain",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a domain name")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ns, errNs := cmd.Flags().GetString("ns-primary")
			if errNs != nil {
				return fmt.Errorf("error parsing 'ns-primary' flag for domain soa : %v", errNs)
			}

			email, errEm := cmd.Flags().GetString("email")
			if errEm != nil {
				return fmt.Errorf("error parsing 'email' flag for domain soa : %v", errEm)
			}

			o.SOAUpdateReq = &govultr.Soa{
				NSPrimary: ns,
				Email:     email,
			}

			if err := o.domainSOAUpdate(); err != nil {
				return fmt.Errorf("error updating domain soa : %v", err)
			}

			o.Base.Printer.Display(printer.Info("domain soa has been updated"), nil)

			return nil
		},
	}

	domainSOAUpdate.Flags().StringP("ns-primary", "n", "", "primary nameserver to store in the SOA record")
	domainSOAUpdate.Flags().StringP("email", "e", "", "administrative email to store in the SOA record")

	domain.AddCommand(
		domainList,
		domainGet,
		domainCreate,
		domainDelete,
		domainDNSSEC,
		domainDNSSECInfo,
		domainSOAInfo,
		domainSOAUpdate,
	)

	// Record
	record := &cobra.Command{
		Use:   "record",
		Short: "Commands to mangage DNS records",
	}

	// Record List
	recordList := &cobra.Command{
		Use:   "list <Domain Name>",
		Short: "List all DNS records",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a domain name")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Base.Options = utils.GetPaging(cmd)

			recs, meta, err := o.recordList()
			if err != nil {
				return fmt.Errorf("error retrieiving domain records : %v", err)
			}

			data := &DNSRecordsPrinter{Records: recs, Meta: meta}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	recordList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	recordList.Flags().IntP(
		"per-page",
		"p",
		utils.PerPageDefault,
		fmt.Sprintf(
			"(optional) Number of items requested per page. Default is %d and Max is 500.",
			utils.PerPageDefault,
		),
	)

	// Record Get
	recordGet := &cobra.Command{
		Use:   "get <Domain Name> <Record ID>",
		Short: "Get a DNS record",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide a domain name and record ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			rec, err := o.recordGet()
			if err != nil {
				return fmt.Errorf("error while getting domain record : %v", err)
			}

			data := &DNSRecordPrinter{Record: *rec}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	// Record Create
	recordCreate := &cobra.Command{
		Use:   "create <Domain Name>",
		Short: "Create a DNS record",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a domain name")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			rType, errTy := cmd.Flags().GetString("type")
			if errTy != nil {
				return fmt.Errorf("error parsing 'type' flag for domain record create : %v", errTy)
			}

			name, errNa := cmd.Flags().GetString("name")
			if errNa != nil {
				return fmt.Errorf("error parsing 'name' flag for domain record create : %v", errNa)
			}

			dt, errDa := cmd.Flags().GetString("data")
			if errDa != nil {
				return fmt.Errorf("error parsing 'data' flag for domain record create : %v", errDa)
			}

			ttl, errTt := cmd.Flags().GetInt("ttl")
			if errTt != nil {
				return fmt.Errorf("error parsing 'ttl' flag for domain record create : %v", errTt)
			}

			priority, errPr := cmd.Flags().GetInt("priority")
			if errPr != nil {
				return fmt.Errorf("error parsing 'priority' flag for domain record create : %v", errPr)
			}

			o.RecordCreateReq = &govultr.DomainRecordCreateReq{
				Name:     name,
				Type:     rType,
				Data:     dt,
				TTL:      ttl,
				Priority: &priority,
			}

			rec, err := o.recordCreate()
			if err != nil {
				return fmt.Errorf("error creating domain record : %v", err)
			}

			data := &DNSRecordPrinter{Record: *rec}
			o.Base.Printer.Display(data, nil)

			return nil
		},
	}

	recordCreate.Flags().StringP("type", "t", "", "type for the record")
	if err := recordCreate.MarkFlagRequired("type"); err != nil {
		fmt.Printf("error marking dns record create 'type' flag required: %v", err)
		os.Exit(1)
	}

	recordCreate.Flags().StringP("name", "n", "", "name of the record")
	if err := recordCreate.MarkFlagRequired("name"); err != nil {
		fmt.Printf("error marking dns record create 'name' flag required: %v", err)
		os.Exit(1)
	}

	recordCreate.Flags().StringP("data", "d", "", "data for the record")
	if err := recordCreate.MarkFlagRequired("data"); err != nil {
		fmt.Printf("error marking dns record create 'data' flag required: %v", err)
		os.Exit(1)
	}

	recordCreate.Flags().IntP("ttl", "l", 0, "ttl for the record")
	recordCreate.Flags().IntP("priority", "p", -1, "only required for MX and SRV")

	// Record Delete
	recordDelete := &cobra.Command{
		Use:     "delete <Domain Name> <Record ID>",
		Short:   "Delete a DNS record",
		Aliases: []string{"destroy"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide a domain name & record ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.recordDelete(); err != nil {
				return fmt.Errorf("error deleting domain record : %v", err)
			}

			o.Base.Printer.Display(printer.Info("domain record has been deleted"), nil)

			return nil
		},
	}

	// Record Update
	recordUpdate := &cobra.Command{
		Use:   "update <Domain Name> <Record ID>",
		Short: "Update DNS record",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide a domain name & record ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			name, errNa := cmd.Flags().GetString("name")
			if errNa != nil {
				return fmt.Errorf("error parsing 'name' flag for domain record update : %v", errNa)
			}

			dt, errDa := cmd.Flags().GetString("data")
			if errDa != nil {
				return fmt.Errorf("error parsing 'data' flag for domain record update : %v", errDa)
			}

			ttl, errTt := cmd.Flags().GetInt("ttl")
			if errTt != nil {
				return fmt.Errorf("error parsing 'ttl' flag for domain record update : %v", errTt)
			}

			priority, errPr := cmd.Flags().GetInt("priority")
			if errPr != nil {
				return fmt.Errorf("error parsing 'priority' flag for domain record update : %v", errPr)
			}

			o.RecordUpdateReq = &govultr.DomainRecordUpdateReq{}

			if cmd.Flags().Changed("name") {
				o.RecordUpdateReq.Name = &name
			}

			if cmd.Flags().Changed("data") {
				o.RecordUpdateReq.Data = dt
			}

			if cmd.Flags().Changed("ttl") {
				o.RecordUpdateReq.TTL = ttl
			}

			if cmd.Flags().Changed("priority") {
				o.RecordUpdateReq.Priority = govultr.IntToIntPtr(priority)
			}

			if err := o.recordUpdate(); err != nil {
				return fmt.Errorf("error updating domain record : %v", errPr)
			}

			o.Base.Printer.Display(printer.Info("domain record has been updated"), nil)

			return nil
		},
	}

	recordUpdate.Flags().StringP("name", "n", "", "name of record")
	recordUpdate.Flags().StringP("data", "d", "", "data for the record")
	recordUpdate.Flags().IntP("ttl", "", 0, "time to live for the record")
	recordUpdate.Flags().IntP("priority", "p", 0, "only required for MX and SRV")

	record.AddCommand(
		recordList,
		recordGet,
		recordCreate,
		recordUpdate,
		recordDelete,
	)

	cmd.AddCommand(
		domain,
		record,
	)

	return cmd
}

type options struct {
	Base                *cli.Base
	DomainCreateReq     *govultr.DomainReq
	DomainDNSSECEnabled string
	SOAUpdateReq        *govultr.Soa
	RecordCreateReq     *govultr.DomainRecordCreateReq
	RecordUpdateReq     *govultr.DomainRecordUpdateReq
}

// domainList ...
func (o *options) domainList() ([]govultr.Domain, *govultr.Meta, error) {
	dms, meta, _, err := o.Base.Client.Domain.List(o.Base.Context, o.Base.Options)
	return dms, meta, err
}

// domainGet ...
func (o *options) domainGet() (*govultr.Domain, error) {
	dm, _, err := o.Base.Client.Domain.Get(o.Base.Context, o.Base.Args[0])
	return dm, err
}

// domainCreate ...
func (o *options) domainCreate() (*govultr.Domain, error) {
	dm, _, err := o.Base.Client.Domain.Create(o.Base.Context, o.DomainCreateReq)
	return dm, err
}

// domainUpdate ...
func (o *options) domainUpdate() error {
	return o.Base.Client.Domain.Update(o.Base.Context, o.Base.Args[0], o.DomainDNSSECEnabled)
}

// domainDelete ...
func (o *options) domainDelete() error {
	return o.Base.Client.Domain.Delete(o.Base.Context, o.Base.Args[0])
}

// domainDNSSECGet ...
func (o *options) domainDNSSECGet() ([]string, error) {
	sec, _, err := o.Base.Client.Domain.GetDNSSec(o.Base.Context, o.Base.Args[0])
	return sec, err
}

// domainSOAGet ...
func (o *options) domainSOAGet() (*govultr.Soa, error) {
	soa, _, err := o.Base.Client.Domain.GetSoa(o.Base.Context, o.Base.Args[0])
	return soa, err
}

// domainSOAUpdate ...
func (o *options) domainSOAUpdate() error {
	return o.Base.Client.Domain.UpdateSoa(o.Base.Context, o.Base.Args[0], o.SOAUpdateReq)
}

// recordList ...
func (o *options) recordList() ([]govultr.DomainRecord, *govultr.Meta, error) {
	rec, meta, _, err := o.Base.Client.DomainRecord.List(o.Base.Context, o.Base.Args[0], o.Base.Options)
	return rec, meta, err
}

// recordGet ...
func (o *options) recordGet() (*govultr.DomainRecord, error) {
	rec, _, err := o.Base.Client.DomainRecord.Get(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
	return rec, err
}

// recordCreate ...
func (o *options) recordCreate() (*govultr.DomainRecord, error) {
	rec, _, err := o.Base.Client.DomainRecord.Create(o.Base.Context, o.Base.Args[0], o.RecordCreateReq)
	return rec, err
}

// recordUpdate ...
func (o *options) recordUpdate() error {
	return o.Base.Client.DomainRecord.Update(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.RecordUpdateReq)
}

// recordDelete ...
func (o *options) recordDelete() error {
	return o.Base.Client.DomainRecord.Delete(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
}
