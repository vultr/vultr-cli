package printer

import "github.com/vultr/govultr"

func DnsRecordsList(records []govultr.DNSRecord) {
	col := columns{"ID", "TYPE", "NAME", "DATA", "PRIORITY", "TTL"}
	display(col)
	for _, r := range records {
		display(columns{r.RecordID, r.Type, r.Name, r.Data, r.Priority, r.TTL})
	}
	flush()
}
