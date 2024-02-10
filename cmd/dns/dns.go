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

type Options struct {
	Base                *cli.Base
	DomainCreateReq     *govultr.DomainReq
	DomainDNSSECEnabled string
	SOAUpdateReq        *govultr.Soa
	RecordReq           *govultr.DomainRecordReq
}

// NewCmdDNS provides the CLI command functionality for DNS
func NewCmdDNS(base *cli.Base) *cobra.Command {
	o := &Options{Base: base}

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
		Short: "Get list of domains",
		Run: func(cmd *cobra.Command, args []string) {
			o.Base.Options = utils.GetPaging(cmd)

			dms, meta, err := o.DomainList()
			if err != nil {
				printer.Error(fmt.Errorf("error retrieving domain list : %v", err))
				os.Exit(1)
			}

			data := &DNSDomainsPrinter{Domains: dms, Meta: meta}
			o.Base.Printer.Display(data, nil)
		},
	}

	domainList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	domainList.Flags().IntP("per-page", "p", utils.PerPageDefault, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

	// Domain Get
	domainGet := &cobra.Command{
		Use:   "get <Domain Name>",
		Short: "Get a domain",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			dm, err := o.DomainGet()
			if err != nil {
				printer.Error(fmt.Errorf("error retrieving domain : %v", err))
				os.Exit(1)
			}

			data := &DNSDomainPrinter{Domain: *dm}
			o.Base.Printer.Display(data, nil)
		},
	}

	// Domain Create
	domainCreate := &cobra.Command{
		Use:     "create",
		Short:   "Create a domain",
		Long:    createLong,
		Example: createExample,
		Run: func(cmd *cobra.Command, args []string) {
			domain, errDo := cmd.Flags().GetString("domain")
			if errDo != nil {
				printer.Error(fmt.Errorf("error parsing 'domain' flag for domain create : %v", errDo))
				os.Exit(1)
			}

			ip, errIP := cmd.Flags().GetString("ip")
			if errIP != nil {
				printer.Error(fmt.Errorf("error parsing 'ip' flag for domain create : %v", errIP))
				os.Exit(1)
			}

			o.DomainCreateReq = &govultr.DomainReq{
				Domain: domain,
				IP:     ip,
			}

			dm, err := o.DomainCreate()
			if err != nil {
				printer.Error(fmt.Errorf("error creating dns domain : %v", err))
				os.Exit(1)
			}

			data := &DNSDomainPrinter{Domain: *dm}
			o.Base.Printer.Display(data, nil)
		},
	}

	domainCreate.Flags().StringP("domain", "d", "", "name of the domain")
	if err := domainCreate.MarkFlagRequired("domain"); err != nil {
		printer.Error(fmt.Errorf("error marking domain create 'domain' flag required: %v", err))
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
		Run: func(cmd *cobra.Command, args []string) {
			if err := o.DomainDelete(); err != nil {
				printer.Error(fmt.Errorf("error delete dns domain : %v", err))
				os.Exit(1)
			}

			o.Base.Printer.Display(printer.Info("dns domain has been deleted"), nil)
		},
	}

	// Domain DNSSEC Update
	domainDNSSEC := &cobra.Command{
		Use:   "dnssec <Domain Name>",
		Short: "enable/disable dnssec",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a domain name")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			enabled, errEn := cmd.Flags().GetBool("enabled")
			if errEn != nil {
				printer.Error(fmt.Errorf("error parsing 'enabled' flag for dnssec : %v", errEn))
				os.Exit(1)
			}

			disabled, errDi := cmd.Flags().GetBool("disabled")
			if errEn != nil {
				printer.Error(fmt.Errorf("error parsing 'disabled' flag for dnssec : %v", errDi))
				os.Exit(1)
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

			if err := o.DomainUpdate(); err != nil {
				printer.Error(fmt.Errorf("error toggling dnssec : %v", err))
				os.Exit(1)
			}

			o.Base.Printer.Display(printer.Info("dns domain DNSSEC has been updated"), nil)
		},
	}

	domainDNSSEC.Flags().BoolP("enabled", "e", true, "enable dnssec")
	domainDNSSEC.Flags().BoolP("disabled", "d", true, "disable dnssec")
	domainDNSSEC.MarkFlagsOneRequired("enabled", "disabled")
	domainDNSSEC.MarkFlagsMutuallyExclusive("enabled", "disabled")

	// Domain DNSSEC Info
	domainDNSSECInfo := &cobra.Command{
		Use:   "dnssec-info <Domain Name>",
		Short: "Get DNS SEC info",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a domain name")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			info, err := o.DomainDNSSECGet()
			if err != nil {
				printer.Error(fmt.Errorf("error getting domain dnssec info : %v", err))
				os.Exit(1)
			}

			data := &DNSSECPrinter{SEC: info}
			o.Base.Printer.Display(data, nil)
		},
	}

	// Domain SOA Info
	domainSOAInfo := &cobra.Command{
		Use:   "soa-info <Domain Name>",
		Short: "Get DNS SOA info",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a domain name")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			info, err := o.DomainSOAGet()
			if err != nil {
				printer.Error(fmt.Errorf("error getting domain soa info : %v", err))
				os.Exit(1)
			}

			data := &DNSSOAPrinter{SOA: *info}
			o.Base.Printer.Display(data, nil)
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
		Run: func(cmd *cobra.Command, args []string) {
			ns, errNs := cmd.Flags().GetString("ns-primary")
			if errNs != nil {
				printer.Error(fmt.Errorf("error parsing 'ns-primary' flag for domain soa : %v", errNs))
				os.Exit(1)
			}

			email, errEm := cmd.Flags().GetString("email")
			if errEm != nil {
				printer.Error(fmt.Errorf("error parsing 'email' flag for domain soa : %v", errEm))
				os.Exit(1)
			}

			o.SOAUpdateReq = &govultr.Soa{
				NSPrimary: ns,
				Email:     email,
			}

			if err := o.DomainSOAUpdate(); err != nil {
				printer.Error(fmt.Errorf("error updating domain soa : %v", err))
				os.Exit(1)
			}

			o.Base.Printer.Display(printer.Info("domain soa has been updated"), nil)
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
		Short: "dns record",
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
		Run: func(cmd *cobra.Command, args []string) {
			o.Base.Options = utils.GetPaging(cmd)

			recs, meta, err := o.RecordList()
			if err != nil {
				printer.Error(fmt.Errorf("error retrieiving domain records : %v", err))
				os.Exit(1)
			}

			data := &DNSRecordsPrinter{Records: recs, Meta: meta}
			o.Base.Printer.Display(data, nil)
		},
	}

	recordList.Flags().StringP("cursor", "c", "", "(optional) Cursor for paging.")
	recordList.Flags().IntP("per-page", "p", utils.PerPageDefault, "(optional) Number of items requested per page. Default is 100 and Max is 500.")

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
		Run: func(cmd *cobra.Command, args []string) {
			rec, err := o.RecordGet()
			if err != nil {
				printer.Error(fmt.Errorf("error while getting domain record : %v", err))
				os.Exit(1)
			}

			data := &DNSRecordPrinter{Record: *rec}
			o.Base.Printer.Display(data, nil)
		},
	}

	// Record Create
	recordCreate := &cobra.Command{
		Use:   "create <Domain Name>",
		Short: "Create a dns record",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide a domain name")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			rType, errTy := cmd.Flags().GetString("type")
			if errTy != nil {
				printer.Error(fmt.Errorf("error parsing 'type' flag for domain record create : %v", errTy))
				os.Exit(1)
			}

			name, errNa := cmd.Flags().GetString("name")
			if errNa != nil {
				printer.Error(fmt.Errorf("error parsing 'name' flag for domain record create : %v", errNa))
				os.Exit(1)
			}

			dt, errDa := cmd.Flags().GetString("data")
			if errDa != nil {
				printer.Error(fmt.Errorf("error parsing 'data' flag for domain record create : %v", errDa))
				os.Exit(1)
			}

			ttl, errTt := cmd.Flags().GetInt("ttl")
			if errTt != nil {
				printer.Error(fmt.Errorf("error parsing 'ttl' flag for domain record create : %v", errTt))
				os.Exit(1)
			}

			priority, errPr := cmd.Flags().GetInt("priority")
			if errPr != nil {
				printer.Error(fmt.Errorf("error parsing 'priority' flag for domain record create : %v", errPr))
				os.Exit(1)
			}

			o.RecordReq = &govultr.DomainRecordReq{
				Name:     name,
				Type:     rType,
				Data:     dt,
				TTL:      ttl,
				Priority: &priority,
			}

			rec, err := o.RecordCreate()
			if err != nil {
				printer.Error(fmt.Errorf("error creating domain record : %v", err))
				os.Exit(1)
			}

			data := &DNSRecordPrinter{Record: *rec}
			o.Base.Printer.Display(data, nil)
		},
	}

	recordCreate.Flags().StringP("type", "t", "", "type for the record")
	if err := recordCreate.MarkFlagRequired("type"); err != nil {
		printer.Error(fmt.Errorf("error marking dns record create 'type' flag required: %v", err))
		os.Exit(1)
	}

	recordCreate.Flags().StringP("name", "n", "", "name of the record")
	if err := recordCreate.MarkFlagRequired("name"); err != nil {
		printer.Error(fmt.Errorf("error marking dns record create 'name' flag required: %v", err))
		os.Exit(1)
	}

	recordCreate.Flags().StringP("data", "d", "", "data for the record")
	if err := recordCreate.MarkFlagRequired("data"); err != nil {
		printer.Error(fmt.Errorf("error marking dns record create 'data' flag required: %v", err))
		os.Exit(1)
	}

	recordCreate.Flags().IntP("ttl", "l", 0, "ttl for the record")
	if err := recordCreate.MarkFlagRequired("ttl"); err != nil {
		printer.Error(fmt.Errorf("error marking dns record create 'ttl' flag required: %v", err))
		os.Exit(1)
	}

	recordCreate.Flags().IntP("priority", "p", 0, "only required for MX and SRV")
	if err := recordCreate.MarkFlagRequired("priority"); err != nil {
		printer.Error(fmt.Errorf("error marking dns record create 'priority' flag required: %v", err))
		os.Exit(1)
	}

	// Record Delete
	recordDelete := &cobra.Command{
		Use:     "delete <Domain Name> <Record ID>",
		Short:   "Delete DNS record",
		Aliases: []string{"destroy"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("please provide a domain name & record ID")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := o.RecordDelete(); err != nil {
				printer.Error(fmt.Errorf("error deleting domain record : %v", err))
				os.Exit(1)
			}

			o.Base.Printer.Display(printer.Info("domain record has been deleted"), nil)
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
		Run: func(cmd *cobra.Command, args []string) {
			name, errNa := cmd.Flags().GetString("name")
			if errNa != nil {
				printer.Error(fmt.Errorf("error parsing 'name' flag for domain record update : %v", errNa))
				os.Exit(1)
			}

			dt, errDa := cmd.Flags().GetString("data")
			if errDa != nil {
				printer.Error(fmt.Errorf("error parsing 'data' flag for domain record update : %v", errDa))
				os.Exit(1)
			}

			ttl, errTt := cmd.Flags().GetInt("ttl")
			if errTt != nil {
				printer.Error(fmt.Errorf("error parsing 'ttl' flag for domain record update : %v", errTt))
				os.Exit(1)
			}

			priority, errPr := cmd.Flags().GetInt("priority")
			if errPr != nil {
				printer.Error(fmt.Errorf("error parsing 'priority' flag for domain record update : %v", errPr))
				os.Exit(1)
			}

			o.RecordReq = &govultr.DomainRecordReq{}

			if cmd.Flags().Changed("name") {
				o.RecordReq.Name = name
			}

			if cmd.Flags().Changed("data") {
				o.RecordReq.Data = dt
			}

			if cmd.Flags().Changed("ttl") {
				o.RecordReq.TTL = ttl
			}

			if cmd.Flags().Changed("priority") {
				o.RecordReq.Priority = govultr.IntToIntPtr(priority)
			}

			if err := o.RecordUpdate(); err != nil {
				printer.Error(fmt.Errorf("error updating domain record : %v", errPr))
				os.Exit(1)
			}

			o.Base.Printer.Display(printer.Info("domain record has been updated"), nil)
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

// DomainList ...
func (o *Options) DomainList() ([]govultr.Domain, *govultr.Meta, error) {
	dms, meta, _, err := o.Base.Client.Domain.List(o.Base.Context, o.Base.Options)
	return dms, meta, err
}

// DomainGet ...
func (o *Options) DomainGet() (*govultr.Domain, error) {
	dm, _, err := o.Base.Client.Domain.Get(o.Base.Context, o.Base.Args[0])
	return dm, err
}

// DomainCreate ...
func (o *Options) DomainCreate() (*govultr.Domain, error) {
	dm, _, err := o.Base.Client.Domain.Create(o.Base.Context, o.DomainCreateReq)
	return dm, err
}

// DomainUpdate ...
func (o *Options) DomainUpdate() error {
	return o.Base.Client.Domain.Update(o.Base.Context, o.Base.Args[0], o.DomainDNSSECEnabled)
}

// DomainDelete ...
func (o *Options) DomainDelete() error {
	return o.Base.Client.Domain.Delete(o.Base.Context, o.Base.Args[0])
}

// DomainDNSSECGet ...
func (o *Options) DomainDNSSECGet() ([]string, error) {
	sec, _, err := o.Base.Client.Domain.GetDNSSec(o.Base.Context, o.Base.Args[0])
	return sec, err
}

// DomainSOAGet ...
func (o *Options) DomainSOAGet() (*govultr.Soa, error) {
	soa, _, err := o.Base.Client.Domain.GetSoa(o.Base.Context, o.Base.Args[0])
	return soa, err
}

// DomainSOAUpdate ...
func (o *Options) DomainSOAUpdate() error {
	return o.Base.Client.Domain.UpdateSoa(o.Base.Context, o.Base.Args[0], o.SOAUpdateReq)
}

// RecordList ...
func (o *Options) RecordList() ([]govultr.DomainRecord, *govultr.Meta, error) {
	rec, meta, _, err := o.Base.Client.DomainRecord.List(o.Base.Context, o.Base.Args[0], o.Base.Options)
	return rec, meta, err
}

// RecordGet ...
func (o *Options) RecordGet() (*govultr.DomainRecord, error) {
	rec, _, err := o.Base.Client.DomainRecord.Get(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
	return rec, err
}

// RecordCreate ...
func (o *Options) RecordCreate() (*govultr.DomainRecord, error) {
	rec, _, err := o.Base.Client.DomainRecord.Create(o.Base.Context, o.Base.Args[0], o.RecordReq)
	return rec, err
}

// RecordUpdate ...
func (o *Options) RecordUpdate() error {
	return o.Base.Client.DomainRecord.Update(o.Base.Context, o.Base.Args[0], o.Base.Args[1], o.RecordReq)
}

// RecordDelete ...
func (o *Options) RecordDelete() error {
	return o.Base.Client.DomainRecord.Delete(o.Base.Context, o.Base.Args[0], o.Base.Args[1])
}
