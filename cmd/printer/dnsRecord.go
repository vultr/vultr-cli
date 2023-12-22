package printer

import "github.com/vultr/govultr/v3"

func DNSRecordsList(records []govultr.DomainRecord, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "TYPE", "NAME", "DATA", "PRIORITY", "TTL"})

	if len(records) == 0 {
		display(columns{"---", "---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range records {
		display(columns{
			records[i].ID,
			records[i].Type,
			records[i].Name,
			records[i].Data,
			records[i].Priority,
			records[i].TTL,
		})
	}

	Meta(meta)
}

func DNSRecord(record *govultr.DomainRecord) {
	defer flush()

	display(columns{"ID", "TYPE", "NAME", "DATA", "PRIORITY", "TTL"})
	display(columns{record.ID, record.Type, record.Name, record.Data, record.Priority, record.TTL})
}
