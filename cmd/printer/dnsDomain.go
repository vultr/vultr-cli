package printer

import "github.com/vultr/govultr/v3"

func SecInfo(info []string) {
	col := columns{"DNSSEC INFO"}
	display(col)
	for _, i := range info {
		display(columns{i})
	}
	flush()
}

func DomainList(domain []govultr.Domain, meta *govultr.Meta) {
	col := columns{"DOMAIN", "DATE CREATED", "DNS SEC"}
	display(col)
	for _, d := range domain {
		display(columns{d.Domain, d.DateCreated, d.DNSSec})
	}

	Meta(meta)
	flush()
}

func Domain(domain *govultr.Domain) {
	col := columns{"DOMAIN", "DATE CREATED", "DNS SEC"}
	display(col)
	display(columns{domain.Domain, domain.DateCreated, domain.DNSSec})

	flush()
}

func SoaInfo(soa *govultr.Soa) {
	col := columns{"NS PRIMARY", "EMAIL"}
	display(col)
	display(columns{soa.NSPrimary, soa.Email})
	flush()
}
