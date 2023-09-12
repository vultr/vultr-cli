package printer

import "github.com/vultr/govultr/v3"

func DnsRecordsList(records []govultr.DomainRecord, meta *govultr.Meta) {
	col := columns{"ID", "TYPE", "NAME", "DATA", "PRIORITY", "TTL"}
	display(col)
	for _, r := range records {
		display(columns{r.ID, r.Type, r.Name, r.Data, r.Priority, r.TTL})
	}

	Meta(meta)
	flush()
}

func DnsRecord(record *govultr.DomainRecord) {
	col := columns{"ID", "TYPE", "NAME", "DATA", "PRIORITY", "TTL"}
	display(col)

	display(columns{record.ID, record.Type, record.Name, record.Data, record.Priority, record.TTL})
	flush()
}
